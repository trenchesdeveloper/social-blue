package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

func (s *server) BasicAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the Basic Authentication credentials
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			s.unauthorizedBasicError(w, r, fmt.Errorf("missing Authorization header"))
			return
		}

		// parse the Authorization header
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Basic" {
			s.unauthorizedBasicError(w, r, fmt.Errorf("invalid Authorization header"))
			return
		}

		// decode the credentials
		creds, err := base64.StdEncoding.DecodeString(parts[1])
		if err != nil {
			s.unauthorizedBasicError(w, r, fmt.Errorf("invalid base64 encoding"))
			return
		}

		// parse the credentials
		pair := strings.Split(string(creds), ":")
		if len(pair) != 2 {
			s.unauthorizedBasicError(w, r, fmt.Errorf("invalid credentials"))
			return
		}

		username := s.config.BASIC_AUTH_USERNAME
		password := s.config.BASIC_AUTH_PASSWORD

		// validate the credentials
		if pair[0] != username || pair[1] != password {
			s.unauthorizedBasicError(w, r, fmt.Errorf("invalid credentials"))
			return
		}

		next.ServeHTTP(w, r)

	})
}
