package server

import "net/http"

func Cors(w http.ResponseWriter, r *http.Request) (error, bool) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, PATCH, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return nil, true
	}
	return nil, false
}
