package main

import (
	"fmt"
	"net/http"
	"strconv"
)

const server = ":8080"

func main() {
	http.HandleFunc("/health", Health)
	http.HandleFunc("/report", Report)
	fmt.Printf("Starting gotcp server on port %s\n", server)
	_ = http.ListenAndServe(server, nil)
}

func Health(w http.ResponseWriter, r *http.Request) {
	params := GetUrlParameters(r, KEY_HOST, KEY_PORT)
	var scanner Scanner
	var err error
	scanner.Host = params[0]
	if len(scanner.Host) < 1 {
		panic("empty " + KEY_HOST)
	}
	scanner.Port, err = strconv.ParseInt(params[1], 10, 64)
	if err != nil {
		panic(err)
	}
	ports := scanner.Scan(scanner.Port, scanner.Port)
	if len(ports) < 1 {
		w.WriteHeader(http.StatusRequestTimeout)
		w.Write([]byte("408 Request Timeout"))
	} else {
		w.Write([]byte("200 OK"))
	}
}

func Report(w http.ResponseWriter, r *http.Request) {
	params := GetUrlParameters(r, KEY_HOST)
	var scanner Scanner
	scanner.Host = params[0]
	if len(scanner.Host) < 1 {
		panic("empty " + KEY_HOST)
	}
	ports := scanner.Scan(1, TCP_MAX_RANGE)
	resp := "gotcp L4 scanner (https://github.com/foroozf001/gotcp)\n"
	resp += "Target host: " + scanner.Host + "\n"
	resp += "Port range: " + fmt.Sprint(1) + "-" + fmt.Sprint(TCP_MAX_RANGE) + "\n"
	resp += "-"
	for _, port := range ports {
		s := fmt.Sprintf("\n%d/tcp", port)
		resp += s
	}
	if len(ports) < 1 {
		resp += "\nNone"
	}
	w.Write([]byte(resp))
}
