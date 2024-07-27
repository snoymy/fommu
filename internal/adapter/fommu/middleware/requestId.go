package middleware

import (
	"app/internal/adapter/handler"
	"app/internal/log"
	"context"
	"net/http"

	"github.com/google/uuid"
)

type RequestIdMiddleware func(http.Handler) http.Handler

func NewRequestIDMiddleware() RequestIdMiddleware {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(handler.Handle(
            func(w http.ResponseWriter, r *http.Request) error {
                log.EnterMethod(r.Context())

                log.Info(r.Context(), "Generate request ID")
                ctx := context.WithValue(r.Context(), "requestId", uuid.New().String())

                log.Info(ctx, "After generate request ID")
                defer log.ExitMethod(ctx)

                next.ServeHTTP(w, r.WithContext(ctx))
                return nil
            },
        ))    
    }
}
