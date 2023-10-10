package jwt

import (
	"errors"
	"time"

	"armada.api.com/internal/config"
	jwt "github.com/golang-jwt/jwt"
)

type JwtCustomClaims struct {
	UserID   string
	Email    string
	password string
	jwt.StandardClaims
}

type jwtConfig struct {
	secretKey string
	issuer    string
	expired   int
}

var jwtConfigEnv = jwtConfig{
	secretKey: config.AppConfig.JwtSecretKey,
	issuer:    config.AppConfig.JwtIsuuer,
	expired:   config.AppConfig.JwtExpired,
}

func GenerateToken(userID string, email string, password string) (t string, err error) {
	claims := &JwtCustomClaims{
		userID,
		email,
		password,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(jwtConfigEnv.expired)).Unix(),
			Issuer:    jwtConfigEnv.issuer,
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err = token.SignedString(jwtConfigEnv.secretKey)

	return
}

func PasreToken(tokenString string) (claims JwtCustomClaims, err error) {
	if token, err := jwt.ParseWithClaims(tokenString, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtConfigEnv.secretKey), nil
	}); err != nil || !token.Valid {
		return JwtCustomClaims{}, errors.New("token is not valid")
	}
	return
}
