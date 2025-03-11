package ctxdata

import (
	"fmt"
	"testing"
	"time"
)

func TestGetJwtToken(t *testing.T) {
	secretKey := "loser"
	uid := "hello"
	iat := time.Now().Unix()
	var expireTime int64 = 86400
	token, err := GetJwtToken(secretKey, iat, expireTime, uid)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(token)
}
