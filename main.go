package main

import (
	config "check-point/src/db"
	"check-point/src/router"
	"fmt"
	"log"
	"net/http"
)

func main() {
	config.Connection()
	r := router.Load()
	fmt.Printf("Executando na porta %d", config.APIConfigInfo.APIPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.APIConfigInfo.APIPort), r))
}
