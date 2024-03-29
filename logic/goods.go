package logic

import (
	"context"
	"github.com/yumo001/blog2/global"
	goods "github.com/yumo001/blog2/pb/goods"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GoodsServer struct {
	goods.UnimplementedGoodsServer
}

func RegisterGoodsServer(s grpc.ServiceRegistrar) {
	s.RegisterService(&goods.Goods_ServiceDesc, GoodsServer{})
}

func (GoodsServer) GoodsAdd(ctx context.Context, in *goods.GoodsAddRequest) (*goods.GoodsAddResponse, error) {
	err := global.MysqlDB.Table("goods").Create(&in.Good).Error
	if err != nil {
		return nil, status.Errorf(codes.Internal, "添加商品失败")
	}

	return nil, status.Errorf(codes.OK, "成功")
}
func (GoodsServer) GoodsList(ctx context.Context, in *goods.GoodsListRequest) (*goods.GoodsListResponse, error) {

	var goodsList []*goods.Good
	err := global.MysqlDB.Table("goods").Find(&goodsList).Error
	if err != nil {
		return nil, status.Errorf(codes.Internal, "获取商品列表失败")
	}
	return &goods.GoodsListResponse{
		Goods: goodsList,
	}, nil
}
