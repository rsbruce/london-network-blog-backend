package resourceroutes

import (
	"encoding/json"
	"log"
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

	var text struct {
		Content string `json:"content"`
	}
	text.Content = content
	json.NewEncoder(w).Encode(text)

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

func (svc *Service) EditUser(w http.ResponseWriter, r *http.Request) {
	var user resourcedata.User
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := svc.AuthData.GetUserId(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	user.ID = id

	err = svc.ResourceData.EditUser(user)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
