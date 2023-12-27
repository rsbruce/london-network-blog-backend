package resourceroutes

import (
	"net/http"
	"rsbruce/blogsite-api/internal/authdata"
	"rsbruce/blogsite-api/internal/resourcedata"
)

type Service struct {
	AuthData     *authdata.Service
	ResourceData *resourcedata.Service
}

func (svc *Service) GetPost(w http.ResponseWriter, r *http.Request) {

}

func (svc *Service) GetFeed(w http.ResponseWriter, r *http.Request) {

}

func (svc *Service) GetTextContent(w http.ResponseWriter, r *http.Request) {

}

func (svc *Service) GetUser(w http.ResponseWriter, r *http.Request) {

}

func (svc *Service) UploadPhoto(w http.ResponseWriter, r *http.Request) {

}

func (svc *Service) CreatePost(w http.ResponseWriter, r *http.Request) {

}

func (svc *Service) EditPost(w http.ResponseWriter, r *http.Request) {

}

func (svc *Service) DeletePost(w http.ResponseWriter, r *http.Request) {

}

func (svc *Service) RestorePost(w http.ResponseWriter, r *http.Request) {

}

func (svc *Service) EditUser(w http.ResponseWriter, r *http.Request) {

}
