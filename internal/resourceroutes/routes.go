package resourceroutes

import (
	"encoding/json"
	"filepath"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"rsbruce/blogsite-api/internal/authdata"
	"rsbruce/blogsite-api/internal/resourcedata"
	"time"

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

	r.ParseMultipartForm(5 << 20) // limit your max input length!
	id, err := svc.AuthData.GetUserId(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	fmt.Println(r.FormValue("testkey"))

	file, _, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()

	uploadPath := fmt.Sprintf("static/%v/display-picture/%v.jpg", id, time.Now().Unix())

	outputFile, err := createPath(uploadPath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = svc.ResourceData.UpdateDisplayPicture(id, uploadPath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(uploadPath)
}

func createPath(p string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(p), 0666); err != nil {
		return nil, err
	}
	return os.Create(p)
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

}

func (svc *Service) DeletePost(w http.ResponseWriter, r *http.Request) {

}

func (svc *Service) RestorePost(w http.ResponseWriter, r *http.Request) {

}

func (svc *Service) EditUser(w http.ResponseWriter, r *http.Request) {

}
