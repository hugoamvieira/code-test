package main

import (
	"log"

	"github.com/hugoamvieira/code-test/server/api"
)

func main() {
	addr := ":5000"
	a := api.New(addr)

	log.Printf("Starting API on %v\n", addr)
	log.Fatalln(a.Start())
}
