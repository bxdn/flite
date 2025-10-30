package flite

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type jsonKey struct{}

func Body[T any](f *Flite) (*T, error) {
	val := f.req.Context().Value(jsonKey{})
	typed, ok := val.(*T)
	if !ok {
		return typed, errors.New("type assertion failed")
	}
	return typed, nil
}

func Json[T any](f *Flite) error {
	ptr := new(T)
	decoder := json.NewDecoder(f.req.Body)
	if e := decoder.Decode(ptr); e != nil {
		log.Println(e)
		if e2 := f.ReturnError("Body is not in the correct JSON schema", http.StatusBadRequest); e2 != nil {
			return fmt.Errorf("Error trying to parse JSON body, then error trying to return that error: %v, %v", e, e2)
		}
		return e
	}
	f.AddContext(jsonKey{}, ptr)
	return nil
}
