package flite

import (
	"encoding/json"
	"log"
	"net/http"
)

func (f *Flite) ReturnText(message string) error {
	_, e := f.Res.Write([]byte(message))
	if e != nil {
		log.Println(e)
		f.ReturnError("internal server error", http.StatusInternalServerError)
	}
	return e
}

func (f *Flite) ReturnJSON(object any) error {
	jsonBytes, e := json.Marshal(object)
	if e != nil {
		log.Println(e)
		f.ReturnError("internal server error", http.StatusInternalServerError)
		return e
	}
	_, e = f.Res.Write(jsonBytes)
	if e != nil {
		log.Println(e)
		f.ReturnError("internal server error", http.StatusInternalServerError)
	}
	return e
}

func (f *Flite) ReturnError(message string, status int) {
	http.Error(f.Res, message, status)
}
