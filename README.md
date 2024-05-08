# Quick gRPC example

## How to build

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
make
```

## How to run

### Server

```
./server
```

### Client

```
# runs client with 1 goroutine x 3 loops of 500ms call
$ time ./client -ms 500 -loops 3 -gos 1
2024/05/08 19:47:01 Greeting: Hello Scott
2024/05/08 19:47:02 [0000] Slow: Slept 500 ms
2024/05/08 19:47:02 [0000] Slow: Slept 500 ms
2024/05/08 19:47:03 [0000] Slow: Slept 500 ms
./client -ms 500 -loops 3 -gos 1  0.01s user 0.01s system 1% cpu 1.525 total


# runs client with 1000 goroutine x 3 loops of 500ms call
$ time ./client -ms 500 -loops 3 -gos 1000
...
./client -ms 500 -loops 3 -gos 1000  0.18s user 0.06s system 15% cpu 1.566 total
```
