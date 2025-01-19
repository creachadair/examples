package main

import (
	"context"
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/creachadair/chirp"
	"github.com/creachadair/chirp/channel"
)

func main() {
	log.SetPrefix("[sub] ")
	flag.Parse()

	pstr, ok := os.LookupEnv("HOST_FD")
	if !ok {
		log.Fatal("Can't find HOST_FD in the environment")
	}
	fd, err := strconv.Atoi(pstr)
	if err != nil {
		log.Fatalf("Invalid HOST_FD: %v", err)
	}
	log.Printf("received HOST_FD %d", fd)

	hr := os.NewFile(uintptr(fd), "host-in")
	hw := os.NewFile(uintptr(fd+1), "host-out")

	p := chirp.NewPeer().Start(channel.IO(hr, hw))
	defer p.Stop()

	log.Printf("running, args are %q", flag.Args())
	rsp, err := p.Call(context.Background(), "ping", []byte("hello"))
	if err != nil {
		log.Printf("WARNING: Host ping failed: %v", err)
	} else {
		log.Printf("host ping responded: %q", rsp.Data)
	}

	log.Print("done, exiting...")
}
