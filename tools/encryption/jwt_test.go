package encryption

import (
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

func TestEncryptGeneratesValidJWT(t *testing.T) {
	encrypt := NewJWTEncryption[jwt.MapClaims](
		JWTRequirements{
			SecretKey:     "secret",
			SigningMethod: JWTSigningMethodHS256,
		},
	)
	token, err := encrypt.Encrypt(jwt.MapClaims{"foo": "bar"})

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestDecryptParsesValidJWT(t *testing.T) {
	encrypt := NewJWTEncryption[jwt.MapClaims](
		JWTRequirements{
			SecretKey:     "secret",
			SigningMethod: jwt.SigningMethodHS256,
		},
	)
	token, _ := encrypt.Encrypt(jwt.MapClaims{"foo": "bar"})
	valid, decrypted, err := encrypt.Decrypt(token)

	assert.NoError(t, err)
	assert.True(t, valid)
	assert.Equal(t, "bar", decrypted["foo"])
}

func TestCustomStructDecryptParsesValidJWT(t *testing.T) {
	type CustomStruct struct {
		Foo string
	}

	encrypt := NewJWTEncryption[CustomStruct](
		JWTRequirements{
			SecretKey:     "secret",
			SigningMethod: jwt.SigningMethodHS256,
		},
	)
	token, _ := encrypt.Encrypt(jwt.MapClaims{"foo": "bar"})
	valid, decrypted, err := encrypt.Decrypt(token)

	assert.NoError(t, err)
	assert.True(t, valid)
	assert.Equal(t, "bar", decrypted.Foo)
}
