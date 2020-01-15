package main

import (
	"debug_tools/pb/echo_server"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var port int

func handler(w http.ResponseWriter, r *http.Request) {

	headers := make(map[string]string)
	for key, values := range r.Header {
		headers[key] = strings.Join(values, ",")
	}

	req := echo_server.HttpRequestDump{
		Url:     r.URL.String(),
		Headers: headers,
		Method:  r.Method,
	}
	if r.Method == "POST" {
		defer r.Body.Close()
		postData, _ := ioutil.ReadAll(r.Body)
		req.PostData = string(postData)
	}
	bytes, err := json.Marshal(req)
	if err != nil {
		_, _ = w.Write([]byte("dump req failed"))
	}
	_, _ = w.Write([]byte(bytes))
}

func main() {
	flag.IntVar(&port, "port", 8080, "port to listen")
	flag.Parse()
	http.HandleFunc("/", handler)
	addr := fmt.Sprintf(":%d", port)
	log.Fatal(http.ListenAndServe(addr, nil))
}
