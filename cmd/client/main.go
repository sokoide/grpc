package main

import (
	"context"
	"flag"
	"log"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "sokoide.com/grpc/proto"
)

const (
	defaultName = "Scott"
)

var (
	addr  = flag.String("addr", "localhost:50051", "the address to connect to")
	name  = flag.String("name", defaultName, "Name to greet")
	ms    = flag.Int64("ms", 100, "milliseconds to sleep")
	loops = flag.Int("loops", 10, "loops")
	gos   = flag.Int("gos", 10, "go routines")
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

		r2, err := c.Slow(ctx, &pb.SlowRequest{Ms: *ms})
		chkerr(err)
		log.Printf("[%04d] Slow: %s", id, r2.GetMessage())
	}
	wg.Done()
}

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	chkerr(err)
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: *name})
	chkerr(err)
	log.Printf("Greeting: %s", r.GetMessage())

	var wg sync.WaitGroup
	for i := 0; i < *gos; i++ {
		wg.Add(1)
		go callSlow(i, &wg, c, *loops)
	}
	wg.Wait()
}
