package students

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/vikramgurjar2/practice-project/internal/types"
	"github.com/vikramgurjar2/practice-project/internal/utils/response"
)

func New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("creating a new student")

		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student)
		if err != nil {
			if errors.Is(err, io.EOF) {
				response.WriteJSON(w, http.StatusBadRequest, err.Error())
				return
			}
		}

		//getting the data from request body
		slog.Info("student data", "student", student)

		response.WriteJSON(w, http.StatusCreated, map[string]string{"success": "ok"})
	}
}
