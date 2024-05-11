package middleware

import (
	"app/internal/appstatus"
	"app/internal/core/usecase"
	"app/internal/handler"
	"net/http"
	"strings"
)

func NewAuthMiddleware(auth *usecase.AuthUsecase) func(next http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(handler.Handle(
            func(w http.ResponseWriter, r *http.Request) error {
                sessionId, err := r.Cookie("session_id")
                if err != nil {
                    return appstatus.InvalidSession("Session not found")
                }

                bearer := r.Header.Get("Authorization")
                accessToken := strings.ReplaceAll(bearer, "Bearer ", "") 

                if err := auth.Exec(r.Context(), sessionId.Value, accessToken); err != nil {
                    return err
                }

                next.ServeHTTP(w, r)
                return nil
            },
        ))    
    }
}
