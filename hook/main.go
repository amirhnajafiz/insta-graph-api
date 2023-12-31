package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

type request struct {
	Host string `json:"host"`
	Data int    `json:"data"`
}

func callback(host string, data []byte) error {
	req, err := http.NewRequest(http.MethodPost, host, bytes.NewReader(data))
	if err != nil {
		return err
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	log.Println(fmt.Sprintf("host=%s, status code=%d %s", host, resp.StatusCode, resp.Status))

	return nil
}

func process(r *request) {
	for i := 0; i < r.Data; i++ {
		time.Sleep(1 * time.Second)
	}

	if err := callback(r.Host, []byte(fmt.Sprintf("server busy time: %ds", r.Data))); err != nil {
		log.Println(fmt.Errorf("error in callback host: %s\n\t%w", r.Host, err))
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	req := new(request)

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)

		log.Println(fmt.Errorf("failed to parse request: %w", err))

		return
	}

	go process(req)

	w.WriteHeader(http.StatusOK)
}

func main() {
	var PortFlag = flag.Int("port", 8080, "http port of hook service")

	flag.Parse()

	mux := http.NewServeMux()

	mux.HandleFunc("/", handler)

	log.Println(fmt.Sprintf("hook server started on %d ...", *PortFlag))

	if err := http.ListenAndServe(fmt.Sprintf(":%d", *PortFlag), mux); err != nil {
		panic(err)
	}
}
