package router

import (
	"google.golang.org/grpc"
	"morris/im-grpc/middleware"
	"morris/im-grpc/proto"
	"morris/im-grpc/server"
)

func InitRouter() *grpc.Server {

	//中间件   日志以及 recover
	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(middleware.LoggerInterceptor),
		grpc.ChainUnaryInterceptor(middleware.RecoverInterceptor),
	}

	//注册服务    todo  添加加密证书支持
	serv := grpc.NewServer(opts...)

	//注册搜索服务
	proto.RegisterSearchServiceServer(serv, &server.SearchService{})

	//注册用户服务

	//注册订单服务

	return serv
}
