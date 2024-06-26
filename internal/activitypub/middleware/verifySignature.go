package middleware

import (
	"app/internal/activitypub/core/usecase"
	"app/internal/handler"
	"fmt"
	"net/http"
)

func NewVerifyMiddleware(verify *usecase.VerifySignatureUsecase) func(next http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(handler.Handle(
            func(w http.ResponseWriter, r *http.Request) error {
                println("middleware")
                fmt.Println(r.Host, r.URL.Path)

                if err := verify.Exec(r.Context(), r); err != nil {
                    return err
                }

                //ctx := context.WithValue(r.Context(), "userId", session.Owner)

                next.ServeHTTP(w, r)
                return nil
            },
        ))    
    }
}
