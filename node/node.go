package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

var (
	inFD  = flag.Int("in", -1, "Input file descriptor")
	outFD = flag.Int("out", -1, "Output file descriptor")
)

func main() {
	flag.Parse()
	if *inFD < 0 {
		log.Fatal("No -in descriptor provided")
	}
	if *outFD < 0 {
		log.Fatal("No -out descriptor provided")
	}
	log.Println("Node started")

	// Create files for the descriptors passed by the parent.
	in := os.NewFile(uintptr(*inFD), "input")
	sin := bufio.NewScanner(in)
	out := os.NewFile(uintptr(*outFD), "output")

	// Send a "request" to the application and wait for its "response".
	fmt.Fprintln(out, "Hello!")
	log.Printf("[node] ⇒ Hello!\n")
	for sin.Scan() {
		log.Printf("[node] ⇐ %s\n", sin.Text())
	}
	if err := sin.Err(); err != nil && err != io.EOF {
		log.Fatalf("[node] Scan failed: %v", err)
	}
	log.Printf("[host] close output: %v", out.Close())
	in.Close()
	log.Println("[host] exit OK")
}
