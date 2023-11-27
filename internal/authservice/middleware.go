package authservice

import (
	"context"
	"net/http"
)

var userCtxKey = &contextKey{"user"}

type contextKey struct {
	name string
}

func (as *AuthService) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie(cookieName)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		// h := r.Header.Get("Authorization")
		// s := strings.Split(h, " ")
		// if err != nil || len(s) != 2 {
		// 	next.ServeHTTP(w, r)
		// 	return
		// }

		//tokenStr := s[1]
		ctx := r.Context()
		tokenStr := c.Value

		jwt, err := as.Auth.ValidateJWT(ctx, tokenStr)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		username, err := jwt.Claims.GetSubject()
		if err != nil {
			http.Error(w, "invalid claims", http.StatusUnauthorized)
			return
		}

		user, err := as.Db.FindUser(ctx, username)
		if err != nil {
			return
		}

		ctx = context.WithValue(r.Context(), userCtxKey, user)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

// ForContext finds the user from the context. REQUIRES Middleware to have run.
func ForContext(ctx context.Context) *User {
	raw, _ := ctx.Value(userCtxKey).(*User)
	return raw
}
