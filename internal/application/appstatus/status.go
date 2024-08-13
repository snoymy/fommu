package appstatus

import (
	"net/http"
)

type AppStatus struct {
    code string
    status int
    msg string
}

func NewAppStatus(code string, status int) func(msgs... string) *AppStatus {
    a := &AppStatus{}
    a.code = code
    a.status = status

    return func(msgs... string) *AppStatus {
        a.msg = ""
        if len(msgs) > 0 {
            a.msg += msgs[0]
            for _, msg := range msgs[1:] {
                a.msg += " " + msg
            }
        }
        if a.status == http.StatusOK {
            return nil
        }
        return a
    }
}

func (a *AppStatus) Error() string {
    return a.msg
}

func (a *AppStatus) Code() string {
    return a.code
}

func (a *AppStatus) Status() int {
    return a.status
}

var (
    Success = NewAppStatus("success", http.StatusOK)
    Duplicate = NewAppStatus("duplicate_data", http.StatusConflict)
    BadValue = NewAppStatus("bad_value", http.StatusBadRequest)
    BadUsername = NewAppStatus("bad_username", http.StatusBadRequest)
    BadEmail = NewAppStatus("bad_email", http.StatusBadRequest)
    BadPassword = NewAppStatus("bad_password", http.StatusBadRequest)
    BadLogin = NewAppStatus("invalid_login", http.StatusUnauthorized)
    NotFound = NewAppStatus("not_found", http.StatusNotFound)
    NotAccept = NewAppStatus("not_accept", http.StatusNotAcceptable)
    NotSupport = NewAppStatus("not_support", http.StatusNotImplemented)
    InvalidToken = NewAppStatus("invalid_token", http.StatusUnauthorized)
    InvalidSession = NewAppStatus("invalid_session", http.StatusUnauthorized)
    InvalidCredential = NewAppStatus("invalid_credential", http.StatusUnauthorized)
    InternalServerError = NewAppStatus("internal_server_error", http.StatusInternalServerError)
)

