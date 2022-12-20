package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"a21hc3NpZ25tZW50/db"
	"a21hc3NpZ25tZW50/model"
)

func isExpired(s model.Session) bool {
	return s.Expiry.Before(time.Now())
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("session_token")
		if err != nil {
			if err == http.ErrNoCookie {
				w.WriteHeader(401)
				error := model.ErrorResponse{
					Error: "http: named cookie not present",
				}
				jsonError, err := json.Marshal(error)
				if err != nil {
					return
				}

				w.Write([]byte(jsonError))
				return
			}
		}

		sessionToken := c.Value

		userSession, exists := db.Sessions[sessionToken]

		if !exists {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if isExpired(userSession) {
			delete(db.Sessions, sessionToken)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "username", userSession.Username)
		next.ServeHTTP(w, r.WithContext(ctx))
	}) // TODO: replace this
}
