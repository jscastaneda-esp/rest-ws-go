package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/jscastaneda-esp/rest-ws-go/server"
	"github.com/jscastaneda-esp/rest-ws-go/services"
)

var (
	NO_AUTH_NEEDED = []string{
		"signup",
		"login",
	}
)

func shouldCheckToken(route string) bool {
	for _, r := range NO_AUTH_NEEDED {
		if strings.Contains(route, r) {
			return false
		}
	}

	return true
}

func CheckAuthMiddleware(s server.Server) func(h http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !shouldCheckToken(r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}

			_, err := services.CheckToken(r.Header.Get("Authorization"), s.Config().JWTSecret)
			if err != nil {
				log.Println("CheckToken:", err)
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
