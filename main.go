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

func cors(fs http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Headers", "*")
		w.Header().Add("Access-Control-Allow-Methods", "*")
		fs.ServeHTTP(w, r)
	}
}

func setupRoutes(r *mux.Router, db *database.Database) {

	handler := transport.NewHttpHandler(db)
	authHandler := auth.NewAuthHandler(db)
	userAuth := authHandler.CanAccessUser
	postAuth := authHandler.CanAccessPost

	fs := http.FileServer(http.Dir("./static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	r.PathPrefix("/").HandlerFunc(handler.HandleCors).Methods("OPTIONS")

	api := r.PathPrefix("/api").Subrouter()

	api.HandleFunc("/text-content/{slug}", handler.GetTextContent).Methods("GET")
	api.HandleFunc("/latest-posts/{handle}", handler.GetLatestForAuthor).Methods("GET")
	api.HandleFunc("/latest-posts", handler.GetLatestAllAuthors).Methods("GET")
	api.HandleFunc("/user/{handle}", handler.GetUserProfile).Methods("GET")
	api.HandleFunc("/login", authHandler.Login).Methods("POST")
	api.HandleFunc("/checkAuth/{id}", authHandler.CheckAuth)
	api.HandleFunc("/post/{slug}", handler.GetPostPage).Methods("GET")
	api.HandleFunc("/slugs/{handle}", handler.GetSlugsForUser).Methods("GET")

	// AUTH ROUTES
	api.HandleFunc("/user/{handle}", userAuth(handler.UpdateUserProfile)).Methods("PUT")
	api.HandleFunc("/new-password", userAuth(handler.UpdatePassword)).Methods("PUT")
	api.HandleFunc("/profile-picture/{id}", handler.UploadProfilePicture).Methods("POST")

	api.HandleFunc("/new-post", postAuth(handler.NewPost)).Methods("POST")
	api.HandleFunc("/post/{id}", postAuth(handler.UpdatePost)).Methods("PUT")
	api.HandleFunc("/post/{id}", postAuth(handler.DeletePost)).Methods("DELETE")
	api.HandleFunc("/post/archive/{id}", postAuth(handler.ArchivePost)).Methods("PUT")
	api.HandleFunc("/post/restore/{id}", postAuth(handler.RestorePost)).Methods("PUT")

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
