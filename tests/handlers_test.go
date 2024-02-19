package tests

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/huilunang/OneCV-Assgn/types"
	"github.com/huilunang/OneCV-Assgn/utils"
)

func TestHandleRootServer(t *testing.T) {
	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)

	if err != nil {
		t.Error(err)
	}

	handler := http.HandlerFunc(makeHTTPHandleFunc(server.HandleRootServer))
	handler.ServeHTTP(rr, req)

	if status := rr.Result().StatusCode; status != http.StatusOK {
		t.Errorf("expected 200 but got %v", status)
	}
	defer rr.Result().Body.Close()

	expected := `"Welcome to OneCV Technical Assessment"`
	b, err := io.ReadAll(rr.Result().Body)
	res := string(b)
	res = strings.TrimSpace(res)

	if err != nil {
		t.Errorf("failed to read response body: %v", err)
	}

	if res != expected {
		t.Errorf("expected %s but we got %s", expected, res)
	}
}

func TestHandleRegisterStudents(t *testing.T) {
	ResetDatabaseState()

	rr := httptest.NewRecorder()
	body := bytes.NewBuffer([]byte(`{
        "teacher": "teacherwario@gmail.com",
        "students": [
            "commonstudentpeach@gmail.com"
        ]
    }`))
	req, err := http.NewRequest("POST", "/api/register", body)

	if err != nil {
		t.Error(err)
	}

	handler := http.HandlerFunc(makeHTTPHandleFunc(server.HandleRegisterStudents))
	handler.ServeHTTP(rr, req)

	if status := rr.Result().StatusCode; status != http.StatusNoContent {
		t.Errorf("expected 204 but got %v", status)
		responseBody, _ := io.ReadAll(rr.Body)
		t.Logf("Response Body: %s", responseBody)
	}
}

func TestHandleGetCommonStudents(t *testing.T) {
	ResetDatabaseState()

	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/api/commonstudents?teacher=teacherrosa@gmail.com", nil)

	if err != nil {
		t.Error(err)
	}

	handler := http.HandlerFunc(makeHTTPHandleFunc(server.HandleGetCommonStudents))
	handler.ServeHTTP(rr, req)

	if status := rr.Result().StatusCode; status != http.StatusOK {
		t.Errorf("expected 200 but got %v", status)
	}
	defer rr.Result().Body.Close()

	expected := `{"students":["studentmario@gmail.com","commonstudentpeach@gmail.com"]}`
	b, err := io.ReadAll(rr.Result().Body)
	res := string(b)
	res = strings.TrimSpace(res)

	if err != nil {
		t.Errorf("failed to read response body: %v", err)
	}

	if res != expected {
		t.Errorf("expected %s but we got %s", expected, res)
	}
}

func TestHandleSuspendStudent(t *testing.T) {
	ResetDatabaseState()

	rr := httptest.NewRecorder()
	body := bytes.NewBuffer([]byte(`{"student": "commonstudentpeach@gmail.com"}`))
	req, err := http.NewRequest("POST", "/api/suspend", body)

	if err != nil {
		t.Error(err)
	}

	handler := http.HandlerFunc(makeHTTPHandleFunc(server.HandleSuspendStudent))
	handler.ServeHTTP(rr, req)

	if status := rr.Result().StatusCode; status != http.StatusNoContent {
		t.Errorf("expected 204 but got %v", status)
		responseBody, _ := io.ReadAll(rr.Body)
		t.Logf("Response Body: %s", responseBody)
	}
}

func TestHandleGetNotifiedStudents(t *testing.T) {
	ResetDatabaseState()

	rr := httptest.NewRecorder()
	body := bytes.NewBuffer([]byte(`{
		"teacher":  "teacherrosa@gmail.com",
		"notification": "Hello students! @studentmario@gmail.com"
	}`))
	req, err := http.NewRequest("POST", "/api/retrievefornotifications", body)

	if err != nil {
		t.Error(err)
	}

	handler := http.HandlerFunc(makeHTTPHandleFunc(server.HandleGetNotifiedStudents))
	handler.ServeHTTP(rr, req)

	if status := rr.Result().StatusCode; status != http.StatusOK {
		t.Errorf("expected 200 but got %v", status)
	}
	defer rr.Result().Body.Close()

	expected := `{"recipients":["commonstudentpeach@gmail.com"]}`
	b, err := io.ReadAll(rr.Result().Body)
	res := string(b)
	res = strings.TrimSpace(res)

	if err != nil {
		t.Errorf("failed to read response body: %v", err)
	}

	if res != expected {
		t.Errorf("expected %s but we got %s", expected, res)
	}
}

type apiFunc func(http.ResponseWriter, *http.Request) error

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			utils.WriteAPIJSON(w, http.StatusBadRequest, types.APIError{Error: err.Error()})
		}
	}
}
