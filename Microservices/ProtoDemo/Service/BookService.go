package Service

import (
	"ProtoDemo/pb"
	"context"
)

type BookService struct {
	pb.UnimplementedBookServiceServer //根据proto文件中的BookService生成的
}

// CreateBook 实现BookService中的两个方法
func (b *BookService) CreateBook(ctx context.Context) {

}
