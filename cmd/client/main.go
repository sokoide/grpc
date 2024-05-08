package main

import (
	"context"
	"flag"
	"log"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	tlsh "sokoide.com/grpc/pkg/tlshelper"
	pb "sokoide.com/grpc/proto"
)

const (
	defaultName = "Scott"
)

type options struct {
	addr  string
	name  string
	ms    int64
	loops int
	gos   int
	ka    bool
	tlso  tlsh.TlsOptions
}

var opts options = options{
	addr:  "localhost:50051",
	name:  "defaultName",
	ms:    100,
	loops: 10,
	gos:   10,
	ka:    true,
	tlso: tlsh.TlsOptions{
		Tls:          "none",
		Cert:         "cert.pem",
		Key:          "key.pem",
		Cacert:       "cacert.pem",
		AllowedUsers: "",
	},
}

var (
// addr  = flag.String("addr", "localhost:50051", "the address to connect to")
)

func chkerr(err error) {
	if err != nil {
		log.Fatalf("err: %v", err)
	}
}

func callSlow(id int, wg *sync.WaitGroup, c pb.GreeterClient, loops int) {
	for i := 0; i < loops; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		r2, err := c.Slow(ctx, &pb.SlowRequest{Ms: opts.ms})
		chkerr(err)
		log.Printf("[%04d] Slow: %s", id, r2.GetMessage())
	}
	wg.Done()
}

func parseFlags() {
	flag.StringVar(&opts.addr, "addr", opts.addr, "the address to connect to")
	flag.StringVar(&opts.name, "name", opts.name, "Name to greet")
	flag.Int64Var(&opts.ms, "ms", opts.ms, "milliseconds to sleep")
	flag.IntVar(&opts.loops, "loops", opts.loops, "loops")
	flag.IntVar(&opts.gos, "gos", opts.gos, "go routines")
	flag.BoolVar(&opts.ka, "keepalive", opts.ka, "keepalive")
	flag.StringVar(&opts.tlso.Tls, "tls", opts.tlso.Tls, "none|oneway|mtls")
	flag.StringVar(&opts.tlso.Cert, "cert", opts.tlso.Cert, "full path of cert")
	flag.StringVar(&opts.tlso.Key, "key", opts.tlso.Key, "full path of key")
	flag.StringVar(&opts.tlso.Cacert, "cacert", opts.tlso.Cacert, "full path of CA cert")
	flag.Parse()
}

func main() {
	parseFlags()

	var do []grpc.DialOption
	if opts.ka {
		ka := grpc.WithKeepaliveParams(
			keepalive.ClientParameters{
				Time:                10 * time.Second,
				Timeout:             5 * time.Second,
				PermitWithoutStream: true,
			},
		)
		do = append(do, ka)
	}
	do = append(do, grpc.WithTransportCredentials(insecure.NewCredentials()))

	var transportSecurityOpt grpc.DialOption

	switch opts.tlso.Tls {
	case "none":
		transportSecurityOpt = grpc.WithInsecure()
	case "oneway":
		creds := tlsh.LoadKeyPairSingle(opts.tlso)
		transportSecurityOpt = grpc.WithTransportCredentials(creds)
	case "mtls":
		creds := tlsh.LoadKeyPairMutual(opts.tlso)
		transportSecurityOpt = grpc.WithTransportCredentials(creds)
	default:
		panic("Error while setting up tls")
	}
	do = append(do, transportSecurityOpt)

	// Set up a connection to the server.
	conn, err := grpc.Dial(
		opts.addr,
		do...,
	)
	chkerr(err)
	defer conn.Close()

	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: opts.name})
	chkerr(err)
	log.Printf("Greeting: %s", r.GetMessage())

	var wg sync.WaitGroup
	for i := 0; i < opts.gos; i++ {
		wg.Add(1)
		go callSlow(i, &wg, c, opts.loops)
	}
	wg.Wait()
}
