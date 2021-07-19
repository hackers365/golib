package ip_locate

import (
	"errors"
	"sync"
	"sync/atomic"
)

type ipStat struct {
	location *IpLocation
	sumQuery uint32
}

var ipCacheMap *sync.Map
var numHit, numMiss int32
var max_capacity int32 = 1000000

func IpLocate(ipAddr string) (*IpLocation, error) {
	if ipCache, ok := ipCacheMap.Load(ipAddr); ok {
		if ipstat, isOk := ipCache.(*ipStat); isOk {
			atomic.AddInt32(&numHit, 1)
			if !isOk || ipstat == nil || ipstat.location == nil {
				return nil, errors.New("failed to find location")
			}
			ipstat.sumQuery++
			return ipstat.location, nil
		}
	}

	location, err := ipLocate(ipAddr)
	ipCacheMap.Store(ipAddr, &ipStat{
		location: location,
		sumQuery: 1,
	})
	atomic.AddInt32(&numMiss, 1)

	return location, err
}
