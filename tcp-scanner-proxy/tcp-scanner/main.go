package main

import (
	"fmt"
	"net"
	"os"
	"sort"
	"strconv"
)

func worker(ports, results chan int) {
	for p := range ports {
		adress := fmt.Sprintf("scanme.nmap.org:%d", p)
		conn, err := net.Dial("tcp", adress)
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
	if len(args) != 2 {
		fmt.Println("Usage: tcp-scanner-proxy <port>")
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

	ports := make(chan int, 100)
	results := make(chan int, 1000)
	var openPorts []int

	numWorkers := 200
	for i := 0; i < numWorkers; i++ {
		go worker(ports, results)
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
