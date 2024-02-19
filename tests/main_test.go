package tests

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"

	"github.com/huilunang/OneCV-Assgn/api"
	"github.com/huilunang/OneCV-Assgn/storage"
)

type StoreWrapper struct {
	*storage.PostgreStore
}

var (
	server *api.APIServer
	store  *StoreWrapper
)

func TestMain(m *testing.M) {
	setUp()
	defer tearDown()

	os.Exit(m.Run())
}

func setUp() {
	godotenv.Load("../.env")

	ps, err := storage.NewPostGreStore(os.Getenv("DB_TEST_CONN_STR"))
	log.Print(os.Getenv("DB_TEST_CONN_STR"))
	if err != nil {
		log.Fatalf("error creating database connection: %v", err)
	}

	if err := ps.CreateStudentTable(); err != nil {
		log.Fatalf("error initializing database: %v", err)
	}
	if err := ps.CreateTeacherTable(); err != nil {
		log.Fatalf("error initializing database: %v", err)
	}

	store = &StoreWrapper{ps}
	server = api.NewAPIServer(":3001", store)
}

func tearDown() {
	if store != nil && store.PostgreStore != nil && store.PostgreStore.Db != nil {
		store.PostgreStore.Db.Close()
	}
}

func GetServer() *api.APIServer {
	return server
}

func ResetDatabaseState() error {
	if err := store.clearTables(); err != nil {
		log.Fatalf("error clearing database: %v", err)
	}

	if err := store.populateTables(); err != nil {
		log.Fatalf("error seeding database: %v", err)
	}

	return nil
}

func (s *StoreWrapper) clearTables() error {
	if _, err := s.Db.Exec("TRUNCATE TABLE student, teacher"); err != nil {
		return err
	}

	return nil
}

func (s *StoreWrapper) populateTables() error {
	timeN := time.Now().UTC()

	queries := []string{
		`
		INSERT INTO student (email, suspend_status, teachers, created_at)
		VALUES ('studentmario@gmail.com', true, ARRAY['teacherrosa@gmail.com', 'teacherwario@gmail.com'], $1),
			   ('commonstudentpeach@gmail.com', false, ARRAY['teacherrosa@gmail.com'], $2)
		`,
		`
		INSERT INTO teacher (email, created_at)
		VALUES ('teacherrosa@gmail.com', $1),
			   ('teacherwario@gmail.com', $2)
		`,
	}

	for _, query := range queries {
		stmt, err := s.Db.Prepare(query)
		if err != nil {
			return err
		}

		if _, err := stmt.Exec(timeN, timeN); err != nil {
			return err
		}
		
		stmt.Close()
	}

	return nil
}
