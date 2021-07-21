package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	hostname, err := os.Hostname()
	panicOnErr(err)

	panicOnErr(http.ListenAndServe(":8080", http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodPost {
			body, err := io.ReadAll(req.Body)
			if err != nil {
				log.Printf("error reading request body: %v", err)
				return
			}
			log.Printf("request received in host %v: %v", hostname, string(body))
			_, _ = resp.Write([]byte(fmt.Sprintf("OK from %v", hostname)))
		}
	})))
}

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}
