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

func (s *server) unauthorizedError(w http.ResponseWriter, r *http.Request, err error) {
	s.logger.Warnf("Unauthorized: %s %s", r.Method, r.URL.Path)

	writeJSONError(w, http.StatusUnauthorized, err.Error())
}

func (s *server) unauthorizedBasicError(w http.ResponseWriter, r *http.Request, err error) {
	s.logger.Warnf("Unauthorized basic error: %s %s", r.Method, r.URL.Path)

	w.Header().Set("WWW-Authenticate", `Basic realm="Restricted", charset="UTF-8"`)

	writeJSONError(w, http.StatusUnauthorized, err.Error())
}

func (s *server) forbiddenError(w http.ResponseWriter, r *http.Request, err error) {
	s.logger.Warnf("Forbidden: %s %s", r.Method, r.URL.Path)

	writeJSONError(w, http.StatusForbidden, err.Error())
}