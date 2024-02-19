package api

import (
	"net/http"

	"github.com/huilunang/OneCV-Assgn/types"
	"github.com/huilunang/OneCV-Assgn/utils"
)

func (s *APIServer) HandleRootServer(w http.ResponseWriter, r *http.Request) error {
	return utils.WriteAPIJSON(w, http.StatusOK, "Welcome to OneCV Technical Assessment")
}

func (s *APIServer) HandleRegisterStudents(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "POST":
		registerStudentsPayload := new(types.RegisterStudentsPayload)
		if err := utils.DecodeAPIJSON(r, registerStudentsPayload); err != nil {
			return err
		}

		if err := s.store.RegisterStudents(registerStudentsPayload); err != nil {
			return err
		}

		return utils.WriteAPIJSON(w, http.StatusNoContent, nil)
	default:
		return utils.WriteAPIJSON(
			w,
			http.StatusMethodNotAllowed,
			map[string]string{"error": "Requested method is not supported for the target endpoint"},
		)
	}
}

func (s *APIServer) HandleGetCommonStudents(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		r.ParseForm()
		teachers := r.Form["teacher"]

		if len(teachers) == 0 {
			return utils.WriteAPIJSON(
				w,
				http.StatusBadRequest,
				map[string]string{"error": "No teachers specified"},
			)
		}

		students, err := s.store.GetCommonStudents(teachers)
		if err != nil {
			return err
		}

		return utils.WriteAPIJSON(w, http.StatusOK, map[string][]string{"students": students})
	default:
		return utils.WriteAPIJSON(
			w,
			http.StatusMethodNotAllowed,
			map[string]string{"error": "Requested method is not supported for the target endpoint"},
		)
	}
}

func (s *APIServer) HandleSuspendStudent(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "POST":
		suspendStudentPayload := new(types.SuspendStudentPayload)
		if err := utils.DecodeAPIJSON(r, suspendStudentPayload); err != nil {
			return err
		}

		if err := s.store.SuspendStudent(suspendStudentPayload.Student); err != nil {
			return utils.WriteAPIJSON(
				w,
				http.StatusBadRequest,
				map[string]string{"error": "Requested student does not exists"},
			)
		}

		return utils.WriteAPIJSON(w, http.StatusNoContent, nil)
	default:
		return utils.WriteAPIJSON(
			w,
			http.StatusMethodNotAllowed,
			map[string]string{"error": "Requested method is not supported for the target endpoint"},
		)
	}
}

func (s *APIServer) HandleGetNotifiedStudents(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "POST":
		getNotifiedStudentsPayload := new(types.GetNotifiedStudentsPayload)
		if err := utils.DecodeAPIJSON(r, getNotifiedStudentsPayload); err != nil {
			return err
		}

		students, err := s.store.GetNotifiedStudents(getNotifiedStudentsPayload)
		if err != nil {
			return err
		}

		return utils.WriteAPIJSON(w, http.StatusOK, map[string][]string{"recipients": students})
	default:
		return utils.WriteAPIJSON(
			w,
			http.StatusMethodNotAllowed,
			map[string]string{"error": "Requested method is not supported for the target endpoint"},
		)
	}
}
