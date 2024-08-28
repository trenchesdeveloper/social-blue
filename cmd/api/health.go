package main

import "net/http"

func (s *server) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))

}
