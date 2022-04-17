package business

import (
	"fmt"
	"github.com/yael-castro/godi/internal/model"
	"strconv"
	"strings"
)

type ScopeParser interface {
	// ParseScope validates and returns the scope build from the string
	ParseScope(string) (interface{}, error)
}

func NewScopeParser() ScopeParser {
	return maskParser{}
}

// maskParser parse scopes (permissions) to a mask
type maskParser struct{}

// ParseScope parses a string of hexadecimal values split by spaces to hash map where each key of hash map contains
// unsigned integer values of 64 bits to be used as bit masks
func (m maskParser) ParseScope(str string) (interface{}, error) {
	if str == "" {
		return nil, nil
	}

	slice := strings.Split(str, " ")

	mask := model.Mask{}

	var err error

	for _, v := range slice {
		slice := strings.Split(v, ":")

		if len(slice) != 2 {
			return nil, fmt.Errorf("%w: malformed permission requested", model.InvalidScope)
		}

		mask[slice[0]], err = strconv.ParseUint(slice[1], 16, 64)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", model.InvalidScope, err.Error())
		}
	}

	return mask, nil
}
