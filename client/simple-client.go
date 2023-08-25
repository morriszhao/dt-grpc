package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"morris/im-grpc/proto"
	"time"
)

func main() {

	//必须要设置transport 证书、否则拨号失败
	conn, err := grpc.Dial("localhost:8888", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("grpc 拨号失败")
		return
	}
	defer conn.Close()

	//实例化一个客户端
	grpcSearchClient := proto.NewSearchServiceClient(conn)

	//发起rpc请求
	timeoutCtx, _ := context.WithTimeout(context.Background(), time.Second)
	ctx, cancelFunc := context.WithTimeout(timeoutCtx, time.Second)
	defer cancelFunc()

	response, err := grpcSearchClient.Search(ctx, &proto.SearchRequest{Request: "通信好了吗"})
	if err != nil {
		log.Fatal("grpc 远程调用失败", err.Error())
		return
	}

	fmt.Println(response.Response)

}
