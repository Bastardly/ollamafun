package session

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
)

func generateSessionID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

type contextKey string


// SessionIDKey passes sessionID to the request context, which can be fetched like so: sessionID := r.Context().Value(SessionIDKey).(string)
const SessionIDKey = contextKey("sessionID")

// Very simple Middleware to pass sessionID
func WithSession(handler http.HandlerFunc, sessionCookieName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var sessionID string
		cookie, err := r.Cookie(sessionCookieName)
		if err != nil || cookie.Value == "" {
			sessionID = generateSessionID()

			http.SetCookie(w, &http.Cookie{
				Name:  sessionCookieName,
				Value: sessionID,
				Path:  "/",
			})
		} else {
			sessionID = cookie.Value
		}

		ctx := context.WithValue(r.Context(), SessionIDKey, sessionID)
		handler(w, r.WithContext(ctx))
	}
}
