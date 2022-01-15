package testingutils

import (
	"context"
	"github.com/golang/mock/gomock"
)

type anyContext struct {
}

func (a anyContext) Matches(x interface{}) bool {
	_, err := x.(context.Context)
	return !err
}

func (a anyContext) String() string {
	panic("checks if argument is a context")
}

func AnyContext() gomock.Matcher {
	return anyContext{}
}
