package business

import (
	"github.com/yael-castro/godi/internal/model"
	"strconv"
	"strings"
)

type ScopeParser interface {
	// ParseScope validates and returns the scope build from the string
	ParseScope(string) (interface{}, error)
}

// _ "implement" constraint for ScopeParser
var _ ScopeParser = MaskParser{}

// MaskParser parse scopes (permissions) to a mask
type MaskParser struct{}

// ParseScope parses a string of hexadecimal values split by spaces to hash map where each key of hash map contains
// unsigned integer values of 64 bits to be used as bit masks
func (m MaskParser) ParseScope(str string) (interface{}, error) {
	slice := strings.Split(str, " ")

	mask := model.Mask{}

	var err error

	for _, v := range slice {
		slice := strings.Split(v, ":")

		if len(slice) != 2 {
			return nil, model.ValidationError("malformed permission requested")
		}

		mask[slice[0]], err = strconv.ParseUint(slice[1], 16, 64)
		if err != nil {
			return nil, model.ValidationError(err.Error())
		}
	}

	return mask, nil
}
