package main

import (
	"flag"
	"io"
	"log"
	"net"
	"os"

	"golang.org/x/sys/unix"
)

var (
	sockPath = flag.String("socket", "", "Socket path to listen on")
)

func main() {
	flag.Parse()
	switch {
	case *sockPath == "":
		log.Fatal("You must provide a -socket path to listen on")
	}

	lst, err := net.Listen("unix", *sockPath)
	if err != nil {
		log.Fatalf("Listen: %v", err)
	}

	conn, err := lst.Accept()
	if err != nil {
		log.Fatalf("Accept: %v", err)
	}
	lst.Close()
	os.Remove(*sockPath)

	log.Printf("Connected: %v", conn)
	cfd, err := connFD(conn.(*net.UnixConn))
	if err != nil {
		log.Fatalf("Get connection descriptor: %v", err)
	}

	log.Printf("File descriptor for connection: %v", cfd)

	buf := make([]byte, unix.CmsgSpace(4))
	if _, _, _, _, err := unix.Recvmsg(int(cfd), nil /* oob */, buf, 0); err != nil {
		log.Fatalf("Receive failed: %v", err)
	}

	cmsg, err := unix.ParseSocketControlMessage(buf)
	if err != nil {
		log.Fatalf("Parse control message: %v", err)
	}

	fds, err := unix.ParseUnixRights(&cmsg[0])
	if err != nil {
		log.Fatalf("Parse Unix rights: %v", err)
	}

	if len(fds) == 0 {
		log.Fatal("No descriptors in control message")
	}

	ffd := fds[0]
	log.Printf("File descriptor: %v", ffd)
	f := os.NewFile(uintptr(ffd), "file")
	defer f.Close()

	nc, _ := io.Copy(os.Stdout, f)
	log.Printf("Copied %d bytes to stdout", nc)
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
