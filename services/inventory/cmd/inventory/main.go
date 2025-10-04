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

	// generate some db
	service.parts = make(map[string]*InventorySvc.Part)
	service.parts["body"] = &InventorySvc.Part{Uuid: "c0eadac2-c9a3-47b9-ad28-9791f75bcba5", Name: "body"}
	service.parts["wing"] = &InventorySvc.Part{Uuid: "f485dcc2-a2f3-4a7c-a3b0-5a1b6652b3a0", Name: "wing"}
	service.parts["engine"] = &InventorySvc.Part{Uuid: "e5c46c2d-961c-469b-a040-6714d5c71320", Name: "wing"}

	InventorySvc.RegisterStarshipInventoryServiceServer(server, service)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	reflection.Register(server)

	log.Println("Starting inventory service")
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
