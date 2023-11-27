package authservice

import (
	"encoding/json"
	"net/http"
)

const cookieName = "session-cookie"

func NewAuthService(db userStorage) AuthService {
	var auth Authenticator
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

	u, _ := as.Db.FindUser(ctx, input.Username)
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

	u, _ := as.Db.FindUser(ctx, input.Username)
	if u != nil {
		http.Error(w, "user already exists", http.StatusBadRequest)
		return
	}
	encryptedPass, err := as.Auth.HashPassword(ctx, input.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	input.Password = string(encryptedPass)
	as.Db.CreateUser(ctx, input)
	w.WriteHeader(http.StatusOK)
}

func (as *AuthService) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "invalid request method", http.StatusMethodNotAllowed)
		return
	}
	_, err := r.Cookie(cookieName)
	if err == nil {
		cookie := http.Cookie{
			Name:     cookieName,
			Value:    "",
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
			MaxAge:   -1,
		}
		http.SetCookie(w, &cookie)
	}
	w.WriteHeader(http.StatusOK)
}
