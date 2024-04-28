package hashx

// IHash  主要是雜湊 (Hash)的interface, such as bcrypt, md5, sha256 ...
type IHash[beforeHash, afterHash any] interface {
	// TODO delete
	// SetHashRequirements 設定工具必須要的參數
	//SetHashRequirements(requirements Requirements)
	// Hash 雜湊資料
	Hash(beforeHash beforeHash) (afterHash afterHash, err error)
	// CompareHash 比較雜湊資料
	CompareHash(beforeHash beforeHash, afterHash afterHash) (err error)
}
