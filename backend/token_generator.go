package main

import (
	"fmt"
	"log"
	"nutrition-platform/middleware"
)

func main() {
	token, err := middleware.GenerateToken("1", "test@example.com", "user", false)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(token)
}
