package logic

import (
	"context"
	"crypto/md5"
	"fmt"
	"github.com/yumo001/blog2/global"
	users "github.com/yumo001/blog2/pb/users"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UsersServer struct {
	users.UnimplementedUsersServer
}

func RegisterUsersServer(s grpc.ServiceRegistrar) {
	s.RegisterService(&users.Users_ServiceDesc, UsersServer{})
}

func (UsersServer) Ping(ctx context.Context, in *users.Request) (*users.Response, error) {
	if in.Ping != "ping" {
		return nil, status.Errorf(codes.InvalidArgument, "不符合所需参数")
	}

	return &users.Response{
		Pong: "pong",
	}, nil
}

func (UsersServer) Register(ctx context.Context, in *users.RegisterRequest) (*users.RegisterResponse, error) {
	var count int64

	err := global.MysqlDB.Table("users").Where("username = ?", in.User.Username).Count(&count).Error
	if err != nil {
		return nil, status.Errorf(codes.Unimplemented, "查询失败")
	}
	if count > 0 {
		return nil, status.Errorf(codes.AlreadyExists, "用户名已存在")
	}

	err = global.MysqlDB.Table("users").Create(&in.User).Error
	if err != nil {
		return nil, status.Errorf(codes.Unimplemented, "注册失败")
	}

	return nil, status.Errorf(codes.OK, "注册成功")
}

func (UsersServer) Login(ctx context.Context, in *users.LoginRequest) (*users.LoginResponse, error) {
	var count int64
	err := global.MysqlDB.Table("users").Where("username = ?", in.User.Username).Count(&count).Error
	if err != nil {
		return nil, status.Errorf(codes.Unimplemented, "查询失败")
	}
	if count <= 0 {
		return nil, status.Errorf(codes.NotFound, "该用户不存在")
	}

	var relPwd string
	err = global.MysqlDB.Table("users").Where("username =?", in.User.Username).Pluck("password", &relPwd).Error
	if err != nil {
		return nil, status.Errorf(codes.Unimplemented, "查询失败")
	}
	if encrypted(in.User.Password) != relPwd {
		return nil, status.Errorf(codes.InvalidArgument, "密码错误")
	}

	return nil, nil
}

// 密码加密
func encrypted(pwd string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(pwd)))
}

func (UsersServer) List(ctx context.Context, in *users.ListRequest) (*users.ListResponse, error) {

	var us []*users.User
	err := global.MysqlDB.Table("users").Find(&us).Error
	if err != nil {
		return nil, status.Errorf(codes.Unimplemented, "查询失败")
	}

	return &users.ListResponse{
		Users: us,
	}, nil

}
