package storage

import (
	"database/sql"
	"fmt"
	"log"
	"regexp"

	"github.com/lib/pq"

	"github.com/huilunang/OneCV-Assgn/types"
	"github.com/huilunang/OneCV-Assgn/utils"
)

type PostgreStore struct {
	Db *sql.DB
}

func NewPostGreStore(connStr string) (*PostgreStore, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgreStore{Db: db}, nil
}

func (s *PostgreStore) Init() error {
	if err := s.CreateStudentTable(); err != nil {
		return err
	}

	if err := s.CreateTeacherTable(); err != nil {
		return err
	}

	if err := s.CreateStudents(); err != nil {
		return err
	}

	if err := s.CreateTeachers(); err != nil {
		return err
	}

	return nil
}

func (s *PostgreStore) CreateStudentTable() error {
	_, err := s.Db.Exec(`CREATE TABLE IF NOT EXISTS student (
		id SERIAL PRIMARY KEY,
		email VARCHAR(50),
		suspend_status BOOLEAN,
		teachers TEXT[],
		created_at TIMESTAMP
	)`)

	return err
}

func (s *PostgreStore) CreateTeacherTable() error {
	_, err := s.Db.Exec(`CREATE TABLE IF NOT EXISTS teacher (
		id SERIAL PRIMARY KEY,
		email VARCHAR(50),
		created_at TIMESTAMP
	)`)

	return err
}

func (s *PostgreStore) CreateStudents() error {
	studentArr := [3]string{"studentmario", "studentluigi", "commonstudentpeach"}
	query := `
	INSERT INTO student (email, suspend_status, teachers, created_at)
	VALUES ($1, $2, $3, $4)
	`
	stmt, err := s.prepareStmt(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, e := range studentArr {
		student := types.NewStudent(fmt.Sprintf("%s@google.com", e))

		_, err := stmt.Exec(
			student.Email,
			student.SuspendSatus,
			pq.StringArray(student.Teachers),
			student.CreatedAt,
		)

		if err != nil {
			return err
		}
	}

	return nil
}

func (s *PostgreStore) CreateTeachers() error {
	teacherArr := [3]string{"teacherrosa", "teacherwario", "teacherboo"}
	query := `
	INSERT INTO teacher (email, created_at)
	VALUES ($1, $2)
	`
	stmt, err := s.prepareStmt(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, e := range teacherArr {
		teacher := types.NewTeacher(fmt.Sprintf("%s@google.com", e))

		_, err := stmt.Exec(
			teacher.Email,
			teacher.CreatedAt,
		)

		if err != nil {
			return err
		}
	}

	return nil
}

func (s *PostgreStore) RegisterStudents(p *types.RegisterStudentsPayload) error {
	teacher := types.NewTeacher(p.Teacher)
	tstmt, err := s.prepareStmt(`SELECT EXISTS (SELECT 1 FROM teacher WHERE email = $1)`)
	if err != nil {
		return err
	}
	defer tstmt.Close()

	_, err = tstmt.Exec(teacher.Email)
	if err != nil {
		return fmt.Errorf("teacher %s is not found", teacher.Email)
	}

	sstmt, err := s.prepareStmt(`(SELECT id, teachers FROM student WHERE email = $1)`)
	if err != nil {
		return err
	}
	defer sstmt.Close()

	for _, studentEmail := range p.Students {
		student := types.NewStudent(studentEmail)
		err := sstmt.QueryRow(student.Email).Scan(&student.Id, (*pq.StringArray)(&student.Teachers))
		if err != nil {
			return fmt.Errorf("student %s is not found", studentEmail)
		}

		student.Teachers = append(student.Teachers, teacher.Email)
		usstmt, err := s.prepareStmt(`UPDATE student SET teachers = $1 WHERE id = $2`)
		if err != nil {
			return err
		}

		_, err = usstmt.Exec(pq.StringArray(student.Teachers), student.Id)
		if err != nil {
			return err
		}
		usstmt.Close()
	}

	return nil
}

func (s *PostgreStore) GetCommonStudents(p []string) ([]string, error) {
	commonStudents := []string{}

	stmt, err := s.prepareStmt(`(SELECT email FROM student WHERE teachers @> $1)`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(pq.StringArray(p))
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var student string
		err := rows.Scan(&student)

		if err != nil {
			return nil, err
		}

		commonStudents = append(commonStudents, student)
	}

	return commonStudents, nil
}

func (s *PostgreStore) SuspendStudent(p string) error {
	stmt, err := s.prepareStmt(`UPDATE student SET suspend_status = $1 WHERE email = $2`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(true, p)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgreStore) GetNotifiedStudents(p *types.GetNotifiedStudentsPayload) ([]string, error) {
	students := []string{}

	tstmt, err := s.prepareStmt(`(SELECT email FROM student WHERE $1 = ANY(teachers) AND suspend_status = $2)`)
	if err != nil {
		return nil, err
	}
	defer tstmt.Close()

	rows, err := tstmt.Query(p.Teacher, false)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var student string
		err := rows.Scan(&student)

		if err != nil {
			return nil, err
		}

		students = append(students, student)
	}

	re := regexp.MustCompile(`([^\s@]+@[^\s]+)`)
	mentions := re.FindAllString(p.Notification, -1)
	sstmt, err := s.prepareStmt(`(SELECT EXISTS (SELECT 1 FROM student WHERE email = $1 AND suspend_status = $2))`)
	if err != nil {
		return nil, err
	}
	defer sstmt.Close()

	for _, e := range mentions {
		if utils.VContains(students, e) {
			continue
		}

		var exists bool
		err := sstmt.QueryRow(e, false).Scan(&exists)
		if err != nil {
			log.Printf("error checking student: %v", err)
			continue
		}

		if exists {
			students = append(students, e)
		}
	}

	return students, nil
}

func (s *PostgreStore) prepareStmt(q string) (*sql.Stmt, error) {
	stmt, err := s.Db.Prepare(q)

	if err != nil {
		return nil, err
	}

	return stmt, nil
}
