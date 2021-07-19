package consul

import (
	"fmt"
	"net"

	"github.com/hashicorp/consul/api"
)

type ServiceInstance interface {

	// return The unique instance ID as registered.
	GetInstanceId() string

	// return The service ID as registered.
	GetServiceId() string

	// return The hostname of the registered service instance.
	GetHost() string

	// return The port of the registered service instance.
	GetPort() int

	// return Whether the port of the registered service instance uses HTTPS.
	IsSecure() bool

	// return The key / value pair metadata associated with the service instance.
	GetMetadata() map[string]string

	GetCheck() *api.AgentServiceCheck
}

type DefaultServiceInstance struct {
	InstanceId string
	ServiceId  string
	Host       string
	Port       int
	Secure     bool
	Metadata   map[string]string
	Check      *api.AgentServiceCheck
}

func NewDefaultServiceInstance(serviceId string, host string, port int) (*DefaultServiceInstance, error) {

	// 如果没有传入 IP 则获取一下，这个方法在多网卡的情况下，并不好用
	/*if len(host) == 0 {
		localIP, err := FindFirstNonLoopbackIP()
		if err != nil {
			return nil, err
		}
		host = localIP
	}*/

	var instanceId string
	var secure bool
	metadata := map[string]string{}
	if len(instanceId) == 0 {
		//instanceId = serviceId + "-" + strconv.FormatInt(time.Now().Unix(), 10) + "-" + strconv.Itoa(rand.Intn(9000)+1000)
		instanceId = fmt.Sprintf("%s-%s:%d", serviceId, host, port)
	}

	check := new(api.AgentServiceCheck)
	check.CheckID = fmt.Sprintf("checkid:%s", instanceId)

	return &DefaultServiceInstance{
		InstanceId: instanceId,
		ServiceId:  serviceId,
		Host:       host,
		Port:       port,
		Secure:     secure,
		Metadata:   metadata,
		Check:      check,
	}, nil
}

func (serviceInstance *DefaultServiceInstance) SetSecure(secure bool) {
	serviceInstance.Secure = secure
}

func (serviceInstance *DefaultServiceInstance) SetHttpCheck(checkUrl, timeout, interval string) {
	serviceInstance.Check.HTTP = checkUrl

	//check.HTTP = fmt.Sprintf("%s://%s:%d/actuator/health", schema, ip, portApi)
	serviceInstance.Check.Timeout = timeout
	serviceInstance.Check.Interval = interval
}

func (serviceInstance *DefaultServiceInstance) SetDeregisterAfter(deregisterAfter string) {
	serviceInstance.Check.DeregisterCriticalServiceAfter = deregisterAfter // 故障检查失败30s后 consul自动将注册服务删除
}

func (serviceInstance *DefaultServiceInstance) SetTtlCheck(ttl string) {
	serviceInstance.Check.TTL = ttl
	serviceInstance.Check.Status = "passing"
}

func (serviceInstance *DefaultServiceInstance) GetInstanceId() string {
	return serviceInstance.InstanceId
}

func (serviceInstance *DefaultServiceInstance) GetServiceId() string {
	return serviceInstance.ServiceId
}

func (serviceInstance *DefaultServiceInstance) GetHost() string {
	return serviceInstance.Host
}

func (serviceInstance *DefaultServiceInstance) GetPort() int {
	return serviceInstance.Port
}

func (serviceInstance *DefaultServiceInstance) IsSecure() bool {
	return serviceInstance.Secure
}

func (serviceInstance *DefaultServiceInstance) GetMetadata() map[string]string {
	return serviceInstance.Metadata
}

func (serviceInstance *DefaultServiceInstance) GetCheck() *api.AgentServiceCheck {
	return serviceInstance.Check
}

func FindFirstNonLoopbackIP() (ipv4 string, err error) {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags & net.FlagUp) != 0 {
			adders, _ := netInterfaces[i].Addrs()

			for _, address := range adders {
				if inet, ok := address.(*net.IPNet); ok && !inet.IP.IsLoopback() {
					fmt.Println(inet)
					if inet.IP.To4() != nil {
						return inet.IP.String(), nil
					}
				}
			}
		}
	}

	return "", fmt.Errorf("not find")
}
