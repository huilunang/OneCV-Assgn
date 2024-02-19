package types

import "time"

type Teacher struct {
	Id        int       `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
}

func NewTeacher(email string) *Teacher {
	return &Teacher{
		Email:     email,
		CreatedAt: time.Now().UTC(),
	}
}
