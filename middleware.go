package flite

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type jsonKey struct{}

func GetTypedBody[T any](ctx context.Context) (*T, error) {
	val := ctx.Value(jsonKey{})
	typed, ok := val.(*T)
	if !ok {
		return typed, errors.New("type assertion failed")
	}
	return typed, nil
}

func Json[T any](f *Flite) (context.Context, error) {
	ptr := new(T)
	decoder := json.NewDecoder(f.r.Body)
	if e := decoder.Decode(ptr); e != nil {
		log.Println(e)
		http.Error(f.w, "bad request", http.StatusBadRequest)
		return f.r.Context(), e
	}
	return context.WithValue(f.r.Context(), jsonKey{}, ptr), nil
}
