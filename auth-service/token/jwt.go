package token

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type jwtCustomClaims struct {
	jwt.StandardClaims
	OAuth2Details
}

type TokenEnhancer interface {
	// 组装 Token 信息
	Enhance(detail *OAuth2Details) (string, error)
	// 从 Token 中还原信息
	Extract(toke string) (*OAuth2Details, error)
}

type JwtTokenEnhancer struct {
	secretKey []byte
}

func NewJwtTokenEnhancer(secretKey string) TokenEnhancer {
	return &JwtTokenEnhancer{
		secretKey: []byte(secretKey),
	}
}

func (enhancer *JwtTokenEnhancer) Enhance(detail *OAuth2Details) (string, error) {

	claims := jwtCustomClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: int64(time.Now().Add(time.Second * time.Duration(detail.Client.AccessTokenValiditySeconds)).Unix()),
			Issuer:    "System",
		},
		OAuth2Details: *detail,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenValue, err := token.SignedString(enhancer.secretKey)

	if err != nil {
		return "", err
	}

	return tokenValue, nil
}

func (enhancer *JwtTokenEnhancer) Extract(tokenValue string) (*OAuth2Details, error) {

	token, err := jwt.ParseWithClaims(tokenValue, &jwtCustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return enhancer.secretKey, nil
	})

	if err == nil {
		claims := token.Claims.(*jwtCustomClaims)
		return &claims.OAuth2Details, nil
	}
	return nil, err
}
