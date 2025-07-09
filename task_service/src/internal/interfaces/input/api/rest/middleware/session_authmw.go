package sessionauthmiddleware

import (
	"fmt"
	"net/http"
	sessionclient "task_service/src/internal/adaptors/grpcclient"
)

func SessionAuthMiddleware(sessionClient *sessionclient.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("sess")
			if err != nil {
				http.Error(w, "session cookie is missing", http.StatusUnauthorized)
				return
			}
			sessionID := cookie.Value
			fmt.Println("Session id : ", sessionID)
			valid, _, err := sessionClient.ValidateSession(r.Context(), sessionID)

			if err != nil || !valid {
				http.Error(w, "invalid session", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
