package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var addr = "127.0.0.1:0"

func init() {
	flag.StringVar(&addr, "addr", addr, "bind address for server")
}

func main() {
	flag.Parse()

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}

	fmt.Printf("Listening on: %s\n", ln.Addr().String())

	if err := http.Serve(ln, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)

		dumpRequest(r)
	})); err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
}

// Mutex to make sure we don't copy multiple request bodies to stdout at once.
var outMu sync.Mutex

func dumpRequest(r *http.Request) {
	outMu.Lock()
	defer outMu.Unlock()

	fmt.Printf("RECEIVED REQUEST AT %s FROM %s\n", time.Now().Format(time.RFC3339), r.RemoteAddr)
	fmt.Println("METHOD: " + r.Method)
	fmt.Println("URL: " + r.URL.String())
	fmt.Println("HEADERS:")
	for k, v := range r.Header {
		fmt.Printf("%s: %s\n", k, strings.Join(v, ","))
	}

	fmt.Println("BODY:")
	io.Copy(os.Stdout, r.Body)
	r.Body.Close()
}
