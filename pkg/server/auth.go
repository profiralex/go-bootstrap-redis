package server

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

func NewBasicApiKeyAuthMiddleware(apiKey string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeaderEncoded := r.Header.Get("Authorization")
			authHeaderEncoded = strings.Replace(authHeaderEncoded, "Basic ", "", -1)
			if len(authHeaderEncoded) == 0 {
				respondError(w, r, fmt.Errorf("unauthorized"), http.StatusUnauthorized)
				return
			}

			authHeaderDecodedBytes, err := base64.StdEncoding.DecodeString(authHeaderEncoded)
			if err != nil {
				respondError(w, r, fmt.Errorf("unauthorized"), http.StatusUnauthorized)
				return
			}

			authHeaderDecoded := string(authHeaderDecodedBytes)
			parts := strings.Split(authHeaderDecoded, ":")
			if len(parts) < 2 {
				respondError(w, r, fmt.Errorf("unauthorized"), http.StatusUnauthorized)
				return
			}

			username := parts[0]
			password := parts[1]
			if username != apiKey || password != "X" {
				respondError(w, r, fmt.Errorf("unauthorized"), http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
