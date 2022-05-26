package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
)

var (
	nodePath = flag.String("node", "", "Node binary path")
	appPath  = flag.String("app", "", "Application binary path")
)

func main() {
	flag.Parse()
	if *nodePath == "" {
		log.Fatal("No -node path is set")
	}
	if *appPath == "" {
		log.Fatal("No -app path is set")
	}

	npr, apw, err := os.Pipe()
	if err != nil {
		log.Fatalf("Pipe: %v", err)
	}
	apr, npw, err := os.Pipe()
	if err != nil {
		log.Fatalf("Pipe: %v", err)
	}

	nproc, err := os.StartProcess(*nodePath, []string{
		filepath.Base(*nodePath),
		"-in", "3",
		"-out", "4",
	}, &os.ProcAttr{
		Files: []*os.File{nil, nil, os.Stderr, npr, npw},
	})
	if err != nil {
		log.Fatalf("Start node: %v", err)
	}
	log.Printf("Started node at pid=%d", nproc.Pid)
	npr.Close()
	npw.Close() // parent

	aproc, err := os.StartProcess(*appPath, []string{
		filepath.Base(*appPath),
		"-in", "3",
		"-out", "4",
	}, &os.ProcAttr{
		Files: []*os.File{nil, nil, os.Stderr, apr, apw},
	})
	if err != nil {
		log.Fatalf("Start app: %v", err)
	}
	log.Printf("Started app at pid=%d", aproc.Pid)
	apr.Close()
	apw.Close() // parent

	if _, err := nproc.Wait(); err != nil {
		log.Printf("WARNING: node process exit: %v", err)
	}
	if _, err := aproc.Wait(); err != nil {
		log.Printf("WARNING: app process exit: %v", err)
	}
}
