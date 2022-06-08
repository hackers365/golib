package http_client

import (
	"fmt"
	_ "net/http"
	"testing"
)

func TestGet(t *testing.T) {
	url := "http://192.168.123.12"
	params := map[string]string{
		"k1": "v1",
		"k2": "v2",
	}
	header := map[string]string{
		"Host": "abc.com",
	}
	httpClient := NewHttpClient()
	status, retData, err := httpClient.Get(url, params, header, 3)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("status: %v, retData: %s\n", status, string(retData))

	params2 := map[string]interface{}{
		"k1": "v1",
		"k2": "v2",
	}
	status, retData, err = httpClient.Post(url, params2, header, 3)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("status: %v, retData: %s\n", status, string(retData))

}
