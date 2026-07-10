package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JwtConfig struct {
	Secret string
	TTL    time.Duration
}

type Claims struct {
	Sub   string
	Email string
	Role  string
}

func New(secret string, ttl time.Duration) *JwtConfig {
	return &JwtConfig{
		Secret: secret,
		TTL:    ttl,
	}
}

func (j *JwtConfig) Generate(userID uuid.UUID, email, role string) (string, error) {
	now := time.Now()

	jwtClaims := jwt.MapClaims{
		"sub":   userID.String(),
		"email": email,
		"role":  role,
		"iat":   now.Unix(),
		"exp":   now.Add(j.TTL).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)

	return token.SignedString([]byte(j.Secret))
}

func (j *JwtConfig) Validate(tokenString string) (*Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.Secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	sub, _ := claims["sub"].(string)
	email, _ := claims["email"].(string)
	role, _ := claims["role"].(string)

	return &Claims{
		Sub:   sub,
		Email: email,
		Role:  role,
	}, nil
}
