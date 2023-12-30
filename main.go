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

func CorsMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}

func JSONMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func setupRoutes(r *mux.Router) {

	r.Use(CorsMiddleWare)

	// PRE FLIGHT REQUESTS
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {}).Methods("OPTIONS")

	// STATIC FILES
	staticFileSubrouter := r.PathPrefix("/static").Subrouter()
	fs := http.FileServer(http.Dir("./static"))
	staticFileSubrouter.NewRoute().Handler(http.StripPrefix("/static/", fs))

	// AUTH
	authSubrouter := r.NewRoute().Subrouter()
	authSubrouter.Use(JSONMiddleware)
	authSubrouter.HandleFunc("/login", authRoutesService.Login).Methods("POST")
	authSubrouter.HandleFunc("/user-handle", authRoutesService.UserHandle).Methods("GET")
	authSubrouter.HandleFunc("/refresh-access", authRoutesService.RefreshAccess).Methods("POST")
	authSubrouter.HandleFunc("/reset-password", authRoutesService.ResetPassword).Methods("POST")

	// RESOURCES
	resourceSubrouter := r.NewRoute().Subrouter()
	resourceSubrouter.Use(JSONMiddleware)
	// CREATE
	resourceSubrouter.HandleFunc("/post", resourceRoutesService.CreatePost).Methods("POST")
	resourceSubrouter.HandleFunc("/post/{slug}/main-image", resourceRoutesService.UpdatePostImage).Methods("POST")
	resourceSubrouter.HandleFunc("/display-picture", resourceRoutesService.UpdateDisplayPicture).Methods("POST")
	// READ
	resourceSubrouter.HandleFunc("/feed", resourceRoutesService.GetFeed).Methods("GET")
	resourceSubrouter.HandleFunc("/feed/{handle}", resourceRoutesService.GetSingleUserFeed).Methods("GET")
	resourceSubrouter.HandleFunc("/personal-feed", resourceRoutesService.GetPersonalFeed).Methods("GET")
	resourceSubrouter.HandleFunc("/post/{handle}/{slug}", resourceRoutesService.GetPost).Methods("GET")
	resourceSubrouter.HandleFunc("/text-content/{slug}", resourceRoutesService.GetTextContent).Methods("GET")
	resourceSubrouter.HandleFunc("/user/{handle}", resourceRoutesService.GetUser).Methods("GET")
	// UPDATE
	resourceSubrouter.HandleFunc("/post/{slug}", resourceRoutesService.EditPost).Methods("PUT")
	resourceSubrouter.HandleFunc("/post/restore/{slug}", resourceRoutesService.RestorePost).Methods("PUT")
	resourceSubrouter.HandleFunc("/post/archive/{slug}", resourceRoutesService.ArchivePost).Methods("PUT")
	resourceSubrouter.HandleFunc("/user", resourceRoutesService.EditUser).Methods("PUT")
	// DELETE
	resourceSubrouter.HandleFunc("/post/{slug}", resourceRoutesService.DeletePost).Methods("DELETE")

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
