package rpcserver

import (
	"context"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	zerr "github.com/zeromicro/x/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func LogInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	resp, err = handler(ctx, req)
	if err == nil {
		return resp, nil
	}

	logx.WithContext(ctx).Errorf("[RPC SEV ERR] %v", err)
	causeErr := errors.Cause(err)
	if e, ok := causeErr.(*zerr.CodeMsg); ok {
		// 下面一个句子说人话就是 err = errors.New(e.Msg) 当然也可以加上 erroCode
		err = status.Error(codes.Code(e.Code), e.Msg)
	}
	return resp, err
}
