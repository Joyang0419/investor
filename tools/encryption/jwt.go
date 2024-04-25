package encryption

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/mitchellh/mapstructure"
)

// JWTRequirements 定義 JWT 特定的設定需求。
type JWTRequirements struct {
	SecretKey     []byte
	SigningMethod jwt.SigningMethod
}

type JWTEncryption[decryptedType any] struct {
	requirements JWTRequirements
}

func (encrypt *JWTEncryption[decryptedType]) SetEncryptionRequirements(requirements JWTRequirements) {
	if requirements.SecretKey == nil {
		panic("[JWTEncryption][SetEncryptionRequirements]SecretKey is required")
	}
	if requirements.SigningMethod == nil {
		panic("[JWTEncryption][SetEncryptionRequirements]SigningMethod is required")
	}
	encrypt.requirements = requirements
}

// NewJWTEncryption 創建一個新的 JWTEncryption 實例。
func NewJWTEncryption[decryptedType any](requirements JWTRequirements) *JWTEncryption[decryptedType] {
	return &JWTEncryption[decryptedType]{requirements: requirements}
}

// Encrypt 生成一個 JWT。
func (encrypt *JWTEncryption[decryptedType]) Encrypt(beforeEncrypt jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(encrypt.requirements.SigningMethod, beforeEncrypt)
	// 使用從 SetEncryptionRequirements 設定的密鑰簽名。
	tokenString, err := token.SignedString(encrypt.requirements.SecretKey)
	if err != nil {
		return "", fmt.Errorf("[JWTEncryption][Encrypt]token.SignedString err: %w", err)
	}
	return tokenString, nil
}

// Decrypt 驗證並解析 JWT。
func (encrypt *JWTEncryption[decryptedType]) Decrypt(beforeDecrypt string) (decrypted decryptedType, err error) {
	token, err := jwt.ParseWithClaims(beforeDecrypt, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return encrypt.requirements.SecretKey, nil
	})
	if err != nil {
		return decrypted, fmt.Errorf("[JWTEncryption][Decrypt]jwt.ParseWithClaims err: %w", err)
	}

	if !token.Valid {
		return decrypted, fmt.Errorf("[JWTEncryption][Decrypt]Invalid token")
	}

	if err = mapstructure.Decode(token.Claims, &decrypted); err != nil {
		return decrypted, fmt.Errorf("[JWTEncryption][Decrypt]mapstructure.Decode err: %w", err)
	}

	return decrypted, nil
}
