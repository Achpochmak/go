package main

import (
	"context"
	"log"

	"HOMEWORK-1/pkg/api/proto/pvz/v1/pvz/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	target = "localhost:50051"
)

func main() {
	conn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
		return
	}
	defer conn.Close()

	client := pvz.NewPVZClient(conn)

	ctx := context.Background()
	ctx = metadata.AppendToOutgoingContext(ctx, "x-my-header", "123")

	_, errAdd := client.AddOrder(ctx, &pvz.AddOrderRequest{
		Order: &pvz.Order{
			Id:         16,
			IdReceiver: 12345,
			Packaging:  "bag",
			Weight:     1.5,
			// Add other fields as necessary
		},
	})
	if errAdd != nil {
		status := status.Code(errAdd)
		if status == codes.InvalidArgument {
			return
		}
		log.Fatal(errAdd)
		return
	}

	err = list(ctx, client)
	if err != nil {
		log.Fatal(err)
		return
	}
	err = deleteOrder(ctx, client)
	if err != nil {
		log.Fatal(err)
		return
	}
	err = find(ctx, client)
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Println("Client done")
}

func list(ctx context.Context, client pvz.PVZClient) error {
	resp, errList := client.ListOrder(ctx, &pvz.ListOrderRequest{})
	if errList != nil {
		status := status.Code(errList)
		if status == codes.InvalidArgument {
			return nil
		}
		log.Fatal(errList)
		return nil
	}
	for _, v := range resp.GetList() {
		log.Println(v)
	}
	return nil
}

func deleteOrder(ctx context.Context, client pvz.PVZClient) error {
	_, errDelete := client.DeleteOrder(ctx, &pvz.DeleteOrderRequest{
		ID: &pvz.ID{
			Id: 16, // specify the correct ID here
		},
	})
	if errDelete != nil {
		status := status.Code(errDelete)
		if status == codes.InvalidArgument {
			return nil
		}
		log.Fatal(errDelete)
		return nil
	}
	return nil
}

func find(ctx context.Context, client pvz.PVZClient) error {
	f, errFind := client.FindOrder(ctx, &pvz.FindOrderRequest{
		Id: &pvz.ID{
			Id: 16, // specify the correct ID here
		},
	})
	if errFind != nil {
		status := status.Code(errFind)
		if status == codes.InvalidArgument {
			return nil
		}
		log.Fatal(errFind)
		return nil
	}
	log.Println(f.GetOrder())
	return nil
}
