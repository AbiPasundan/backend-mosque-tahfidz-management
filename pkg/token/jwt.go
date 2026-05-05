package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Payload struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type Maker interface {
	CreateToken(userID string, role string) (string, error)
	VerifyToken(token string) (*Payload, error)
}

type jwtMaker struct {
	secret []byte
}

func NewJWTMaker(secret string) Maker {
	return &jwtMaker{secret: []byte(secret)}
}

func (m *jwtMaker) CreateToken(userID string, role string) (string, error) {
	payload := Payload{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return token.SignedString(m.secret)
}

func (m *jwtMaker) VerifyToken(tokenStr string) (*Payload, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Payload{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return m.secret, nil
	})
	if err != nil {
		return nil, err
	}
	payload, ok := token.Claims.(*Payload)
	if !ok || !token.Valid {
		return nil, jwt.ErrInvalidKey
	}
	return payload, nil
}
