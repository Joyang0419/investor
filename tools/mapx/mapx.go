package mapx

import (
	"reflect"

	"golang.org/x/exp/constraints"

	"tools/slicex"
)

// CombineMaps combines multiple maps into one map.
func CombineMaps[TypeKey constraints.Ordered](maps ...map[TypeKey]interface{}) map[TypeKey]interface{} {
	if slicex.IsEmpty(maps) {
		return nil
	}

	result := make(map[TypeKey]interface{})
	for _, m := range maps {
		for k, v := range m {
			value, keyExist := result[k]
			if !keyExist {
				result[k] = v
				continue
			}

			// 以下代碼, 是目前想到是容器, ex: map or slice, 不是Replace value, 而是合併
			if checkIsMap(value) {
				existingVal := reflect.ValueOf(value)
				newVal := reflect.ValueOf(v)
				for _, mapKey := range newVal.MapKeys() {
					existingVal.SetMapIndex(mapKey, newVal.MapIndex(mapKey))
				}
				result[k] = existingVal.Interface()
				continue
			}

			if checkIsSlice(value) {
				existingVal := reflect.ValueOf(value)
				newVal := reflect.ValueOf(v)
				combinedSlice := reflect.AppendSlice(existingVal, newVal)
				result[k] = combinedSlice.Interface()
				continue
			}

			result[k] = v
		}
	}

	return result
}

func checkIsMap(data any) bool {
	return reflect.TypeOf(data).Kind() == reflect.Map
}

func checkIsSlice(data any) bool {
	return reflect.TypeOf(data).Kind() == reflect.Slice
}
