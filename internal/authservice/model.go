package authservice

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
)

type userStorage interface {
	FindUser(ctx context.Context, username string) (*User, error)
	CreateUser(ctx context.Context, input SignupInput) (*User, error)
}

type authenticator interface {
	HashPassword(ctx context.Context, pass string) ([]byte, error)
	ComparePasswords(ctx context.Context, plainPass string, hashedPass []byte) error
	GenerateJWT(ctx context.Context, user *User) (string, error)
	ValidateJWT(ctx context.Context, tokenString string) (*jwt.Token, error)
}

type AuthService struct {
	Db   userStorage
	Auth authenticator
}

type User struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	ID       string `json:"id,omitempty"`
	Email    string `json:"email,omitempty"`
}

type LoginInput struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type SignupInput struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Email    string `json:"email,omitempty"`
}
