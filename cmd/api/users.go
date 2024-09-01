package main

import (
	"database/sql"
	"errors"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func (s *server) getUserHandler(w http.ResponseWriter, r *http.Request) {
	// get the user id from the url
	userID := chi.URLParam(r, "userID")

	parsedID, err := strconv.ParseInt(userID, 10, 64)

	// fetch the user from the store
	user, err := s.store.GetUserByID(r.Context(), parsedID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.notFoundError(w, r)
		}
		s.internalServerError(w, r, err)
		return
	}

	// return the user
	jsonRespose(w, http.StatusOK, user)
}
