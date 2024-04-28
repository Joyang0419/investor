package hashx

// HashRequirements 雜湊需求
type HashRequirements struct {
	HashFunc    func(data []byte, cost int) ([]byte, error)
	CompareFunc func(hashedPassword, password []byte) error
	Cost        int
}

// implement IHash with bcryptHash
type bcryptHash struct {
	requirements HashRequirements
}

// SetHashRequirements 設定工具必須要的參數
func (b *bcryptHash) SetHashRequirements(requirements HashRequirements) {
	if requirements.HashFunc == nil {
		panic("[bcryptHash][SetHashRequirements]HashFunc is required")
	}
	if requirements.CompareFunc == nil {
		panic("[bcryptHash][SetHashRequirements]CompareFunc is required")
	}
	if requirements.Cost == 0 {
		panic("[bcryptHash][SetHashRequirements]Cost is required")
	}
	b.requirements = requirements
}

// Hash 雜湊資料
func (b *bcryptHash) Hash(beforeHash []byte) (afterHash []byte, err error) {
	afterHash, err = b.requirements.HashFunc(beforeHash, b.requirements.Cost)
	if err != nil {
		return nil, err
	}

	return afterHash, nil
}

// CompareHash 比較雜湊資料
func (b *bcryptHash) CompareHash(beforeHash []byte, afterHash []byte) (err error) {
	err = b.requirements.CompareFunc(afterHash, beforeHash)
	if err != nil {
		return err
	}

	return nil
}

// NewBcrypt create a new bcryptHash instance
func NewBcrypt() IHash[HashRequirements, []byte, []byte] {
	return new(bcryptHash)
}
