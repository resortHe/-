package service_test

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"proto_demo/pb"
	"proto_demo/sample"
	"proto_demo/serializer"
	"proto_demo/service"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func TestClientCreateLaptop(t *testing.T) {
	t.Parallel()
	laptopstore := service.NewInMemoryLaptopStore()
	serverAddress := startTestLaptopServer(t, laptopstore, nil, nil)
	laptopClient := newTestLaptopClient(t, serverAddress)

	laptop := sample.NewLaptop()
	expectedID := laptop.Id
	req := &pb.CreateLaptopRequest{
		Laptop: laptop,
	}

	res, err := laptopClient.CreateLaptop(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, expectedID, res.Id)

	other, err := laptopstore.Find(res.Id)
	require.NoError(t, err)
	require.NotNil(t, other)

	requireSameLaptop(t, laptop, other)

}
func TestClientSearchLaptop(t *testing.T) {
	t.Parallel()

	filter := &pb.Filter{
		MaxPriceUsd: 2000,
		MinCpuCores: 4,
		MinCpuGhz:   2.2,
		MinRam:      &pb.Memory{Value: 8, Unit: pb.Memory_GIGABYTE},
	}
	laptopstore := service.NewInMemoryLaptopStore()

	expectedIDs := make(map[string]bool)

	for i := 0; i < 6; i++ {
		laptop := sample.NewLaptop()

		switch i {
		case 0:
			laptop.PriceUsd = 2500
		case 1:
			laptop.Cpu.NumberCores = 2
		case 2:
			laptop.Cpu.MinGhz = 2.2
		case 3:
			laptop.Ram = &pb.Memory{Value: 4096, Unit: pb.Memory_MEGABYTE}
		case 4:
			laptop.PriceUsd = 1999
			laptop.Cpu.NumberCores = 4
			laptop.Cpu.MinGhz = 2.5
			laptop.Cpu.MaxGhz = 4.5
			laptop.Ram = &pb.Memory{Value: 16, Unit: pb.Memory_GIGABYTE}
			expectedIDs[laptop.Id] = true
		case 5:
			laptop.PriceUsd = 2000
			laptop.Cpu.NumberCores = 6
			laptop.Cpu.MinGhz = 2.8
			laptop.Cpu.MaxGhz = 5.0
			laptop.Ram = &pb.Memory{Value: 64, Unit: pb.Memory_GIGABYTE}
			expectedIDs[laptop.Id] = true
		}

		err := laptopstore.Save(laptop)
		require.NoError(t, err)
	}
	serverAddress := startTestLaptopServer(t, laptopstore, nil, nil)
	laptopClient := newTestLaptopClient(t, serverAddress)

	req := &pb.SearchLaptopRequest{Filter: filter}
	stream, err := laptopClient.SearchLaptop(context.Background(), req)
	require.NoError(t, err)
	found := 0
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		require.NoError(t, err)
		require.Contains(t, expectedIDs, res.GetLaptop().GetId())
		found += 1
	}
	require.Equal(t, len(expectedIDs), found)

}
func TestClientUploadImage(t *testing.T) {
	t.Parallel()

	testImageFolder := "../tmp"
	laptopstore := service.NewInMemoryLaptopStore()
	imagestore := service.NewDiskImageStore(testImageFolder)

	laptop := sample.NewLaptop()
	err := laptopstore.Save(laptop)
	require.NoError(t, err)

	serverAddress := startTestLaptopServer(t, laptopstore, imagestore, nil)
	laptopClient := newTestLaptopClient(t, serverAddress)

	imagePath := fmt.Sprintf("%s/laptop.jpg", testImageFolder)
	file, err := os.Open(imagePath)
	require.NoError(t, err)
	defer file.Close()

	stream, err := laptopClient.UploadImage(context.Background())
	require.NoError(t, err)

	imageTyep := filepath.Ext(imagePath)
	req := &pb.UploadImageRequest{
		Data: &pb.UploadImageRequest_Info{
			Info: &pb.ImageInfo{
				LaptopId:  laptop.GetId(),
				ImageType: imageTyep,
			},
		},
	}

	err = stream.Send(req)
	require.NoError(t, err)
	reader := bufio.NewReader(file)
	buffer := make([]byte, 1024)
	size := 0

	for {

		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		require.NoError(t, err)
		size += n
		req := &pb.UploadImageRequest{
			Data: &pb.UploadImageRequest_ChunkData{
				ChunkData: buffer[:n],
			},
		}
		err = stream.Send(req)
		require.NoError(t, err)
	}
	res, err := stream.CloseAndRecv()
	require.NoError(t, err)
	require.NotZero(t, res.GetId())
	require.EqualValues(t, size, res.GetSize())

	savedImagePath := fmt.Sprintf("%s/%s%s", testImageFolder, res.GetId(), imageTyep)
	require.FileExists(t, savedImagePath)
	require.NoError(t, os.Remove(savedImagePath))

}
func TestClientRateImage(t *testing.T) {
	t.Parallel()

	laptopstore := service.NewInMemoryLaptopStore()
	ratingstore := service.NewInMemoryRatingStore()

	laptop := sample.NewLaptop()
	err := laptopstore.Save(laptop)
	require.NoError(t, err)

	serverAddress := startTestLaptopServer(t, laptopstore, nil, ratingstore)
	laptopClient := newTestLaptopClient(t, serverAddress)

	stream, err := laptopClient.RateLaptop(context.Background())
	require.NoError(t, err)

	scores := []float64{8, 7.5, 10}
	averages := []float64{8, 7.75, 8.5}

	n := len(scores)
	for i := 0; i < n; i++ {
		req := &pb.RateLaptopRequest{
			LaptopId: laptop.GetId(),
			Score:    scores[i],
		}

		err := stream.Send(req)
		require.NoError(t, err)
	}

	err = stream.CloseSend()
	require.NoError(t, err)

	for idx := 0; ; idx++ {
		res, err := stream.Recv()
		if err == io.EOF {
			require.Equal(t, n, idx)
			return
		}

		require.NoError(t, err)
		require.Equal(t, laptop.GetId(), res.GetLaptopId())
		require.Equal(t, uint32(idx+1), res.GetRatedCount())
		require.Equal(t, averages[idx], res.GetAverageScore())

	}

}
func startTestLaptopServer(t *testing.T, laptopstore service.LaptopStore, imagestore service.ImageStore, ratingstore service.RatingStore) string {
	laptopServer := service.NewLaptopService(laptopstore, imagestore, ratingstore)
	grpcService := grpc.NewServer()
	pb.RegisterLaptopServiceServer(grpcService, laptopServer)

	l, err := net.Listen("tcp", ":0") // random available port
	require.NoError(t, err)

	go grpcService.Serve(l)

	return l.Addr().String()

}
func newTestLaptopClient(t *testing.T, serverAddress string) pb.LaptopServiceClient {
	cc, err := grpc.Dial(serverAddress, grpc.WithInsecure()) //不安全连接
	require.NoError(t, err)
	return pb.NewLaptopServiceClient(cc)
}
func requireSameLaptop(t *testing.T, laptop1 *pb.Laptop, laptop2 *pb.Laptop) {
	s, err := serializer.ProtobufToJSON(laptop1)
	require.NoError(t, err)
	s2, err := serializer.ProtobufToJSON(laptop2)
	require.NoError(t, err)
	require.Equal(t, s, s2)

}
