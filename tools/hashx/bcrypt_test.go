package hashx

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestBcryptHash_Hash(t *testing.T) {
	b := NewBcrypt(bcrypt.DefaultCost)
	h, err := b.Hash([]byte("123456"))

	assert.NoError(t, err)
	assert.NotEmpty(t, h)
}

func TestBcryptHash_CompareHash(t *testing.T) {
	b := NewBcrypt(bcrypt.DefaultCost)

	h, err := b.Hash([]byte("123456"))

	err = b.CompareHash([]byte("123456"), h)

	assert.NoError(t, err)
}
