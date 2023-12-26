package authdata

import (
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

func NewService(dbConn *sqlx.DB) *Service {
	return &Service{DbConn: dbConn}
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
	userClaims := UserClaims{
		ID: id,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Minute * 120).Unix(),
		},
	}
	signedAccessToken, err := NewAccessToken(userClaims)
	if err != nil {
		return "", "", err
	}

	refreshClaims := jwt.StandardClaims{
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Hour * 48).Unix(),
	}

	signedRefreshToken, err := NewRefreshToken(refreshClaims)
	if err != nil {
		return "", "", err
	}

	return signedAccessToken, signedRefreshToken, nil
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
