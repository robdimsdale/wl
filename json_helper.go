package wundergo

import "encoding/json"

type JSONHelper interface {
	Marshal(v interface{}) ([]byte, error)
	Unmarshal(data []byte, v interface{}) (interface{}, error)
}

type DefaultJSONHelper struct {
}

func NewDefaultJSONHelper() *DefaultJSONHelper {
	return &DefaultJSONHelper{}
}

func (h DefaultJSONHelper) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (h DefaultJSONHelper) Unmarshal(data []byte, v interface{}) (interface{}, error) {
	err := json.Unmarshal(data, v)
	return v, err
}
