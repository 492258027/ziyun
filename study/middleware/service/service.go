package service

import (
	"errors"
	"strings"
)

// 把业务逻辑服务抽象为接口
type IStringServicer interface {
	Uppercase(string) (string, error)
	Count(string) int
}

// 定义一个类
type StringService struct{}

func (StringService) Count(s string) int {
	return len(s)
}

func (StringService) Uppercase(s string) (string, error) {
	if s == "" {
		return "", errors.New("empty string")
	}
	return strings.ToUpper(s), nil
}

type ServiceMiddleware func(IStringServicer) IStringServicer
