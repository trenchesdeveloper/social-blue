package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
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

func (s *server) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			s.unauthorizedError(w, r, fmt.Errorf("missing Authorization header"))
			return
		}

		// parse the Authorization header
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			s.unauthorizedError(w, r, fmt.Errorf("invalid Authorization header"))
			return
		}

		// validate the token
		token := parts[1]
		jwtToken, err := s.authenticator.ValidateToken(token)
		if err != nil {
			s.logger.Error(err)
			s.unauthorizedError(w, r, fmt.Errorf("invalid token"))
			return
		}

		claims := jwtToken.Claims.(jwt.MapClaims)

		userID, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)

		if err != nil {
			s.unauthorizedError(w, r, err)
			return
		}
		ctx := r.Context()
		user, err := s.store.GetUserByID(ctx, userID)
		if err != nil {
			s.unauthorizedError(w, r, err)
			return
		}
		ctx = context.WithValue(ctx, userKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
