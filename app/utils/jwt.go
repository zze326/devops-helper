package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"reflect"
	"time"
)

func GenJwtToken(secretKey string, seconds int64, userID uint, username, userRealName string) (jwtToken string, accessExpire, refreshAfter int64, err error) {
	iat := time.Now().Unix()
	accessExpire = iat + seconds
	refreshAfter = iat + (seconds / 2)
	claims := make(jwt.MapClaims)
	claims["exp"] = accessExpire
	claims["iat"] = iat
	claims["user_id"] = userID
	claims["username"] = username
	claims["user_real_name"] = userRealName
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	jwtToken, err = token.SignedString([]byte(secretKey))
	return
}

func GetFromClaims[T any](key string, claims jwt.Claims) *T {
	v := reflect.ValueOf(claims)
	if v.Kind() == reflect.Map {
		for _, k := range v.MapKeys() {
			value := v.MapIndex(k)

			if fmt.Sprintf("%s", k.Interface()) == key {
				v := value.Interface().(T)
				return &v
			}
		}
	}
	return nil
}
