package middleware

import (
	"fmt"

	"github.com/ellofae/deanery-gateway/config"
	"github.com/golang-jwt/jwt/v5"
)

type AccessTokenData struct {
	Expiry     int64
	IssuedAt   int64
	RecordCode string
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
			State:      claims["state"].(string),
		}, nil
	} else {
		return nil, err
	}
}

// func AuthenticateUser(c *gin.Context) {
// 	storage := repository.SessionStorage()

// 	session, err := storage.Get(c.Request, "session")
// 	if err != nil {
// 		response_errors.NewErrorResponse(c, http.StatusInternalServerError, "Unable to get session")
// 		return
// 	}

// 	sessionValue, ok := session.Values["access_token"]
// 	if !ok {
// 		response_errors.NewErrorResponse(c, http.StatusUnauthorized, "Authorization data field is empty")
// 		return
// 	}

// 	jwtString := strings.Split(sessionValue.(string), "Bearer ")
// 	if len(jwtString) < 2 {
// 		response_errors.NewErrorResponse(c, http.StatusInternalServerError, "Must provide Authorization data with format `Bearer {token}`")
// 		return
// 	}

// 	token_claims, err := ParseToken(jwtString[1])
// 	if err != nil {
// 		response_errors.NewErrorResponse(c, http.StatusInternalServerError, "Incorrect access token provided")
// 		return
// 	}

// 	expiry := token_claims.Expiry
// 	if expiry < time.Now().Unix() {
// 		response_errors.NewErrorResponse(c, http.StatusInternalServerError, "Token expired")
// 		return
// 	}

// 	session.Values["user_id"] = token_claims.UserID
// 	err = session.Save(c.Request, c.Writer)
// 	if err != nil {
// 		response_errors.NewHTTPResponse(c, http.StatusInternalServerError, "Unable to save session data", err)
// 		return
// 	}

// 	c.Next()
// }
