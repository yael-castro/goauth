package repository

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/yael-castro/godi/internal/model"
	"reflect"
	"strconv"
	"testing"
)

func TestClientFinder_Find(t *testing.T) {
	tdt := []struct {
		clientId           string
		expectedClient     model.Client
		expectedErr        error
		skipInitialization bool
	}{
		{
			clientId: "abc",
			expectedClient: model.Client{
				Id: "abc",
				AllowedOrigins: []string{
					"http://localhost:8080",
					"http://localhost",
				},
			},
		},
		{
			clientId: "xyz",
			expectedClient: model.Client{
				Id: "xyz",
				AllowedOrigins: []string{
					"http://localhost:8080",
					"http://localhost",
				},
			},
			expectedErr:        redis.Nil,
			skipInitialization: true,
		},
	}

	client, err := NewRedisClient(defaultRedisConfiguration)
	if err != nil {
		t.Fatal(err)
	}

	finder := ClientFinder{client}

	for _, v := range tdt {
		if v.skipInitialization {
			continue
		}

		cmd := client.Set(
			context.TODO(),
			finder.secretKey(v.clientId),
			v.expectedClient.Secret,
			0,
		)

		if err := cmd.Err(); err != nil {
			t.Fatal(err)
		}

		intCmd := client.LPush(
			context.TODO(),
			finder.listKey(v.clientId),
			v.expectedClient.AllowedOrigins,
		)

		if err := intCmd.Err(); err != nil {
			t.Fatal(err)
		}
	}

	t.Cleanup(func() {
		for _, v := range tdt {
			client.Del(context.TODO(), finder.secretKey(v.clientId))
			client.Del(context.TODO(), finder.listKey(v.clientId))
		}

		_ = client.Close()
	})

	for i, v := range tdt {
		t.Run(strconv.Itoa(i+1), func(t *testing.T) {
			gotData, err := finder.Find(v.clientId)
			if !errors.Is(err, v.expectedErr) {
				t.Fatalf(`expected error "%v" got "%v""`, v.expectedErr, err)
			}

			if err != nil {
				t.Skip(err)
			}

			gotClient := gotData.(model.Client)

			if compareArrays(v.expectedClient.AllowedOrigins, gotClient.AllowedOrigins) {
				v.expectedClient.AllowedOrigins = nil
				gotClient.AllowedOrigins = nil
			}

			if !reflect.DeepEqual(v.expectedClient, gotClient) {
				t.Fatalf(`expected data "%v" got "%v"`, v.expectedClient, gotClient)
			}

			t.Logf("%+v", gotData)
		})
	}
}

func compareArrays(arr1, arr2 []string) bool {
	values := make(map[string]struct{})

	for _, v := range arr1 {
		values[v] = struct{}{}
	}

	counter := 0

	for _, v := range arr2 {
		if _, ok := values[v]; ok {
			counter++
		}
	}

	return counter == len(arr1)
}
