// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type BaseResponse interface {
	IsBaseResponse()
	GetCode() int
	GetMessage() string
}

type Mutation struct {
}

type MutationEchoOutput struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"Data"`
}

func (MutationEchoOutput) IsBaseResponse()         {}
func (this MutationEchoOutput) GetCode() int       { return this.Code }
func (this MutationEchoOutput) GetMessage() string { return this.Message }

type Query struct {
}

type QueryEchoOutput struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"Data"`
}

func (QueryEchoOutput) IsBaseResponse()         {}
func (this QueryEchoOutput) GetCode() int       { return this.Code }
func (this QueryEchoOutput) GetMessage() string { return this.Message }
