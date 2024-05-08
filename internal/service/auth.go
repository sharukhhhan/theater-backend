package service

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"theater/pkg/hasher"
	"time"
)

type TokenClaims struct {
	jwt.StandardClaims
}

type AuthService struct {
	hasher        hasher.PasswordHasher
	signKey       string
	tokenTTL      time.Duration
	adminUsername string
	adminPassword string
}

func NewAuthService(hasher hasher.PasswordHasher, signKey string, tokenTTL time.Duration, adminUsername, adminPassword string) *AuthService {
	return &AuthService{
		hasher:        hasher,
		signKey:       signKey,
		tokenTTL:      tokenTTL,
		adminUsername: adminUsername,
		adminPassword: adminPassword,
	}
}

func (s *AuthService) Login(username, password string) (string, error) {
	hashedPassword := s.hasher.Hash(password)
	if username != s.adminUsername || hashedPassword != s.adminPassword {
		return "", ErrInvalidUsernameOrPassword
	}

	token, err := s.GenerateToken()
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) GenerateToken() (string, error) {
	claims := TokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(s.tokenTTL).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(s.signKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (s *AuthService) ParseToken(tokenString string) error {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.signKey), nil
	})

	if err != nil {
		return err
	}

	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		if claims.ExpiresAt < time.Now().Unix() {
			return errors.New("token is expired")
		}
	} else {
		return errors.New("invalid token")
	}

	return nil
}
