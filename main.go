package main

import (
	"flag"
	"fmt"
	"log"

	"l4d2serverquery-go/router"
	"l4d2serverquery-go/service"
)

func main() {
	port := flag.Int("port", 44316, "Port to run the server on")
	flag.Parse()

	service.InitClient()

	if err := webServer(*port); err != nil {
		log.Fatal(err)
	}
}

func webServer(port int) error {
	engine := router.Router()
	return engine.Run(fmt.Sprintf(":%d", port))
}
