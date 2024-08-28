package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/jwt"
	"github.com/clerk/clerk-sdk-go/v2/user"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/hello", HelloHandler)
	mux.Handle("/protected", ProtectRouteMiddleware(http.HandlerFunc(ProtectedHandler)))

	clerk.SetKey("")

	slog.Info("Starting server on port 8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		slog.Error("Error starting server", "error", err)
	}
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Protected route")
}

func ProtectRouteMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Protecting route")
		slog.Info("Request headers", "headers", r.Header)
		sessionToken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		slog.Info("Session token", "token", sessionToken)
		claims, err := jwt.Verify(r.Context(), &jwt.VerifyParams{
			Token: sessionToken,
		})
		if err != nil {
			slog.Error("Error verifying JWT", "error", err)
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"access": "unauthorized"}`))
			return
		}

		usr, err := user.Get(r.Context(), claims.Subject)
		if err != nil {
			slog.Error("Error getting user from Clerk")
			return
		}
		fmt.Fprintf(w, `{"user_id": "%s", "user_banned": "%t"}`, usr.ID, usr.Banned)
		next.ServeHTTP(w, r)
	})
}
