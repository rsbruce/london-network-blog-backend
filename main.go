package main

import (
	"fmt"
	"log"
	"os"

	"net/http"

	"rsbruce/blogsite-api/internal/auth"
	"rsbruce/blogsite-api/internal/database"
	"rsbruce/blogsite-api/internal/transport"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
)

func setupRoutes(r *mux.Router, db *database.Database) {

	r.Use(transport.CorsMiddleWare)

	handler := transport.NewHttpHandler(db)
	authHandler := auth.NewAuthHandler(db)

	fs := http.FileServer(http.Dir("./static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	r.PathPrefix("/").HandlerFunc(handler.PreFlight).Methods("OPTIONS")

	r.HandleFunc("/login", authHandler.Login).Methods("POST")
	r.HandleFunc("/check-auth", authHandler.CheckAuth).Methods("GET")

	r.HandleFunc("/latest-posts", handler.GetLatestAllAuthors).Methods("GET")
	r.HandleFunc("/latest-posts/{handle}", handler.GetLatestForAuthor).Methods("GET")
	r.HandleFunc("/post/{handle}/{slug}", handler.GetPostPage).Methods("GET")
	r.HandleFunc("/slugs/{handle}", handler.GetSlugsForUser).Methods("GET")
	r.HandleFunc("/text-content/{slug}", handler.GetTextContent).Methods("GET")
	r.HandleFunc("/user/{handle}", handler.GetUserProfile).Methods("GET")

	r.HandleFunc("/user", handler.UpdateUserProfile).Methods("PUT")
	r.HandleFunc("/reset-password", handler.UpdatePassword).Methods("PUT")
	r.HandleFunc("/profile-picture", handler.UploadProfilePicture).Methods("POST")
	r.HandleFunc("/post", handler.NewPost).Methods("POST")
	r.HandleFunc("/post", handler.UpdatePost).Methods("PUT")
	r.HandleFunc("/post/{id}", handler.DeletePost).Methods("DELETE")
	r.HandleFunc("/post/archive/{id}", handler.ArchivePost).Methods("PUT")
	r.HandleFunc("/post/restore/{id}", handler.RestorePost).Methods("PUT")

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
