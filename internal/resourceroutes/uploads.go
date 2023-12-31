package resourceroutes

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
)

func (svc *Service) UpdatePostImage(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	slug := params["slug"]

	r.ParseMultipartForm(5 << 20) // limit your max input length!
	id, err := svc.AuthData.GetUserId(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()

	uploadPath := fmt.Sprintf("static/%v/%v.jpg", id, slug)

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

	err = svc.ResourceData.UpdatePostImage(id, slug, uploadPath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(uploadPath)
}

func (svc *Service) UpdateDisplayPicture(w http.ResponseWriter, r *http.Request) {

	r.ParseMultipartForm(5 << 20) // limit your max input length!
	id, err := svc.AuthData.GetUserId(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

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
