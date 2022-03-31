package model

import "encoding/json"

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
