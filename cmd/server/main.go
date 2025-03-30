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
	tlsh "sokoide.com/grpc/pkg/tlshelper"
	pb "sokoide.com/grpc/proto"
)

var ()

type options struct {
	port int
	ka   bool
	tlso tlsh.TlsOptions
}

var opts options = options{
	port: 50051,
	ka:   true,
	tlso: tlsh.TlsOptions{
		Tls:          "none",
		Cert:         "cert.pem",
		Key:          "key.pem",
		Cacert:       "cacert.pem",
		AllowedUsers: "",
	},
}

func parseFlags() {
	flag.IntVar(&opts.port, "port", opts.port, "The server port")
	flag.BoolVar(&opts.ka, "keepalive", opts.ka, "keepalive")
	flag.StringVar(&opts.tlso.Tls, "tls", opts.tlso.Tls, "none|oneway|mtls")
	flag.StringVar(&opts.tlso.Cert, "cert", opts.tlso.Cert, "full path of cert")
	flag.StringVar(&opts.tlso.Key, "key", opts.tlso.Key, "full path of key")
	flag.StringVar(&opts.tlso.Cacert, "cacert", opts.tlso.Cacert, "full path of CA cert")
	flag.StringVar(&opts.tlso.AllowedUsers, "allowedUsers", opts.tlso.AllowedUsers, "comma separated allowed users")
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

func (s *server) Push(ctx context.Context, in *pb.PushRequest) (*pb.PushReply, error) {
	log.Printf("Received: %d bytes", len(in.GetData()))
	return &pb.PushReply{Message: fmt.Sprintf("%d bytes received", len(in.GetData()))}, nil
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

	switch opts.tlso.Tls {
	case "oneway":
		creds := tlsh.LoadKeyPairSingle(opts.tlso)
		transportSecurityOpt := grpc.Creds(creds)
		so = append(so, transportSecurityOpt)
	case "mtls":
		creds := tlsh.LoadKeyPairMutual(opts.tlso)
		transportSecurityOpt := grpc.Creds(creds)
		so = append(so, transportSecurityOpt)
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
