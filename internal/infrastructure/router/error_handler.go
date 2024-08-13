package router

import (
	"app/internal/application/appstatus"
	"encoding/json"
	"net/http"
)

func HandleError(e error, w http.ResponseWriter, r *http.Request) {
    if apperr, ok := e.(*appstatus.AppStatus); ok {
        body := map[string]interface{}{
            "error_code": apperr.Code(),
            "description": apperr.Error(),
        }

        bytes, err := json.Marshal(body)
        if err != nil {
            w.WriteHeader(500)
            w.Write([]byte(err.Error()))
        }

        w.Header().Add("Content-Type", "application/json")
        w.WriteHeader(apperr.Status())
        w.Write(bytes)
    } else {
        w.WriteHeader(500)
        w.Write([]byte(e.Error()))
    }
}
