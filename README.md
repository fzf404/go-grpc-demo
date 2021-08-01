# Go语言的gRPC例子

> 单向流与双向流加密传输

## 服务端

```bash
cd ./server

# 直接运行
go run main.go

# 或
# 下载依赖包
go mod download

# 编译proto
protoc --proto_path=./protos ./protos/*.proto --go_out=plugins=grpc:.

# 生成ssl证书
openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cery.pem -days 365 -nodes -subj /CN=localhost

# 运行
go run main.go
```

## 客户端

```bash

cd ./client

# 直接运行
go run main.go

# 或
# 拷贝证书
cp ../server/cert.pem .
# 设置证书环境变量
export GODEBUG=x509ignoreCN=0

# 编译proto
protoc --proto_path=./protos ./protos/*.proto --go_out=plugins=grpc:.

# 运行
go run main.go
```