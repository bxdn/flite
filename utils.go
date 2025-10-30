package flite

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

func (f *Flite) ReturnText(message string) error {
	if f.done {
		return errors.New("Repsonse already finalized!")
	}
	_, e := f.res.Write([]byte(message))
	if e != nil {
		log.Println(e)
		if e2 := f.ReturnError("internal server error", http.StatusInternalServerError); e2 != nil {
			return fmt.Errorf("Error trying to return text, then error trying to return that error: %v, %v", e, e2)
		}
		return e
	}
	log.Printf("200 - %s - %s", f.req.Method, f.req.RequestURI)
	f.done = true
	return nil
}

func (f *Flite) ReturnJSON(object any) error {
	if f.done {
		return errors.New("Repsonse already finalized!")
	}
	jsonBytes, e := json.Marshal(object)
	if e != nil {
		log.Println(e)
		if e2 := f.ReturnError("internal server error", http.StatusInternalServerError); e2 != nil {
			return fmt.Errorf("Error trying to return text, then error trying to return that error: %v, %v", e, e2)
		}
		return e
	}
	_, e = f.res.Write(jsonBytes)
	if e != nil {
		log.Println(e)
		if e2 := f.ReturnError("internal server error", http.StatusInternalServerError); e2 != nil {
			return fmt.Errorf("Error trying to return text, then error trying to return that error: %v, %v", e, e2)
		}
		return e
	}
	log.Printf("200 - %s - %s", f.req.Method, f.req.RequestURI)
	f.done = true
	return nil
}

func (f *Flite) ReturnError(message string, status int) error {
	if f.done {
		return errors.New("Repsonse already finalized!")
	}
	http.Error(f.res, message, status)
	log.Printf("%d - %s - %s", status, f.req.Method, f.req.RequestURI)
	f.done = true
	return nil
}

func (f *Flite) Return() error {
	if f.done {
		return errors.New("Repsonse already finalized!")
	}
	if f.res.status == 0 {
		f.res.WriteHeader(http.StatusOK)
	}
	log.Printf("%d - %s - %s", f.res.status, f.req.Method, f.req.RequestURI)
	f.done = true
	return nil
}
