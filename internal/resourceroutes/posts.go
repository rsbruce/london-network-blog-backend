package resourceroutes

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"rsbruce/blogsite-api/internal/resourcedata"

	"github.com/gorilla/mux"
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
	var updated_post resourcedata.Post
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&updated_post)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !svc.AuthData.CanEditPost(r, updated_post.ID) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = svc.ResourceData.UpdatePost(updated_post)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (svc *Service) ArchivePost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	slug := params["slug"]

	userId, err := svc.AuthData.GetUserId(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = svc.ResourceData.ArchivePost(userId, slug)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}
}

func (svc *Service) RestorePost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	slug := params["slug"]

	userId, err := svc.AuthData.GetUserId(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = svc.ResourceData.RestorePost(userId, slug)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}
}

func (svc *Service) DeletePost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	slug := params["slug"]

	userId, err := svc.AuthData.GetUserId(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = svc.ResourceData.DeletePost(userId, slug)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}
}
