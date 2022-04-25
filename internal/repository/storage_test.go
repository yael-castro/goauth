package repository

import (
	"errors"
	"github.com/yael-castro/goauth/internal/model"
	"net/url"
	"reflect"
	"strconv"
	"testing"
)

// testCase is a common test case for crud operations
type testCase struct {
	// id record identifier
	id string
	// input record to be created or expected data
	input interface{}
	// expectedErr error that is expected
	expectedErr error
}

// TestStorage_Create
// Check the functionality of the Create method for different Storage interface implementations (StateStorage and OwnerStorage)
// TODO validate the duplicated records
func TestStorage_Create(t *testing.T) {
	client, err := NewRedisClient(defaultRedisConfiguration)
	if err != nil {
		t.Fatal(err)
	}

	tdt := []struct {
		storage Storage
		tests   []testCase
	}{
		{
			storage: StateStorage{Client: client},
			tests: []testCase{
				{
					id: "abc",
					input: model.Authorization{
						State:               "ABC",
						ClientId:            "3aad9943-714d-4576-9c6f-bb45b142666c",
						Scope:               "http://localhost/private/,http://localhost/private2/",
						ResponseType:        "code",
						CodeChallenge:       "FFF",
						CodeChallengeMethod: "PLAIN",
						RedirectURL: func() *url.URL {
							uri, _ := url.Parse("http://localhost/callback")
							return uri
						}(),
					},
				},
				{
					id: "xyz",
					input: model.Authorization{
						State: "ABC",
						RedirectURL: func() *url.URL {
							uri, _ := url.Parse("http://localhost/callback")
							return uri
						}(),
					},
				},
			},
		},
		{
			storage: OwnerStorage{Client: client},
			tests: []testCase{
				{id: "abc", input: model.Owner{Id: "abc", Password: "xyz"}},
				{id: "xyz", input: model.Owner{Id: "xyz", Password: "abc"}},
			},
		},
		{
			storage: SessionStorage{Client: client},
			tests: []testCase{
				{
					id: "abc",
					input: model.Session{
						UserAgent: "Go/1.17",
						Owner: model.Owner{
							Id: "Golang",
						},
					},
				},
				{
					id: "xyz",
					input: model.Session{
						UserAgent: "Go/1.17",
						Owner: model.Owner{
							Id: "Go",
						},
					},
				},
			},
		},
	}

	t.Cleanup(func() {
		_ = client.Close()
	})

	// Here starts the unit tests
	for _, v := range tdt {
		storage := v.storage
		t.Run(reflect.TypeOf(v.storage).String(), func(t *testing.T) {
			for i, v := range v.tests {
				t.Run(strconv.Itoa(i+1), func(t *testing.T) {
					t.Cleanup(func() {
						_ = storage.Delete(v.id)
					})

					err := storage.Create(v.id, v.input)
					if !errors.Is(err, v.expectedErr) {
						t.Fatal(err)
					}

					if err != nil {
						t.Skip(err)
					}

					t.Logf("%+v", v.input)
				})
			}
		})
	}
}

// TestStorage_Obtain
// Checks the functionality of the Obtain method for different Storage interface implementations (StateStorage and OwnerStorage)
func TestStorage_Obtain(t *testing.T) {
	client, err := NewRedisClient(defaultRedisConfiguration)
	if err != nil {
		t.Fatal(err)
	}

	tdt := []struct {
		storage Storage
		tests   []testCase
	}{
		{
			storage: StateStorage{Client: client},
			tests: []testCase{
				{
					id: "abc",
					input: model.Authorization{
						State:               "",
						ClientId:            "def",
						ClientSecret:        "def",
						Scope:               "qwerty",
						ResponseType:        "code",
						CodeChallenge:       "abc",
						CodeChallengeMethod: "S256",
					},
				},
			},
		},
		{
			storage: OwnerStorage{Client: client},
			tests: []testCase{
				{
					id: "xyz",
					input: model.Owner{
						Id:       "xyz",
						Password: "foo",
					},
				},
			},
		},
		{
			storage: SessionStorage{Client: client},
			tests: []testCase{
				{
					id: "",
					input: model.Session{
						UserAgent: "Go/1.17",
						Owner:     model.Owner{Id: "foo"},
					},
				},
			},
		},
	}

	// Here starts unit tests
	for _, v := range tdt {
		storage := v.storage
		t.Run(reflect.TypeOf(storage).String(), func(t *testing.T) {

			for i, v := range v.tests {
				err := storage.Create(v.id, v.input)
				if err != nil {
					t.Fatal(err)
				}

				t.Cleanup(func() {
					_ = storage.Delete(v.id)
				})

				t.Run(strconv.Itoa(i+1), func(t *testing.T) {
					gotData, err := storage.Obtain(v.id)
					if !errors.Is(err, v.expectedErr) {
						t.Fatal(err)
					}

					if err != nil {
						t.Skip(err)
					}

					if !reflect.DeepEqual(gotData, v.input) {
						t.Error(`mismatch expected data`)
						t.Error(`got data`)
						t.Errorf(`%+v`, gotData)
						t.Error(`expected data`)
						t.Errorf("%+v", v.input)
						t.Fatal()
					}

					t.Logf("%+v", gotData)
				})
			}
		})
	}

}

// TestStorage_Delete
// Checks the functionality of the Delete method for different Storage interface implementations (StateStorage and OwnerStorage)
// TODO makes a better testing for each implementation
func TestStorage_Delete(t *testing.T) {
	client, err := NewRedisClient(defaultRedisConfiguration)
	if err != nil {
		t.Fatal(err)
	}

	tdt := []struct {
		input       string
		expectedErr error
	}{
		{input: "abc"},
		{input: "abc"},
	}

	t.Cleanup(func() {
		_ = client.Close()
	})

	storages := []Storage{
		StateStorage{Client: client},
		OwnerStorage{Client: client},
		SessionStorage{Client: client},
	}

	// Here starts unit tests
	for _, storage := range storages {
		t.Run(reflect.TypeOf(storage).String(), func(t *testing.T) {
			for i, v := range tdt {
				t.Run(strconv.Itoa(i+1), func(t *testing.T) {
					err := storage.Delete(v.input)
					if !errors.Is(err, v.expectedErr) {
						t.Fatal(err)
					}
				})
			}
		})
	}
}
