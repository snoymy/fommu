package middlewares

import (
	"app/internal/application/fommu/usecases"
	"app/internal/application/appstatus"
	"app/internal/adapter/handler"
	"app/internal/log"
	"context"
	"net/http"
	"strings"
)

type AuthMiddleware func(http.Handler) http.Handler

func NewAuthMiddleware(auth *usecases.AuthUsecase) AuthMiddleware {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(handler.Handle(
            func(w http.ResponseWriter, r *http.Request) error {
                ctx := r.Context()
                log.EnterMethod(ctx)
                defer log.ExitMethod(ctx)

                log.Info(ctx, "Get session id from cookie")
                sessionId, err := r.Cookie("session_id")
                if err != nil {
                    log.Warn(ctx, "Response with error: " + err.Error())
                    return appstatus.InvalidSession("Session not found")
                }

                bearer := r.Header.Get("Authorization")
                accessToken := strings.ReplaceAll(bearer, "Bearer ", "") 

                session, err := auth.Exec(ctx, sessionId.Value, accessToken)
                if err != nil {
                    log.Info(ctx, "Response with error: " + err.Error())
                    return err
                }

                log.Info(ctx, "Set userId to context")
                ctx = context.WithValue(r.Context(), "userId", session.Owner)

                next.ServeHTTP(w, r.WithContext(ctx))
                return nil
            },
        ))    
    }
}
