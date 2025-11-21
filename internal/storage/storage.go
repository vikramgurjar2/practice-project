package storage

import "github.com/vikramgurjar2/practice-project/internal/types"

type Storage interface {
	CreateStudent(name string, age int, email string) (int64, error)
	GetStudentById(id int) (types.Student, error)
	IsEmailExists(email string) (bool, error)
}
