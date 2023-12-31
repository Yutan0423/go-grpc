package main

import (
	"context"
	"fmt"
	"go-proto/generated/pb"
	"io"
	"log"
	"os"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	fileServiceClient := pb.NewFileServiceClient(conn)
	// callListFiles(fileServiceClient)
	// callDownload(fileServiceClient)
	callUpload(fileServiceClient)
}

func callListFiles(c pb.FileServiceClient) {
	res, err := c.ListFiles(context.Background(), &pb.ListFilesRequest{})
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(res.GetFilenames())
}

func callDownload(c pb.FileServiceClient) {
	req := &pb.DownloadRequest{Filename: "name.txt"}
	stream, err := c.Download(context.Background(), req)
	if err != nil {
		log.Fatalln(err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err)
		}

		log.Printf("Response from Download(bytes): %v", res.GetData())
		log.Printf("Response from Download(string): %v", string(res.GetData()))
	}
}

func callUpload(c pb.FileServiceClient) {
	filename := "sport.txt"
	path := "/Users/yutan/Desktop/go-proto/storage/" + filename

	file, err := os.Open(path)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	stream, err := c.Upload(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
	buf := make([]byte, 5)
	for {
		n, err := file.Read(buf)
		if n == 0 || err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err)
		}

		req := &pb.UploadRequest{Data: buf[:n]}
		sendErr := stream.Send(req)
		if sendErr != nil {
			log.Fatalln(err)
		}

		// time.Sleep(1 * time.Second)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("received data size: %v", res.GetSize())
}
