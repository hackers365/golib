package consul

import (
	// "fmt"
	log "github.com/Sirupsen/logrus"
	"time"
)

const FAILED_TIMES = 5

type Registry interface {
	RegisterWithHttp(serviceName string, ip string, port int, checkUrl string, interval string, timeout string, deRegisterTime string) error
	RegisterWithTtl(serviceName string, ip string, port int, ttl string, keepaliveTime int, deRegisterTime string) error
	Deregister()
	GetService(serviceName string) ([]ServiceInstance, error)
	Watch(serviceName string, handler func(string, []ServiceInstance)) error
	GetAndWatch(serviceName string, handler func(string, []ServiceInstance)) ([]ServiceInstance, error)
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

	registerAgainChan := make(chan bool)
	serviceInsChan := make(chan *DefaultServiceInstance)

	//keepalive
	go func() {
		for {
			select {
			case serviceInstance := <-serviceInsChan:
				checkId := serviceInstance.GetCheck().CheckID
				failedTimes := 0
				for {
					if failedTimes == FAILED_TIMES {
						registerAgainChan <- true
						break
					}
					err := r.srcRegistry.TtlKeepalive(checkId, "pass")
					if err != nil {
						log.Error("keep alive error:", err.Error())
						failedTimes++
						time.Sleep(1 * time.Second)
						continue
					}
					failedTimes = 0
					time.Sleep(time.Duration(keepaliveTime) * time.Second)
				}
			}
		}

	}()

	serviceInsChan <- serviceInstanceInfo

	go func() {
		for {
			select {
			case <-registerAgainChan:
				for {
					err = r.srcRegistry.Register(serviceInstanceInfo)
					if err != nil {
						log.Error("register again error: ", err.Error())
						time.Sleep(1 * time.Second)
						continue
					}
					serviceInsChan <- serviceInstanceInfo
					break
				}
			}

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
