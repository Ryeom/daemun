package internal

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func GetLocalIP() string {
	addr, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addr {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func IsRunning(host, port string) bool {
	address := net.JoinHostPort(host, port)
	conn, err := net.DialTimeout("tcp", address, 3*time.Second)
	if err != nil {
		fmt.Println(err)
	} else if conn != nil {
		defer conn.Close()
		fmt.Printf("%s:%s is opened \n", host, port)
		return true
	}
	return false
}

func SendHttpRequest(method, uri string, headers string, param string, bodyObj interface{}) (string, string) {
	var err error
	var req *http.Request
	var header map[string][]string
	if headers == "" || method == "" || uri == "" || bodyObj == nil {
		return "", ""
	}
	err = json.Unmarshal([]byte(headers), &header)
	if err != nil {
		fmt.Println("header json.Unmarshal Err", method, uri, header, param, bodyObj)
		return "", ""
	}
	if header["Content-Type"] == nil {
		return "", ""
	}
	switch header["Content-Type"][0] {
	case "application/x-www-form-urlencoded":
		var body *bytes.Buffer
		body = bytes.NewBuffer([]byte(bodyObj.(string)))
		req, err = http.NewRequest(method, uri, body)
	case "application/json":
		str := bodyObj.(string)
		var body *strings.Reader
		body = strings.NewReader(str)
		req, err = http.NewRequest(method, uri, body)
	default:
		return "", ""
	}

	if err != nil {
		return "", ""
	}
	// header
	for key, val := range header {
		req.Header.Add(key, val[0])
	}
	if method == "GET" || param != "" {
		q := req.URL.Query()
		var queryParam map[string]string
		err = json.Unmarshal([]byte(param), &queryParam)
		if err != nil {
			fmt.Println("파람 마샬 에러", err)
			return "", ""
		}
		for k, v := range queryParam {
			value, _ := url.QueryUnescape(v)
			q.Add(k, value)
		}
		// query
		req.URL.RawQuery = q.Encode()
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", ""
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		str := string(respBody)
		fmt.Println("respBody ioutil.ReadAll Err", err, str)
	}
	defer resp.Body.Close()
	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(bytes.NewReader(respBody))
		if err != nil {
			// header: gzip , body: plain
			reader = resp.Body
		}
		bytes, _ := ioutil.ReadAll(reader)
		respBody = bytes
		defer reader.Close()
	default:
		reader = resp.Body
	}
	return resp.Status, string(respBody)
}
