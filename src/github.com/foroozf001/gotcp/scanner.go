package main

import (
	"fmt"
	"log"
	"net"
	"sort"
)

type Scanner struct {
	Host string
	Port int64
}

const ULIMIT = 1024

func (s *Scanner) Scan(first int64, last int64) []int64 {
	ports := make(chan int64, ULIMIT)
	results := make(chan int64)
	var open []int64
	for i := 0; i < cap(ports); i++ {
		go Worker(s.Host, ports, results)
	}
	go func() {
		for i := first; i <= last; i++ {
			ports <- i
		}
	}()
	for i := first; i <= last; i++ {
		port := <-results
		if port != 0 {
			open = append(open, port)
		}
	}
	close(ports)
	close(results)
	sort.Slice(open, func(i, j int) bool { return open[i] < open[j] }) // sort for type int64
	return open
}

func Worker(host string, ports, results chan int64) {
	for p := range ports {
		address := fmt.Sprintf("%s:%d", host, p)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			results <- 0
			continue
		}
		log.Printf("dial tcp %s:%d: connect: connection succeeded\n", host, p)
		conn.Close()
		results <- p
	}
}
