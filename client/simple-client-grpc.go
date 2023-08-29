package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"log"
	"morris/im-grpc/etcdd"
	"morris/im-grpc/proto"
	"time"
)

func main() {

	//注册自己的 解析器
	etcdResolverBuilder := etcdd.NewEtcdResolverBuilder()
	resolver.Register(etcdResolverBuilder)

	/**
	1、etcd::// 协议 表示使用etcd      /grpc/user_service/ 前缀key
	2、必须要设置transport 证书、否则拨号失败
	*/
	conn, err := grpc.Dial("etcd:///grpc/user_service/", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("grpc 拨号失败")
		return
	}
	defer conn.Close()

	//实例化一个客户端
	grpcSearchClient := proto.NewSearchServiceClient(conn)

	//发起rpc请求
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelFunc()

	response, err := grpcSearchClient.Search(ctx, &proto.SearchRequest{Request: "通信好了吗"})
	if err != nil {
		log.Fatal("grpc 远程调用失败", err.Error())
		return
	}

	fmt.Println(response.Response)

}
