package dao

import (
	"context"
	"errors"
	"math"
)

// AddService 创建一个结构体实现服务中的两个方法就是实现了这个接口
type AddService struct{}

// 为实现服务中可能出现的问题创建错误
var (
	// ErrTwoZeroes  Sum方法的业务规则不能对两个0求和
	ErrTwoZeroes = errors.New("can't sum two zeroes")
	// ErrIntOverflow Sum参数越界
	ErrIntOverflow = errors.New("integer overflow")
	// ErrTwoEmptyStrings Concat方法业务规则规定参数不能是两个空字符串.
	ErrTwoEmptyStrings = errors.New("can't concat two empty strings")
	// ErrMaxSizeExceeded Concat方法的参数超出范围
	ErrMaxSizeExceeded = errors.New("result exceeds maximum size")
)

const maxLen = 20

// Sum 对两个数字求和  注：不对addService中的数据修改，可以使用值传递
func (add AddService) Sum(_ context.Context, a, b int) (int, error) {
	if a == 0 && b == 0 {
		return 0, ErrTwoZeroes
	}
	if (b > 0 && a > (math.MaxInt-b)) || (b < 0 && a < (math.MinInt-b)) {
		return 0, ErrIntOverflow
	}
	return a + b, nil
}

// Concat 连接两个字符串，实现AddService。
func (add AddService) Concat(_ context.Context, a, b string) (string, error) {
	if a == "" && b == "" {
		return "", ErrTwoEmptyStrings
	}
	if len(a)+len(b) > maxLen {
		return "", ErrMaxSizeExceeded
	}
	return a + b, nil
}
