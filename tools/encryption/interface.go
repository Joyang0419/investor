package encryption

// IEncryption  主要是加密 (Encryption)的interface, such as JWT, AES, RSA ...
type IEncryption[beforeEncrypt, encrypted, beforeDecrypt, decrypted any] interface {
	// Encrypt 加密資料
	Encrypt(beforeEncrypt beforeEncrypt) (encrypted encrypted, err error)
	// Decrypt 解密資料
	Decrypt(beforeDecrypt beforeDecrypt) (decrypted decrypted, err error)
}
