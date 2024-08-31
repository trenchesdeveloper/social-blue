package main

import "net/http"

func (s *server) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":      "ok",
		"environment": s.config.Environment,
		"version":     version,
	}
	if err := writeJSON(w, http.StatusOK, data); err != nil {
		s.internalServerError(w, r, err)
	}

}
