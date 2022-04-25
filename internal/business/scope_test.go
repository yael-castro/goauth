package business

import (
	"errors"
	"github.com/yael-castro/goauth/internal/model"
	"reflect"
	"strconv"
	"testing"
)

// TestMaskParser_ParseScope check the correct functionality of the MaskParser structure and more specifically the method ParseScope
func TestMaskParser_ParseScope(t *testing.T) {
	parser := NewScopeParser()

	tdt := []struct {
		input       string
		output      interface{}
		expectedErr error
	}{
		// Empty scope
		{
			input: "",
		},
		// Success build scope
		{
			input: "read:aa write:BB update:cc delete:DD",
			output: model.Mask{
				"read":   0xAA,
				"write":  0xbb,
				"update": 0xCC,
				"delete": 0xdd,
			},
		},
		// Malformed scope
		{
			input:       "r:1_w:2_u:3_d:4",
			expectedErr: model.InvalidScope,
		},
		// Invalid scope
		{
			input:       "r:w",
			expectedErr: model.InvalidScope,
		},
		// Scope out of range
		{
			input:       "r:80000000000000001",
			expectedErr: model.InvalidScope,
		},
	}

	for i, v := range tdt {
		t.Run(strconv.Itoa(i+1), func(t *testing.T) {
			scope, err := parser.ParseScope(v.input)
			if !errors.Is(err, v.expectedErr) {
				t.Fatalf(`unexpected error "%v"`, err)
			}

			if err != nil {
				t.Skip(err)
			}

			if !reflect.DeepEqual(v.output, scope) {
				t.Fatalf(`expected scope "%+v" got "%+v"`, v.output, scope)
			}

			t.Log(scope, uint(1<<63))
		})
	}
}
