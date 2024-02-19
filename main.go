package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/huilunang/OneCV-Assgn/api"
	"github.com/huilunang/OneCV-Assgn/storage"
)

func main() {
	godotenv.Load(".env")
	
	store, err := storage.NewPostGreStore(os.Getenv("DB_CONN_STR"))
	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	server := api.NewAPIServer(":3000", store)
	server.Run()
}