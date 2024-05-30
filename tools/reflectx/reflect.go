package reflectx

import (
	"reflect"
)

// IsStruct 判斷傳入變數是否為結構體
func IsStruct(v interface{}) bool {
	// 使用反射獲取傳入變數的類型
	t := reflect.TypeOf(v)
	return t.Kind() == reflect.Struct
}

func IsStructPtr(v interface{}) bool {
	t := reflect.TypeOf(v)
	return t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct
}
