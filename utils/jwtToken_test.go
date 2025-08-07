package utils

import (
	"os"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {

	id := 1
	name := "razi"

	tokenstring, err := GenerateToken(id, name)

	assert.NoError(t, err)          // if error, then is logged
	assert.NotEmpty(t, tokenstring) // if empty string , then logged

	token, err := jwt.Parse(tokenstring, func(Token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	assert.NoError(t, err)
	assert.True(t, token.Valid) // returns booleans value

}
