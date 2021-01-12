package main

import (
	"log"
	"os"
)

func main() {
	for _, env := range os.Environ() {
		log.Printf("%v", env)
	}
}
