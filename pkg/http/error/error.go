package error

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Error RespErr `json:"error"`
}

type RespErr struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func InternalError(w http.ResponseWriter, err error) {
	httpError(w, http.StatusInternalServerError, err)
}

func httpError(w http.ResponseWriter, code int, err error) {
	w.WriteHeader(code)

	json.NewEncoder(w).Encode(
		Response{
			Error: RespErr{
				Message: err.Error(),
				Code:    code,
			},
		},
	)
}
