package api

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/huilunang/OneCV-Assgn/storage"
	"github.com/huilunang/OneCV-Assgn/types"
	"github.com/huilunang/OneCV-Assgn/utils"
)

type APIServer struct {
	listenAddr string
	store      storage.Storage
}

func NewAPIServer(listenAddr string, store storage.Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/", makeHTTPHandleFunc(s.HandleRootServer))
	router.HandleFunc("/api/register", makeHTTPHandleFunc(s.HandleRegisterStudents))
	router.HandleFunc("/api/commonstudents", makeHTTPHandleFunc(s.HandleGetCommonStudents))
	router.HandleFunc("/api/suspend", makeHTTPHandleFunc(s.HandleSuspendStudent))
	router.HandleFunc("/api/retrievefornotifications", makeHTTPHandleFunc(s.HandleGetNotifiedStudents))

	log.Println("API server running on port: ", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			utils.WriteAPIJSON(w, http.StatusBadRequest, types.APIError{Error: err.Error()})
		}
	}
}
