package resourceroutes

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"rsbruce/blogsite-api/internal/resourcedata"
)

func (svc *Service) GetPost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	authorHandle, ah := params["handle"]
	slug, s := params["slug"]

	if !(ah && s) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	post, err := svc.ResourceData.GetPost(authorHandle, slug)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(post)

}

func (svc *Service) CreatePost(w http.ResponseWriter, r *http.Request) {
	var post resourcedata.Post
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&post)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := svc.AuthData.GetUserId(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	post.AuthorID = id
	post.CreatedAt = time.Now().Format(time.DateTime)

	err = svc.ResourceData.CreatePost(post)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (svc *Service) EditPost(w http.ResponseWriter, r *http.Request) {
	var post resourcedata.Post
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&post)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := svc.AuthData.GetUserId(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	post.AuthorID = id
	post.UpdatedAt = time.Now().Format(time.DateTime)
	err = svc.ResourceData.UpdatePost(post)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
	}

}

func (svc *Service) DeletePost(w http.ResponseWriter, r *http.Request) {

}

func (svc *Service) RestorePost(w http.ResponseWriter, r *http.Request) {

}
