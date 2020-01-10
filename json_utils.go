package goserv

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// ConvertToInteger converts a raw unmarshaled value to a float64. When unmarshaling to an interface, the unmarshaler may unmarshal a number to a json.Number or a float64. The float is then floored by this function and returned. This conversion function encapsulates this behavior.
func ConvertToInteger(unmarshaledVal interface{}) (int, error) {
	f, err := ConvertToFloat(unmarshaledVal)
	if err != nil {
		return 0, err
	}
	return int(f), nil
}

// ConvertToFloat converts a raw unmarshaled value to a float64. When unmarshaling to an interface, the unmarshaler may unmarshal a number to a json.Number or a float64. This conversion function encapsulates this behavior.
func ConvertToFloat(unmarshaledVal interface{}) (float64, error) {
	if f, ok := unmarshaledVal.(float64); ok {
		return f, nil
	}
	if f, ok := unmarshaledVal.(int); ok {
		return float64(f), nil
	}
	if num, ok := unmarshaledVal.(json.Number); ok {
		f, err := num.Float64()
		if err != nil {
			return 0, err
		}
		return f, nil
	}
	return 0, fmt.Errorf("conversion called with an interface with type %v", reflect.TypeOf(unmarshaledVal))
}
