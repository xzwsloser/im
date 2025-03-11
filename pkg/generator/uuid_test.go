package generator

import (
	"fmt"
	"testing"
)

func TestGeneratorUUid(t *testing.T) {
	uuid := GeneratorUUid()
	fmt.Println(uuid)
	fmt.Println(len(uuid))
}
