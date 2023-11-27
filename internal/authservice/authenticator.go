package authservice

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

const signingSecret string = "hmacSampleSecret"

type Authenticator struct{}

func (a *Authenticator) HashPassword(ctx context.Context, pass string) ([]byte, error) {
	var passwordBytes = []byte(pass)
	return bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
}

// Check if two passwords match
func (a *Authenticator) ComparePasswords(ctx context.Context, plainPass string, hashedPass []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPass, []byte(plainPass))
}

func (a *Authenticator) GenerateJWT(ctx context.Context, u *User) (string, error) {
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Second * 30).Unix(),
		"sub": &u.Username,
	})
	return token.SignedString([]byte(signingSecret))
}

func (a *Authenticator) ValidateJWT(ctx context.Context, tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(signingSecret), nil
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to parse token")
		return nil, err
	}
	return token, nil
}
