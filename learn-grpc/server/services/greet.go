package services

import (
	"context"
	"fmt"
	"io"
	"learn-grpc/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GreetService struct {
	pb.GreetServiceServer
}

func NewGreetService() *GreetService {
	return &GreetService{}
}

func (s *GreetService) Greet(ctx context.Context, in *pb.GreetRequest) (*pb.GreetResponse, error) {
	return &pb.GreetResponse{Result: fmt.Sprintf("Hi %v", in.FirstName)}, nil
}

func (s *GreetService) GreetManyTimes(in *pb.GreetRequest, stream grpc.ServerStreamingServer[pb.GreetResponse]) error {
	for i := range 10 {
		stream.Send(&pb.GreetResponse{Result: fmt.Sprintf("Hi %v %d", in.FirstName, i)})
	}
	return nil
}

func (s *GreetService) LongGreet(stream grpc.ClientStreamingServer[pb.GreetRequest, pb.GreetResponse]) error {
	result := "Hi ..."

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.GreetResponse{Result: result})
		}

		if err != nil {
			fmt.Printf("LongGreet: failed streaming: %v\n", err)
			return status.Errorf(codes.Internal, "%v", err)
		}

		fmt.Printf("LongGreet: receiving req long greet: %v\n", req)
		result += " " + req.FirstName // append name to result
	}
}

func (s *GreetService) GreetEveryOne(stream grpc.BidiStreamingServer[pb.GreetRequest, pb.GreetResponse]) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			fmt.Printf("GreetEveryOne: failed reading client streaming: %v\n", err)
			return status.Errorf(codes.Internal, "%v", err)
		}

		res := "Helo " + req.FirstName
		if err = stream.Send(&pb.GreetResponse{Result: res}); err != nil {
			fmt.Printf("GreetEveryOne: failed sending stream: %v\n", err)
			return status.Errorf(codes.Internal, "%v", err)
		}
	}
}
