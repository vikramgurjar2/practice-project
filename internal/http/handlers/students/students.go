package students

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/vikramgurjar2/practice-project/internal/types"
	"github.com/vikramgurjar2/practice-project/internal/utils/response"
)

func New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("creating a new student")

		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student)

		if errors.Is(err, io.EOF) {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
			return
		}

		if err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid json: %w", err)))
			return
		}

		//getting the data from request body
		slog.Info("student data", "student", student)

		//request validation

		if err := validator.New().Struct(student); err != nil {

			validateErr := err.(validator.ValidationErrors)
			response.WriteJSON(w, http.StatusBadRequest, response.ValidationError(validateErr))
			return
		}

		response.WriteJSON(w, http.StatusCreated, map[string]string{"success": "ok"})
	}
}
