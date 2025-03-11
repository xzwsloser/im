package ctxdata

import "github.com/golang-jwt/jwt/v4"

/**
@Author: loser
@Description: to create jwt token
*/

const IdentifyKey = "userInfo"

// @Description: get the jwt key
func GetJwtToken(secretKey string, iat, seconds int64, uid string) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims[IdentifyKey] = uid

	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims

	return token.SignedString([]byte(secretKey))
}
