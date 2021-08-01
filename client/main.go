package main

import (
	"context"
	"fmt"
	"grpc-client/pb"
	"io"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const port = ":5000"

func main() {
	// 加载证书
	creds, err := credentials.NewClientTLSFromFile("cert.pem", "")
	if err != nil {
		log.Fatalln(err.Error())
	}
	options := []grpc.DialOption{grpc.WithTransportCredentials(creds)}
	// 获得tcp连接
	conn, err := grpc.Dial("localhost"+port, options...)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer conn.Close()
	// 初始化客户端
	client := pb.NewEmployeeServiceClient(conn)
	// 调用各个方法
	GetByNo(client)
	GetAll(client)
	AddPhoto(client)
	SaveAll(client)

}

// 通过序号获得职员信息
func GetByNo(client pb.EmployeeServiceClient) {
	// 发送请求
	res, err := client.GetByNo(context.Background(), &pb.GetByNoRequest{No: 404})
	if err != nil {
		log.Fatalln(err.Error())
	}
	fmt.Println(res.Employee)
}

// 获得全部职员信息
func GetAll(client pb.EmployeeServiceClient) {
	// 创建请求，获得流
	stream, err := client.GetAll(context.Background(), &pb.GetAllRequest{})
	if err != nil {
		log.Fatalln(err.Error())
	}
	for {
		// 持续处理流
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err.Error())
		}
		fmt.Println(res.Employee)
	}

}

// 上传图像
func AddPhoto(client pb.EmployeeServiceClient) {
	imgFile, err := os.Open("demo.jpg")
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer imgFile.Close()
	// 使用metadata发送用户序号
	md := metadata.New(map[string]string{"no": "404"})
	// 创建上下文信息
	context := context.Background()
	context = metadata.NewOutgoingContext(context, md)
	// 创建请求，获得发送流
	stream, err := client.AddPhoto(context)
	if err != nil {
		log.Fatalln(err.Error())
	}

	for {
		chunk := make([]byte, 128*1024)
		chunkSize, err := imgFile.Read(chunk)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err.Error())
		}
		// 判断是否到达文件末尾，减小ChunkSize
		if chunkSize < len(chunk) {
			chunk = chunk[:chunkSize]
		}
		// 持续发送
		stream.Send(&pb.AddPhotoRequest{Data: chunk})

	}
	// 关闭并获得返回信息
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalln(err.Error())
	}
	fmt.Println(res)
}

// 批量上传职员信息
func SaveAll(client pb.EmployeeServiceClient) {
	// 教职员信息
	var employees = []pb.Employee{
		{
			Id:        4,
			No:        407,
			FirstName: "qi",
			LastName:  "xuan",
			MonthSalary: &pb.MonthSalary{
				Basic: 3000,
				Bonus: 100,
			},
			Status: pb.EmployeeStatus_NORMAL,
			LastModfied: &timestamppb.Timestamp{
				Seconds: time.Now().Unix(),
			},
		},
		{
			Id:        5,
			No:        408,
			FirstName: "meng",
			LastName:  "zi",
			MonthSalary: &pb.MonthSalary{
				Basic: 1000,
				Bonus: 50,
			},
			Status: pb.EmployeeStatus_NORMAL,
			LastModfied: &timestamppb.Timestamp{
				Seconds: time.Now().Unix(),
			},
		},
	}
	// 创建发送流
	stream, err := client.SaveAll(context.Background())
	if err != nil {
		log.Fatalln(err.Error())
	}
	// 创建管道
	finshChannel := make(chan struct{})
	// 匿名函数，用于处理服务端发送的返回流
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				finshChannel <- struct{}{}
				break
			}
			if err != nil {
				log.Fatalln(err.Error())
			}
			fmt.Println(res.Employee)
		}
	}() // 新的轻量级线程
	// 批量发送
	for _, e := range employees {
		err := stream.Send(&pb.EmployeeRequest{Employee: &e})
		if err != nil {
			log.Fatalln(err.Error())
		}
	}
	stream.CloseSend()
	<-finshChannel

}
