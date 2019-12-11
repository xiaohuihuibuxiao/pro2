package Baseinfo

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"log"
	"time"
)

type UserClaim struct {
	Uname string `json:"username"`
	jwt.StandardClaims
}

const (
	key               = "dashaiduhakdkadhq132u489274927498(&(*&(*^(E" //未使用
	symmetricalsecret = "spacemanagementsystembasedonsensors"         //使用中
)

//生成token--对称密钥
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
		//fmt.Println(claims["sub"])
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
	//ok, errt, user := Authtoken(token)//对称
	ok, errt, user := TokenCheck_asymmetricalkey(token) //非对称验证
	if !ok || errt != nil {
		return errt, user
	}
	return nil, user
}

//------对称加密生成token---------
func TokenGen_symmetricalkey() string {
	sec := []byte(symmetricalsecret)
	//hs256
	token_obj := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaim{Uname: "shenyi"})
	token, _ := token_obj.SignedString(sec)
	return token
}

func TokenCheck_symmetricalkey(token string, secret []byte) {
	uc := UserClaim{}
	getToken, _ := jwt.ParseWithClaims(token, &uc, func(token *jwt.Token) (i interface{}, e error) {
		return secret, nil
	})
	if getToken.Valid {
		fmt.Println(getToken.Claims.(*UserClaim).Uname)
		fmt.Println(getToken.Claims.(*UserClaim).ExpiresAt)
	}
}

//=-----非对称加密token------
func TokenGen_asymmetricalkey(userid string) (string, error) {
	priKeyBytes, err := ioutil.ReadFile("./pem/private.pem")
	if err != nil {
		log.Fatal("私钥文件读取失败")
		return "", err
	}
	priKey, err := jwt.ParseRSAPrivateKeyFromPEM(priKeyBytes)
	if err != nil {
		log.Fatal("私钥文件不正确")
		return "", err
	}

	pubKeyBytes, err := ioutil.ReadFile("./pem/public.pem")
	if err != nil {
		log.Fatal("公钥文件读取失败")
		return "", err
	}
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pubKeyBytes)
	if err != nil {
		log.Fatal("公钥文件不正确")
		return "", err
	}
	user := UserClaim{Uname: userid}
	user.ExpiresAt = time.Now().Add(time.Duration(Expiredtime) * time.Second).Unix() //TODO 在config中配置
	token_obj := jwt.NewWithClaims(jwt.SigningMethodRS256, user)
	token, _ := token_obj.SignedString(priKey)

	//--校验token时使用pubkey
	uc := UserClaim{}
	getToken, _ := jwt.ParseWithClaims(token, &uc, func(token *jwt.Token) (i interface{}, e error) {
		return pubKey, nil
	})
	if getToken.Valid {
		fmt.Println(getToken.Claims.(*UserClaim).Uname)
	}
	return token, nil
}

func TokenCheck_asymmetricalkey(token string) (bool, error, string) {
	pubKeyBytes, err := ioutil.ReadFile("./pem/public.pem")
	if err != nil {
		log.Fatal("公钥文件读取失败")
		return false, err, ""
	}
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pubKeyBytes)
	if err != nil {
		log.Fatal("公钥文件不正确")
		return false, err, ""
	}
	uc := UserClaim{}
	getToken, err := jwt.ParseWithClaims(token, &uc, func(token *jwt.Token) (i interface{}, e error) {
		return pubKey, nil
	})
	if getToken != nil && getToken.Valid {
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return false, errors.New("invalid token"), ""
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {

			return false, errors.New("token is expired"), ""
		} else {
			return false, errors.New("Couldn't handle this token:" + err.Error()), ""
		}
	} else {
		return false, errors.New("unresolved token err:" + err.Error()), ""
	}
	return true, nil, getToken.Claims.(*UserClaim).Uname
}
