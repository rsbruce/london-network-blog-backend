package main

import (
	"fmt"
	"log"
	"os"

	"net/http"

	"rsbruce/blogsite-api/internal/database"
	"rsbruce/blogsite-api/internal/models"
	"rsbruce/blogsite-api/internal/transport"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
)

func setupRoutes(r *mux.Router, db *database.Database) {

	textContentService := models.NewTextContentService(db)
	textContentHandler := transport.NewTextContentHandler(textContentService)
	r.HandleFunc("/text-content/{slug}", textContentHandler.Get)

	feedItemService := models.NewFeedItemService(db)
	feedItemHandler := transport.NewPostFeedItemHandler(feedItemService)
	r.HandleFunc("/latest-posts/{handle}", feedItemHandler.GetLatestForAuthor)
	r.HandleFunc("/latest-posts", feedItemHandler.GetLatestAllAuthors)

	userService := models.NewUserService(db)
	userProfileHandler := transport.NewUserProfileHandler(userService, feedItemService)
	r.HandleFunc("/user/{handle}", userProfileHandler.Get).Methods("GET")
	r.HandleFunc("/user/{handle}", userProfileHandler.Update).Methods("POST")

	r.PathPrefix("/").HandlerFunc(corsHandler).Methods("OPTIONS")

	postService := models.NewPostService(db)
	postServiceHandler := transport.NewPostHandler(postService)
	r.HandleFunc("/post/{slug}", postServiceHandler.GetPostPage).Methods("GET")
	r.HandleFunc("/new-post", postServiceHandler.NewPost).Methods("POST")

	slugsService := models.NewSlugsService(db)
	slugsHandler := transport.NewSlugsHandler(slugsService)
	r.HandleFunc("/slugs/{handle}", slugsHandler.Get)

}

func corsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "*")
	w.WriteHeader(http.StatusOK)
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
