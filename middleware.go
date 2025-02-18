package flite

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"reflect"
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

func Json(typ reflect.Type) RequestNode {
	return func(w http.ResponseWriter, r *http.Request) (context.Context, error) {
		ptr := reflect.New(typ).Interface()
		decoder := json.NewDecoder(r.Body)
		if e := decoder.Decode(ptr); e != nil {
			log.Println(e)
			http.Error(w, "bad request", http.StatusBadRequest)
			return r.Context(), e
		}
		return context.WithValue(r.Context(), jsonKey{}, ptr), nil
	}
}
