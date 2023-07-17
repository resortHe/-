package service

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log"
	"proto_demo/pb"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const maxImageSize = 1 << 20

type LaptopServer struct {
	laptopStore LaptopStore
	imageStore  ImageStore
	ratingStore RatingStore
}

func NewLaptopService(laptopstore LaptopStore, imagestore ImageStore, ratingstore RatingStore) *LaptopServer {
	return &LaptopServer{laptopstore, imagestore, ratingstore}
}

func (service *LaptopServer) CreateLaptop(
	ctx context.Context,
	req *pb.CreateLaptopRequest,
) (*pb.CreateLaptopResponse, error) {
	laptop := req.GetLaptop()
	log.Printf("receive a create-laptop request with id :%s", laptop.Id)

	if len(laptop.Id) > 0 {
		_, err := uuid.Parse(laptop.Id)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "laptop ID is not a valid UUID:%v", err)
		}
	} else {
		id, err := uuid.NewRandom()
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "cannot generate a new laptop ID:%v", err)
		}
		laptop.Id = id.String()
	}
	//semo heavy processing set timeout
	//time.Sleep(6 * time.Second)

	if err := contextError(ctx); err != nil {
		return nil, err
	}

	err := service.laptopStore.Save(laptop)
	if err != nil {
		code := codes.Internal
		if errors.Is(err, ErrAlreadyExists) {
			code = codes.AlreadyExists
		}
		return nil, status.Errorf(code, "cannot save laptop to the store:%v", err)
	}
	log.Printf("saved laptop with id :%s", laptop.Id)
	res := &pb.CreateLaptopResponse{
		Id: laptop.Id,
	}
	return res, nil
}
func contextError(ctx context.Context) error {
	switch ctx.Err() {
	case context.Canceled:
		return logError(status.Error(codes.Canceled, "request is canceled"))
	case context.DeadlineExceeded:
		return logError(status.Error(codes.DeadlineExceeded, "deadline is exceeded"))
	default:
		return nil
	}

}
func (server *LaptopServer) SearchLaptop(
	req *pb.SearchLaptopRequest,
	stream pb.LaptopService_SearchLaptopServer) error {
	filter := req.GetFilter()
	log.Printf("receive a search-laptop request with filter:%v", filter)

	err := server.laptopStore.Search(
		stream.Context(),
		filter,
		func(laptop *pb.Laptop) error {
			res := &pb.SearchLaptopResponse{Laptop: laptop}
			err := stream.Send(res)
			if err != nil {
				return err
			}
			log.Printf("sent laptop with id : %s", laptop.GetId())
			return nil
		},
	)
	if err != nil {
		return status.Errorf(codes.Internal, "unexpected err:%v", err)
	}
	return nil

}
func (server *LaptopServer) UploadImage(stream pb.LaptopService_UploadImageServer) error {
	req, err := stream.Recv()
	if err != nil {
		return logError(status.Errorf(codes.Unknown, "cannot receive image info"))
	}
	lagtopID := req.GetInfo().GetLaptopId()
	imageTyep := req.GetInfo().GetImageType()
	log.Printf("receive an uploda-image request for laptop %s with image type %s", lagtopID, imageTyep)

	laptop, err := server.laptopStore.Find(lagtopID)
	if err != nil {
		return logError(status.Errorf(codes.Internal, "cannot find laptop %s", err))
	}
	if laptop == nil {
		return logError(status.Errorf(codes.InvalidArgument, "laptop %s doesn't exist", lagtopID))
	}
	imageData := bytes.Buffer{}
	imageSize := 0
	for {
		//
		if err := contextError(stream.Context()); err != nil {
			return err
		}
		log.Print("waiting to receive more data")

		req, err := stream.Recv()
		if err == io.EOF {
			log.Print("no more data")
			break
		}
		if err != nil {
			return logError(status.Errorf(codes.Unknown, "cannot receive chunk data :%v", err))
		}
		chunk := req.GetChunkData()
		size := len(chunk)

		log.Printf("received a chunk with size:%d", size)
		imageSize += size
		if imageSize > maxImageSize {
			return logError(status.Errorf(codes.InvalidArgument, "image is too large :%d > %d", imageSize, maxImageSize))
		}
		//time.Sleep(time.Second)

		_, err = imageData.Write(chunk)
		if err != nil {
			return logError(status.Errorf(codes.Internal, "cannot write chunk data %v", err))
		}

	}

	imageID, err := server.imageStore.Save(lagtopID, imageTyep, imageData)
	if err != nil {
		return logError(status.Errorf(codes.Internal, "cannot save image to the store %v", err))
	}
	res := &pb.UploadImageResponse{
		Id:   imageID,
		Size: uint32(imageSize),
	}

	err = stream.SendAndClose(res)
	if err != nil {
		return logError(status.Errorf(codes.Unknown, "cannot send response %v", err))
	}
	log.Printf("saved image with id: %s,size:%d", imageID, imageSize)

	return nil
}
func (service *LaptopServer) RateLaptop(stream pb.LaptopService_RateLaptopServer) error {
	for {
		err := contextError(stream.Context())
		if err != nil {
			return err
		}
		req, err := stream.Recv()
		if err == io.EOF {
			log.Printf("no more data")
			break
		}
		if err != nil {
			return logError(status.Errorf(codes.Unknown, "cannot receive stream request:%v", err))
		}

		laptopID := req.GetLaptopId()
		score := req.GetScore()

		log.Printf("recevied a rate-laptop request:id = %s , score= %.2f", laptopID, score)

		found, err := service.laptopStore.Find(laptopID)
		if err != nil {
			return logError(status.Errorf(codes.Internal, "cannot find laptop %v", err))
		}
		if found == nil {
			return logError(status.Errorf(codes.NotFound, "laptopID %s is not found", laptopID))
		}

		rating, err := service.ratingStore.Add(laptopID, score)
		if err != nil {
			return logError(status.Errorf(codes.Internal, "cannot add rating to the store: %v", err))
		}
		res := &pb.RateLaptopResponse{
			LaptopId:     laptopID,
			RatedCount:   rating.Count,
			AverageScore: rating.Sum / float64(rating.Count),
		}
		err = stream.Send(res)
		if err != nil {
			return logError(status.Errorf(codes.Unknown, "cannot send stream response %v", err))
		}
	}
	return nil
}
func logError(err error) error {
	if err != nil {
		log.Print(err)
	}
	return err
}
