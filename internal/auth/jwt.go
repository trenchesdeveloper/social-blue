package auth

import "github.com/golang-jwt/jwt/v5"

type JWTAuthenticator struct {
	secretKey string
	aud       string
	iss       string
	ExpiresAt int64
}

func NewJWTAuthenticator(secretKey, aud, iss string) *JWTAuthenticator {
	return &JWTAuthenticator{
		secretKey: secretKey,
		aud:       aud,
		iss:       iss,
		ExpiresAt: 3600,
	}
}

func (j *JWTAuthenticator) GenerateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

func (j *JWTAuthenticator) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secretKey), nil
	})
}