package grpcBegin

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"strconv"
)

func NewBolgServer(port int, f func(s *grpc.Server)) {
	lis, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(port))
	if err != nil {
		log.Fatal("监听端口"+strconv.Itoa(port)+"失败", err)
		return
	}

	s := grpc.NewServer()
	healthpb.RegisterHealthServer(s, health.NewServer())

	//反射接口
	reflection.Register(s)
	log.Println("grpc服务启动成功...")
	f(s)

	err = s.Serve(lis)
	if err != nil {
		log.Fatal("grpc服务启动失败", err)
		return
	}
}
