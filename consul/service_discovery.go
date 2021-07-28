package consul

import (
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
)

type ServiceRegistry interface {
	Register(serviceInstance ServiceInstance) error
	Deregister()

	GetInstances(serviceId string) ([]ServiceInstance, error)
	GetServices() ([]string, error)
	TtlKeepalive(checkId string, note string) error
	WatchPlan(serviceName string, handler func([]ServiceInstance)) error
}

type consulServiceRegistry struct {
	addr string
	token string
	serviceInstances     map[string]map[string]ServiceInstance
	client               api.Client
	localServiceInstance ServiceInstance
}

func NewConsulServiceRegistry(addr string, token string) (ServiceRegistry, error) {
	config := api.DefaultConfig()
	config.Address = addr
	config.Token = token
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	instance := &consulServiceRegistry{client: *client}
	instance.addr = addr
	instance.token = token
	return instance, nil
}

/**
 * 服务发现的watch
 */
func (c *consulServiceRegistry) WatchPlan(serviceName string, handler func(string, []ServiceInstance)) error {
	watchConfig := make(map[string]interface{})

  watchConfig["type"] = "service"
  watchConfig["service"] = serviceName
  watchConfig["handler_type"] = "script"
  watchConfig["passingonly"] = true

  watchPlan, err := watch.Parse(watchConfig)
  if err != nil {
  	return err
  }

  cb := func(lastIndex uint64, result interface{}) {
      serviceList := result.([]*api.ServiceEntry)
			ret := make([]ServiceInstance, len(serviceList))
			for index, service := range serviceList {
				ret[index] = getServiceInstance(service)
			}
			handler(serviceName, ret)
  }

  watchPlan.Handler = cb

  return watchPlan.Run(c.addr)
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
		for index, service := range serviceList {
			result[index] = getServiceInstance(service)
		}
		return result, nil
	}
	return nil, nil
}

func getServiceInstance(service *api.ServiceEntry) ServiceInstance {
	return &DefaultServiceInstance{
		InstanceId: service.Service.ID,
		ServiceId:  service.Service.Service,
		Host:       service.Service.Address,
		Port:       service.Service.Port,
		Metadata:   service.Service.Meta,
	}
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

