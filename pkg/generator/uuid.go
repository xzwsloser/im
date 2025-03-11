package generator

import "github.com/google/uuid"

/*
	@Author: loser
	@Description: generator a new uuid
*/

// @Desciption: generate a new uuid
func GeneratorUUid() string {
	return uuid.New().String()
}
