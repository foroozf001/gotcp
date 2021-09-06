package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

const SERVER_PORT = ":8080"
const URL_HOST = "host"
const URL_PORT = "port"
const MAX_RANGE = 65535

func main() {
	http.HandleFunc("/health", Health)
	http.HandleFunc("/report", Report)
	fmt.Printf("Starting gotcp server on port %s\n", SERVER_PORT)
	_ = http.ListenAndServe(SERVER_PORT, nil)
}

func Health(w http.ResponseWriter, r *http.Request) {
	var scanner Scanner
	var err error
	params := GetUrlParameters(r, URL_HOST, URL_PORT)
	scanner.Host = params[0]
	if len(scanner.Host) < 1 {
		panic("empty " + URL_HOST)
	}
	scanner.Port, err = strconv.ParseInt(params[1], 10, 64)
	if err != nil {
		panic(err)
	}
	ports := scanner.Scan(scanner.Port, scanner.Port)
	if len(ports) < 1 {
		w.WriteHeader(http.StatusRequestTimeout)
		w.Write([]byte("408 Request Timeout\n"))
	} else {
		w.Write([]byte("200 OK\n"))
	}
}

func Report(w http.ResponseWriter, r *http.Request) {
	var scanner Scanner
	params := GetUrlParameters(r, URL_HOST)
	scanner.Host = params[0]
	if len(scanner.Host) < 1 {
		panic("empty " + URL_HOST)
	}
	ports := scanner.Scan(1, MAX_RANGE)
	resp := "Target host: " + scanner.Host + "\n"
	resp += "Port range: " + fmt.Sprint(1) + "-" + fmt.Sprint(MAX_RANGE) + "\n"
	resp += "-"
	for _, port := range ports {
		s := fmt.Sprintf("\n%d/tcp", port)
		resp += s
	}
	if len(ports) < 1 {
		resp += "\nNone"
	}
	w.Write([]byte(resp + "\n"))
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
	return values
}
