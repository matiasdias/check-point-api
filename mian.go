package main

import (
	"check-point/api/src/db"
	"fmt"
)

func main() {
	db.Connection()
	fmt.Printf("Executando na porta 3000")

}
