package auth

import (
	"context"

	"golang.org/x/crypto/bcrypt"
)

type Authenticator struct{}

func (a *Authenticator) HashPassword(ctx context.Context, pass string) ([]byte, error) {
	var passwordBytes = []byte(pass)
	return bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
}

// Check if two passwords match
func (a *Authenticator) ComparePasswords(ctx context.Context, plainPass string, hashedPass []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPass, []byte(plainPass))
}
