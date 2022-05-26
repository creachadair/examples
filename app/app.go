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
	log.Println("App started")

	in := os.NewFile(uintptr(*inFD), "input")
	sin := bufio.NewScanner(in)
	out := os.NewFile(uintptr(*outFD), "output")

	for sin.Scan() {
		log.Printf("[app] ⇐ %s\n", sin.Text())
		fmt.Fprintf(out, "OK %s\n", sin.Text())
		log.Printf("[app] ⇒ OK %s\n", sin.Text())
		break
	}
	if err := sin.Err(); err != nil && err != io.EOF {
		log.Fatalf("[app] Scan failed: %v", err)
	}
	log.Printf("[app] close output: %v", out.Close())
	in.Close()
	log.Println("[app] exit OK")
}
