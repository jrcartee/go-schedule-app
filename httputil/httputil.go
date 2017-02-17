package httputil

import (
	"encoding/json"
	"log"
	"net/http"
)

func EncodeJSONResponse(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	return json.NewEncoder(w).Encode(data)
}

func DecodeJSONRequest(r *http.Request, i interface{}) error {
	decoder := json.NewDecoder(r.Body)
	return decoder.Decode(i)
}

func InternalError(rw http.ResponseWriter, err error) {
	log.Printf("Unexpected Error:\n%s\n", err.Error())
	http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
}
