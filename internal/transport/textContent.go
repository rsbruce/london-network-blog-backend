package transport

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"rsbruce/blogsite-api/internal/service/textContent"
)

type TextContentHandler struct {
	Get func(w http.ResponseWriter, r *http.Request)
}

func NewTextContentHandler(service *textContent.Service) TextContentHandler {
	var handler TextContentHandler
	handler.Get = func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		slug := params["slug"]

		textContent, err := service.Store.GetTextContent(slug)
		if err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(w).Encode(textContent)
	}

	return handler
}
