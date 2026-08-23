package security

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	secret            []byte
	accessExpiration  time.Duration
	refreshExpiration time.Duration
}

type Claims struct {
	UserID    int64  `json:"user_id"`
	SessionID int64  `json:"session_id,omitempty"`
	Role      string `json:"role,omitempty"`
	TokenType string `json:"token_type"`

	jwt.RegisteredClaims
}

func NewJWTManager(
	secret string,
	accessExpiration time.Duration,
	refreshExpiration time.Duration,
) *JWTManager {
	return &JWTManager{
		secret:            []byte(secret),
		accessExpiration:  accessExpiration,
		refreshExpiration: refreshExpiration,
	}
}

func (m *JWTManager) GenerateAccessToken(
	userID int64,
	role string,
) (string, error) {

	now := time.Now()

	claims := Claims{
		UserID:    userID,
		Role:      role,
		TokenType: "access",

		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(
				now.Add(m.accessExpiration),
			),
			IssuedAt: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString(m.secret)
}

func (m *JWTManager) GenerateRefreshToken(
	userID int64,
	sessionID int64,
) (string, error) {

	now := time.Now()

	claims := Claims{
		UserID:    userID,
		SessionID: sessionID,
		TokenType: "refresh",

		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(
				now.Add(m.refreshExpiration),
			),
			IssuedAt: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString(m.secret)
}

func (m *JWTManager) parseToken(
	tokenString string,
) (*Claims, error) {

	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {

			if token.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf(
					"unexpected signing method: %v",
					token.Method.Alg(),
				)
			}

			return m.secret, nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf(
			"invalid token: %w",
			err,
		)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, fmt.Errorf(
			"invalid token claims",
		)
	}

	if !token.Valid {
		return nil, fmt.Errorf(
			"invalid token",
		)
	}

	return claims, nil
}

func (m *JWTManager) ValidateAccessToken(
	tokenString string,
) (*Claims, error) {

	claims, err := m.parseToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.TokenType != "access" {
		return nil, fmt.Errorf(
			"token is not an access token",
		)
	}

	return claims, nil
}

func (m *JWTManager) ValidateRefreshToken(
	tokenString string,
) (*Claims, error) {

	claims, err := m.parseToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.TokenType != "refresh" {
		return nil, fmt.Errorf(
			"token is not a refresh token",
		)
	}

	if claims.SessionID == 0 {
		return nil, fmt.Errorf(
			"refresh token missing session ID",
		)
	}

	return claims, nil
}

func (m *JWTManager) AccessExpiration() time.Duration {
	return m.accessExpiration
}

func (m *JWTManager) RefreshExpiration() time.Duration {
	return m.refreshExpiration
}
