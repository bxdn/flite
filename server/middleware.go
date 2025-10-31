package server

import "net/http"

func Cors(w http.ResponseWriter, r *http.Request) (error, bool) {
	if r.Method != http.MethodOptions {
		return nil, false
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, PATCH, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.WriteHeader(http.StatusNoContent)
	return nil, true
}
