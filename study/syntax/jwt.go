package main

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

var (
	SecretKey = "xuheng"
)

type UserDetails struct {
	// 用户标识
	UserId int64
	// 用户名 唯一
	Username string
	// 用户密码
	Password string
	// 用户具有的权限
	Authorities []string // 具备的权限
}

type jwtCustomClaims struct {
	jwt.StandardClaims
	UserDetails
}

type TokenEnhancer interface {
	// 组装 Token 信息
	Enhance(user *UserDetails) (string, error)
	// 从 Token 中还原信息
	Extract(toke string) (*UserDetails, error)
}

type JwtTokenEnhancer struct {
	secretKey []byte
}

func NewJwtTokenEnhancer(secretKey string) TokenEnhancer {
	return &JwtTokenEnhancer{
		secretKey: []byte(secretKey),
	}
}

func (enhancer *JwtTokenEnhancer) Enhance(user *UserDetails) (string, error) {

	claims := jwtCustomClaims{
		StandardClaims: jwt.StandardClaims{
			//ExpiresAt:int64(time.Now().Add(time.Hour * 72).Unix()),
			ExpiresAt: int64(time.Now().Add(time.Second * 10).Unix()),
			Issuer:    "System",
		},
		UserDetails: *user,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenValue, err := token.SignedString(enhancer.secretKey)

	if err != nil {
		return "", err
	}

	return tokenValue, nil
}

func (enhancer *JwtTokenEnhancer) Extract(tokenValue string) (*UserDetails, error) {

	token, err := jwt.ParseWithClaims(tokenValue, &jwtCustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return enhancer.secretKey, nil
	})

	if err == nil {
		claims := token.Claims.(*jwtCustomClaims)
		return &claims.UserDetails, nil
	}
	return nil, err
}

func main() {
	user := UserDetails{111111, "username", "passwd", nil}
	jt := NewJwtTokenEnhancer("xuheng")
	token, err := jt.Enhance(&user)
	if err == nil {
		println(token)
	}

	time.Sleep(2 * time.Second)

	ut, err := jt.Extract(token)

	if err == nil {
		println(ut.UserId, ut.Username, ut.Password, ut.Authorities)
	} else {
		println(err.Error())
	}
}
