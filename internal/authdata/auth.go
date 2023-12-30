package authdata

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/jmoiron/sqlx"
	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	DbConn *sqlx.DB
}

type UserClaims struct {
	ID int64
	jwt.StandardClaims
}

type LoginCredentials struct {
	Handle   string
	Password string
}

func (svc *Service) CheckPassword(creds LoginCredentials) (int64, error) {

	password := []byte(creds.Password)

	var storedHash []byte
	var id int64
	row := svc.DbConn.QueryRow(`SELECT password, id FROM user WHERE handle = ?`, creds.Handle)
	err := row.Scan(&storedHash, &id)
	if err != nil {
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword(storedHash, password)

	return id, err
}

func (svc *Service) GenerateTokensFromId(id int64) (string, string, error) {
	signedAccessToken, err := NewAccessToken(id)
	if err != nil {
		return "", "", err
	}

	signedRefreshToken, err := NewRefreshToken()
	if err != nil {
		return "", "", err
	}

	return signedAccessToken, signedRefreshToken, nil
}

func (svc *Service) GetUserId(r *http.Request) (int64, error) {
	accessToken := r.Header.Get("Authorization")
	if accessToken == "" {
		return 0, errors.New("Empty access token")
	}

	userClaims, err := ParseAccessToken(accessToken)
	if err != nil {
		return 0, err
	}
	return userClaims.ID, nil
}

func (svc *Service) GetHandleFromAccessToken(accessToken string) (string, error) {
	var handle string

	userClaims, err := ParseAccessToken(accessToken)
	if err != nil {
		return "", err
	}
	id := userClaims.ID

	row := svc.DbConn.QueryRow(`SELECT handle FROM user WHERE id = ?`, id)
	err = row.Scan(&handle)
	if err != nil {
		return "", err
	}

	return handle, nil
}

func NewAccessToken(id int64) (string, error) {
	claims := UserClaims{
		ID: id,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return accessToken.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
}

func NewRefreshToken() (string, error) {
	claims := jwt.StandardClaims{
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Hour * 48).Unix(),
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return refreshToken.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
}

func ParseAccessToken(accessToken string) (*UserClaims, error) {
	if accessToken == "" {
		return nil, errors.New("Empty access token")
	}

	parsedAccessToken, err := jwt.ParseWithClaims(accessToken, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("TOKEN_SECRET")), nil
	})

	return parsedAccessToken.Claims.(*UserClaims), err
}

func RefreshTokenIsValid(refreshToken string) bool {
	parsedRefreshToken, err := jwt.ParseWithClaims(refreshToken, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("TOKEN_SECRET")), nil
	})

	return (err == nil) && (parsedRefreshToken.Claims.(*jwt.StandardClaims).Valid() == nil)
}

func (svc *Service) GenerateTokensFromRefresh(r *http.Request) (string, string, error) {
	accessToken := r.Header.Get("Authorization")
	if accessToken == "" {
		return "", "", errors.New("Empty access token")
	}

	decoder := json.NewDecoder(r.Body)
	var refresh struct {
		Token string `json:"token"`
	}
	err := decoder.Decode(&refresh)
	if err != nil {
		return "", "", errors.New("No refresh token found in refresh body")
	}

	if !RefreshTokenIsValid(refresh.Token) {
		return "", "", errors.New("Invalid refresh token")
	}

	accessClaims, _ := ParseAccessToken(accessToken)
	if accessClaims == nil {
		return "", "", errors.New("Invalid access token claims")
	}

	newRefreshToken, err := NewRefreshToken()
	newAccessToken, err := NewAccessToken(accessClaims.ID)
	if err != nil {
		return "", "", err
	}
	return newAccessToken, newRefreshToken, nil
}

func (svc *Service) UpdatePassword(userID int64, newPassword string) error {

	password := []byte(newPassword)
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, password)
	if err != nil {
		return err
	}

	_, err = svc.DbConn.Exec(`UPDATE user SET password = ? WHERE id = ?`, hashedPassword, userID)

	return err
}
