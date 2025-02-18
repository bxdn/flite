package flite

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type jsonKey struct{}

func Deserialize[T any](ctx context.Context) (*T, error) {
	val := ctx.Value(jsonKey{})
	typed, ok := val.(*T)
	if !ok {
		return typed, errors.New("type assertion failed")
	}
	return typed, nil
}

func Json[T any](w http.ResponseWriter, r *http.Request) (context.Context, error) {
	ptr := new(T)
	decoder := json.NewDecoder(r.Body)
	if e := decoder.Decode(ptr); e != nil {
		log.Println(e)
		http.Error(w, "bad request", http.StatusBadRequest)
		return r.Context(), e
	}
	return context.WithValue(r.Context(), jsonKey{}, ptr), nil
}
