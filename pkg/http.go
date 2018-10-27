package proxge

import (
	"encoding/json"
	"log"
	"net/http"
)

func writeError(code int, err string, w http.ResponseWriter) {
	w.WriteHeader(code)

	e := struct {
		Error string `json:"error"`
	}{Error: err}
	werr := json.NewEncoder(w).Encode(e)
	if werr != nil {
		log.Printf("failed to write error %s: %v\n", err, werr)
	}
}

func asJson(h http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		h.ServeHTTP(writer, request)
	}
}
