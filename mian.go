package main

import (
	"check-point/api/src/db"
	"check-point/api/src/routers"
	"fmt"
	"log"
	"net/http"
)

func main() {
	db.Connection()
	r := routers.Load()
	fmt.Printf("Executando na porta %d", db.Door)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", db.Door), r))

}
