package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

const (
	udpAddr = "0.0.0.0:8080"
	socketPath = "/tmp/structura.sock"
)

func main() {
	fmt.Println("=== Structura Go Gateway starting ===")

	ipcConn, err := net.Dial("unix", socketPath)
	if err != nil {
		log.Fatalf("[Go Gateway] Failed to connect to IPC socket (%s): %v\nMake sure C# Core is runnig first!", socketPath, err)
	}

	defer ipcConn.Close()
	fmt.Printf("[GO Gateway] Connected to C# Core via IPC (%s)\n", socketPath)

	addr, err := net.ResolveUDPAddr("udp", udpAddr)
	if err != nil {
		log.Fatalf("Failed to resolve UDP address: %v", err)
	}

	udpConn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatalf("Failed to listen on UDP: %v", err)
	}

	defer udpConn.Close()
	fmt.Printf("[Go Gateway] Listening for UDP packets on %s...\n", udpAddr)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		buf := make([]byte, 65535)
		for {
			n, _, err := udpConn.ReadFromUDP(buf)
			if err != nil {
				log.Printf("UDP Read Error: %n\n", err)
				continue
			}
			_, err = ipcConn.Write(buf[:n])
			if err != nil {
				log.Printf("[Go Gateway] IPC Write Error: %v\n", err)
			}
		}
	}()

	<-sigChan
	fmt.Println("\n[Go Gateway] Shutting down gracefully..")
}

