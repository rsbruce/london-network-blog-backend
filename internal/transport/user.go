package transport

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"rsbruce/blogsite-api/internal/auth"
	"rsbruce/blogsite-api/internal/database"
	"strconv"
	"time"

	"github.com/gorilla/sessions"

	"github.com/gorilla/mux"
)

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

func (handler *HttpHandler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	params := mux.Vars(r)
	handle := params["handle"]

	author, err := handler.DB_conn.GetUser(handle)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ResponseMessage{Message: fmt.Sprintf("Could not get user with handle: %s", handle)})
		return
	}

	feed_item_posts, err := handler.DB_conn.GetAllFeedItemPostsForAuthor(handle)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ResponseMessage{Message: fmt.Sprintf("Failed to retrieve post items for author with handle: %s.", handle)})
		return
	}

	user_profile := database.UserProfile{
		User:        author,
		LatestPosts: feed_item_posts,
	}

	json.NewEncoder(w).Encode(user_profile)
}

func (handler *HttpHandler) UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	var request struct {
		User        database.User  `json:"user"`
		Auth_tokens auth.TokenPair `json:"auth_tokens"`
	}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseMessage{Message: "Invalid JSON payload for this route."})
		return
	}

	auth_tokens, err := request.Auth_tokens.GetNewTokenPair()
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseMessage{Message: "Bad tokens"})
		return
	}

	access_token := auth.ParseAccessToken(auth_tokens.AccessToken)
	request.User.ID = access_token.ID

	user, err := handler.DB_conn.UpdateUser(request.User)

	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ResponseMessage{Message: "Could not update user."})
		return
	}

	json.NewEncoder(w).Encode(user)
}

func (handler *HttpHandler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)

	var user_auth database.UserAuth
	err := decoder.Decode(&user_auth)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseMessage{Message: "Invalid JSON payload for this route."})
		return
	}

	err = handler.DB_conn.UpdatePassword(user_auth)

	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ResponseMessage{Message: "Could not update password."})
		return
	}

	json.NewEncoder(w).Encode(ResponseMessage{Message: "Success"})

}

func (handler *HttpHandler) UploadProfilePicture(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	r.ParseMultipartForm(32 << 20) // limit your max input length!
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		json.NewEncoder(w).Encode(ResponseMessage{Message: "Failed"})
		panic(err)
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		json.NewEncoder(w).Encode(ResponseMessage{Message: "Failed"})
		panic(err)
	}
	defer file.Close()

	uploadPath := fmt.Sprintf("static/%v/display-picture/%v.jpg", id, time.Now().Unix())

	outputFile, err := createPath(uploadPath)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, file)
	if err != nil {
		panic(err)
	}

	err = handler.DB_conn.UpdateDisplayPicture(int64(id), uploadPath)
	if err != nil {
		panic(err)
	}

	json.NewEncoder(w).Encode(ResponseMessage{Message: "Success"})
}

func createPath(p string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(p), 0666); err != nil {
		return nil, err
	}
	return os.Create(p)
}
