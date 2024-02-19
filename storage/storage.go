package storage

import "github.com/huilunang/OneCV-Assgn/types"

type Storage interface {
	RegisterStudents(*types.RegisterStudentsPayload) error
	GetCommonStudents([]string) ([]string, error)
	SuspendStudent(string) error
	GetNotifiedStudents(*types.GetNotifiedStudentsPayload) ([]string, error)
}
