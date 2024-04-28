package mapx_test

import (
	"testing"

	"tools/mapx"
)

func TestCombineMapsSingleMap(t *testing.T) {
	map1 := map[int]interface{}{1: "one", 2: "two"}
	result := mapx.CombineMaps[int](map1)

	if len(result) != 2 || result[1] != "one" || result[2] != "two" {
		t.Errorf("CombineMaps failed with single map")
	}
}

func TestCombineMapsMultipleMaps(t *testing.T) {
	map1 := map[int]interface{}{1: "one", 2: "two"}
	map2 := map[int]interface{}{3: "three", 4: "four"}
	result := mapx.CombineMaps[int](map1, map2)

	if len(result) != 4 || result[1] != "one" || result[2] != "two" || result[3] != "three" || result[4] != "four" {
		t.Errorf("CombineMaps failed with multiple maps")
	}
}

func TestCombineMapsOverlappingKeys(t *testing.T) {
	map1 := map[int]interface{}{1: "one", 2: "two"}
	map2 := map[int]interface{}{2: "two", 3: "three"}
	result := mapx.CombineMaps[int](map1, map2)

	if len(result) != 3 || result[1] != "one" || result[2] != "two" || result[3] != "three" {
		t.Errorf("CombineMaps failed with overlapping keys")
	}
}

func TestCombineMapsEmptyMaps(t *testing.T) {
	map1 := map[int]interface{}{}
	map2 := map[int]interface{}{}
	result := mapx.CombineMaps[int](map1, map2)

	if len(result) != 0 {
		t.Errorf("CombineMaps failed with empty maps")
	}
}
