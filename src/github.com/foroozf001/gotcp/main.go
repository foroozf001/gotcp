package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

const portNumber = ":8080"

func main() {
	http.HandleFunc("/healthz", healthz)
	http.HandleFunc("/reportz", reportz)
	fmt.Printf("Starting Portcat server on port %s\n", portNumber)
	_ = http.ListenAndServe(portNumber, nil)
}

func extract(w http.ResponseWriter, r *http.Request, m map[string]string) {
	for key := range m {
		keys, ok := r.URL.Query()[key]
		if !ok || len(keys[0]) < 1 {
			log.Println("missing " + key)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("400 Bad Request"))
			return
		}
		m[key] = string(keys[0])
	}
}

func healthz(w http.ResponseWriter, r *http.Request) {
	tp := map[string]string{"host": "default", "port": "-1"}
	extract(w, r, tp)
	host := tp["host"]
	if len(host) < 1 {
		panic("empty host")
	}
	port, err := strconv.ParseInt(tp["port"], 10, 64)
	if err != nil {
		panic(err)
	}
	ports := scan(host, port, port)
	if len(ports) < 1 {
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte("502 Bad Gateway"))
	} else {
		w.Write([]byte("200 OK"))
	}
}

func reportz(w http.ResponseWriter, r *http.Request) {
	rp := map[string]string{"host": "default", "startPort": "-1", "endPort": "-2"}
	extract(w, r, rp)
	host := rp["host"]
	if len(host) < 1 {
		panic("empty host")
	}
	startPort, err := strconv.ParseInt(rp["startPort"], 10, 64)
	if err != nil {
		panic(err)
	}
	endPort, err := strconv.ParseInt(rp["endPort"], 10, 64)
	if err != nil {
		panic(err)
	}
	ports := scan(host, startPort, endPort)
	resp := "gotcp scanner (https://github.com/foroozf001/gotcp)\n"
	resp += "Target host: " + host + "\n"
	resp += "Port range: " + fmt.Sprint(startPort) + "-" + fmt.Sprint(endPort) + "\n"
	resp += "-"
	for _, port := range ports {
		s := fmt.Sprintf("\ntcp/%d\t%s", port, "UP")
		resp += s
	}
	if len(ports) < 1 {
		resp += "\nNone"
	}
	fmt.Fprintln(w, resp)
}
