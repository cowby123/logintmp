package main

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

//Retdata 载荷，可以加一些自己需要的信息
type Retdata struct {
	Name  string `json:"Name"`
	Token string `json:"Token"`
	Data  string `json:"Data"`
}

//CustomClaims 载荷，可以加一些自己需要的信息
type CustomClaims struct {
	ID   string `json:"userId"`
	Name string `json:"name"`
	Data string `json:"Data"`
	jwt.StandardClaims
}

// JWTAuth 中间件，检查token
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var Data Retdata
		err := c.BindJSON(&Data)
		token := Data.Token
		if token == "" {
			c.JSON(http.StatusOK, gin.H{
				"Static":  -1,
				"Message": "沒有token，沒有權限",
				"Data":    "",
			})
			c.Abort()
			return
		}

		log.Print("get token: ", token)

		j := NewJWT()
		// parseToken 解析token包含的信息
		claims, err := j.ParseToken(token)
		if err != nil {
			if err == errTokenExpired {
				c.JSON(http.StatusOK, gin.H{
					"status": -1,
					"msg":    "授權過期",
				})
				c.Abort()
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"status": -1,
				"msg":    err.Error(),
			})
			c.Abort()
			return
		}
		// 继续交由下一个路由处理,并将解析出的信息传递下去
		claims.Data = Data.Data
		c.Set("claims", claims)
	}
}

// JWT 签名结构
type JWT struct {
	SigningKey []byte
}

// 一些常量
var (
	errTokenExpired     error  = errors.New("token is expired")
	errTokenNotValidYet error  = errors.New("token not active yet")
	errTokenMalformed   error  = errors.New("that's not even a token")
	errTokenInvalid     error  = errors.New("couldn't handle this token: ")
	SignKey             string = "legbone"
)

//NewJWT 新建一个jwt实例
func NewJWT() *JWT {
	return &JWT{
		[]byte(GetSignKey()),
	}
}

//GetSignKey 获取signKey
func GetSignKey() string {
	return SignKey
}

//SetSignKey 这是SignKey
func SetSignKey(key string) string {
	SignKey = key
	return SignKey
}

// CreateToken 生成一个token
func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

//ParseToken 解析Tokne
func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, errTokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, errTokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, errTokenNotValidYet
			} else {
				return nil, errTokenInvalid
			}
		}
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errTokenInvalid
}

//RefreshToken 更新token
func (j *JWT) RefreshToken(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
		return j.CreateToken(*claims)
	}
	return "", errTokenInvalid
}
