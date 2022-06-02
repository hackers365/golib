package http_client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

type HttpClient interface {
	Get(url string, params map[string]string, header map[string]string) (int, []byte, error)
	Post(requestUrl string, params map[string]interface{}, header map[string]string) (int, []byte, error)
}

type httpClient struct {
	instance *http.Client
}

var client *httpClient

func GetHttpClient() HttpClient {
	return client
}

func NewHttpClient(timeout int) HttpClient {
	//init http client
	if timeout == 0 {
		timeout = 10
	}
	iClient := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	//init http transport
	transport := &http.Transport{
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

func (h *httpClient) Get(url string, params map[string]string, header map[string]string) (int, []byte, error) {
	//实例化req
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, nil, err
	}
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

func (h *httpClient) Post(requestUrl string, params map[string]interface{}, header map[string]string) (int, []byte, error) {
	bytesParams, _ := json.Marshal(params)
	body := bytes.NewBuffer(bytesParams)
	//实例化req
	req, err := http.NewRequest("POST", requestUrl, body)
	if err != nil {
		return 0, nil, err
	}
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
