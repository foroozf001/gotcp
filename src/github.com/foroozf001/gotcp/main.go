package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

const SERVER_PORT = 8080
const URL_HOST = "host"
const URL_PORT = "port"

func main() {
	http.HandleFunc("/health", Health)
	http.HandleFunc("/report", Report)
	fmt.Printf("Starting GOTCP server on port %d\n", SERVER_PORT)
	_ = http.ListenAndServe(":"+fmt.Sprint(SERVER_PORT), nil)
}

func Health(w http.ResponseWriter, r *http.Request) {
	defer TimeTrack(time.Now())

	var params []string
	var err error
	var port int64
	var ports []int64

	params = GetUrlParameters(r, URL_HOST, URL_PORT)
	var scanner = Scanner{
		Host:     params[0],
		Protocol: "tcp",
		Timeout:  200 * time.Millisecond,
		Ulimit:   1024,
	}

	if !scanner.HasValidHost() {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(http.StatusText(http.StatusBadRequest)))
		return
	}

	port, err = strconv.ParseInt(params[1], 10, 64)
	if err != nil {
		panic(err)
	}

	ports = scanner.Scan(port, port)

	if len(ports) < 1 {
		w.WriteHeader(http.StatusRequestTimeout)
		w.Write([]byte(http.StatusText(http.StatusRequestTimeout)))
		return
	} else {
		w.Write([]byte(http.StatusText(http.StatusOK)))
	}
}

func Report(w http.ResponseWriter, r *http.Request) {
	defer TimeTrack(time.Now())

	var params []string
	var resp string
	var ports []int64

	params = GetUrlParameters(r, URL_HOST, URL_PORT)
	var scanner = Scanner{
		Host:     params[0],
		Protocol: "tcp",
		Timeout:  200 * time.Millisecond,
		Ulimit:   1024,
	}

	if !scanner.HasValidHost() {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(http.StatusText(http.StatusBadRequest)))
		return
	}

	ports = scanner.Scan(1, 65535)

	for _, port := range ports {
		s := fmt.Sprintf("%d/tcp\n", port)
		resp += s
	}
	if len(ports) < 1 {
		resp += "None\n"
	}

	w.Write([]byte(resp))
}

func GetUrlParameters(r *http.Request, s ...string) []string {
	var keys = []string{}
	var values = []string{}
	keys = append(keys, s...)
	for _, key := range keys {
		keys, ok := r.URL.Query()[key]
		if !ok {
			log.Println("invalid " + key)
		}
		values = append(values, string(keys[0]))
	}
	return values
}

func TimeTrack(start time.Time) {
	elapsed := time.Since(start)
	log.Printf("elapsed time: %s", elapsed)
}
