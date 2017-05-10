package main

import (
	"fmt"
	"log"
	"os"
)

const help = `Usage: grpc gen protos...
`

func main() {
	if err := gen(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}
