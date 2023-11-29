package main

import (
	"fmt"
	"log"
	"os"

	"net/http"

	"rsbruce/blogsite-api/internal/database"
	"rsbruce/blogsite-api/internal/transport"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
)

func setupRoutes(r *mux.Router, db *database.Database) {

	handler := transport.NewHttpHandler(db)

	r.PathPrefix("/").HandlerFunc(handler.HandleCors).Methods("OPTIONS")

	r.HandleFunc("/text-content/{slug}", handler.GetTextContent)

	r.HandleFunc("/latest-posts/{handle}", handler.GetLatestForAuthor)
	r.HandleFunc("/latest-posts", handler.GetLatestAllAuthors)

	r.HandleFunc("/user/{handle}", handler.GetUserProfile).Methods("GET")
	r.HandleFunc("/user/{handle}", handler.UpdateUserProfile).Methods("PUT")

	r.HandleFunc("/post/{slug}", handler.GetPostPage).Methods("GET")
	r.HandleFunc("/new-post", handler.NewPost).Methods("POST")

	r.HandleFunc("/slugs/{handle}", handler.GetSlugsForUser)

}

func main() {
	file, err := os.OpenFile("log.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Println("App started")

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	defer func() {
		file.Close()
	}()

	db, err := database.NewDatabase()
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Client.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

	r := mux.NewRouter()

	setupRoutes(r, db)

	serveAddress := ":" + os.Getenv("SERVE_PORT")
	if err := http.ListenAndServe(serveAddress, r); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
