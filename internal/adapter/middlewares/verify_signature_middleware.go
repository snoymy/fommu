package middlewares

import (
	"app/internal/application/activitypub/usecases"
	"app/internal/infrastructure/router"
	"fmt"
	"net/http"
)

type VerifyMiddleware func(http.Handler) http.Handler

func NewVerifyMiddleware(verify *usecases.VerifySignatureUsecase) VerifyMiddleware {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(router.Handle(
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
