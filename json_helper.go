package wundergo

import "encoding/json"

// JSONHelper provides a wrapper around JSON marshalling/unmarshalling.
type JSONHelper interface {
	Marshal(v interface{}) ([]byte, error)
	Unmarshal(data []byte, v interface{}) (interface{}, error)
}

// DefaultJSONHelper is an implementation of JSONHelper.
type DefaultJSONHelper struct {
}

// Marshal is a wrapper around json.Marshal.
func (h DefaultJSONHelper) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// Unmarshal is a wrapper around json.Unmarshal.
func (h DefaultJSONHelper) Unmarshal(data []byte, v interface{}) (interface{}, error) {
	err := json.Unmarshal(data, v)
	return v, err
}
