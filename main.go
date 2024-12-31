package main

import (
	"log"

	"github.com/pangmaowang/network-toolbox/cmd"
)

func main() {

	err := cmd.Execute()
	if err != nil {
		log.Printf("Error executing command: %v", err)
		log.Fatal(err)
	}
}
