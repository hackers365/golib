package consul

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
)

type ServiceRegistry interface {
	Register(serviceInstance ServiceInstance) error
	Deregister()

	GetInstances(serviceId string) ([]ServiceInstance, error)
	GetServices() ([]string, error)
	TtlKeepalive(checkId string, note string) error
	WatchPlan(wType, serviceName string, handler watch.HandlerFunc) error
}

type consulServiceRegistry struct {
	host string
	port int
	token string
	serviceInstances     map[string]map[string]ServiceInstance
	client               api.Client
	localServiceInstance ServiceInstance
}

func NewConsulServiceRegistry(host string, port int, token string) (ServiceRegistry, error) {
	if len(host) < 3 {
		return nil, fmt.Errorf("check host")
	}

	if port <= 0 || port > 65535 {
		return nil, fmt.Errorf("check port, port should between 1 and 65535")
	}

	config := api.DefaultConfig()
	config.Address = host + ":" + strconv.Itoa(port)
	config.Token = token
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	instance := &consulServiceRegistry{client: *client}
	instance.host = host
	instance.port = port
	instance.token = token
	return instance, nil
}

/**
 * 服务发现的watch
 */
func (c *consulServiceRegistry) WatchPlan(wType, serviceName string, handler watch.HandlerFunc) error {
	watchConfig := make(map[string]interface{})

  watchConfig["type"] = wType
  watchConfig["service"] = serviceName
  watchConfig["handler_type"] = "script"

  watchPlan, err := watch.Parse(watchConfig)
  if err != nil {
  	fmt.Println(err)
  	return err
  }
  watchPlan.Handler = handler

  return watchPlan.Run(fmt.Sprintf("%s:%d", c.host, c.port))
}

/**
服务注册
*/
func (c *consulServiceRegistry) Register(serviceInstance ServiceInstance) error {
	// 创建注册到consul的服务到
	registration := new(api.AgentServiceRegistration)
	registration.ID = serviceInstance.GetInstanceId()
	registration.Name = serviceInstance.GetServiceId()
	registration.Port = serviceInstance.GetPort()
	var tags []string
	if serviceInstance.IsSecure() {
		tags = append(tags, "secure=true")
	} else {
		tags = append(tags, "secure=false")
	}
	if serviceInstance.GetMetadata() != nil {
		var tags []string
		for key, value := range serviceInstance.GetMetadata() {
			tags = append(tags, key+"="+value)
		}
	}
	registration.Tags = tags

	registration.Address = serviceInstance.GetHost()
	registration.Check = serviceInstance.GetCheck()

	// 注册服务到consul
	err := c.client.Agent().ServiceRegister(registration)
	if err != nil {
		return err
	}

	if c.serviceInstances == nil {
		c.serviceInstances = map[string]map[string]ServiceInstance{}
	}

	services := c.serviceInstances[serviceInstance.GetServiceId()]
	if services == nil {
		services = map[string]ServiceInstance{}
	}
	services[serviceInstance.GetInstanceId()] = serviceInstance

	c.serviceInstances[serviceInstance.GetServiceId()] = services

	c.localServiceInstance = serviceInstance

	return nil
}

/**
 * check为ttl类型下keepalive
 */
func (c *consulServiceRegistry) TtlKeepalive(checkId string, note string) error {
	return c.client.Agent().PassTTL(checkId, note)
}

/**
服务剔除
*/
func (c *consulServiceRegistry) Deregister() {
	if c.serviceInstances == nil {
		return
	}

	services := c.serviceInstances[c.localServiceInstance.GetServiceId()]

	if services == nil {
		return
	}

	delete(services, c.localServiceInstance.GetInstanceId())

	if len(services) == 0 {
		delete(c.serviceInstances, c.localServiceInstance.GetServiceId())
	}

	_ = c.client.Agent().ServiceDeregister(c.localServiceInstance.GetInstanceId())

	c.localServiceInstance = nil
}

/**
  获取实例
*/
func (c *consulServiceRegistry) GetInstances(serviceId string) ([]ServiceInstance, error) {
	serviceList, _, _ := c.client.Health().Service(serviceId, "", true, nil)
	if len(serviceList) > 0 {
		result := make([]ServiceInstance, len(serviceList))
		for index, sever := range serviceList {
			s := &DefaultServiceInstance{
				InstanceId: sever.Service.ID,
				ServiceId:  sever.Service.Service,
				Host:       sever.Service.Address,
				Port:       sever.Service.Port,
				Metadata:   sever.Service.Meta,
			}
			result[index] = s
		}
		return result, nil
	}
	return nil, nil
}

/**
服务列表
*/
func (c *consulServiceRegistry) GetServices() ([]string, error) {
	services, _, _ := c.client.Catalog().Services(nil)
	result := make([]string, len(services))
	index := 0
	for serviceName, _ := range services {
		result[index] = serviceName
		index++
	}
	return result, nil
}

