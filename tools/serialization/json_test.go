package serialization

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestJsonMarshal(t *testing.T) {
	person := Person{
		Name: "John",
		Age:  30,
	}

	// 序列化
	jsonBytes, err := JsonMarshal(person)
	if err != nil {
		t.Fatalf("JsonMarshal failed: %v", err)
	}

	expectedJSONStr := `{"name":"John","age":30}`
	assert.JSONEq(t, expectedJSONStr, string(jsonBytes))
}

func TestJsonUnmarshal(t *testing.T) {
	jsonStr := `{"name":"John","age":30}`

	// 反序列化

	person, err := JsonUnmarshal[Person](jsonStr)

	expectedPerson := Person{
		Name: "John",
		Age:  30,
	}

	assert.Equal(t, expectedPerson, person)
	assert.NoError(t, err)
}
