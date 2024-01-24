package main

import (
	"check-point/src/config"
	"check-point/src/router"
	"fmt"
	"log"
	"net/http"
)

func main() {
	config.Connection()
	r := router.Load()
	fmt.Printf("Executando na porta %d", config.Door)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Door), r))

}
