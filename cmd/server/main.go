package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	pb "sokoide.com/grpc/proto"
)

var ()

type options struct {
	port int
	ka   bool
}

var opts options = options{
	port: 50051,
	ka:   true,
}

func parseFlags() {
	flag.IntVar(&opts.port, "port", opts.port, "The server port")
	flag.BoolVar(&opts.ka, "keepalive", opts.ka, "keepalive")
	flag.Parse()
}

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func (s *server) Slow(ctx context.Context, in *pb.SlowRequest) (*pb.SlowReply, error) {
	ms := in.GetMs()
	log.Printf("sleeping %v ms...", ms)
	time.Sleep(time.Millisecond * time.Duration(ms))
	return &pb.SlowReply{Message: fmt.Sprintf("Slept %d ms", ms)}, nil
}

func main() {
	parseFlags()

	var so []grpc.ServerOption

	if opts.ka {
		kaPolicy := grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             5 * time.Second,
			PermitWithoutStream: true,
		})
		so = append(so, kaPolicy)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", opts.port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(so...)
	pb.RegisterGreeterServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
