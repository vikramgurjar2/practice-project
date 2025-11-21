package students

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/vikramgurjar2/practice-project/internal/storage"
	"github.com/vikramgurjar2/practice-project/internal/types"
	"github.com/vikramgurjar2/practice-project/internal/utils/response"
)

func New(storage storage.Storage) http.HandlerFunc {
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
		///////////----------->
		//check if email already exists
		exists, err := storage.IsEmailExists(student.Email)

		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(fmt.Errorf("internal server error: %w", err)))
			return
		}

		if exists {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("email already exists")))
			return
		}

		lastid, err := storage.CreateStudent(student.Name, student.Age, student.Email)
		if err != nil {
			slog.Error("failed to shutdown server", "error", err)
			return
		}
		slog.Info("student created successfully", "id", lastid)

		response.WriteJSON(w, http.StatusCreated, map[string]string{"success": "ok"})
	}
}

func GetById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("getting student by id")
		//id from url
		idstr := r.URL.Path[len("/api/students/"):]

		//convert id to int
		id, err := strconv.Atoi(idstr)
		if err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("id is of different format")))
		}
		///get student by id
		student, err := storage.GetStudentById(id)
		if err != nil {
			response.WriteJSON(w, http.StatusNotFound, response.GeneralError(fmt.Errorf("student not found")))
			return
		}
		response.WriteJSON(w, http.StatusFound, student)
	}
}

func GetStudent(s storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("getting students list")
		//fetch students from storage
		students, err := s.GetStudents()
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(fmt.Errorf("internal server error: %w", err)))
			return
		}
		response.WriteJSON(w, http.StatusOK, students)
	}
}
