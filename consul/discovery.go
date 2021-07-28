package consul

import (
  "fmt"
  "sync"
  "sync/atomic"
)

type Discovery interface {
  GetInstance(serviceName string) (string, error)
}

type discovery struct {
  srcRegistry ServiceRegistry
  service2Instances *sync.Map
  service2Index *sync.Map
}

func NewDiscovery(consulAddr string, token string) (Discovery, error) {
  //服务注册
  registryDiscoveryClient, err := NewConsulServiceRegistry(consulAddr, token)
  if err != nil {
    return nil, err
  }

  reg := &discovery{
    srcRegistry: registryDiscoveryClient,
    service2Instances: &sync.Map{},
    service2Index: &sync.Map{},
  }

  return reg, nil
}

func (r *discovery) GetInstance(serviceName string) (string, error) {
  var serviceList []ServiceInstance
  var instance string
  var err error
  if val, ok := r.service2Instances.Load(serviceName); ok {
    serviceList = val.([]ServiceInstance)
  } else {
    serviceList, err = r.watchAndSave(serviceName)
    if err != nil {
      return "", err
    }
  }

  if len(serviceList) == 0 {
    return "", fmt.Errorf("not found serviceList")
  }

  if data, ok := r.service2Index.Load(serviceName); ok {
    if indexAddr, ok := data.(*uint64); ok {
      index := atomic.AddUint64(indexAddr, 1)
      if index > 1000000000 {
        atomic.StoreUint64(indexAddr, 0)
        index = 0
      }
      rIndex := index % uint64(len(serviceList))
      sInstance := serviceList[rIndex]
      instance = getUrlFromInstance(sInstance)
    }
  }

  return instance, nil
}

func (r *discovery) watchAndSave(serviceName string) ([]ServiceInstance, error) {
  handler := func(svcName string, instances []ServiceInstance) {
    r.service2Instances.Store(svcName, instances)
  }
  
  serviceList, err := r.srcRegistry.GetInstances(serviceName)
  if err != nil {
    return nil, err
  }

  var index uint64
  r.service2Instances.Store(serviceName, serviceList)
  r.service2Index.Store(serviceName, &index)
  go r.srcRegistry.WatchPlan(serviceName, handler)
  return serviceList, err
}

func getUrlFromInstance(instance ServiceInstance) string {
  schema := "http"
  if instance.IsSecure() {
    schema = "https"
  }

  return fmt.Sprintf("%s://%s:%d", schema, instance.GetHost(), instance.GetPort())
}

