package main

import (
	"net/http"
)

const (
	ErrorNotFound = "sql.ErrNoRows"
)

func (s *server) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	s.logger.Errorw("Internal server error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeJSONError(w, http.StatusInternalServerError, "Something went wrong")
}

func (s *server) notFoundError(w http.ResponseWriter, r *http.Request) {
	s.logger.Warnf("Not found: %s %s", r.Method, r.URL.Path)
	writeJSONError(w, http.StatusNotFound, "Not found")
}

func (s *server) badRequestError(w http.ResponseWriter, r *http.Request, err error) {
	s.logger.Warnf("Bad request: %s %s", r.Method, r.URL.Path)
	writeJSONError(w, http.StatusBadRequest, err.Error())
}
