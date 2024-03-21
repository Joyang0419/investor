package serialization

import (
	"fmt"

	"github.com/bytedance/sonic"
)

// why choose sonic: https://cloud.tencent.com/developer/article/2296705

func JsonMarshal(input any) (jsonBytes []byte, err error) {
	return sonic.Marshal(input)
}

func JsonUnmarshal[T any](jsonStr string) (result T, err error) {
	if err = sonic.Unmarshal([]byte(jsonStr), &result); err != nil {
		return result, fmt.Errorf("[JsonUnmarshal]sonic.Unmarshal err: %w", err)
	}

	return result, nil
}
