package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"proto_demo/pb"
	"proto_demo/service"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func seedUsers(userStore service.UserStore) error {
	err := createUser(userStore, "admin1", "secret", "admin")
	if err != nil {
		return err
	}
	return createUser(userStore, "user1", "secret", "user")
}
func createUser(userstore service.UserStore, username, passwrod, role string) error {
	user, err := service.NewUser(username, passwrod, role)
	if err != nil {
		return err
	}
	return userstore.Save(user)
}

const (
	secreKey      = "secret"
	tokenDuration = 15 * time.Minute
)

func accessibleRoles() map[string][]string {
	const laptopServicePath = "/techschool.pcbook.LaptopService/"
	return map[string][]string{
		laptopServicePath + "CreateLaptop": {"admin"},
		laptopServicePath + "UploadImage":  {"admin"},
		laptopServicePath + "RateLaptop":   {"admin", "user"},
	}
}

func main() {
	port := flag.Int("port", 0, "the server port")
	flag.Parse()
	fmt.Println(*port)
	log.Printf("start server on port %d", *port)
	userStore := service.NewInMemoryUserStore()
	err := seedUsers(userStore)
	if err != nil {
		log.Fatal("cannot seed users")
	}
	jwtmanager := service.NewJwtManager(secreKey, tokenDuration)
	authServer := service.NewAuthServer(userStore, jwtmanager)

	laptopStore := service.NewInMemoryLaptopStore()
	imageStore := service.NewDiskImageStore("img")
	ratingStore := service.NewInMemoryRatingStore()
	laptopServer := service.NewLaptopService(laptopStore, imageStore, ratingStore)

	interceptor := service.NewAuthInterceptor(jwtmanager, accessibleRoles())
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.Unary()),   //一元拦截器
		grpc.StreamInterceptor(interceptor.Stream()), //流拦截器

	)
	pb.RegisterAuthServiceServer(grpcServer, authServer)
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
