package consul

import (
  "fmt"
  "sync"
  "sync/atomic"
)

type Discovery interface {
  WatchAndSave(serviceName string) error
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

func (r *discovery) WatchAndSave(serviceName string) error {
  handler := func(instances []ServiceInstance) {
    r.service2Instances.Store(serviceName, instances)
  }
  
  serviceList, err := r.srcRegistry.GetInstances(serviceName)
  if err != nil {
    return err
  }

  var index uint64
  r.service2Instances.Store(serviceName, serviceList)
  r.service2Index.Store(serviceName, &index)
  go r.srcRegistry.WatchPlan(serviceName, handler)
  return err
}

func (r *discovery) GetInstance(serviceName string) (string, error) {
  var serviceList []ServiceInstance
  var instance string
  if val, ok := r.service2Instances.Load(serviceName); ok {
    serviceList = val.([]ServiceInstance)
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

func getUrlFromInstance(instance ServiceInstance) string {
  schema := "http"
  if instance.IsSecure() {
    schema = "https"
  }

  return fmt.Sprintf("%s://%s:%d", schema, instance.GetHost(), instance.GetPort())
}

