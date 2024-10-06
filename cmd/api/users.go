package main

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/lib/pq"
	db "github.com/trenchesdeveloper/social-blue/internal/db/sqlc"
	"github.com/trenchesdeveloper/social-blue/internal/dto"
	"net/http"
	"strconv"
)

type UserContextKey string

const userKey UserContextKey = "user"

// GetUser godoc
//
//	@Summary		Get a user
//	@Description	Get a user by ID
//	@ID				get-user
//	@Produce		json
//	@Param			userID	path		int	true	"User ID"
//	@Success		200		{object}	dto.UserResponseDto
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Router			/users/{userID} [get]
func (s *server) getUserHandler(w http.ResponseWriter, r *http.Request) {
	// get the user from the context
	user, err := s.getUserFromContext(r.Context())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.notFoundError(w, r)
			return
		}
		s.internalServerError(w, r, err)
		return
	}

	userResponse := dto.UserResponseDto{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}

	// return the user
	jsonRespose(w, http.StatusOK, userResponse)
}

type FollowUser struct {
	UserID int64 `json:"user_id"`
}

// FollowUser godoc
//
//	@Summary		Follow a user
//	@Description	Follow a user by ID
//	@ID				follow-user
//	@Produce		json
//	@Param			userID	path		int	true	"User ID"
//	@Success		200		{object}	db.GetUserByIDRow
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//
//	@Security		ApiKeyAuth
//
//	@Router			/users/{userID}/follow [put]
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
		s.logger.Info(err)
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
		s.logger.Info(err)
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
				return
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
