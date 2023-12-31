package authservice

import (
	"context"
	"net/http"
	"strings"
)

var userCtxKey = &contextKey{"user"}

type contextKey struct {
	name string
}

func (as *AuthService) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var tokenStr string
		c, err := r.Cookie(cookieName)
		if err == nil {
			tokenStr = c.Value
		}

		if tokenStr == "" {
			header := r.Header.Get("Authorization")
			if str := strings.TrimPrefix(header, "Bearer "); str != "" {
				tokenStr = str
			}
		}

		if tokenStr == "" {
			next.ServeHTTP(w, r)
			return
		}

		ctx := r.Context()

		jwt, err := as.Auth.ValidateJWT(ctx, tokenStr)
		if err != nil {
			http.Error(w, `{"error":"invalid token"}`, http.StatusUnauthorized)
			return
		}

		id, err := jwt.Claims.GetSubject()
		if err != nil {
			http.Error(w, `{"error":"invalid claims"}`, http.StatusUnauthorized)
			return
		}

		user, err := as.Db.FindUser(ctx, id, "")
		if err != nil {
			http.Error(w, `{"error":"invalid claims"}`, http.StatusUnauthorized)
			return
		}

		ctx = context.WithValue(r.Context(), userCtxKey, user)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

// GetUser gets the user from the context. Requires middleware to have run.
func GetUser(ctx context.Context) *User {
	raw, _ := ctx.Value(userCtxKey).(*User)
	return raw
}
