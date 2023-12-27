package main

import (
	"fmt"
	"log"
	"os"

	"net/http"

	"rsbruce/blogsite-api/internal/authdata"
	"rsbruce/blogsite-api/internal/authroutes"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
)

var authDataService *authdata.Service
var authRoutesService *authroutes.Service

func setupRoutes(r *mux.Router) {

	fs := http.FileServer(http.Dir("./static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	r.HandleFunc("/login", authRoutesService.Login).Methods("POST")
	r.HandleFunc("/user-handle", authRoutesService.UserHandle).Methods("GET")

	// r.HandleFunc("/latest-posts", handler.GetLatestAllAuthors).Methods("GET")
	// r.HandleFunc("/latest-posts/{handle}", handler.GetLatestForAuthor).Methods("GET")
	// r.HandleFunc("/post/{handle}/{slug}", handler.GetPostPage).Methods("GET")
	// r.HandleFunc("/slugs/{handle}", handler.GetSlugsForUser).Methods("GET")
	// r.HandleFunc("/text-content/{slug}", handler.GetTextContent).Methods("GET")
	// r.HandleFunc("/user/{handle}", handler.GetUserProfile).Methods("GET")

	// r.HandleFunc("/user", handler.UpdateUserProfile).Methods("PUT")
	// r.HandleFunc("/reset-password", handler.UpdatePassword).Methods("PUT")
	// r.HandleFunc("/profile-picture", handler.UploadProfilePicture).Methods("POST")
	// r.HandleFunc("/post", handler.NewPost).Methods("POST")
	// r.HandleFunc("/post", handler.UpdatePost).Methods("PUT")
	// r.HandleFunc("/post/{id}", handler.DeletePost).Methods("DELETE")
	// r.HandleFunc("/post/archive/{id}", handler.ArchivePost).Methods("PUT")
	// r.HandleFunc("/post/restore/{id}", handler.RestorePost).Methods("PUT")

}

func NewDbConnection() (*sqlx.DB, error) {
	log.Println("Setting up new database connection")

	cfg := mysql.Config{
		User:   os.Getenv("DB_USER"),
		Passwd: os.Getenv("DB_PASS"),
		Net:    "tcp",
		Addr:   os.Getenv("DB_HOST"),
		DBName: os.Getenv("DB_NAME"),
	}

	db, err := sqlx.Connect("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("could not connect to database: %w", err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		return nil, fmt.Errorf("could not connect to database: %w", err)
	}

	return db, nil
}

func main() {
	file, err := os.OpenFile("log.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Println("App started")

	defer func() {
		file.Close()
	}()

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	authDb, err := NewDbConnection()
	// resourceDb, err := NewDbConnection()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected!")

	authDataService = &authdata.Service{DbConn: authDb}
	authRoutesService = &authroutes.Service{AuthDataService: authDataService}

	r := mux.NewRouter()
	setupRoutes(r)

	serveAddress := ":" + os.Getenv("SERVE_PORT")
	if err := http.ListenAndServe(serveAddress, r); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
