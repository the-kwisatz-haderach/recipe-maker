package authservice

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/the-kwisatz-haderach/recipemaker/internal/config"
)

const cookieName = "session-cookie"

func NewAuthService(db userStorage) AuthService {
	var auth = Authenticator{
		signingSecret:           config.Config.JWT_SIGNING_SECRET,
		shouldValidateJwt:       config.Config.VALIDATE_JWT,
		tokenExpirationDuration: time.Hour * 24,
	}
	return AuthService{Db: db, Auth: &auth}
}

func (as *AuthService) LoginHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method != http.MethodPost {
		http.Error(w, "invalid request method", http.StatusMethodNotAllowed)
		return
	}
	if r.Body == nil {
		http.Error(w, "missing request body", http.StatusBadRequest)
		return
	}
	var input LoginInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	u, _ := as.Db.FindUser(ctx, "", input.Username)
	if u == nil {
		http.Error(w, "invalid login credentials", http.StatusUnauthorized)
		return
	}
	err = as.Auth.ComparePasswords(ctx, input.Password, []byte(u.Password))
	if err != nil {
		http.Error(w, "invalid login credentials", http.StatusUnauthorized)
		return
	}
	tokenStr, err := as.Auth.GenerateJWT(ctx, u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	cookie := http.Cookie{
		Name:     cookieName,
		Value:    tokenStr,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
		Expires:  time.Now().Add(time.Hour * 24),
		Path:     "/",
	}
	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusOK)
}

func (as *AuthService) SignupHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method != http.MethodPost {
		http.Error(w, "invalid request method", http.StatusMethodNotAllowed)
		return
	}
	if r.Body == nil {
		http.Error(w, "missing request body", http.StatusBadRequest)
		return
	}
	var input SignupInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	u, _ := as.Db.FindUser(ctx, "", input.Username)
	if u != nil {
		http.Error(w, "user already exists", http.StatusConflict)
		return
	}
	encryptedPass, err := as.Auth.HashPassword(ctx, input.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	input.Password = string(encryptedPass)
	if _, err = as.Db.CreateUser(ctx, input); err != nil {
		http.Error(w, "unable to create user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (as *AuthService) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "invalid request method", http.StatusMethodNotAllowed)
		return
	}
	cookie := http.Cookie{
		Name:     cookieName,
		Value:    "",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
		MaxAge:   -1,
		Path:     "/",
	}
	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusOK)
}
