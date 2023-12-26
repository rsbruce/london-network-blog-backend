package authroutes

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"rsbruce/blogsite-api/internal/authdata"
)

type Service struct {
	router      *mux.Router
	dataService *authdata.Service
}

func NewService(router *mux.Router, ads *authdata.Service) *Service {
	return &Service{
		router:      router,
		dataService: ads,
	}
}

func (svc *Service) UserHandle(w http.ResponseWriter, r *http.Request) {
	// var auth_response struct {
	// 	Message string `json:"message"`
	// 	Handle  string `json:"handle"`
	// }

	// access_token := r.Header.Get("Authorization")
	// user_claims, err := ParseAccessToken(access_token)
	// if err != nil {
	// 	log.Print(err)
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	auth_response.Message = "Failed"
	// 	json.NewEncoder(w).Encode(auth_response)
	// 	return
	// }
	// id := user_claims.ID
	// handle, err := ah.DB_conn.UserHandleFromId(id)
	// if err != nil {
	// 	log.Print(err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	json.NewEncoder(w).Encode("Authenticated, but could not find handle")
	// 	return
	// }

	// auth_response.Message = "Success"
	// auth_response.Handle = handle

	// json.NewEncoder(w).Encode(auth_response)
}

func (svc *Service) Login(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var credentials authdata.LoginCredentials
	err := decoder.Decode(&credentials)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Invalid JSON payload for this route.")
		return
	}

	id, err := svc.dataService.CheckPassword(credentials)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Could not authenticate")
		return
	}

	accessToken, refreshToken, err := svc.dataService.GenerateTokensFromId(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("error creating refresh token")
		return
	}

	var auth_response struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
	auth_response.AccessToken = accessToken
	auth_response.RefreshToken = refreshToken

	json.NewEncoder(w).Encode(auth_response)
}
