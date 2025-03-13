# 构建 rpc 结构
goctl rpc protoc apps/social/rpc/social.proto --go_out=apps/social/rpc --go-grpc_out=apps/social/rpc --zrpc_out=apps/social/rpc

# 构建 models 结构,注意包名称
goctl model mysql ddl -src="./deploy/sql/social.sql" -dir="./apps/social/socialmodels/" -c

# 构建社交服务
goctl api go -api apps/social/api/social.api -dir apps/social/api -style gozero