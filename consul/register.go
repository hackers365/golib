package consul

import (
	"time"
)

type Registry interface {
  RegisterWithHttp(serviceName string, ip string, port int, checkUrl string, interval string, timeout string, deRegisterTime string) error
  RegisterWithTtl(serviceName string, ip string, port int, ttl string, keepaliveTime int, deRegisterTime string) error
  Deregister()
  GetService(serviceName string) ([]ServiceInstance, error)
  Watch(serviceName string, handler func([]ServiceInstance)) error
  GetAndWatch(serviceName string, handler func([]ServiceInstance)) ([]ServiceInstance, error)
}

type registry struct {
  srcRegistry ServiceRegistry
}

func NewRegistry(consulAddr string, token string) (Registry, error) {
	//服务注册
	registryDiscoveryClient, err := NewConsulServiceRegistry(consulAddr, token)
	if err != nil {
		return nil, err
	}

	reg := &registry{
		srcRegistry: registryDiscoveryClient,
	}

	return reg, nil
}

func (r *registry) RegisterWithHttp(serviceName string, ip string, port int, checkUrl string, interval string, timeout string, deRegisterTime string) error {
	serviceInstanceInfo, err := NewDefaultServiceInstance(serviceName, ip, port)
	if err != nil {
		return err
	}
	
	serviceInstanceInfo.SetHttpCheck(checkUrl, timeout, interval)
	serviceInstanceInfo.SetDeregisterAfter(deRegisterTime)

	err = r.srcRegistry.Register(serviceInstanceInfo)
	if err != nil {
		return err
	}
	return nil
}

func (r *registry) RegisterWithTtl(serviceName string, ip string, port int, ttl string, keepaliveTime int, deRegisterTime string) error {
	serviceInstanceInfo, err := NewDefaultServiceInstance(serviceName, ip, port)
	if err != nil {
		return err
	}

	serviceInstanceInfo.SetTtlCheck(ttl)
	serviceInstanceInfo.SetDeregisterAfter(deRegisterTime)

	err = r.srcRegistry.Register(serviceInstanceInfo)
	if err != nil {
		return err
	}

	//keepalive
	go func() {
		checkId := serviceInstanceInfo.GetCheck().CheckID
		for {
			err := r.srcRegistry.TtlKeepalive(checkId, "pass")
			if err != nil {
				continue
			}
			time.Sleep(time.Duration(keepaliveTime) * time.Second)
		}
	}()

	return nil
}

func (r *registry) Deregister() {
	r.srcRegistry.Deregister()
}

func (r *registry) GetService(serviceName string) ([]ServiceInstance, error) {
	return r.srcRegistry.GetInstances(serviceName)
}
 
func (r *registry) Watch(serviceName string, handler func(string, []ServiceInstance)) error {
	return r.srcRegistry.WatchPlan(serviceName, handler)
}

func (r *registry) GetAndWatch(serviceName string, handler func(string, []ServiceInstance)) ([]ServiceInstance, error) {
	serviceList, err := r.GetService(serviceName)
	if err != nil {
		return nil, err
	}
	go r.Watch(serviceName, handler)
	return serviceList, err
}
