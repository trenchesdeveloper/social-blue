package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	db "github.com/trenchesdeveloper/social-blue/internal/db/sqlc"
	"github.com/trenchesdeveloper/social-blue/internal/dto"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type UserWithToken struct {
	db.CreateUserRow
	Token string `json:"token"`
}

// @Summary		Register a new user
// @Description	Register a new user
// @Tags			users
// @Accept			json
// @Produce		json
// @Param			input	body		dto.RegisterUserDto	true	"User data"
// @Success		200		{object}	UserWithToken
// @Router			/auth/register [post]
func (s *server) registerUserHandler(w http.ResponseWriter, r *http.Request) {

	var input dto.RegisterUserDto
	if err := readJSON(w, r, &input); err != nil {
		s.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(input); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Validation failed: "+err.Error())
		return
	}

	// hash the user password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	if err != nil {
		s.internalServerError(w, r, err)
		return
	}

	// create a plain token
	plainToken := uuid.New().String()

	hash := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hash[:])

	user, err := s.store.CreateAndInviteUser(r.Context(), hashToken, s.mailConfig.EXP, db.CreateUserParams{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Username:  input.Username,
		Email:     input.Email,
		Password:  hashedPassword,
	})

	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"`:
			s.badRequestError(w, r, db.ErrDuplicateUsername)
			return
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			s.badRequestError(w, r, db.ErrDuplicateEmail)
			return

		default:
			s.internalServerError(w, r, err)
			return
		}
	}

	// send email to user
	userWithToken := UserWithToken{
		CreateUserRow: user,
		Token:         plainToken,
	}

	jsonResponse(w, http.StatusCreated, userWithToken)
}

// @Summary		Activate a user
// @Description	Activate a user
// @Tags			users
// @Accept			json
// @Produce		json
// @Param			token	path		string	true	"Activation token"
// @Success		200		{object}	nil
// @Failure		400		{object}	error
// @Failure		404		{object}	error
// @Router			/auth/activate/{token} [put]
func (s *server) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	// get the token from the request
	token := chi.URLParam(r, "token")

	s.logger.Info("Token: ", token)

	// check if the token is empty
	if token == "" {
		s.badRequestError(w, r, db.ErrInvalidToken)
		return
	}

	// hash the token
	hash := sha256.Sum256([]byte(token))
	hashToken := hex.EncodeToString(hash[:])

	// get the user by the token
	_, err := s.store.ActivateUser(r.Context(), hashToken)

	if err != nil {
		switch {
		case errors.Is(err, db.ErrInvalidToken):
			s.badRequestError(w, r, err)
			return
		case errors.Is(err, sql.ErrNoRows):
			s.notFoundError(w, r)
			return
		default:
			s.internalServerError(w, r, err)
			return
		}
	}

	// send email to user

	jsonResponse(w, http.StatusOK, nil, "User activated successfully")
}
