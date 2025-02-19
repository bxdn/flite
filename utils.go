package flite

import (
	"encoding/json"
	"log"
	"net/http"
)

func (f *Flite) ReturnText(message string) {
	_, e := f.Res.Write([]byte(message))
	if e != nil {
		log.Println(e)
		http.Error(f.Res, "internal server error", http.StatusInternalServerError)
	}
}

func (f *Flite) ReturnJSON(object any) {
	jsonBytes, err := json.Marshal(object)
	if err != nil {
		log.Println(err)
		http.Error(f.Res, "internal server error", http.StatusInternalServerError)
		return
	}
	_, e := f.Res.Write(jsonBytes)
	if e != nil {
		log.Println(e)
		http.Error(f.Res, "internal server error", http.StatusInternalServerError)
	}
}
