package ml

import (
	"fmt"
)

func readValue[T any](query map[string]interface{}, key string) (T, error) {
	var result T
	v, ok := query[key]
	if !ok {
		return result, fmt.Errorf("required field '%s' is missing", key)
	}
	result, ok = v.(T)
	if !ok {
		return result, fmt.Errorf("field '%s' has type %T but expected string", key, v)
	}
	return result, nil
}

func readOptionalValue[T any](query map[string]interface{}, key string) (*T, error) {
	var result T
	v, ok := query[key]
	if !ok {
		return nil, fmt.Errorf("required field '%s' is missing", key)
	}
	result, ok = v.(T)
	if !ok {
		return nil, fmt.Errorf("field '%s' has type %T but expected string", key, v)
	}
	return &result, nil
}
