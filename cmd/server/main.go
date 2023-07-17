package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"proto_demo/pb"
	"proto_demo/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func unaryInterceptor(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	log.Print("-->unary interceptor:", info.FullMethod)
	return handler(ctx, req)
}

func streaminterceptor(srv interface{},
	stream grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	log.Print("-->stream interceptor:", info.FullMethod)
	return handler(srv, stream)

}

func main() {
	port := flag.Int("port", 0, "the server port")
	flag.Parse()
	fmt.Println(*port)
	log.Printf("start server on port %d", *port)
	laptopStore := service.NewInMemoryLaptopStore()
	imageStore := service.NewDiskImageStore("img")
	ratingStore := service.NewInMemoryRatingStore()
	laptopServer := service.NewLaptopService(laptopStore, imageStore, ratingStore)
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(unaryInterceptor),   //一元拦截器
		grpc.StreamInterceptor(streaminterceptor), //流拦截器

	)
	pb.RegisterLaptopServiceServer(grpcServer, laptopServer)
	//evans 反射  evans -r repl -p 8080启动evans
	reflection.Register(grpcServer)

	address := fmt.Sprintf("0.0.0.0:%d", *port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}

}
