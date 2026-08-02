package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/creachadair/chirp"
	"github.com/creachadair/chirp/channel"
)

func main() {
	log.SetPrefix("[host] ")
	flag.Parse()
	if flag.NArg() == 0 {
		log.Fatal("You must provide a command and arguments")
	}
	sp, cp, err := channel.NewPipe()
	if err != nil {
		log.Fatalf("Pipe 2: %v", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	p := chirp.NewPeer().Handle("ping", handlePing).Start(sp)
	defer p.Stop()

	cmd := exec.CommandContext(ctx, flag.Arg(0), flag.Args()[1:]...)
	cmd.Env = append(os.Environ(), "HOST_FD=3")
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	cr, cw := cp.Files()
	cmd.ExtraFiles = []*os.File{cr, cw}
	if err := cmd.Start(); err != nil {
		log.Fatalf("Start subprocess: %v", err)
	}
	cp.Close()
	log.Printf("subprocess started pid=%d", cmd.Process.Pid)
	log.Printf("subprocess exited: %v", cmd.Wait())
}

func handlePing(ctx context.Context, req *chirp.Request) ([]byte, error) {
	log.Printf("received ping request (%q)", req.Data)
	return []byte("ok"), nil
}
