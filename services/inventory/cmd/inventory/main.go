package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

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

	partid, ok := is.parts[req.Uuid]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "partid not found")
	}
	part := &InventorySvc.InventoryPartResponse{
		Uuid: partid.Uuid,
	}

	return part, nil
}

func (is *inventoryService) ListParts(_ context.Context, req *InventorySvc.ListPartsRequest) (*InventorySvc.ListInventoryPartsResponse, error) {
	is.mu.Lock()
	defer is.mu.Unlock()

	filter := req.Filter
	fmt.Printf("Filter: %v\n", filter)

	listParts := make([]*InventorySvc.Part, 0)
	for key, value := range is.parts {
		fmt.Printf("%v -> %V\n", key, value)
		listParts = append(listParts, value)
	}
	// TODO: Implement filtering

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
	service.parts["c0eadac2-c9a3-47b9-ad28-9791f75bcba5"] = &InventorySvc.Part{Uuid: "c0eadac2-c9a3-47b9-ad28-9791f75bcba5", Name: "body", Category: []InventorySvc.Category{InventorySvc.Category_CATEGORY_PORTHOLE}}
	service.parts["f485dcc2-a2f3-4a7c-a3b0-5a1b6652b3a0"] = &InventorySvc.Part{Uuid: "f485dcc2-a2f3-4a7c-a3b0-5a1b6652b3a0", Name: "wing", Category: []InventorySvc.Category{InventorySvc.Category_CATEGORY_WING}}
	service.parts["e5c46c2d-961c-469b-a040-6714d5c71320"] = &InventorySvc.Part{Uuid: "e5c46c2d-961c-469b-a040-6714d5c71320", Name: "engine", Category: []InventorySvc.Category{InventorySvc.Category_CATEGORY_ENGINE}}

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
