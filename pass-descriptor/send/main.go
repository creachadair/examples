package main

import (
	"flag"
	"log"
	"net"
	"os"

	"golang.org/x/sys/unix"
)

var (
	filePath = flag.String("file", "", "File path to send")
	sockPath = flag.String("socket", "", "Socket path to send to")
)

func main() {
	flag.Parse()
	switch {
	case *filePath == "":
		log.Fatal("You must provide a -file path to send")
	case *sockPath == "":
		log.Fatal("You must provide a -socket path to send to")
	}

	f, err := os.Open(*filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	ffd := f.Fd()
	log.Printf("File descriptor for %q: %v", *filePath, ffd)

	conn, err := net.Dial("unix", *sockPath)
	if err != nil {
		log.Fatalf("Dial: %v", err)
	}
	defer conn.Close()

	cfd, err := connFD(conn.(*net.UnixConn))
	if err != nil {
		log.Fatalf("Get connection descriptor: %v", err)
	}

	log.Printf("File descriptor for connection: %v", cfd)

	rights := unix.UnixRights(int(ffd))
	if err := unix.Sendmsg(int(cfd), nil, rights, nil, 0); err != nil {
		log.Fatalf("Send descriptor: %v", err)
	}
}

func connFD(uc *net.UnixConn) (uintptr, error) {
	rc, err := uc.SyscallConn()
	if err != nil {
		return 0, err
	}
	var ucfd uintptr
	if err := rc.Control(func(fd uintptr) {
		ucfd = fd
	}); err != nil {
		return 0, err
	}
	return ucfd, nil
}
