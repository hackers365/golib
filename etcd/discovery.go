package etcd

import (
  "fmt"
  "sync"
  "context"
  "sync/atomic"

  clientv3 "go.etcd.io/etcd/client/v3"
)

type Discovery interface {
  GetInstance(prefix string) (string, error)
}

type ParseCb func(string)string

type discovery struct {
  etcdClient *clientv3.Client
  service2Instances *sync.Map
  service2Index *sync.Map
  parseCb ParseCb
}

func NewDiscovery(etcdClient *clientv3.Client, cb ParseCb) (Discovery, error) {
  reg := &discovery{
    etcdClient: etcdClient,
    service2Instances: &sync.Map{},
    service2Index: &sync.Map{},
    parseCb: cb,
  }

  return reg, nil
}

func (r *discovery) GetInstance(prefix string) (string, error) {
  var serviceName2List *node2Url
  var instance string
  if val, ok := r.service2Instances.Load(prefix); ok {
    serviceName2List = val.(*node2Url)
  } else {
    //初始化
    serviceName2List = r.GetAndSaveInstance(prefix)

    var index uint64
    r.service2Index.Store(prefix, &index)

    go r.watchAndSave(prefix)    
  }

  if serviceName2List.Len() == 0 {
    return "", fmt.Errorf("not found serviceList")
  }

  if data, ok := r.service2Index.Load(prefix); ok {
    if indexAddr, ok := data.(*uint64); ok {
      index := atomic.AddUint64(indexAddr, 1)
      if index > 1000000000 {
        atomic.StoreUint64(indexAddr, 0)
        index = 0
      }
      rIndex := index % uint64(serviceName2List.Len())
      sInstance := serviceName2List.GetUrl(rIndex)
      //sInstance := serviceList[rIndex]
      instance = getUrlFromInstance(sInstance)
    }
  }

  return instance, nil
}

func (r *discovery) watchAndSave(prefix string) {
  watchChan := r.etcdClient.Watch(context.Background(), prefix, clientv3.WithPrefix())

  for _ = range watchChan {
    r.GetAndSaveInstance(prefix)
  }
}

func (r *discovery) GetAndSaveInstance(prefix string) *node2Url {
  serviceList := r.GetAllInstance(prefix)

  fmt.Println(serviceList)

  serviceInfo := &node2Url{
    urlList: serviceList,
  }
  r.service2Instances.Store(prefix, serviceInfo)
  return serviceInfo
}

func (r *discovery) GetAllInstance(prefix string) []string {
  kv := clientv3.NewKV(r.etcdClient)
  rangeResp, err := kv.Get(context.TODO(), prefix, clientv3.WithPrefix())
  if err != nil {
    return []string{}
  }
  
  serviceList := []string{}
  for _, kv := range rangeResp.Kvs {
    serviceList = append(serviceList, r.parseCb(string(kv.Value)))
  }
  
  return serviceList
}

func getUrlFromInstance(instance string) string {
  return instance
}

