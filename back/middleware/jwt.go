package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"netdisk/models"
	"netdisk/utils"
	"time"
)

type MyClaims struct {
	models.User
	jwt.StandardClaims
}

var (
	mySigningKey = []byte("woshihuaiche")
)

func AuthMiddle() gin.HandlerFunc {
	return func(context *gin.Context) {
		token := context.GetHeader("X-Token")
		if user, err := ParseToken(token); err != nil {
			utils.ClientFailWithMessage(context, err.Error())
			context.Abort()
		} else {
			context.Set("userID", &user.ID)
			context.Set("IsVIP", user.IsVIP)
			context.Next()
		}
	}
}

func GetToken(user models.User) (string, error) {
	c := MyClaims{
		User: user,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * 12).Unix(),
			Issuer:    "huaichekk",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	s, err := token.SignedString(mySigningKey)
	if err != nil {
		return "", err
	}
	return s, err
}

func ParseToken(token string) (*models.User, error) {
	claims, err := jwt.ParseWithClaims(token, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	})
	if err != nil {
		return nil, err
	}
	return &claims.Claims.(*MyClaims).User, nil
}
