package main

import (
	"log"
	"morris/im-grpc/etcdd"
	"morris/im-grpc/router"
	"net"
)

func main() {

	//监听端口
	listener, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatal("tcp监听端口8888失败:", err.Error())
		return
	}

	//路由启动
	serv := router.InitRouter()

	//etcd 服务注册
	etcdd.InitRegisterEtcd()

	//启动grpc 服务
	err = serv.Serve(listener)
	if err != nil {
		log.Fatal("grpc监听端口8888失败:", err.Error())
	}
}
