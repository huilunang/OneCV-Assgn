package types

import "time"

type Student struct {
	Id           int       `json:"id"`
	Email        string    `json:"email"`
	SuspendSatus bool      `json:"status"`
	Teachers     []string  `json:"teachers"`
	CreatedAt    time.Time `json:"createdAt"`
}

func NewStudent(email string) *Student {
	return &Student{
		Email:        email,
		SuspendSatus: false,
		Teachers:     []string{},
		CreatedAt:    time.Now().UTC(),
	}
}
