package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

const SERVER_PORT = 8080
const URL_HOST = "host"
const URL_PORT = "port"
const TCP_RANGE = 65535
const BAD_REQUEST = "400 Bad Request"
const GATEWAY_TIMEOUT = "504 Gateway Time-out"
const OK = "200 OK"

func main() {
	http.HandleFunc("/health", Health)
	http.HandleFunc("/report", Report)
	fmt.Printf("Starting GOTCP server on port %d\n", SERVER_PORT)
	_ = http.ListenAndServe(":"+fmt.Sprint(SERVER_PORT), nil)
}

func Health(w http.ResponseWriter, r *http.Request) {
	var scanner Scanner
	var err error
	params := GetUrlParameters(r, URL_HOST, URL_PORT)
	scanner.Port, err = strconv.ParseInt(params[1], 10, 64)
	if err != nil {
		panic(err)
	}
	scanner.Host = params[0]
	if !scanner.HasValidHost() || !scanner.HasValidPort() {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(BAD_REQUEST))
		return
	}
	ports := scanner.Scan(scanner.Port, scanner.Port)
	if len(ports) < 1 {
		w.WriteHeader(http.StatusGatewayTimeout)
		w.Write([]byte(GATEWAY_TIMEOUT))
		return
	} else {
		w.Write([]byte(OK))
	}
}

func Report(w http.ResponseWriter, r *http.Request) {
	var scanner Scanner
	var resp string
	params := GetUrlParameters(r, URL_HOST)
	scanner.Host = params[0]
	if !scanner.HasValidHost() {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(BAD_REQUEST))
		return
	}
	ports := scanner.Scan(1, TCP_RANGE)
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
