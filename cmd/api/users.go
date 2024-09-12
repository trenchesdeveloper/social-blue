package main

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/lib/pq"
	db "github.com/trenchesdeveloper/social-blue/internal/db/sqlc"
	"log"
	"net/http"
	"strconv"
)

type UserContextKey string

const userKey UserContextKey = "user"

func (s *server) getUserHandler(w http.ResponseWriter, r *http.Request) {
	// get the user from the context
	user, err := s.getUserFromContext(r.Context())
	if err != nil {
		s.internalServerError(w, r, err)
		return
	}

	// return the user
	jsonRespose(w, http.StatusOK, user)
}

type FollowUser struct {
	UserID int64 `json:"user_id"`
}

func (s *server) followUserHandler(w http.ResponseWriter, r *http.Request) {
	// get the user from the context
	followerUser, err := s.getUserFromContext(r.Context())
	if err != nil {
		s.internalServerError(w, r, err)
		return
	}
	var payload FollowUser
	if err = readJSON(w, r, &payload); err != nil {
		s.badRequestError(w, r, err)
		return
	}

	// follow the user
	err = s.store.FollowUser(r.Context(), db.FollowUserParams{
		UserID:     payload.UserID,
		FollowerID: followerUser.ID,
	})

	if err != nil {
		log.Println(err)
		if pqError, ok := err.(*pq.Error); ok && pqError.Code == "23505" {
			s.badRequestError(w, r, errors.New("user already followed"))
			return
		}
		s.internalServerError(w, r, err)
		return
	}

	// return the user
	jsonRespose(w, http.StatusOK, followerUser)
}

func (s *server) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	// get the user from the context
	followerUser, err := s.getUserFromContext(r.Context())
	if err != nil {
		s.internalServerError(w, r, err)
		return
	}
	var payload FollowUser
	if err = readJSON(w, r, &payload); err != nil {
		s.badRequestError(w, r, err)
		return
	}

	// follow the user
	err = s.store.UnFollowUser(r.Context(), db.UnFollowUserParams{
		UserID:     payload.UserID,
		FollowerID: followerUser.ID,
	})

	if err != nil {
		log.Println(err)
		s.internalServerError(w, r, err)
		return
	}

	// return the user
	jsonRespose(w, http.StatusOK, followerUser)
}

func (s *server) userContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		// set the user in the context
		ctx := context.WithValue(r.Context(), userKey, user)

		// call the next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *server) getUserFromContext(ctx context.Context) (db.GetUserByIDRow, error) {
	user, ok := ctx.Value(userKey).(db.GetUserByIDRow)
	if !ok {
		return db.GetUserByIDRow{}, errors.New("could not get user from context")
	}

	return user, nil
}
