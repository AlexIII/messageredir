package middleware

import (
	"context"
	"log"
	"messageredir/cmd/messageredir/app"
	"messageredir/cmd/messageredir/db/repo"
	"net/http"
	"strings"
)

type contextKey string

const UserKey contextKey = "user"

// AuthMiddleware checks user authorization and fetches user information.
func UserAuth(config *app.Config, db repo.DbRepo, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.Split(r.URL.Path, "/")
		var token *string
		if len(path) > 1 && len(path[1]) > (config.UserTokenLength-1)*8/6 { // Approximate check for token length
			token = &path[1]
		}
		if token == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		user := db.GetUserByToken(*token)
		if user == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		log.Println("[middleware] UserAuth:", user)

		ctx := context.WithValue(r.Context(), UserKey, user)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
