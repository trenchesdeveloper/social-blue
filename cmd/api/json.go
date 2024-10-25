package main

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"net/http"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())
}

func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxByte := 1 << 20 // 1MB
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxByte))
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(data); err != nil {
		return err
	}
	return nil
}

func writeJSONError(w http.ResponseWriter, status int, message string) error {
	data := map[string]string{"error": message}
	return writeJSON(w, status, data)

}

func jsonResponse(w http.ResponseWriter, status int, data any, msg ...string) error {
	type response struct {
		Data    any    `json:"data,omitempty"`
		Message string `json:"message,omitempty"`
	}

	var message string
	if len(msg) > 0 {
		message = msg[0]
	}

	var responseData any
	if data != nil {
		responseData = data
	}

	return writeJSON(w, status, response{Data: responseData, Message: message})
}
