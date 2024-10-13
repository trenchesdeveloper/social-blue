package main

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/google/uuid"
	db "github.com/trenchesdeveloper/social-blue/internal/db/sqlc"
	"github.com/trenchesdeveloper/social-blue/internal/dto"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func (s *server) registerUserHandler(w http.ResponseWriter, r *http.Request) {

	var input dto.RegisterUserDto
	if err := readJSON(w, r, &input); err != nil {
		s.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(input); err != nil {
		s.badRequestError(w, r, err)
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

	jsonRespose(w, http.StatusOK, user)
}
