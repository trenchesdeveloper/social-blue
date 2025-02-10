package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	db "github.com/trenchesdeveloper/social-blue/internal/db/sqlc"
	"github.com/trenchesdeveloper/social-blue/internal/dto"
	"golang.org/x/crypto/bcrypt"
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
		RoleID:   1,
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

	confirmationUrl := fmt.Sprintf("%s/confirm/%s", s.config.FrontendURL, plainToken)

	isProdEnv := s.config.Environment == "production"

	vars := struct{
		Username string
		ConfirmationURL string
	}{
		Username: user.Username,
		ConfirmationURL: confirmationUrl,
	}

	if isProdEnv {
		err = s.mailer.Send("user_invitation.tmpl", user.Username, user.Email, vars, false)
		if err != nil {
			s.logger.Error("Failed to send welcome email: ", err)

			//  rollback the user creation
			err = s.store.DeleteUserAndInvitation(r.Context(), user.ID)
			if err != nil {
				s.logger.Error("Failed to rollback user creation: ", err)
			}
		}
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

// @Summary		Login a user
// @Description	Login a user
// @Tags			users
// @Accept			json
// @Produce		json
// @Param			input	body		dto.LoginUserDto	true	"User data"
// @Success		200		{object}	string
// @Failure		400		{object}	error
// @Failure		401		{object}	error
// @Router			/auth/login [post]
func (s *server) loginHandler(w http.ResponseWriter, r *http.Request) {
	var input dto.LoginUserDto
	if err := readJSON(w, r, &input); err != nil {
		s.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(input); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Validation failed: "+err.Error())
		return
	}

	user, err := s.store.GetActiveUserByEmail(r.Context(), input.Email)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			s.unauthorizedError(w, r, db.ErrInvalidCredentials)
			return
		default:
			s.internalServerError(w, r, err)
			return
		}
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(input.Password))
	if err != nil {
		s.unauthorizedError(w, r, db.ErrInvalidCredentials)
		return
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 24 * 3).Unix(),
		"iat":     time.Now().Unix(),
		"nbf":     time.Now().Unix(),
		"iss":     s.config.APP_NAME,
		"aud":     s.config.APP_NAME,
	}

	token, err := s.authenticator.GenerateToken(claims)
	if err != nil {
		s.internalServerError(w, r, err)
		return
	}

	jsonResponse(w, http.StatusOK, token)
}