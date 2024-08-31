package main

import (
	"encoding/json"
	"net/http"
)

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
