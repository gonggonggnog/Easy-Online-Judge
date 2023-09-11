package helper

import (
	"crypto/md5"
	"fmt"
	"github.com/dgrijalva/jwt-go"
)

type UserClaims struct {
	Identity string `json:"identity"`
	Name     string `json:"name"`
	IsAdmin  int    `json:"is_admin"`
	jwt.StandardClaims
}

var key = "my-oj-project"

func GetMd5(string2 string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(string2)))
}

func GenerateToken(identity, name string, isAdmin int) (string, error) {
	UserClaim := &UserClaims{
		Identity:       identity,
		Name:           name,
		IsAdmin:        isAdmin,
		StandardClaims: jwt.StandardClaims{},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaim)
	return token.SignedString(key)
}
