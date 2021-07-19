package consul

import (
	"testing"
  "time"
  //"github.com/hashicorp/consul/api"
)


func TestConsul(t *testing.T) {
  host := "192.168.208.214"
  port := 8500
  registryDiscoveryClient, _ := NewConsulServiceRegistry(host, port, "")

  ip := "192.168.208.214"

  isSecure := false
  portApi := 8095
  

  /*check := new(api.AgentServiceCheck)
  check.TTL = "30s"
  check.DeregisterCriticalServiceAfter = "1m" // 故障检查失败30s后 consul自动将注册服务删除
  */
  check := nil

  /*schema := "http"
  if isSecure {
    schema = "https"
  }
  check.HTTP = fmt.Sprintf("%s://%s:%d/actuator/health", schema, ip, portApi)
  check.Timeout = "5s"
  check.Interval = "5s"*/

  serviceInstanceInfo, _ := NewDefaultServiceInstance("service-ttl-check", ip, portApi,
    isSecure, map[string]string{"user": "zyn"}, "", nil)

  for i := 0; i < 10; i++ {
    registryDiscoveryClient.Register(serviceInstanceInfo)
    time.Sleep(15 * time.Second)
  }
}
