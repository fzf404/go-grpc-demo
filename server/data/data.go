package data

import (
	"grpc-server/pb"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

var Employees = []pb.Employee{
	{
		Id:        1,
		No:        404,
		FirstName: "nmd",
		LastName:  "fzf",
		MonthSalary: &pb.MonthSalary{
			Basic: 5000,
			Bonus: 200,
		},
		Status: pb.EmployeeStatus_NORMAL,
		LastModfied: &timestamppb.Timestamp{
			Seconds: time.Now().Unix(),
		},
	},
	{
		Id:        2,
		No:        405,
		FirstName: "fzf",
		LastName:  "404",
		MonthSalary: &pb.MonthSalary{
			Basic: 6000,
			Bonus: 500,
		},
		Status: pb.EmployeeStatus_NORMAL,
		LastModfied: &timestamppb.Timestamp{
			Seconds: time.Now().Unix(),
		},
	},
	{
		Id:        3,
		No:        406,
		FirstName: "big",
		LastName:  "binys",
		MonthSalary: &pb.MonthSalary{
			Basic: 5000,
			Bonus: 1000,
		},
		Status: pb.EmployeeStatus_NORMAL,
		LastModfied: &timestamppb.Timestamp{
			Seconds: time.Now().Unix(),
		},
	},
}
