package transport

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"rsbruce/blogsite-api/internal/models"
)

type PostPageHandler struct {
	GetPostPage func(w http.ResponseWriter, r *http.Request)
}

func NewPostPageHandler(post_svc *models.PostService) PostPageHandler {
	var handler PostPageHandler
	handler.GetPostPage = func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		slug := params["slug"]

		post_with_user, err := post_svc.Store.GetPostWithUser(slug)
		if err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(w).Encode(post_with_user)
	}

	return handler
}
