package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	flag.Parse()
	log.Println("Node started")

	sin := bufio.NewScanner(os.Stdin)

	// Send a "request" to the application and wait for its "response".
	fmt.Println("Hello!")
	log.Printf("[node] ⇒ Hello!\n")
	for sin.Scan() {
		log.Printf("[node] ⇐ %s\n", sin.Text())
	}
	if err := sin.Err(); err != nil && err != io.EOF {
		log.Fatalf("[node] Scan failed: %v", err)
	}
	log.Printf("[host] close output: %v", os.Stdout.Close())
	log.Println("[host] exit OK")
}
