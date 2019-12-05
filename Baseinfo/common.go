package Baseinfo

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

const (
	key = "dashaiduhakdkadhq132u489274927498(&(*&(*^(E" //TODO key目前是随意设置的 如果有必要在这里修改
)

//生成token
func Gentoken(user string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["sub"] = user //用于在controller中确定用户
	//claims["exp"] = time.Now().Add(time.Hour.Truncate(72)) //设置过期时间为72小时后
	claims["exp"] = time.Now().Add(30 * time.Minute) //设置30分钟以后token过期
	claims["iat"] = time.Now().Unix()                //用作和exp对比的时间
	token.Claims = claims

	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

//验证token有效性 是否超时
func Authtoken(tokenString string) (bool, error, string) {
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//验证是否是给定的加密算法
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("fail to gen authtoken ")
		}
		return []byte(key), nil
	})
	if !token.Valid {
		return false, errors.New("token is invalid"), ""
	} else {
		claims := token.Claims.(jwt.MapClaims)
		fmt.Println(claims["sub"])
		local, _ := time.LoadLocation("Local")
		t, _ := time.ParseInLocation("2006-01-02 15:04:05", claims["exp"].(string)[:10]+" "+claims["exp"].(string)[11:19], local)
		if t.Unix() < time.Now().Unix() {
			return false, errors.New("token is expired"), claims["sub"].(string)
		} else {
			return true, nil, claims["sub"].(string)
		}
	}

}

//登陆接口中校验token的方法
func Logintokenauth(token string) (error, string) {
	if token == "" {
		return errors.New("lack of token"), ""
	}
	ok, errt, user := Authtoken(token)
	if !ok || errt != nil {
		return errt, user
	}
	return nil, user
}
