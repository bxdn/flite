package flite

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type jsonKey struct{}

func (f *F[T]) Body() *T {
	var zero T
    switch any(zero).(type) {
	case Never:
		panic("Cannot parse body of endpoint with [Never] type!")
	}
	val := f.req.Context().Value(jsonKey{})
	return val.(*T)
}

func DeserializeBody[T any]() func(f *F[T]) error {
	var zero T
    switch any(zero).(type) {
    case string:
        return Text
	case Never:
		return func(f *F[T]) error{return nil}
    default:
        return Json
    }
}

func Json[T any](f *F[T]) error {
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

func Text[T any](f *F[T]) error {
	bodyBytes, err := io.ReadAll(f.req.Body)
	if err != nil {
		return f.ReturnError("Error reading text request body", http.StatusInternalServerError)
	}
	defer f.req.Body.Close()
	bodyString := string(bodyBytes)
	f.AddContext(jsonKey{}, &bodyString)
	return nil
}
