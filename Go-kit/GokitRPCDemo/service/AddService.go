package service

import "context"

// AddServiceInter 服务从业务逻辑开始，将服务建模为一个接口
type AddServiceInter interface {
	Sum(ctx context.Context, a, b int) (int, error)
	Concat(ctx context.Context, a, b string) (string, error)
}
