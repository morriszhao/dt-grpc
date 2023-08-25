package main

import (
	"log"
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

	serv := router.InitRouter()
	err = serv.Serve(listener)
	if err != nil {
		log.Fatal("grpc监听端口8888失败:", err.Error())
	}
}
