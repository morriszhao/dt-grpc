package middleware

import (
	"context"
	"google.golang.org/grpc"
	"log"
)

// LoggerInterceptor grpc 中间件  一元拦截器
func LoggerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	log.Printf("grpc method: %s, %v", info.FullMethod, req)
	resp, err := handler(ctx, req)

	log.Printf("grpc method: %s, %v", info.FullMethod, resp)
	return resp, err
}
