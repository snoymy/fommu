package handler

import (
    "net/http"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

var ErrorHandler func(error, http.ResponseWriter, *http.Request)

func defaultErrorHandler(e error, w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(500)
    w.Write([]byte(e.Error()))
}

func Handle(f HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        err := f(w, r)

        if err != nil {
            if ErrorHandler != nil {
                ErrorHandler(err, w, r)
            } else {
                defaultErrorHandler(err, w, r)
            }
        }
    }
}
