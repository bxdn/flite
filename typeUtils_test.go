package flite

import (
	"context"
	"testing"
)

type test struct {
	X string `json:"x"`
	Y int    `json:"y"`
}

func TestGetFromContext(t *testing.T) {
	x := 5
	ctx := context.WithValue(context.Background(), jsonKey{}, &x)
	val, e := GetTypedBody[int](ctx)
	if e != nil {
		t.Error(e)
	}
	if *val != 5 {
		t.Errorf("expected 5, got %d", *val)
	}
}

func TestGetFromContextError(t *testing.T) {
	x := 5
	ctx := context.WithValue(context.Background(), jsonKey{}, &x)
	_, e := GetTypedBody[string](ctx)
	if e == nil {
		t.Errorf("expected an error, but did not get one")
	}
}

func TestJSONMapToType(t *testing.T) {
	obj := map[string]any{"x": "test", "y": 1}
	ptr, e := JsonMapToType[test](obj)
	if e != nil {
		t.Error(e)
	}
	expected := test{"test", 1}
	if *ptr != expected {
		t.Errorf("expected test{'test', 1}, got %v", *ptr)
	}
}

func TestJSONMapToTypeDecodeError(t *testing.T) {
	obj := map[string]any{"x": "test", "y": 1}
	_, e := JsonMapToType[int](obj)
	if e == nil {
		t.Errorf("expected an error, but did not get one")
	}
}

func TestJSONMapToTypeEncodeError(t *testing.T) {
	obj := map[string]any{
		"hello": map[bool]any{
			false: "world",
		},
	}
	_, e := JsonMapToType[any](obj)
	if e == nil {
		t.Errorf("expected an error, but did not get one")
	}
}
