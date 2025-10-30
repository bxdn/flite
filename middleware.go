package flite

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type jsonKey struct{}

func (f *Flite[T]) Body() *T {
	val := f.req.Context().Value(jsonKey{})
	typed, ok := val.(*T)
	if !ok {
		panic("Endpoint not configured properly, Json required as middleware to use this method!")
	}
	return typed
}

func (e *endpoint[T]) Json(f *Flite[T]) error {
	ptr := new(T)
	decoder := json.NewDecoder(f.req.Body)
	if e := decoder.Decode(ptr); e != nil {
		if e2 := f.ReturnError("Body is not in the correct JSON schema", http.StatusBadRequest); e2 != nil {
			return fmt.Errorf("Error returning bad request error: %v", e2)
		}
		return nil
	}
	f.AddContext(jsonKey{}, ptr)
	return nil
}
