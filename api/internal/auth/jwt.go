package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Service struct {
	secret []byte
	exp    time.Duration
}

func New(secret string, exp time.Duration) *Service {
	return &Service{secret: []byte(secret), exp: exp}
}

func (s *Service) Generate(userID, email string) (string, error) {
	claims := jwt.MapClaims{
		"sub":   userID,
		"email": email,
		"exp":   time.Now().Add(s.exp).Unix(),
		"iat":   time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

func (s *Service) Validate(tokenStr string) (jwt.MapClaims, error) {
	tok, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("método de assinatura inválido")
		}
		return s.secret, nil
	})
	if err != nil || !tok.Valid {
		return nil, errors.New("token inválido")
	}
	if claims, ok := tok.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}
	return nil, errors.New("claims inválidas")
}
