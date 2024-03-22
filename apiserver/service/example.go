package service

import (
	"apiserver/handler"
)

// todo 這邊還是需要service 層 因為之後是GraphQL, 要在這邊組裝response

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
