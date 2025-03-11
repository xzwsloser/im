package encrypt

import (
	"crypto/md5"
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
)

/**
@Author: loser
@Descirption: to gen password using hash
*/

// @Description: create the md5 secret key
func Md5(str []byte) string {
	h := md5.New()
	h.Write(str)
	return hex.EncodeToString(h.Sum(nil))
}

// @Description: user the hash to gen password
func GenPasswordHash(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
}

// @Description: to verify the password
func ValidatePasswordHash(password string, hash string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return false
	}
	return true
}
