package types

type APIError struct {
	Error string `json:"error"`
}

type RegisterStudentsPayload struct {
	Teacher  string   `json:"teacher"`
	Students []string `json:"students"`
}

type SuspendStudentPayload struct {
	Student string `json:"student"`
}

type GetNotifiedStudentsPayload struct {
	Teacher      string `json:"teacher"`
	Notification string `json:"notification"`
}
