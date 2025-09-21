package main

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/neogan74/go-pet-store/api"
)

const (
	port = ":5000"
	// timout for http server
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

type PetStorage struct {
	mu       sync.RWMutex
	petStore map[int64]api.Pet
}

func main() {

	r := chi.NewRouter()

	// just get api handler
	petstoreAPI, err := api.NewPetstore()
	if err != nil {     
		log.Fatalln(err)
	}
	// serve the api at /api
	log.Println("Serving petstore api on http://127.0.0.1:8344/api/")
	_ = http.ListenAndServe(":8344", petstoreAPI)
}
