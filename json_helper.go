package wundergo

import "encoding/json"

type JSONHelper interface {
	Marshal(v interface{}) ([]byte, error)
	Unmarshal(data []byte, v interface{}) (interface{}, error)
}

type defaultJSONHelper struct {
}

func newDefaultJSONHelper() *defaultJSONHelper {
	return &defaultJSONHelper{}
}

func (h defaultJSONHelper) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (h defaultJSONHelper) Unmarshal(data []byte, v interface{}) (interface{}, error) {
	err := json.Unmarshal(data, v)
	return v, err
}
