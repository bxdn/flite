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
	decoder := json.NewDecoder(f.Req.Body)
	if e := decoder.Decode(ptr); e != nil {
		log.Println(e)
		http.Error(f.Res, "bad request", http.StatusBadRequest)
		return f.Req.Context(), e
	}
	return context.WithValue(f.Req.Context(), jsonKey{}, ptr), nil
}
