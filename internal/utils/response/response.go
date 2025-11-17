package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func GeneralError(err error) Response {
	return Response{
		Status:  "error",
		Message: err.Error(),
	}
}

func WriteJSON(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)

}

func ValidationError(err validator.ValidationErrors) Response {

	var errMsg []string
	for _, e := range err {
		switch e.ActualTag() {
		case "required":
			errMsg = append(errMsg, fmt.Sprintf("%s is required", e.Field()))
		default:
			errMsg = append(errMsg, fmt.Sprintf("%s is not valid", e.Field()))
		}
	}
	return Response{
		Status:  "error",
		Message: strings.Join(errMsg, ","),
	}
}
