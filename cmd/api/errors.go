package main

import (
	"log"
	"net/http"
)

func (s *server) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Internal server error: %s path=%s: error: %s", r.Method, r.URL.Path, err)
	writeJSONError(w, http.StatusInternalServerError, "Something went wrong")
}

func (s *server) notFoundError(w http.ResponseWriter, r *http.Request) {
	log.Printf("Not found: %s path=%s", r.Method, r.URL.Path)
	writeJSONError(w, http.StatusNotFound, "Not found")
}

func (s *server) badRequestError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Bad request: %s path=%s: error: %s", r.Method, r.URL.Path, err)
	writeJSONError(w, http.StatusBadRequest, err.Error())
}
