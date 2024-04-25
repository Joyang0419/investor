package encryption

// IEncryption  主要是加密 (Encryption)的interface, such as JWT, AES, RSA ...
type IEncryption[Requirements, beforeEncrypt, encrypted, beforeDecrypt, decrypted any] interface {
	// SetEncryptionRequirements 設定工具必須要的參數
	SetEncryptionRequirements(requirements Requirements)
	// Encrypt 加密資料
	Encrypt(beforeEncrypt beforeEncrypt) (encrypted encrypted, err error)
	// Decrypt 解密資料
	Decrypt(beforeDecrypt beforeDecrypt) (decrypted decrypted, err error)
}
