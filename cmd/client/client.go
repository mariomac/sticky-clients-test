package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	hostname, err := os.Hostname()
	panicOnErr(err)

	serverURL, ok := os.LookupEnv("SERVER_URL")
	if !ok {
		serverURL = "http://" + fetchNodeHost() + ":30080"
	}
	log.Printf("pinging to server: %v", serverURL)

	for {
		pingServer(serverURL, hostname)
		time.Sleep(5 * time.Second)
	}
}

func fetchNodeHost() string {
	currentPodName, err := os.Hostname()
	panicOnErr(err)
	cfg, err := rest.InClusterConfig()
	panicOnErr(err)
	clientSet, err := kubernetes.NewForConfig(cfg)
	panicOnErr(err)
	currentPod, err := clientSet.CoreV1().Pods("default").
		Get(context.Background(), currentPodName, v1.GetOptions{})
	panicOnErr(err)
	log.Printf("pod %v is in node %v", currentPodName, currentPod.Spec.NodeName)
	currentNode, err := clientSet.CoreV1().Nodes().
		Get(context.Background(), currentPod.Spec.NodeName, v1.GetOptions{})
	panicOnErr(err)
	for _, addr := range currentNode.Status.Addresses {
		if addr.Type == "InternalIP" {
			return addr.Address
		}
	}
	for _, addr := range currentNode.Status.Addresses {
		if addr.Type == "InternalDNS" {
			return addr.Address
		}
	}
	for _, addr := range currentNode.Status.Addresses {
		return addr.Address
	}
	panic("can't find any address")
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
