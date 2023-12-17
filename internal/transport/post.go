package transport

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"rsbruce/blogsite-api/internal/auth"
	"rsbruce/blogsite-api/internal/database"
	"strconv"

	"github.com/gorilla/mux"
)

func (handler *HttpHandler) GetPostPage(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	slug := params["slug"]
	handle := params["handle"]

	post_with_user, err := handler.DB_conn.GetPostWithUser(slug, handle)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ResponseMessage{Message: fmt.Sprintf("Post not found with slug: %s", slug)})
		return
	}

	json.NewEncoder(w).Encode(post_with_user)
}

func (handler *HttpHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)

	var updated_post database.Post
	err := decoder.Decode(&updated_post)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseMessage{Message: "Invalid JSON payload for this route."})
		return
	}

	access_token := r.Header.Get("Authorization")
	user_claims, err := auth.ParseAccessToken(access_token)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ResponseMessage{Message: "Access token invalid"})
		return
	}
	token_author_id := user_claims.ID
	true_author_id, err := handler.DB_conn.GetAuthorIdFromPostId(updated_post.ID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseMessage{Message: "Could not find author for this post"})
		return
	}

	if token_author_id == true_author_id {
		post, err := handler.DB_conn.UpdatePost(updated_post)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ResponseMessage{Message: "Could not add post."})
			return
		}

		json.NewEncoder(w).Encode(post)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ResponseMessage{Message: "Not authorized to edit this post"})
	}

}

func (handler *HttpHandler) ArchivePost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseMessage{Message: "ID in path must be an integer"})
		return
	}

	if handler.canEditPost(r, id) {
		err = handler.DB_conn.ArchivePost(id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ResponseMessage{Message: "Could not restore post"})
			return
		}

		json.NewEncoder(w).Encode(ResponseMessage{Message: "Success"})
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ResponseMessage{Message: "Not authorized to edit this post"})
	}
}

func (handler *HttpHandler) RestorePost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseMessage{Message: "ID in path must be an integer"})
		return
	}

	if handler.canEditPost(r, id) {
		err = handler.DB_conn.RestorePost(id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ResponseMessage{Message: "Could not restore post"})
			return
		}

		json.NewEncoder(w).Encode(ResponseMessage{Message: "Success"})
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ResponseMessage{Message: "Not authorized to edit this post"})
	}
}

func (handler *HttpHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseMessage{Message: "ID in path must be an integer"})
		return
	}

	if handler.canEditPost(r, id) {
		err = handler.DB_conn.DeletePost(id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ResponseMessage{Message: "Could not delete post"})
			return
		}

		json.NewEncoder(w).Encode(ResponseMessage{Message: "Success"})
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ResponseMessage{Message: "Not authorized to edit this post"})
	}
}

func (handler *HttpHandler) NewPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)

	var new_post database.Post
	err := decoder.Decode(&new_post)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseMessage{Message: "Invalid JSON payload for this route."})
		return
	}

	access_token := r.Header.Get("Authorization")
	user_claims, err := auth.ParseAccessToken(access_token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ResponseMessage{Message: "Access token invalid"})
		return
	}
	new_post.Author_id = user_claims.ID

	post, err := handler.DB_conn.NewPost(new_post)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ResponseMessage{Message: "Could not add post."})
		return
	}

	json.NewEncoder(w).Encode(post)
}

func (handler *HttpHandler) canEditPost(r *http.Request, post_id int64) bool {
	access_token := r.Header.Get("Authorization")
	user_claims, err := auth.ParseAccessToken(access_token)
	if err != nil {
		log.Println(err)
		return false
	}
	token_author_id := user_claims.ID
	true_author_id, err := handler.DB_conn.GetAuthorIdFromPostId(post_id)
	if err != nil {
		log.Println(err)
		return false
	}

	return token_author_id == true_author_id
}
