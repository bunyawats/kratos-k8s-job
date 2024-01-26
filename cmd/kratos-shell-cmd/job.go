package main

import (
	"context"
	"fmt"
	_ "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	transgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	v1 "kratos-k8s-job/api/helloworld/v1"
	"log"
)

func runCommand(context.Context) error {
	msg := "message from K8S"
	fmt.Printf("Hello Kratos Application: %v\n\n", msg)

	callGRPC(msg)

	done <- true
	return nil
}

func callGRPC(msg string) {
	conn, err := transgrpc.DialInsecure(
		context.Background(),
		transgrpc.WithEndpoint("127.0.0.1:9000"),
		transgrpc.WithMiddleware(
			recovery.Recovery(),
		),
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := v1.NewGreeterClient(conn)
	reply, err := client.SayHello(context.Background(), &v1.HelloRequest{Name: msg})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[grpc] SayHello %+v\n", reply)

}
