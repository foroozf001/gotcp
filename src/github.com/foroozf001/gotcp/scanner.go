package main

import (
	"fmt"
	"log"
	"net"
	"sort"
	"time"
)

type Scanner struct {
	Host     string
	Protocol string
	Timeout  time.Duration
	Ulimit   int
}

func (s *Scanner) HasValidHost() bool {
	return len(s.Host) > 0
}

func (s *Scanner) Scan(first int64, last int64) []int64 {
	ports := make(chan int64, s.Ulimit)
	results := make(chan int64)
	var open []int64
	for i := 0; i < cap(ports); i++ {
		go Worker(s.Protocol, s.Host, s.Timeout, ports, results)
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

func Worker(protocol string, host string, timeout time.Duration, ports, results chan int64) {
	for p := range ports {
		address := fmt.Sprintf("%s:%d", host, p)
		conn, err := net.DialTimeout(protocol, address, timeout)
		if err != nil {
			results <- 0
			continue
		}
		log.Printf("dial %s %s:%d: connect: connection succeeded\n", protocol, host, p)
		conn.Close()
		results <- p
	}
}
