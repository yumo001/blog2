package main

import (
	"google.golang.org/grpc"

	"github.com/yumo001/blog2/grpcBegin"
	"github.com/yumo001/blog2/initialize"
	"github.com/yumo001/blog2/logic"
)

func init() {
	initialize.Viper()
	initialize.Nacos()
	initialize.Consul()
}

func main() {
	grpcBegin.NewBolgServer(8081, func(s *grpc.Server) {
		logic.RegisterUsersServer(s)
		logic.RegisterGoodsServer(s)
	})
}
