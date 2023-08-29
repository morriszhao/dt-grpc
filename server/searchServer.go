package server

import (
	"context"
	"morris/im-grpc/proto"
	"time"
)

type SearchService struct {
	proto.UnimplementedSearchServiceServer
}

func (s *SearchService) Search(c context.Context, request *proto.SearchRequest) (*proto.SearchResponse, error) {

	//模拟正常业务操作
	time.Sleep(time.Second * 1)

	//如果已经超时  就不用返回了
	deadline, ok := c.Deadline()
	if !ok {
		//没有设置超时时间
		requestName := request.GetRequest()
		return &proto.SearchResponse{Response: "nihao, 收到客户端请求是数据：" + requestName}, nil
	}

	if time.Now().Before(deadline) {
		//未到超时时间
		requestName := request.GetRequest()
		return &proto.SearchResponse{Response: "nihao, 收到客户端请求是数据：" + requestName}, nil
	}

	return nil, nil

}
func getTimeoutFromContext(c context.Context) time.Duration {
	deadline, ok := c.Deadline()
	if !ok {
		return time.Duration(0)
	}

	return time.Until(deadline)
}
