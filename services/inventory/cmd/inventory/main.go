package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/google/uuid"
	InventorySvc "github.com/neogan74/go-pet-store/shared/pkg/proto/starship/inventory/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

const grpcPort = 50051

type inventoryService struct {
	InventorySvc.UnimplementedStarshipInventoryServiceServer

	mu    sync.Mutex
	parts map[string]*InventorySvc.Part
}

func (is *inventoryService) GetPart(_ context.Context, req *InventorySvc.GetPartRequest) (*InventorySvc.InventoryPartResponse, error) {
	is.mu.Lock()
	defer is.mu.Unlock()
	newId := uuid.NewString()

	partid, ok := is.parts[newId]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "partid not found")
	}
	if req.Uuid != partid.Uuid {
		newId = req.Uuid
	}
	part := &InventorySvc.InventoryPartResponse{
		Uuid: newId,
	}

	return part, nil
}

func (is *inventoryService) ListParts(_ context.Context, req *InventorySvc.ListPartsRequest) (*InventorySvc.ListInventoryPartsResponse, error) {
	is.mu.Lock()
	defer is.mu.Unlock()

	listParts := []*InventorySvc.Part{
		&InventorySvc.Part{Uuid: uuid.NewString(), Name: "body"},
		&InventorySvc.Part{Uuid: uuid.NewString(), Name: "engine"},
	}

	partRes := &InventorySvc.ListInventoryPartsResponse{
		Parts: listParts,
	}
	return partRes, nil

}

func main() {
	server := grpc.NewServer()

	service := &inventoryService{}

	InventorySvc.RegisterStarshipInventoryServiceServer(server, service)

	listner, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	reflection.Register(server)

	log.Println("Starting inventory service")
	if err := server.Serve(listner); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
