package main

import (
	"log"
	"net/http"

	"github.com/neogan74/go-pet-store/api"
)

func main() {
	// just get api handler
	petstoreAPI, err := api.NewPetstore()
	if err != nil {
		log.Fatalln(err)
	}
	// serve the api at /api
	log.Println("Serving petstore api on http://127.0.0.1:8344/api/")
	_ = http.ListenAndServe(":8344", petstoreAPI)
}
