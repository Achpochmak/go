package middleware

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func Logging(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		header := md.Get("x-my-header")
		log.Printf("[interceptor.Logging] header: %v", header)
	}
	raw, _ := protojson.Marshal((req).(proto.Message))
	log.Printf("[interceptor.Logging] start: %v, %v", info.FullMethod, string(raw))
	resp, err = handler(ctx, req)
	if err != nil {
		log.Printf("[interceptor.Logging] error:%v", err.Error())
		return nil, err
	}
	log.Print("[interceptor.Logging] end")
	return
}
