package encrypt

import (
	"fmt"
	"testing"
)

func TestGenPasswordHash(t *testing.T) {
	pass := "898989"
	res, err := GenPasswordHash([]byte(pass))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(res))
}

func TestValidatePasswordHash(t *testing.T) {
	pass := "$2a$10$op47lZK2jAEW5BZk2HwVCuV/t8wMGoauib4C1Tw0zhByGrax5JVPm"
	realPass := "898989"
	if !ValidatePasswordHash(realPass, pass) {
		t.Fatal("failed!")
	}
	fmt.Println("测试成功...")
}
