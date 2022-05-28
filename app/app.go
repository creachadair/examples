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
	log.Println("App started")

	sin := bufio.NewScanner(os.Stdin)

	// Accept a "request" from the node and send a "response".
	// Then shut down the "service" and exit.
	for sin.Scan() {
		log.Printf("[app] ⇐ %s\n", sin.Text())
		fmt.Printf("OK %s\n", sin.Text()) // stdout goes back to caller
		log.Printf("[app] ⇒ OK %s\n", sin.Text())
		break
	}
	if err := sin.Err(); err != nil && err != io.EOF {
		log.Fatalf("[app] Scan failed: %v", err)
	}

	log.Printf("[app] close output: %v", os.Stdout.Close())
	log.Println("[app] exit OK")
}
