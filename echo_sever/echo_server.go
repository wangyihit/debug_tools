package main

import (
	"bytes"
	"compress/gzip"
	"debug_tools/pb/echo_server"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var port int
var fakeServer bool
var jsUrl string
func GZipData(data []byte) (compressedData []byte, err error) {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)

	_, err = gz.Write(data)
	if err != nil {
		return nil, errors.New("write failed")
	}

	if err = gz.Flush(); err != nil {
		return nil, errors.New("flush failed")
	}
	if err = gz.Close(); err != nil {
		return nil, errors.New("close failed")
	}
	compressedData = b.Bytes()
	return compressedData, nil
}


func echoHandler(w http.ResponseWriter, r *http.Request) {

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


type Proxy struct {
}

func NewProxy() *Proxy { return &Proxy{} }

func (p *Proxy) ServeHTTP(wr http.ResponseWriter, r *http.Request) {
	var resp *http.Response
	var err error
	var req *http.Request
	client := &http.Client{}
	schema := "https"
	if s, ok := r.Header["schema"]; ok {
		schema = s[0]
	}
	// http://abc:110/path?quesry#hash
	url := fmt.Sprintf("%s://%s%s", schema, r.Host,  r.RequestURI)
	req, err = http.NewRequest(r.Method, url, r.Body)
	if err != nil {
		fmt.Printf("err: %s\n", err.Error())
		return
	}
	for name, value := range r.Header {
		if value != nil && name != "" {
			req.Header.Set(name, strings.Join(value, " "))
		}
	}
	resp, err = client.Do(req)
	r.Body.Close()

	// combined for GET/POST
	if err != nil {
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	// 微信公众号请求
	if r.Host == "mp.weixin.qq.com" && strings.Index(r.RequestURI, "/mp/profile_ext?action=home") == 0 {
		jsTail := fmt.Sprintf( `<script type="text/javascript" src="%s"></script>`, jsUrl)
		contentEncodingKey := "Content-Encoding"
		var reader io.ReadCloser
		switch resp.Header.Get(contentEncodingKey) {
		case "gzip":
			reader, err = gzip.NewReader(resp.Body)
			defer reader.Close()
		default:
			reader = resp.Body
		}
		var byteBuffer bytes.Buffer
		respData, _ := ioutil.ReadAll(reader)
		byteBuffer.Write(respData)
		byteBuffer.Write([]byte(jsTail))
		compressedData, _ := GZipData(byteBuffer.Bytes())

		for k, v := range resp.Header {
			if k == "Content-Length" {
				v[0] = fmt.Sprintf("%d", len(compressedData))
			}else {
				wr.Header().Set(k, strings.Join(v, ""))
			}
		}
		wr.Write(compressedData)

	} else{
		for k, v := range resp.Header {
			wr.Header().Set(k, strings.Join(v, ""))
		}
		wr.WriteHeader(resp.StatusCode)
		io.Copy(wr, resp.Body)
	}
}

func main() {
	flag.IntVar(&port, "port", 8080, "port to listen")
	flag.BoolVar(&fakeServer, "fs", true, "run fake server")
	flag.StringVar(&jsUrl, "jsurl", "http://10.0.0.7/js/jump.baidu.js", "js url to inject to wechat page")
	flag.Parse()
	addr := fmt.Sprintf(":%d", port)
	if fakeServer {
		err := http.ListenAndServe(addr, NewProxy())
		if err != nil {
			log.Fatal("ListenAndServe: ", err.Error())
		}
	}else{
		http.HandleFunc("/", echoHandler)

		log.Fatal(http.ListenAndServe(addr, nil))
	}


}
