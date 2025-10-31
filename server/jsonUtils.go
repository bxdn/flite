package flite

import (
	"encoding/json"
)

func JsonMapToType[T any](obj map[string]any) (*T, error) {
	data, e := json.Marshal(obj)
	if e != nil {
		return nil, e
	}
	ptr := new(T)
	e = json.Unmarshal(data, ptr)
	return ptr, e
}

func TypeToJsonMap(obj any) (map[string]any, error) {
	data, e := json.Marshal(obj)
	if e != nil {
		return nil, e
	}
	m := map[string]any{}
	e = json.Unmarshal(data, &m)
	return m, e
}
