package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/creachadair/chirp"
	"github.com/creachadair/chirp/channel"
)

func main() {
	log.SetPrefix("[sub]  ")
	log.SetFlags(log.Ltime | log.Lmicroseconds)
	flag.Parse()
	log.Printf("subprocess here, pid=%d, args=%q", os.Getpid(), flag.Args())

	pstr, ok := os.LookupEnv("HOST_FD")
	if !ok {
		log.Fatal("Can't find HOST_FD in the environment")
	}
	fd, err := strconv.Atoi(pstr)
	if err != nil {
		log.Fatalf("Invalid HOST_FD: %v", err)
	}
	log.Printf("received HOST_FD %d", fd)

	p := chirp.NewPeer().Start(channel.ConnectPipe(
		os.NewFile(uintptr(fd), "host-in"),
		os.NewFile(uintptr(fd+1), "host-out"),
	))
	defer p.Stop()

	rsp, err := p.Call(context.Background(), "ping", fmt.Appendf(nil, "hello %v", flag.Args()))
	if err != nil {
		log.Printf("WARNING: Host ping failed: %v", err)
	} else {
		log.Printf("host ping responded: %q", rsp.Data)
	}

	log.Print("done, exiting...")
}
