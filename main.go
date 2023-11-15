package main

import (
	"database/sql"
	"encoding/json"
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

var db *sql.DB
var dtbs *database.Database

func singlePostHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	fmt.Println(params["slug"])
	postWithUser, err := postBySlugWithUser(params["slug"])
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(postWithUser)
}
func latestPostsHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := latestPosts()
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(posts)
}
func singleUserProfileHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	user, err := userByHandle(params["handle"])
	if err != nil {
		log.Fatal(err)
	}

	latestPosts, err := postItemsByUserHandle(params["handle"])
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(UserProfile{User: *user, LatestPosts: latestPosts})
}
func aboutHandler(w http.ResponseWriter, r *http.Request) {
	about, err := getAboutContent()
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(about)
}

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
	r.HandleFunc("/user/{handle}", userProfileHandler.Get)

	r.HandleFunc("/post/{slug}", singlePostHandler)

}

func postItemsByUserHandle(handle string) ([]FeedItemPost, error) {
	var feedItemPosts []FeedItemPost
	rows, err := db.Query(
		`SELECT post.title, post.subtitle, post.slug, post.created_at
        FROM post
        INNER JOIN user
        ON post.author_id = user.id
        WHERE user.handle = ?`, handle)
	if err != nil {
		return nil, fmt.Errorf("postItemsByUserHandle %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		feedItemPost := FeedItemPost{}
		err = rows.Scan(&feedItemPost.Title, &feedItemPost.Subtitle, &feedItemPost.Slug, &feedItemPost.Created_at)
		if err != nil {
			return nil, fmt.Errorf("postItemsByUserHandle %v", err)
		}
		feedItemPosts = append(feedItemPosts, feedItemPost)
	}

	return feedItemPosts, nil
}

func getAboutContent() (string, error) {
	rows, err := db.Query("SELECT content FROM text_content WHERE slug = \"about\"")
	var about string
	if err != nil {
		return "", fmt.Errorf("postsBySlug %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&about); err != nil {
			return "", fmt.Errorf("postsBySlug %v", err)
		}
	}

	return about, nil
}

func postBySlugWithUser(slug string) (*PostWithUser, error) {
	var postwithUser PostWithUser

	rows, err := db.Query(
		`SELECT post.title, post.subtitle, post.content, post.created_at, user.display_name, user.display_picture, user.handle 
        FROM post 
        INNER JOIN user on post.author_id = user.id 
        WHERE post.slug = ?`, slug)
	if err != nil {
		return nil, fmt.Errorf("postBySlugWithUser %q: %v", slug, err)
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(
			&postwithUser.Post.Title,
			&postwithUser.Post.Subtitle,
			&postwithUser.Post.Content,
			&postwithUser.Post.Created_at,
			&postwithUser.User.Display_name,
			&postwithUser.User.Display_picture,
			&postwithUser.User.Handle); err != nil {
			return nil, fmt.Errorf("postBySlugWithUser %q: %v", slug, err)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("postBySlugWithUser %q: %v", slug, err)
	}

	fmt.Println(postwithUser)

	return &postwithUser, nil
}

func postsBySlug(slug string) ([]Post, error) {
	var posts []Post

	rows, err := db.Query("SELECT id, author_id, title, subtitle, content, created_at FROM post WHERE slug = ?", slug)
	if err != nil {
		return nil, fmt.Errorf("postsBySlug %q: %v", slug, err)
	}
	defer rows.Close()
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.Author_id, &post.Title, &post.Subtitle, &post.Content, &post.Created_at); err != nil {
			return nil, fmt.Errorf("postsBySlug %q: %v", slug, err)
		}
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("postsBySlug %q: %v", slug, err)
	}

	fmt.Println(posts)

	return posts, nil
}
func latestPosts() ([]FeedItem, error) {
	var feedItems []FeedItem
	var author_display_picture sql.NullString

	rows, err := db.Query(
		`SELECT post.title, post.subtitle, post.created_at, post.slug, user.display_name, user.display_picture, user.handle 
                FROM post 
                INNER JOIN user on post.author_id = user.id
                ORDER BY created_at DESC
                LIMIT 5;`)
	if err != nil {
		return nil, fmt.Errorf("latestPosts: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		feedItem := FeedItem{Post: FeedItemPost{}, User: FeedItemUser{}}
		if err := rows.Scan(&feedItem.Post.Title, &feedItem.Post.Subtitle, &feedItem.Post.Created_at, &feedItem.Post.Slug, &feedItem.User.Display_name, &author_display_picture, &feedItem.User.Handle); err != nil {
			return nil, fmt.Errorf("latestPosts %v", err)
		}
		if author_display_picture.Valid {
			feedItem.User.Display_picture = author_display_picture.String
		}
		feedItems = append(feedItems, feedItem)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("latestPosts %v", err)
	}

	return feedItems, nil
}
func userByHandle(handle string) (*User, error) {
	var user User

	rows, err := db.Query(
		`SELECT id, handle, blurb, display_name, display_picture, created_at
            FROM user
            WHERE handle = ?`, handle)
	if err != nil {
		return nil, fmt.Errorf("usersByHandle %q: %v", handle, err)
	}
	defer rows.Close()

	for rows.Next() {

		err = rows.Scan(
			&user.ID,
			&user.Handle,
			&user.Blurb,
			&user.Display_name,
			&user.Display_picture,
			&user.Created_at)
		if err != nil {
			return nil, fmt.Errorf("usersByHandle %q: %v", handle, err)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("usersByHandle %q: %v", handle, err)
	}

	fmt.Println(user)

	return &user, nil
}

type Post struct {
	ID             []byte `json:"id"`
	Author_id      []byte `json:"author_id"`
	Slug           string `json:"slug"`
	Title          string `json:"title"`
	Subtitle       string `json:"subtitle"`
	Content        string `json:"content"`
	Main_image_url string `json:"main_image"`
	Created_at     string `json:"created_at"`
}
type User struct {
	ID              []byte `json:"id"`
	Handle          string `json:"handle"`
	Blurb           string `json:"blurb"`
	Display_name    string `json:"display_name"`
	Display_picture string `json:"display_picture"`
	User_role       string `json:"user_role"`
	Created_at      string `json:"created_at"`
}
type PostWithUser struct {
	Post Post `json:"post"`
	User User `json:"user"`
}
type FeedItemPost struct {
	Title      string `json:"title"`
	Subtitle   string `json:"subtitle"`
	Created_at string `json:"created_at"`
	Slug       string `json:"slug"`
}
type FeedItemUser struct {
	Handle          string `json:"handle"`
	Display_picture string `json:"display_picture"`
	Display_name    string `json:"display_name"`
}
type FeedItem struct {
	Post FeedItemPost `json:"post"`
	User FeedItemUser `json:"user"`
}
type UserProfile struct {
	User        User           `json:"user"`
	LatestPosts []FeedItemPost `json:"posts"`
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

	dtbs, err := database.NewDatabase()
	if err != nil {
		log.Fatal(err)
	}

	db = dtbs.Client.DB

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

	r := mux.NewRouter()

	setupRoutes(r, dtbs)

	serveAddress := ":" + os.Getenv("SERVE_PORT")

	if err := http.ListenAndServe(serveAddress, r); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
