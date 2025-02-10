package main

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/lib/pq"
	db "github.com/trenchesdeveloper/social-blue/internal/db/sqlc"
	"github.com/trenchesdeveloper/social-blue/internal/dto"
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
//	@Success		200		{object}	UserWithRole
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
	jsonResponse(w, http.StatusOK, userResponse)
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
//	@Success		200		{object}	UserWithRole
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

	// get the user id from the request
	followedID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)

	if err != nil {
		s.badRequestError(w, r, err)
		return
	}
	// follow the user
	err = s.store.FollowUser(r.Context(), db.FollowUserParams{
		UserID:     followedID,
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
	jsonResponse(w, http.StatusOK, followerUser)
}

func (s *server) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	// get the user from the context

	followerUser, err := s.getUserFromContext(r.Context())
	if err != nil {
		s.internalServerError(w, r, err)
		return
	}

	// get the user id from the request
	unFollowUserID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)

	if err != nil {
		s.badRequestError(w, r, err)
		return
	}


	// follow the user
	err = s.store.UnFollowUser(r.Context(), db.UnFollowUserParams{
		UserID:     unFollowUserID,
		FollowerID: followerUser.ID,
	})

	if err != nil {
		s.logger.Info(err)
		s.internalServerError(w, r, err)
		return
	}

	// return the user
	jsonResponse(w, http.StatusOK, followerUser)
}

func (s *server) getUserFromContext(ctx context.Context) (UserWithRole, error) {
	user, ok := ctx.Value(userKey).(UserWithRole)
	if !ok {
		return UserWithRole{}, errors.New("could not get user from context")
	}

	return user, nil
}
