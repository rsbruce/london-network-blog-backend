package resourceroutes

import (
	"encoding/json"
	"net/http"
	"rsbruce/blogsite-api/internal/authdata"
	"rsbruce/blogsite-api/internal/resourcedata"

	"github.com/gorilla/mux"
)

type Service struct {
	AuthData     *authdata.Service
	ResourceData *resourcedata.Service
}

func (svc *Service) GetTextContent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	slug := params["slug"]

	content, err := svc.ResourceData.GetTextContent(slug)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}

	json.NewEncoder(w).Encode(content)

}

func (svc *Service) GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	handle := params["handle"]

	user, err := svc.ResourceData.GetUser(handle)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}

	json.NewEncoder(w).Encode(user)
}

func (svc *Service) UploadPhoto(w http.ResponseWriter, r *http.Request) {

}

func (svc *Service) CreatePost(w http.ResponseWriter, r *http.Request) {
	var post resourcedata.Post
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&post)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	
	accessToken := r.Header.Get("Authorization")
	if accessToken == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	id, err := svc.AuthData.GetIdFromAccessToken(accessToken)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	post.ID = id

	// err = svc.ResourceData.CreatePost(post)

	// if err != nil {
	// 	w.WriteHeader(http.StatusBadRequest)
	// }

	json.NewEncoder(w).Encode(post)
}

func (svc *Service) EditPost(w http.ResponseWriter, r *http.Request) {

}

func (svc *Service) DeletePost(w http.ResponseWriter, r *http.Request) {

}

func (svc *Service) RestorePost(w http.ResponseWriter, r *http.Request) {

}

func (svc *Service) EditUser(w http.ResponseWriter, r *http.Request) {

}
