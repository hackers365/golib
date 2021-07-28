package consul

import (
	"fmt"
	_"github.com/gin-gonic/gin"
	"testing"
	"time"
	_"github.com/hashicorp/consul/api"
)

// go test -v service_register_test.go conf_center.go service_discovery.go service_instance.go

func TestServiceDiscovery(t *testing.T) {
	discovery, err := NewDiscovery("192.168.208.214:8500", "")
	if err != nil {
		fmt.Println("NewDiscovery error")
		return
	}

	/*err = discovery.WatchAndSave("UploadService")
	if err != nil {
		fmt.Println(err)
		return
	}*/
	
	go register()
	for i := 0; i < 10; i++ {
		go get(discovery)
	}

	

	time.Sleep(2 * time.Second)
}

func get(discovery Discovery) {
	for i := 0; i < 10; i++ {
		instance, err := discovery.GetInstance("UploadService")
		if err != nil {
			fmt.Println(err)
			return
		}
		_ = instance
		fmt.Println(instance)
		//time.Sleep(2 * time.Second)
	}
}

func register() {
	instance, err := NewRegistry("192.168.208.214:8500", "")
	if err != nil {
		fmt.Println("NewRegistry error")
		return
	}

	instance.RegisterWithTtl("UploadService", "192.168.208.209", 1000, "30s", 10, "20s")
	instance.RegisterWithTtl("UploadService", "192.168.208.209", 1001, "30s", 10, "20s")
	instance.RegisterWithTtl("UploadService", "192.168.208.209", 1002, "30s", 10, "20s")
}
