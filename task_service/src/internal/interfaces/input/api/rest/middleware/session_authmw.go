package sessionauthmiddleware

import (
	"context"
	"fmt"
	"net/http"
	sessionclient "task_service/src/internal/adaptors/grpcclient"
)

const UserIdKey string = "userId"

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
			valid, userId, err := sessionClient.ValidateSession(r.Context(), sessionID)

			if err != nil || !valid {
				fmt.Println("Error in session auth middleware ", err)
				// fmt.Println("valid ", valid)

				http.Error(w, "invalid session", http.StatusUnauthorized)
				return
			}

			// store userid in context
			ctx := context.WithValue(r.Context(), UserIdKey, userId)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
