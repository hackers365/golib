package http_client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

type HttpClient interface {
	Get(url string, params map[string]string, header map[string]string, timeout int) (int, []byte, error)
	Post(requestUrl string, params map[string]interface{}, header map[string]string, timeout int) (int, []byte, error)
}

type httpClient struct {
	instance *http.Client
}

var client *httpClient
var transport *http.Transport

func GetHttpClient() HttpClient {
	return client
}

func NewHttpClient() HttpClient {
	//init http client
	iClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	//init http transport
	transport = &http.Transport{
		//Proxy: http.ProxyURL(torProxyUrl),
		DialContext: (&net.Dialer{
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		//MaxConnsPerHost: 1,
		//DisableKeepAlives: true,
	}

	iClient.Transport = transport
	client = &httpClient{instance: iClient}
	return client
}

func (h *httpClient) Get(url string, params map[string]string, header map[string]string, timeout int) (int, []byte, error) {
	ctx, cancel := context.WithCancel(context.TODO())
	time.AfterFunc(time.Duration(timeout)*time.Second, func() {
		cancel()
	})

	//实例化req
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, nil, err
	}
	req = req.WithContext(ctx)

	//添加params
	query := req.URL.Query()
	for k, v := range params {
		query.Add(k, v)
	}
	req.URL.RawQuery = query.Encode()

	//添加header
	for k, v := range header {
		req.Header.Add(k, v)
	}

	if host, ok := header["Host"]; ok {
		req.Host = host
	}

	return h.do(req)
}

func (h *httpClient) Post(requestUrl string, params map[string]interface{}, header map[string]string, timeout int) (int, []byte, error) {
	ctx, cancel := context.WithCancel(context.TODO())
	time.AfterFunc(time.Duration(timeout)*time.Second, func() {
		cancel()
	})

	bytesParams, _ := json.Marshal(params)
	body := bytes.NewBuffer(bytesParams)
	//实例化req
	req, err := http.NewRequest("POST", requestUrl, body)
	if err != nil {
		return 0, nil, err
	}
	req = req.WithContext(ctx)
	//添加header
	for k, v := range header {
		req.Header.Add(k, v)
	}

	if host, ok := header["Host"]; ok {
		req.Host = host
	}
	req.Header.Set("Content-Type", "application/json")
	return h.do(req)
}

func (h *httpClient) do(req *http.Request) (int, []byte, error) {
	resp, err := h.instance.Do(req)
	if err != nil {
		return 0, nil, fmt.Errorf("http client.do():%s", err.Error())
	}

	defer resp.Body.Close()
	// check status code
	if resp.StatusCode != http.StatusOK {
		return resp.StatusCode, nil, fmt.Errorf("status code:%d", resp.StatusCode)
	}

	// read from response
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, fmt.Errorf("ReadAll():%s", err.Error())
	}
	return resp.StatusCode, bytes, nil

}
