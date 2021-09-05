package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"sort"
)

type Scanner struct {
	Host string
	Port int64
}

const KEY_HOST = "host"
const KEY_PORT = "port"
const TCP_MAX_RANGE = 65535

func (s *Scanner) Scan(first int64, last int64) []int64 {
	ulimit := 1024
	ports := make(chan int64, ulimit)
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

func GetUrlParameters(r *http.Request, s ...string) []string {
	keys := []string{}
	keys = append(keys, s...)
	values := []string{}
	for _, key := range keys {
		keys, ok := r.URL.Query()[key]
		if !ok {
			log.Println("invalid " + key)
		}
		values = append(values, string(keys[0]))
	}
	fmt.Println(values)
	return values
}
