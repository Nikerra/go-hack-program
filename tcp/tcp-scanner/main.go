package main

import (
	"fmt"
	"net"
	"os"
	"sort"
	"strconv"
	"time"
)

func worker(ports, results chan int, address string) {
	for p := range ports {
		addressForScan := fmt.Sprintf(address+":%d", p)
		conn, err := net.DialTimeout("tcp", addressForScan, 1*time.Second)
		if err != nil {
			results <- 0
			continue
		}
		conn.Close()
		results <- p
	}
}

func main() {
	args := os.Args
	if len(args) != 3 {
		fmt.Println("Usage: tcp-scanner-proxy <port> or <address for scan>")
		os.Exit(1)
	}

	portForScan, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Printf("Invalid port number: %s\n", portForScan)
		os.Exit(1)
	}

	if portForScan <= 0 || portForScan > 65535 {
		fmt.Printf("Port must be between 1 and 65535, got: %d\n", portForScan)
		os.Exit(1)
	}
	addressForScan := args[2]

	ports := make(chan int, 100)
	results := make(chan int, 1000)
	var openPorts []int

	numWorkers := 200
	for i := 0; i < numWorkers; i++ {
		go worker(ports, results, addressForScan)
	}

	go func() {
		for i := 1; i <= portForScan; i++ {
			ports <- i
		}
		close(ports)
	}()

	for i := 0; i < portForScan; i++ {
		port := <-results
		if port != 0 {
			openPorts = append(openPorts, port)
		}
	}

	close(results)
	sort.Ints(openPorts)
	for _, port := range openPorts {
		fmt.Printf("open port: %d\n", port)
	}
}
