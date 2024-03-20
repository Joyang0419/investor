package service

import (
	"apiserver/handler"
)

type ExampleService struct {
}

// 故意的，compile時, 如果沒有implement IExampleService, 會直接不能compile
var _ handler.IExampleService = NewExampleService()

func NewExampleService() *ExampleService {
	return &ExampleService{}
}

func (e *ExampleService) Helloworld() string {
	return "hello world"
}
