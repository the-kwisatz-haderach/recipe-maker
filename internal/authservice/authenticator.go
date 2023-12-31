package authservice

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type Authenticator struct {
	signingSecret           string
	shouldValidateJwt       bool
	tokenExpirationDuration time.Duration
}

type CustomJWTClaims struct {
	User string `json:"user"`
	jwt.RegisteredClaims
}

const issuerClaim = "recipe-maker-auth-service"

// HashPassword hashes a plaintext password.
func (a *Authenticator) HashPassword(ctx context.Context, pass string) ([]byte, error) {
	var passwordBytes = []byte(pass)
	return bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
}

// ComparePasswords checks if two passwords match and returns an error otherwise.
func (a *Authenticator) ComparePasswords(ctx context.Context, plainPass string, hashedPass []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPass, []byte(plainPass))
}

func (a *Authenticator) GenerateJWT(ctx context.Context, u *User) (string, error) {
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(a.tokenExpirationDuration)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		Issuer:    issuerClaim,
		Subject:   u.ID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(a.signingSecret))
}

// ValidateJWT parses and validates a JWT, returning the parsed token struct if valid.
func (a *Authenticator) ValidateJWT(ctx context.Context, tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(a.signingSecret), nil
	})
	if !a.shouldValidateJwt && token != nil {
		return token, nil
	}
	if err != nil {
		log.Error().Err(err).Msg("failed to parse token")
		return nil, err
	}
	return token, nil
}
