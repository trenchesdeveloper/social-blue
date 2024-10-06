package main

import (
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

	user, err := s.store.CreateUser(r.Context(), db.CreateUserParams{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Username:  input.Username,
		Email:     input.Email,
		Password:  hashedPassword,
	})

	if err != nil {
		s.internalServerError(w, r, err)
		return
	}

	jsonRespose(w, http.StatusOK, user)
}
