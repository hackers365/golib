package consul

import (
	"fmt"
	_"github.com/gin-gonic/gin"
	"testing"
	"time"
	"encoding/json"
	"github.com/hashicorp/consul/api"
)

// go test -v service_register_test.go conf_center.go service_discovery.go service_instance.go

func TestConsulServiceRegistry(t *testing.T) {
	host := "192.168.208.214"
	port := 8500
	registryDiscoveryClient, _ := NewConsulServiceRegistry(host, port, "")

	/*check := new(api.AgentServiceCheck)
	schema := "http"
	if isSecure {
		schema = "https"
	}
	check.HTTP = fmt.Sprintf("%s://%s:%d/actuator/health", schema, ip, portApi)
	check.Timeout = "5s"
	check.Interval = "5s"
	check.TTL = "30s"
	check.DeregisterCriticalServiceAfter = "20s" // 故障检查失败30s后 consul自动将注册服务删除
	*/
	serviceInstanceInfo, _ := NewDefaultServiceInstance("service-ttl-checks", "192.168.208.214", 80)

	//url := "http://192.168.208.214:9991"
	//serviceInstanceInfo.SetHttpCheck(url, "20s", "5s")

	serviceInstanceInfo.SetTtlCheck("30s")
	serviceInstanceInfo.SetDeregisterAfter("300s")

  handler := func(lastIndex uint64, result interface{}) {
      services := result.([]*api.ServiceEntry)
      str, _ := json.Marshal(services)
      fmt.Println("notify:" + string(str))
  }


	go registryDiscoveryClient.WatchPlan("service", "service-ttl-checks", handler)

	/*serviceInstanceInfo, _ := NewDefaultServiceInstance("go-user-server", ip, portApi,
		isSecure, map[string]string{"user": "zyn"}, "", check)*/

	err := registryDiscoveryClient.Register(serviceInstanceInfo)
	if err != nil {
		fmt.Println(err)
		return
	}

	checkId := serviceInstanceInfo.GetCheck().CheckID
	for {
		srvList, _ := registryDiscoveryClient.GetInstances("service-ttl-checks")
		fmt.Println(srvList)

		err := registryDiscoveryClient.TtlKeepalive(checkId, "pass")
		fmt.Println(err)
		time.Sleep(10 * time.Second)
	}

	/*r := gin.Default()
	// 健康检测接口，只要是 200 就认为成功了
	r.GET("/actuator/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	err = r.Run(":8095")
	if err != nil {
		registryDiscoveryClient.Deregister()
	}*/
}

func TestConsulServiceDiscovery2(t *testing.T) {
	host := "192.168.208.214"
	port := 8500
	token := ""
	registryDiscoveryClient, err := NewConsulServiceRegistry(host, port, token)
	if err != nil {
		panic(err)
	}

	t.Log(registryDiscoveryClient.GetServices())

	t.Log(registryDiscoveryClient.GetInstances("service-ttl-checks"))
}
