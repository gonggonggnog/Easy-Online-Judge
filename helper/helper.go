package helper

import (
	"blog/define"
	"crypto/md5"
	"crypto/tls"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/jordan-wright/email"
	"net/smtp"
)

type UserClaims struct {
	Identity string `json:"identity"`
	Name     string `json:"name"`
	IsAdmin  int    `json:"is_admin"`
	jwt.StandardClaims
}

var key = []byte("my-gin-oj-project")

//获取md5信息

func GetMd5(string2 string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(string2)))
}

//生成token

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

//发送邮箱验证码

func SendEmail(toUserEmail, code string) error {
	e := email.NewEmail()
	e.From = "Get <13320064423@163.com>"
	e.To = []string{toUserEmail}
	e.Subject = "验证码已发送，请查收"
	e.HTML = []byte("感谢您使用此平台，您的验证码为：<b>" + code + "</b><br>请在十分钟内验证，验证码打死也不要告诉别人哦")
	return e.SendWithTLS("smtp.163.com:465",
		smtp.PlainAuth("", "13320064423@163.com", define.MailPassword, "smtp.163.com"),
		&tls.Config{InsecureSkipVerify: true, ServerName: "smtp.163.com"})
}
