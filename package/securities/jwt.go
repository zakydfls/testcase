package securities

import (
	"fmt"
	"testcase/internal/modules/user/entities"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTPayload struct {
	UserID   uuid.UUID              `json:"user_id"`
	Username string                 `json:"username"`
	Role     entities.RoleEnum      `json:"role"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type JWTClaims struct {
	*JWTPayload
	jwt.RegisteredClaims
}

type JWTManager struct {
	secretKey        string
	refreshSecretKey string
	tokenExpiry      time.Duration
	refreshExpiry    time.Duration
	apiKey           *string
}

func NewJWTManager(secretKey string, refreshSecretKey string, tokenExpiry, refreshExpiry time.Duration, apiKey string) *JWTManager {
	return &JWTManager{
		secretKey:        secretKey,
		refreshSecretKey: refreshSecretKey,
		tokenExpiry:      tokenExpiry,
		refreshExpiry:    refreshExpiry,
		apiKey:           &apiKey,
	}
}

func (jm *JWTManager) GenerateToken(jwtPayload *JWTPayload) (string, error) {
	claims := &JWTClaims{
		JWTPayload: jwtPayload,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(jm.tokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "pln-api",
			Subject:   fmt.Sprintf("user:%d", jwtPayload.UserID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jm.secretKey))
}

func (jm *JWTManager) GenerateRefreshToken(jwtPayload *JWTPayload) (string, error) {
	claims := &JWTClaims{
		JWTPayload: jwtPayload,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(jm.refreshExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "pln-api",
			Subject:   fmt.Sprintf("refresh:%d", jwtPayload.UserID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jm.refreshSecretKey))
}

func (jm *JWTManager) GenerateTokenPair(jwtPayload *JWTPayload) (string, string, error) {
	accessToken, err := jm.GenerateToken(jwtPayload)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := jm.GenerateRefreshToken(jwtPayload)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (jm *JWTManager) ValidateToken(tokenString, tokenType string) (*JWTClaims, error) {
	var secret string
	if tokenType == "refresh" {
		secret = jm.refreshSecretKey
	} else {
		secret = jm.secretKey
	}
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func (jm *JWTManager) ExtractClaims(tokenString string) (*JWTClaims, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, &JWTClaims{})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token claims")
}

func (jm *JWTManager) ValidateAndExtract(tokenString, tokenType string) (*JWTClaims, error) {
	var secret string
	if tokenType == "refresh" {
		secret = jm.refreshSecretKey
	} else {
		secret = jm.secretKey
	}
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("token expired")
	}
	return claims, nil
}

func (jm *JWTManager) IsTokenExpired(tokenString string) bool {
	claims, err := jm.ExtractClaims(tokenString)
	if err != nil {
		return true
	}

	return claims.ExpiresAt.Before(time.Now())
}

func (jm *JWTManager) GetTokenExpiry() time.Duration {
	return jm.tokenExpiry
}

func (jm *JWTManager) GetRefreshExpiry() time.Duration {
	return jm.refreshExpiry
}
