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

	// Create a pair of connected pipes:
	//
	//   apr -- reads by the app
	//   apw -- writes from the app
	//   npr -- reads by the node
	//   apw -- writes by the node
	//
	npr, apw, err := os.Pipe()
	if err != nil {
		log.Fatalf("Pipe: %v", err)
	}
	apr, npw, err := os.Pipe()
	if err != nil {
		log.Fatalf("Pipe: %v", err)
	}

	// Start the node process and pass in its pipe ends.
	// The flags tell it which descriptors to expect.
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

	// Clean up the pipe ends on the parent process.
	npr.Close()
	npw.Close()

	// Start the app process and pass in its pipe ends.
	// The flags tell it which descriptors to expect.
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

	// Clean up the pipe ends on the parent process.
	apr.Close()
	apw.Close()

	// Wait for both children to exit.
	if _, err := nproc.Wait(); err != nil {
		log.Printf("WARNING: node process exit: %v", err)
	}
	if _, err := aproc.Wait(); err != nil {
		log.Printf("WARNING: app process exit: %v", err)
	}
}
