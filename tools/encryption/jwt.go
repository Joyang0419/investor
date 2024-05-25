package encryption

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/mitchellh/mapstructure"
)

var (
	JWTSigningMethodHS256 = jwt.SigningMethodHS256
	JWTSigningMethodHS384 = jwt.SigningMethodHS384
	JWTSigningMethodHS512 = jwt.SigningMethodHS512
)

type JWTMapClaims = jwt.MapClaims

// JWTRequirements 定義 JWT 特定的設定需求。
type JWTRequirements struct {
	SecretKey      string
	SigningMethod  jwt.SigningMethod
	ExpireDuration time.Duration // Token 過期時間
}

type JWTEncryption[decryptedType any] struct {
	requirements JWTRequirements
}

// NewJWTEncryption 創建一個新的 JWTEncryption 實例。
func NewJWTEncryption[decryptedType any](requirements JWTRequirements) JWTEncryption[decryptedType] {
	return JWTEncryption[decryptedType]{requirements: requirements}
}

// Encrypt 生成一個 JWT。
func (encrypt *JWTEncryption[decryptedType]) Encrypt(beforeEncrypt jwt.MapClaims) (string, error) {
	beforeEncrypt["exp"] = time.Now().Add(encrypt.requirements.ExpireDuration).Unix()
	token := jwt.NewWithClaims(encrypt.requirements.SigningMethod, beforeEncrypt)
	tokenString, err := token.SignedString([]byte(encrypt.requirements.SecretKey))
	if err != nil {
		return "", fmt.Errorf("[JWTEncryption][Encrypt]token.SignedString err: %w", err)
	}
	return tokenString, nil
}

// Decrypt 驗證並解析 JWT。
func (encrypt *JWTEncryption[decryptedType]) Decrypt(beforeDecrypt string) (valid bool, decrypted decryptedType, err error) {
	token, err := jwt.ParseWithClaims(beforeDecrypt, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(encrypt.requirements.SecretKey), nil
	})
	if err != nil {
		return false, decrypted, fmt.Errorf("[JWTEncryption][Decrypt]jwt.ParseWithClaims err: %w", err)
	}

	if !token.Valid {
		return false, decrypted, nil
	}

	if err = mapstructure.Decode(token.Claims, &decrypted); err != nil {
		return false, decrypted, fmt.Errorf("[JWTEncryption][Decrypt]mapstructure.Decode err: %w", err)
	}

	return true, decrypted, nil
}
