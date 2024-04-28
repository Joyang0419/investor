package hashx

import (
	"golang.org/x/crypto/bcrypt"
)

// HashRequirements 雜湊需求
type HashRequirements struct {
	// TODO 你這邊這樣寫，變成這個實作，誰實作了，HashFunc, CompareFunc, Cost
	//HashFunc    func(data []byte, cost int) ([]byte, error)
	//CompareFunc func(hashedPassword, password []byte) error
	// TODO COST 註解 这个工作因子决定了哈希函数的复杂性，值越大，哈希函数的计算就越复杂，生成哈希值所需的时间就越长。
	Cost int
}

// implement IHash with bcryptHash
// TODO 因果關係, 不對, 所以我隨便實作 HashFunc, CompareFunc
type bcryptHash struct {
	// Cost决定哈希函数的复杂性，值越大，哈希函数的计算就越复杂，生成哈希值所需的时间就越长。
	Cost int
	//requirements HashRequirements
}

// TODO delete
// SetHashRequirements 設定工具必須要的參數
//func (b *bcryptHash) SetHashRequirements(requirements HashRequirements) {
//	//if requirements.HashFunc == nil {
//	//	panic("[bcryptHash][SetHashRequirements]HashFunc is required")
//	//}
//	//if requirements.CompareFunc == nil {
//	//	panic("[bcryptHash][SetHashRequirements]CompareFunc is required")
//	//}
//	if requirements.Cost == 0 {
//		panic("[bcryptHash][SetHashRequirements]Cost is required")
//	}
//	b.requirements = requirements
//}

// Hash 雜湊資料
func (b *bcryptHash) Hash(beforeHash []byte) (afterHash []byte, err error) {
	return bcrypt.GenerateFromPassword(beforeHash, b.Cost)

	// TODO delete
	//afterHash, err = b.requirements.HashFunc(beforeHash, b.requirements.Cost)
	//if err != nil {
	//	return nil, err
	//}
	//
	//return afterHash, nil
}

// CompareHash 比較雜湊資料
func (b *bcryptHash) CompareHash(beforeHash []byte, afterHash []byte) (err error) {
	return bcrypt.CompareHashAndPassword(afterHash, beforeHash)
	// TODO delete
	//err = b.requirements.CompareFunc(afterHash, beforeHash)
	//if err != nil {
	//	return err
	//}
	//
	//return nil
}

func NewBcrypt(cost int) IHash[[]byte, []byte] {
	return &bcryptHash{Cost: cost}
}
