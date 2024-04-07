package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ellofae/deanery-gateway/config"
	"github.com/ellofae/deanery-gateway/core/session"
	"github.com/golang-jwt/jwt/v5"
)

type AccessTokenData struct {
	Expiry     int64
	IssuedAt   int64
	RecordCode string
	Role       string
	State      string
}

var jwtSecretKey string

func InitJWTSecretKey(cfg *config.Config) {
	jwtSecretKey = cfg.Authentication.JwtSecretToken
}

func ParseToken(tokenString string) (*AccessTokenData, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(jwtSecretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("unable to parse token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		expiry := claims["expiry"].(float64)
		issued_at := claims["issued_at"].(float64)

		return &AccessTokenData{
			Expiry:     int64(expiry),
			IssuedAt:   int64(issued_at),
			RecordCode: claims["record_code"].(string),
			Role:       claims["role"].(string),
			State:      claims["state"].(string),
		}, nil
	} else {
		return nil, err
	}
}

func AuthenticateMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		storage := session.SessionStorage()

		switch r.URL.Path {
		case "/users/login":
			next.ServeHTTP(w, r)
			return
		}

		session, err := storage.Get(r, "session")
		if err != nil {
			http.Error(w, "Unable to get the session, error: %v", http.StatusInternalServerError)
			return
		}

		sessionValue, ok := session.Values["access_token"]
		if !ok {
			http.Error(w, "Authorization data field is empty", http.StatusUnauthorized)
			return
		}

		jwtString := strings.Split(sessionValue.(string), "Bearer ")
		if len(jwtString) < 2 {
			http.Error(w, "Must provide Authorization data with format Bearer {token}", http.StatusBadRequest)
			return
		}

		tokenClaims, err := ParseToken(jwtString[1])
		if err != nil {
			http.Error(w, "Incorrect access token provided", http.StatusBadRequest)
			return
		}

		expiry := tokenClaims.Expiry
		if expiry < time.Now().Unix() {
			http.Error(w, "Token expired", http.StatusUnauthorized)
			return
		}

		session.Values["record_code"] = tokenClaims.RecordCode
		session.Values["role"] = tokenClaims.Role
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, "Unable to save session data", http.StatusInternalServerError)
			return
		}

		next.ServeHTTP(w, r)
	}
}
