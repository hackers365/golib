package etcd

import (
  "testing"
  "time"
  "encoding/json"
  "fmt"
  clientv3 "go.etcd.io/etcd/client/v3"
)

type endpoint struct {
  Id string `json:"id"`
  Address []string `json:"address"`
  Port int `json:"port"`
  MetaData interface{} `json:"metadata"`
}

func TestGetInstance(t *testing.T) {
	prefix := "/onething/app_manager"
	client, err := clientv3.New(clientv3.Config{
      Endpoints:   []string{"192.168.208.214:2379"},
      DialTimeout: 5 * time.Second,
  })

  if err != nil {
  	panic(err)
  }

  cb := func(s string) string {
  	fmt.Println(s)
  	var info endpoint
  	err := json.Unmarshal([]byte(s), &info)
  	if err != nil {
  		fmt.Println(err)
  		return ""
  	}
  	if len(info.Address) > 0 {
  		return info.Address[0]
  	}
  	return ""
  }
  dis, err := NewDiscovery(client, ParseCb(cb))
  if err != nil {
  	panic(err)
  }

  for {
	  url, err := dis.GetInstance(prefix)
	  if err != nil {
	  	_ = err
	  }
	  _ = url
	  //time.Sleep(1 * time.Millisecond)
  }

}
