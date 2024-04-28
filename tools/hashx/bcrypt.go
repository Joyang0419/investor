package hashx

import (
	"golang.org/x/crypto/bcrypt"
)

// implement IHash with bcryptHash
type bcryptHash struct {
	Cost int
}

// Hash 雜湊資料
func (b *bcryptHash) Hash(beforeHash []byte) (afterHash []byte, err error) {
	return bcrypt.GenerateFromPassword(beforeHash, b.Cost)
}

// CompareHash 比較雜湊資料
func (b *bcryptHash) CompareHash(beforeHash []byte, afterHash []byte) (err error) {
	return bcrypt.CompareHashAndPassword(afterHash, beforeHash)
}

func NewBcrypt(cost int) IHash[[]byte, []byte] {
	return &bcryptHash{Cost: cost}
}
