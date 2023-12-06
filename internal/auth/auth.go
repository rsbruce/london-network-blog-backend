package auth

import (
	"rsbruce/blogsite-api/internal/database"
	"rsbruce/blogsite-api/internal/transport"

	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

type AuthHandler struct {
	store   *sessions.CookieStore
	DB_conn *database.Database
}

type ResponseMessage transport.ResponseMessage

type AuthResponse struct {
	Status string `json:"status"`
	Auth   string `json:"auth"`
	UserId int64  `json:"userId"`
}
type AuthCheck struct {
	Message string `json:"message"`
	Handle  string `json:"handle"`
}

func NewAuthHandler(db *database.Database) *AuthHandler {
	return &AuthHandler{
		store:   sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY"))),
		DB_conn: db,
	}
}

func (ah *AuthHandler) CanAccessUser(original func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		decoder := json.NewDecoder(r.Body)
		var user database.User
		err := decoder.Decode(&user)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ResponseMessage{Message: "Invalid JSON payload for this route."})
			return
		}

		session, err := ah.store.Get(r, "ks_session")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		id := session.Values["id"].(int64)
		authenticated := session.Values["authenticated"].(bool)

		if id == user.ID && authenticated {
			original(w, r)
		} else {
			http.Error(w, "Invalid credentials", http.StatusBadRequest)
		}
	}
}

func (ah *AuthHandler) CanAccessPost(original func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		decoder := json.NewDecoder(r.Body)
		var post database.Post
		err := decoder.Decode(&post)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ResponseMessage{Message: "Invalid JSON payload for this route."})
			return
		}

		session, err := ah.store.Get(r, "ks_session")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		id := session.Values["id"].(int64)
		authenticated := session.Values["authenticated"].(bool)

		if id == post.Author_id && authenticated {
			original(w, r)
		} else {
			http.Error(w, "Invalid credentials", http.StatusBadRequest)
		}
	}
}

func (ah *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)

	var user_auth database.UserAuth
	err := decoder.Decode(&user_auth)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseMessage{Message: "Invalid JSON payload for this route."})
		return
	}

	id, err := ah.DB_conn.CheckPassword(user_auth)

	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ResponseMessage{Message: "Could not authenticate"})
		return
	}

	session, err := ah.store.Get(r, "ks_session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["id"] = id
	session.Values["authenticated"] = true
	session.Values["timestamp"] = time.Now().Unix()

	err = session.Save(r, w)

	json.NewEncoder(w).Encode(AuthResponse{
		Status: "Success",
		Auth:   w.Header().Get("Set-Cookie"),
		UserId: id,
	})
}

func (ah *AuthHandler) CheckAuth(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	requestId, err := strconv.Atoi(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(AuthCheck{Message: "Failed"})
		return
	}

	session, err := ah.store.Get(r, "ks_session")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(AuthCheck{Message: "Failed"})
		return
	}
	id, ok := session.Values["id"].(int64)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(AuthCheck{Message: "Failed"})
		return
	}

	authenticated, ok := session.Values["authenticated"].(bool)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(AuthCheck{Message: "Failed"})
		return
	}

	handle, err := ah.DB_conn.UserHandleFromId(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(AuthCheck{Message: "Failed"})
		return
	}

	if id == int64(requestId) && authenticated {
		json.NewEncoder(w).Encode(AuthCheck{Message: "Success", Handle: handle})
	} else {
		http.Error(w, "Invalid credentials", http.StatusBadRequest)
	}
}
