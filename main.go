package main

import (
	"fmt"
	"log"
	"os"

	"net/http"

	"rsbruce/blogsite-api/internal/authdata"
	"rsbruce/blogsite-api/internal/authroutes"
	"rsbruce/blogsite-api/internal/resourcedata"
	"rsbruce/blogsite-api/internal/resourceroutes"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
)

var authDataService *authdata.Service
var authRoutesService *authroutes.Service
var resourceDataService *resourcedata.Service
var resourceRoutesService *resourceroutes.Service

func setupRoutes(r *mux.Router) {

	fs := http.FileServer(http.Dir("./static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	r.HandleFunc("/login", authRoutesService.Login).Methods("POST")
	r.HandleFunc("/user-handle", authRoutesService.UserHandle).Methods("GET")
	r.HandleFunc("/refresh-access", authRoutesService.RefreshAccess).Methods("GET")
	r.HandleFunc("/reset-password", authRoutesService.ResetPassword).Methods("POST")

	// CREATE
	r.HandleFunc("/post", resourceRoutesService.CreatePost).Methods("POST")
	r.HandleFunc("/post/{slug}/main-image", resourceRoutesService.UploadPhoto).Methods("POST")
	r.HandleFunc("/display-picture", resourceRoutesService.UploadPhoto).Methods("POST")
	// READ
	r.HandleFunc("/feed", resourceRoutesService.GetFeed).Methods("GET")
	r.HandleFunc("/feed/{handle}", resourceRoutesService.GetSingleUserFeed).Methods("GET")
	r.HandleFunc("/personal-feed", resourceRoutesService.GetPersonalFeed).Methods("GET")
	r.HandleFunc("/post/{handle}/{slug}", resourceRoutesService.GetPost).Methods("GET")
	r.HandleFunc("/text-content/{slug}", resourceRoutesService.GetTextContent).Methods("GET")
	r.HandleFunc("/user/{handle}", resourceRoutesService.GetUser).Methods("GET")
	// UPDATE
	r.HandleFunc("/post/{slug}", resourceRoutesService.EditPost).Methods("PUT")
	r.HandleFunc("/post/restore/{slug}", resourceRoutesService.RestorePost).Methods("PUT")
	r.HandleFunc("/user", resourceRoutesService.EditUser).Methods("PUT")
	// DELETE
	r.HandleFunc("/post/{slug}", resourceRoutesService.DeletePost).Methods("DELETE")

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
	resourceDb, err := NewDbConnection()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected!")

	authDataService = &authdata.Service{DbConn: authDb}
	resourceDataService = &resourcedata.Service{DbConn: resourceDb}

	authRoutesService = &authroutes.Service{AuthData: authDataService}
	resourceRoutesService = &resourceroutes.Service{
		AuthData:     authDataService,
		ResourceData: resourceDataService,
	}

	r := mux.NewRouter()
	setupRoutes(r)

	serveAddress := ":" + os.Getenv("SERVE_PORT")
	if err := http.ListenAndServe(serveAddress, r); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
