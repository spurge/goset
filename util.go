package goset

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

func isMap(val interface{}) bool {
	return reflect.TypeOf(val).String()[0:3] == "map"
}

func toMap(val interface{}) map[string]interface{} {
	return val.(map[string]interface{})
}

func parse(data []byte) (map[string]interface{}, error) {
	var parsed map[string]interface{}

	if err := json.Unmarshal(data, &parsed); err != nil {
		return nil, err
	}

	return parsed, nil
}

func merge(values ...*map[string]interface{}) map[string]interface{} {
	merged := make(map[string]interface{})

	for _, val := range values {
		for k, v := range *val {
			e := merged[k]

			if e != nil && isMap(e) && isMap(v) {
				eMap := toMap(e)
				vMap := toMap(v)

				merged[k] = merge(&eMap, &vMap)
			} else {
				merged[k] = v
			}
		}
	}

	return merged
}

func extract(source map[string]interface{}, path string) (interface{}, error) {
	frags := strings.Split(path, ".")
	val, exists := source[frags[0]]

	if !exists {
		return nil, fmt.Errorf("%s not found", frags[0])
	}

	if len(frags) == 1 {
		return val, nil
	}

	switch val.(type) {
	case map[string]interface{}:
		return extract(val.(map[string]interface{}), strings.Join(frags[1:], "."))
	default:
		return nil, fmt.Errorf("%s is not a map", frags[0])
	}
}
