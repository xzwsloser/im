package resultx

import (
	"context"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	zrpcErr "github.com/zeromicro/x/errors"
	"google.golang.org/grpc/status"
	"im-chat/pkg/xerr"
	"net/http"
)

/**
@Author: loser
@Description: 封装 http 请求返回结果
*/

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

// 返回数据最好范围 * 类型,便于处理 nil 值
func Success(data any) *Response {
	return &Response{
		Code: 200,
		Msg:  "",
		Data: data,
	}
}

func Fail(code int, err string) *Response {
	return &Response{
		Code: code,
		Msg:  err,
		Data: nil,
	}
}

func OkHandler(_ context.Context, v any) any {
	return Success(v)
}

/*
*
注意这里错误处理的逻辑: 首先对于返回的错误,
使用 Cause 获取到原始的错误,并且利用类型断言进行类型的转换,
最后得到错误码和错误信息即可
并且最后返回的对象都是 Response 形式的,
这里的 http.Response 的第一个参数其实是设置状态码
*/
func ErrHandler(name string) func(ctx context.Context, err error) (int, any) {
	return func(ctx context.Context, err error) (int, any) {
		errcode := xerr.SERVER_COMMON_ERROR
		errmsg := xerr.ErrMsg(errcode)

		causeErr := errors.Cause(err)
		if e, ok := causeErr.(*zrpcErr.CodeMsg); ok {
			errcode = e.Code
			errmsg = e.Msg
		} else {
			if gstatus, ok := status.FromError(causeErr); ok {
				errcode = int(gstatus.Code())
				errmsg = gstatus.Message()
			}
		}

		logx.WithContext(ctx).Errorf("[%s] err %v", name, err)
		return http.StatusBadRequest, Fail(errcode, errmsg)
	}
}
