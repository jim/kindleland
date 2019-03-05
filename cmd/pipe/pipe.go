package main

import (
	"log"
	"os"
)

func main() {
	if _, _, err := os.Pipe(); err != nil {
		log.Fatal(err)
	}
}
