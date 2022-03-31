package repository

import (
	"errors"
	"github.com/yael-castro/godi/internal/model"
	"net/url"
	"reflect"
	"strconv"
	"testing"
)

func TestStateStorage_Create(t *testing.T) {
	tdt := []struct {
		model.Authorization
		expectedErr error
	}{
		{
			Authorization: model.Authorization{
				State:               "ABC",
				ClientId:            "3aad9943-714d-4576-9c6f-bb45b142666c",
				Scope:               "http://localhost/private/,http://localhost/private2/",
				ResponseType:        "code",
				CodeChallenge:       "FEDCBA",
				CodeChallengeMethod: "PLAIN",
				RedirectURL: func() *url.URL {
					uri, _ := url.Parse("http://localhost/callback")
					return uri
				}(),
			},
		},
		{
			Authorization: model.Authorization{
				State: "ABC",
				RedirectURL: func() *url.URL {
					uri, _ := url.Parse("http://localhost/callback")
					return uri
				}(),
			},
		},
	}

	// Here starts unit tests
	client, err := NewRedisClient(defaultRedisConfiguration)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		client.Close()
	})

	storage := StateStorage{Client: client}

	for i, v := range tdt {
		t.Run(strconv.Itoa(i+1), func(t *testing.T) {
			t.Cleanup(func() {
				storage.Delete(v.Authorization.State)
			})

			err := storage.Create(v.Authorization)
			if !errors.Is(err, v.expectedErr) {
				t.Fatal(err)
			}

			if err != nil {
				t.Skip(err)
			}

			t.Logf("%+v", v.Authorization)
		})
	}
}

func TestStateStorage_Obtain(t *testing.T) {
	tdt := []struct {
		expectedData model.Authorization
		expectedErr  error
	}{
		{
			expectedData: model.Authorization{
				State:               "1234",
				ClientId:            "3aad9943-714d-4576-9c6f-bb45b142666c",
				Scope:               "http://localhost/private/,http://localhost/private2/",
				ResponseType:        "code",
				CodeChallenge:       "abc",
				CodeChallengeMethod: "PLAIN",
				RedirectURL: func() *url.URL {
					uri, _ := url.Parse("http://localhost/callback")
					return uri
				}(),
			},
		},
	}

	// Here starts unit tests
	client, err := NewRedisClient(defaultRedisConfiguration)
	if err != nil {
		t.Fatal(err)
	}

	storage := StateStorage{Client: client}

	for _, v := range tdt {
		err := storage.Create(v.expectedData)
		if err != nil {
			t.Fatal(err)
		}
	}

	t.Cleanup(func() {
		for _, v := range tdt {
			storage.Delete(v.expectedData.State)
		}
	})

	for i, v := range tdt {
		t.Run(strconv.Itoa(i+1), func(t *testing.T) {
			gotData, err := storage.Obtain(v.expectedData.State)
			if !errors.Is(err, v.expectedErr) {
				t.Fatal(err)
			}

			if err != nil {
				t.Skip(err)
			}

			if !reflect.DeepEqual(gotData, v.expectedData) {
				t.Error(`mismatch expected data`)
				t.Error(`got data`)
				t.Errorf(`%+v`, gotData)
				t.Error(`expected data`)
				t.Errorf("%+v", v.expectedData)
				t.Fatal()
			}

			t.Logf("%+v", gotData)
		})
	}
}

func TestStateStorage_Delete(t *testing.T) {
	tdt := []struct {
		state       string
		expectedErr error
	}{
		{state: "abc"},
		{state: "abc"},
	}

	// Here starts unit tests
	client, err := NewRedisClient(defaultRedisConfiguration)
	if err != nil {
		t.Fatal(err)
	}

	storage := StateStorage{Client: client}

	for i, v := range tdt {
		t.Run(strconv.Itoa(i+1), func(t *testing.T) {
			err := storage.Delete(v.state)
			if !errors.Is(err, v.expectedErr) {
				t.Fatal(err)
			}
		})
	}
}
