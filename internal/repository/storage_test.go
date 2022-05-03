package repository

import (
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/yael-castro/goauth/internal/model"
	"reflect"
	"strconv"
	"testing"
)

// StorageCase test case for Storage interface
type StorageCase[K Key, V any] struct {
	Id    K
	Input V
	// expectedErr expected error
	expectedErr error
	// skipSetup indicates if the initial configuration must be skipped,
	// configurations like a pre-creation of records
	skipSetup bool
}

// StorageTest test for Storage interface
type StorageTest[K Key, V any] struct {
	// Storage interface instance to test
	Storage[K, V]
	// Cases test cases for integration or unit tests
	Cases []StorageCase[K, V]
}

// TestStorage_Create check if the method Create of different implementations of Storage
// interface works correctly, using the function testStorageCreate
// TODO: check the creation of duplicate records
func TestStorage_Create(t *testing.T) {
	redisClient, err := NewRedisClient(defaultRedisConfiguration)
	if err != nil {
		t.Fatal(err)
	}

	// Testing OwnerStorage
	testStorageCreate(t, StorageTest[string, model.Owner]{
		Storage: OwnerStorage{
			Client: redisClient,
		},
		Cases: []StorageCase[string, model.Owner]{
			{
				Id: "abc",
				Input: model.Owner{
					Id:       "abc",
					Password: "1234",
				},
			},
		},
	})

	// Testing SessionStorage
	testStorageCreate(t, StorageTest[string, model.Session]{
		Storage: SessionStorage{
			Client: redisClient,
		},
		Cases: []StorageCase[string, model.Session]{
			{
				Id: "abc",
				Input: model.Session{
					UserAgent: "Go/1.18",
				},
			},
		},
	})

	// Testing SessionStorage
	testStorageCreate(t, StorageTest[string, model.Authorization]{
		Storage: StateStorage{
			Client: redisClient,
		},
		Cases: []StorageCase[string, model.Authorization]{
			{
				Id: "abc",
				Input: model.Authorization{
					ResponseType: "code",
				},
			},
		},
	})
}

// TestStorage_Obtain check if the method Obtain of different implementations of Storage
// interface works correctly, using the function testStorageObtain
func TestStorage_Obtain(t *testing.T) {
	redisClient, err := NewRedisClient(defaultRedisConfiguration)
	if err != nil {
		t.Fatal(err)
	}

	// Testing OwnerStorage
	testStorageObtain(t, StorageTest[string, model.Owner]{
		Storage: OwnerStorage{
			Client: redisClient,
		},
		Cases: []StorageCase[string, model.Owner]{
			{
				Id: "abc",
				Input: model.Owner{
					Id:       "abc",
					Password: "1234",
				},
			},
			{
				Id:          "abc",
				skipSetup:   true,
				expectedErr: redis.Nil,
			},
		},
	})

	// Testing SessionStorage
	testStorageObtain(t, StorageTest[string, model.Session]{
		Storage: SessionStorage{Client: redisClient},
		Cases: []StorageCase[string, model.Session]{
			{
				Id: "abc",
				Input: model.Session{
					UserAgent: "Go/1.17",
				},
			},
			{
				Id:          "abc",
				skipSetup:   true,
				expectedErr: redis.Nil,
			},
		},
	})

	// Testing SessionStorage
	testStorageObtain(t, StorageTest[string, model.Authorization]{
		Storage: StateStorage{Client: redisClient},
		Cases: []StorageCase[string, model.Authorization]{
			{
				Id: "abc",
				Input: model.Authorization{
					ResponseType: "code",
				},
			},
			{
				Id:          "abc",
				skipSetup:   true,
				expectedErr: redis.Nil,
			},
		},
	})
}

// TestStorage_Delete check if the method Delete of different implementations of Storage
// interface works correctly, using the function testStorageDelete
func TestStorage_Delete(t *testing.T) {
	redisClient, err := NewRedisClient(defaultRedisConfiguration)
	if err != nil {
		t.Fatal(err)
	}

	// Testing OwnerStorage
	testStorageDelete(t, StorageTest[string, model.Owner]{
		Storage: OwnerStorage{Client: redisClient},
		Cases: []StorageCase[string, model.Owner]{
			{
				Id: "abc",
				Input: model.Owner{
					Id:       "abc",
					Password: "1234",
				},
			},
			{
				Id:        "abc",
				skipSetup: true,
			},
		},
	})

	// Testing SessionStorage
	testStorageDelete(t, StorageTest[string, model.Session]{
		Storage: SessionStorage{Client: redisClient},
		Cases: []StorageCase[string, model.Session]{
			{
				Id: "abc",
				Input: model.Session{
					UserAgent: "Go/1.18",
				},
			},
			{
				Id:        "abc",
				skipSetup: true,
			},
		},
	})

	// Testing StateStorage
	testStorageDelete(t, StorageTest[string, model.Authorization]{
		Storage: StateStorage{Client: redisClient},
		Cases: []StorageCase[string, model.Authorization]{
			{
				Id: "abc",
				Input: model.Authorization{
					ResponseType: "code",
				},
			},
			{
				Id:        "abc",
				skipSetup: true,
			},
		},
	})
}

// testStorageCreate util function to test the method Create of the Storage interface
//
// Notes:
//
// - Before starting each test case, a temporary record is created unless otherwise
//   indicated with a false value in the skipSetup field.
//
// - Always to end each test case the temporary record is deleted
//
func testStorageCreate[K Key, V any](t *testing.T, test StorageTest[K, V]) {
	t.Run(reflect.TypeOf(test.Storage).String(), func(t *testing.T) {
		for i, v := range test.Cases {
			t.Run(strconv.Itoa(i+1), func(t *testing.T) {
				t.Logf(`Creating record with id "%v"`, v.Id)
				err := test.Storage.Create(v.Id, v.Input)
				if !errors.Is(err, v.expectedErr) {
					t.Fatalf(`Unexpected error "%v"`, err)
				}

				t.Cleanup(func() {
					t.Logf(`Deleting record with id "%v"`, v.Id)
					test.Storage.Delete(v.Id)
				})
			})
		}
	})
}

// testStorageObtain test the method Obtain of the Storage interface using the received implementation of the Storage interface
// and the test cases
//
// Notes:
//
// - Before starting each test case, a temporary record is created unless otherwise
//   indicated with a false value in the skipSetup field.
//
// - Always to end each test case the temporary record is deleted
//
func testStorageObtain[K Key, V any](t *testing.T, test StorageTest[K, V]) {
	t.Run(reflect.TypeOf(test.Storage).String(), func(t *testing.T) {
		for i, v := range test.Cases {
			t.Run(strconv.Itoa(i+1), func(t *testing.T) {
				if !v.skipSetup {
					t.Logf(`Creating record with id "%v"`, v.Id)
					_ = test.Storage.Create(v.Id, v.Input)
				}

				output, err := test.Obtain(v.Id)
				if !errors.Is(err, v.expectedErr) {
					t.Fatalf(`"Error "%v" was expected but "%v" was obtained`, v.expectedErr, err)
				}

				if err != nil {
					t.Skip(err)
				}

				if !reflect.DeepEqual(output, v.Input) {
					t.Fatalf(`"%v" was expected but "%v" was obtained`, output, v.Input)
				}

				t.Logf("%+v", output)
				t.Cleanup(func() {
					t.Logf(`Deleting record with id "%v"`, v.Id)
					test.Storage.Delete(v.Id)
				})
			})
		}
	})
}

// testStorageDelete test the method Delete of the Storage interface
//
// Notes:
//
// - Before starting each test case, a temporary record is created unless otherwise
//   indicated with a false value in the skipSetup field.
//
// - Always to end each test case the temporary record is deleted
//
func testStorageDelete[K Key, V any](t *testing.T, test StorageTest[K, V]) {
	t.Run(reflect.TypeOf(test.Storage).String(), func(t *testing.T) {
		for i, v := range test.Cases {
			t.Run(strconv.Itoa(i+1), func(t *testing.T) {
				t.Logf(`Creating record with id "%v"`, v.Id)
				if !v.skipSetup {
					_ = test.Storage.Create(v.Id, v.Input)
				}

				t.Logf(`Deleting record with id "%v"`, v.Id)
				err := test.Delete(v.Id)
				if !errors.Is(err, v.expectedErr) {
					t.Fatalf(`"Error "%v" was expected but "%v" was obtained`, v.expectedErr, err)
				}

				//t.Cleanup(func() {
				//	t.Logf(`Deleting record with id "%v"`, v.Id)
				//	test.Storage.Delete(v.Id)
				//})
			})
		}
	})
}
