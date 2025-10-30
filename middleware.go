package flite

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type jsonKey struct{}

func GetTypedBody[T any](f *Flite) (*T, error) {
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
		http.Error(f.res, "bad request", http.StatusBadRequest)
		return e
	}
	f.AddContext(jsonKey{}, ptr)
	return nil
}
