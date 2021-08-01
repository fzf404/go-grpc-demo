package main

import (
	"errors"
	"fmt"
	"grpc-server/data"
	"grpc-server/pb"
	"io"
	"log"
	"net"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

const port = ":5000"

func main() {
	// 监听端口
	listen, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalln(err.Error())
	}
	// 证书加载
	creds, err := credentials.NewServerTLSFromFile("cert.pem", "key.pem")
	if err != nil {
		log.Fatalln(err.Error())
	}
	options := []grpc.ServerOption{grpc.Creds(creds)}
	// 初始化server
	server := grpc.NewServer(options...)
	pb.RegisterEmployeeServiceServer(server, new(employeeService))
	log.Println("gRPC Server Started.")
	log.Println("Listen in", port)
	// 持续监听
	server.Serve(listen)
}

// 定义方法实现结构体
type employeeService struct{}

// 通过序号获得职员信息
func (s *employeeService) GetByNo(ctx context.Context, req *pb.GetByNoRequest) (*pb.EmployeeResponse, error) {
	// 遍历Data
	for _, e := range data.Employees {

		if req.No == e.No {
			// 返回响应
			return &pb.EmployeeResponse{
				Employee: &e,
			}, nil
		}
	}
	return nil, errors.New("emploee not found")
}

// 获得全部职员信息
func (s *employeeService) GetAll(req *pb.GetAllRequest, stream pb.EmployeeService_GetAllServer) error {

	for _, e := range data.Employees {
		stream.Send(&pb.EmployeeResponse{
			Employee: &e,
		})
		// 客户加钱优化
		time.Sleep(2 * time.Second)
	}

	return nil
}

// 上传图像
func (s *employeeService) AddPhoto(stream pb.EmployeeService_AddPhotoServer) error {
	// 获得metadata中的职员信息
	md, ok := metadata.FromIncomingContext(stream.Context())
	if ok {
		// 打印职员信息
		fmt.Printf("Employee: %s\n", md["no"][0])
	}

	img := []byte{}
	for {
		// 持续获得stream传来的二进制流
		data, err := stream.Recv()
		if err == io.EOF {
			fmt.Printf("File Size: %d\n", len(img))
			return stream.SendAndClose(&pb.AddPhotoResponse{IsOK: true})

		}
		if err != nil {
			return err
		}
		fmt.Printf("File Receivd: %d\n", len(data.Data))
		img = append(img, data.Data...)
	}
}

// 未实现
func (s *employeeService) Save(context.Context, *pb.EmployeeRequest) (*pb.EmployeeResponse, error) {
	return nil, nil
}

// 批量增加职员信息
func (s *employeeService) SaveAll(stream pb.EmployeeService_SaveAllServer) error {
	for {
		// 持续获得stream传来的信息
		empReq, err := stream.Recv()
		if err == io.EOF {
			break

		}
		if err != nil {
			return err
		}
		data.Employees = append(data.Employees, *empReq.Employee)
		stream.Send(&pb.EmployeeResponse{Employee: empReq.Employee})
	}
	// 遍历并打印当前全部Employee信息
	for _, emp := range data.Employees {
		fmt.Println(emp)
	}
	return nil
}

// 未实现
func (c *employeeService) CreateToken(context.Context, *pb.TokenRequest) (*pb.TokenResponse, error) {
	return nil, nil
}
