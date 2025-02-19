package flite

import (
	"encoding/json"
	"log"
	"net/http"
)

func (f *flite) ReturnText(message string) {
	_, e := f.w.Write([]byte(message))
	if e != nil {
		log.Println(e)
		http.Error(f.w, "internal server error", http.StatusInternalServerError)
	}
}

func (f *flite) ReturnJSON(object any) {
	jsonBytes, err := json.Marshal(object)
	if err != nil {
		log.Println(err)
		http.Error(f.w, "internal server error", http.StatusInternalServerError)
		return
	}
	_, e := f.w.Write(jsonBytes)
	if e != nil {
		log.Println(e)
		http.Error(f.w, "internal server error", http.StatusInternalServerError)
	}
}
