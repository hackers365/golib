package http_client

import (
	"fmt"
	_ "net/http"
	"testing"
)

func TestGet(t *testing.T) {
	url := "http://192.168.208.214"
	params := map[string]string{
		"k1": "v1",
		"k2": "v2",
	}
	header := map[string]string{
		"Host": "sjbapi-miner_manage_system.onethingpcs.com",
	}
	httpClient := NewHttpClient()
	status, retData, err := httpClient.Get(url, params, header, 10)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("status: %v, retData: %s\n", status, string(retData))

	status, retData, err = httpClient.Post(url, params, header, 10)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("status: %v, retData: %s\n", status, string(retData))

}
