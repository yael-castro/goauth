package business

import (
	"errors"
	"github.com/yael-castro/godi/internal/model"
	"reflect"
	"strconv"
	"testing"
)

// TestMaskParser_ParseScope check the correct functionality of the MaskParser structure and more specifically the method ParseScope
func TestMaskParser_ParseScope(t *testing.T) {
	parser := MaskParser{}

	tdt := []struct {
		input       string
		output      interface{}
		expectedErr error
	}{
		{
			input: "read:aa write:BB update:cc delete:DD",
			output: model.Mask{
				"read":   0xAA,
				"write":  0xbb,
				"update": 0xCC,
				"delete": 0xdd,
			},
		},
		{
			input:       "r:1_w:2_u:3_d:4",
			expectedErr: model.ValidationError("malformed permission requested"),
		},
		{
			input:       "r:w",
			expectedErr: model.ValidationError(`strconv.ParseUint: parsing "w": invalid syntax`),
		},
		{
			input:       "r:80000000000000001",
			expectedErr: model.ValidationError(`strconv.ParseUint: parsing "80000000000000001": value out of range`),
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
