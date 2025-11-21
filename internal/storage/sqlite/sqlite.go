package sqlite

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/vikramgurjar2/practice-project/internal/config"
	"github.com/vikramgurjar2/practice-project/internal/types"
)

type Sqlite struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*Sqlite, error) {

	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		return nil, err
	}

	// Create table
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS students (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT,
            age INTEGER,
            email TEXT
        )
    `)
	if err != nil {
		return nil, err
	}

	return &Sqlite{Db: db}, nil

}

func (s *Sqlite) CreateStudent(name string, age int, email string) (int64, error) {
	smt, err := s.Db.Prepare("INSERT INTO students (name, age, email) VALUES (?, ?, ?)")
	if err != nil {
		return 0, err
	}
	result, err := smt.Exec(name, age, email)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil

}

func (s *Sqlite) GetStudentById(id int) (types.Student, error) {
	var student types.Student
	row := s.Db.QueryRow(`select id, name, age, email from students where id =? `, id)
	//scan the result into the student struct
	err := row.Scan(&student.Id, &student.Name, &student.Age, &student.Email)
	if err != nil {
		return types.Student{}, err
	}
	return student, nil
}

func (s *Sqlite) IsEmailExists(email string) (bool, error) {
	var count int
	row := s.Db.QueryRow(`select count(*) from students where email = ?`, email)
	err := row.Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
func (s *Sqlite) GetStudents() ([]types.Student, error) {
	rows, err := s.Db.Query(`select id, name, age, email from students`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []types.Student
	for rows.Next() {
		var student types.Student
		err := rows.Scan(&student.Id, &student.Name, &student.Age, &student.Email)
		if err != nil {
			return nil, err
		}
		students = append(students, student)
	}

	return students, nil
}
