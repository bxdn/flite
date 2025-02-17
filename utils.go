package flite

import (
	"encoding/json"
	"log"
	"net/http"
)

func ReturnText(w http.ResponseWriter, message string) {
	_, e := w.Write([]byte(message))
	if e != nil {
		log.Println(e)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func ReturnJSON(w http.ResponseWriter, object any) {
	jsonBytes, err := json.Marshal(object)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	_, e := w.Write(jsonBytes)
	if e != nil {
		log.Println(e)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
