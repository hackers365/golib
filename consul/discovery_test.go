package consul

import (
	"fmt"
	_"github.com/gin-gonic/gin"
	"testing"
	"time"
	"encoding/json"
	_"github.com/hashicorp/consul/api"
)

// go test -v service_register_test.go conf_center.go service_discovery.go service_instance.go

func TestServiceDiscovery(t *testing.T) {
	instance, err := NewRegistry("10.0.0.4:8500", "")
	if err != nil {
		fmt.Println("NewRegistry error")
		return
	}

	//watch
	handler := func(i []ServiceInstance) {
		b, _ := json.Marshal(i)
		fmt.Println(string(b))
	}
	serviceList, err := instance.GetAndWatch("UploadService", handler)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("getService: %v", serviceList)
	err = instance.RegisterWithTtl("UploadService", "10.0.0.4", 80, "10s", 100, "30s")
	if err != nil {
		fmt.Println(err)
		return
	}

	time.Sleep(100 * time.Second)
}
