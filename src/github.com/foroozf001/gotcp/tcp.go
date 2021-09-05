package main

import (
	"fmt"
	"log"
	"net"
	"sort"
)

func scan(host string, sp int64, ep int64) []int64 {
	ulimit := 1024
	ports := make(chan int64, ulimit)
	results := make(chan int64)
	var openports []int64
	for i := 0; i < cap(ports); i++ {
		go worker(host, ports, results)
	}
	go func() {
		for i := sp; i <= ep; i++ {
			ports <- i
		}
	}()
	for i := sp; i <= ep; i++ {
		port := <-results
		if port != 0 {
			openports = append(openports, port)
		}
	}
	close(ports)
	close(results)
	sort.Slice(openports, func(i, j int) bool { return openports[i] < openports[j] }) // sort for type int64
	return openports
}

func worker(host string, ports, results chan int64) {
	for p := range ports {
		address := fmt.Sprintf("%s:%d", host, p)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			results <- 0
			log.Println(err)
			continue
		}
		log.Printf("dial tcp %s:%d: connect: connection succeeded\n", host, p)
		conn.Close()
		results <- p
	}
}
