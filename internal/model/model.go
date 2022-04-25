package model

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type Map = map[string]interface{}

// BinaryJSON json serializer that implements encoding.BinaryMarshaler
type BinaryJSON struct {
	// I embed data to later be serialized
	I interface{}
}

// MarshalBinary returns the serialized data in json encode of I
func (b BinaryJSON) MarshalBinary() ([]byte, error) {
	return json.Marshal(b.I)
}

// Mask defines a mask that contains multiple bit masks
type Mask map[string]uint64

// NewIP constructor for IP
// Parse an IP from string
func NewIP(str string) (ip IP, err error) {
	parts := strings.Split(str, ".")

	if len(parts) > 4 {
		err = fmt.Errorf("invalid number of ip segments")
		return
	}

	for i, part := range parts {
		var segment int64

		segment, err = strconv.ParseInt(part, 10, 8)
		if err != nil {
			return
		}

		ip[i] = byte(segment)
	}

	return
}

// _ "implement" constraint for IP
var _ fmt.Stringer = IP{}

// IP data type for ip address v4
type IP [4]byte

// String transforms an IP to string
//
// Example:
//
//     xxx.xxx.xxx.xxx
func (i IP) String() string {
	itoa := strconv.Itoa

	return itoa(int(i[0])) + "." + itoa(int(i[1])) + "." + itoa(int(i[2])) + "." + itoa(int(i[3]))
}
