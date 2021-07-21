package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	hostname, err := os.Hostname()
	panicOnErr(err)

	serverURL := os.Getenv("SERVER_URL")

	for {
		pingServer(serverURL, hostname)
		time.Sleep(5 * time.Second)
	}
}

func pingServer(serverURL, clientHostname string) {
	l := logger{serverURL: serverURL, clientHostname: clientHostname}
	req, err := http.NewRequest(http.MethodPost, serverURL,
		bytes.NewBufferString("ping from "+clientHostname))
	if err != nil {
		l.Printf("error creating request: %v", err)
		return
	}
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		l.Printf("error invoking server: %v", err)
		return
	}
	msg, err := io.ReadAll(resp.Body)
	if err != nil {
		l.Printf("error reading response: %v", err)
		return
	}

	l.Printf("server responded: %v", string(msg))
}

type logger struct {
	serverURL      string
	clientHostname string
}

func (l *logger) Printf(format string, v ...interface{}) {
	log.Printf("hostname=%v server=%v %s", l.clientHostname, l.serverURL,
		fmt.Sprintf(format, v...))
}

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}
