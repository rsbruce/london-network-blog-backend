package auth

import (
	"rsbruce/blogsite-api/internal/database"

	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/sessions"
	_ "github.com/joho/godotenv/autoload"
)

type UserClaims struct {
	ID int64
	jwt.StandardClaims
}

type TokenPair struct {
	AccessToken  string `json:"access_token" binding:"required"`
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type AuthHandler struct {
	store   *sessions.CookieStore
	DB_conn *database.Database
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

func (ah *AuthHandler) CheckAuth(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	var auth_response struct {
		Message string `json:"message"`
		Handle  string `json:"handle"`
	}

	access_token := r.Header.Get("Authorization")
	user_claims, err := ParseAccessToken(access_token)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		auth_response.Message = "Failed"
		json.NewEncoder(w).Encode(auth_response)
		return
	}
	id := user_claims.ID
	handle, err := ah.DB_conn.UserHandleFromId(id)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Authenticated, but could not find handle")
		return
	}

	auth_response.Message = "Success"
	auth_response.Handle = handle

	json.NewEncoder(w).Encode(auth_response)
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
		json.NewEncoder(w).Encode("Invalid JSON payload for this route.")
		return
	}

	id, err := ah.DB_conn.CheckPassword(user_auth)

	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Could not authenticate")
		return
	}

	userClaims := UserClaims{
		ID: id,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Minute * 120).Unix(),
		},
	}

	signedAccessToken, err := NewAccessToken(userClaims)
	if err != nil {
		log.Fatal("error creating access token")
	}

	refreshClaims := jwt.StandardClaims{
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Hour * 48).Unix(),
	}

	signedRefreshToken, err := NewRefreshToken(refreshClaims)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("error creating refresh token")
		return
	}

	handle, err := ah.DB_conn.UserHandleFromId(userClaims.ID)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Authenticated, but could not find handle")
		return
	}
	var auth_response struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		Handle       string `json:"handle"`
	}
	auth_response.AccessToken = signedAccessToken
	auth_response.RefreshToken = signedRefreshToken
	auth_response.Handle = handle

	json.NewEncoder(w).Encode(auth_response)
}

func (token_pair *TokenPair) GetNewTokenPair() (TokenPair, error) {
	var err error
	userClaims, err := ParseAccessToken(token_pair.AccessToken)
	if err != nil {
		return TokenPair{}, errors.New("Access Token invalid")
	}
	refreshClaims, err := ParseRefreshToken(token_pair.RefreshToken)
	if err != nil {
		return TokenPair{}, errors.New("Refresh Token invalid")
	}

	// refresh token is expired
	if refreshClaims.Valid() != nil {
		token_pair.RefreshToken, err = NewRefreshToken(*refreshClaims)
		if err != nil {
			return TokenPair{}, errors.New("Error creating refresh token")
		}
	}

	// access token is expired
	if userClaims.StandardClaims.Valid() != nil && refreshClaims.Valid() == nil {
		token_pair.AccessToken, err = NewAccessToken(*userClaims)
		if err != nil {
			return TokenPair{}, errors.New("Error creating access token")
		}
	}

	return *token_pair, nil
}

func NewAccessToken(claims UserClaims) (string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return accessToken.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
}

func NewRefreshToken(claims jwt.StandardClaims) (string, error) {
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return refreshToken.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
}

func ParseAccessToken(accessToken string) (*UserClaims, error) {
	parsedAccessToken, err := jwt.ParseWithClaims(accessToken, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("TOKEN_SECRET")), nil
	})

	return parsedAccessToken.Claims.(*UserClaims), err
}

func ParseRefreshToken(refreshToken string) (*jwt.StandardClaims, error) {
	parsedRefreshToken, err := jwt.ParseWithClaims(refreshToken, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("TOKEN_SECRET")), nil
	})

	return parsedRefreshToken.Claims.(*jwt.StandardClaims), err
}
