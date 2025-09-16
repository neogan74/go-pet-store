package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/runtime/middleware/untyped"
	"github.com/go-openapi/swag"
)

// NewPetstore creates a new petstore api handler
func NewPetstore() (http.Handler, error) {
	spec, err := loads.Analyzed(json.RawMessage([]byte(swaggerJSON)), "")
	if err != nil {
		return nil, err
	}
	api := untyped.NewAPI(spec)

	api.RegisterOperation("get", "/pets", getAllPets)
	api.RegisterOperation("post", "/pets", createPet)
	api.RegisterOperation("delete", "/pets/{id}", deletePet)
	api.RegisterOperation("get", "/pets/{id}", getPetByID)

	return middleware.Serve(spec, api), nil
}

var getAllPets = runtime.OperationHandlerFunc(func(data interface{}) (interface{}, error) {
	fmt.Println("getAllPets")
	fmt.Printf("%#v\n", data)
	return pets, nil
})

var createPet = runtime.OperationHandlerFunc(func(data interface{}) (interface{}, error) {
	fmt.Println("createPet")
	fmt.Printf("%#v\n", data)
	body := data.(map[string]interface{})["pet"]
	var pet Pet
	if err := swag.FromDynamicJSON(body, &pet); err != nil {
		return nil, err
	}
	addPet(pet)
	return body, nil
})

var deletePet = runtime.OperationHandlerFunc(func(data interface{}) (interface{}, error) {
	fmt.Println("deletePet")
	fmt.Printf("%#v\n", data)
	id := data.(map[string]interface{})["id"].(int64)
	removePet(id)
	return nil, nil
})

var getPetByID = runtime.OperationHandlerFunc(func(data interface{}) (interface{}, error) {
	fmt.Println("getPetByID")
	fmt.Printf("%#v\n", data)
	id := data.(map[string]interface{})["id"].(int64)
	return petByID(id)
})

// Tag the tag model
type Tag struct {
	ID   int64
	Name string
}

// Pet the pet model
type Pet struct {
	ID        int64    `json:"id"`
	Name      string   `json:"name"`
	PhotoURLs []string `json:"photoUrls,omitempty"`
	Status    string   `json:"status,omitempty"`
	Tags      []Tag    `json:"tags,omitempty"`
}

var pets = []Pet{
	{ID: 1, Name: "Dog", PhotoURLs: []string{}, Status: "available", Tags: nil},
	{ID: 2, Name: "Cat", PhotoURLs: []string{}, Status: "pending", Tags: nil},
	{ID: 3, Name: "Parrot", PhotoURLs: []string{}, Status: "delivering", Tags: nil},
}

var (
	petsLock        = &sync.Mutex{}
	lastPetID int64 = 2
)

// newPetID generates a new pet ID
func newPetID() int64 {
	return atomic.AddInt64(&lastPetID, 1)
}

// addPet adds a new pet to the store
func addPet(pet Pet) {
	petsLock.Lock()
	defer petsLock.Unlock()
	pet.ID = newPetID()
	pets = append(pets, pet)
}

// removePet removes a pet from the store by ID
func removePet(id int64) {
	petsLock.Lock()
	defer petsLock.Unlock()
	var newPets []Pet
	for _, pet := range pets {
		if pet.ID != id {
			newPets = append(newPets, pet)
		}
	}
	pets = newPets
}

// petByID finds a pet by ID
func petByID(id int64) (*Pet, error) {
	for _, pet := range pets {
		if pet.ID == id {
			return &pet, nil
		}
	}
	return nil, errors.NotFound("not found: pet %d", id)
}

