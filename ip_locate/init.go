/* init.go -   */
/*
modification history
--------------------
2017/08/17, by Lei Hong, Create
*/
/*
DESCRIPTION

*/
package ip_locate

import (
	"sync"
	"sync/atomic"
	"time"

	log "github.com/Sirupsen/logrus"

	maxminddb "github.com/oschwald/maxminddb-golang"
)

var ipdb *maxminddb.Reader
var provAreaConf map[uint16]uint16 //省份->大区 映射关系

func Init(dbPath string) error {
	db, err := maxminddb.Open(dbPath)
	if err != nil {
		return err
	}

	ipdb = db
	ipCacheMap = &sync.Map{}

	go monitorIpCacheMap()

	conf := map[uint16][]uint16{
		1: []uint16{13, 9, 17, 14, 2, 19},
		2: []uint16{12, 29, 18, 22, 20},
		4: []uint16{1, 3, 10, 15, 16, 11, 32},
		5: []uint16{31, 30, 26, 25, 21},
		6: []uint16{24, 7, 27},
		7: []uint16{8, 5, 23, 28, 4},
		8: []uint16{6, 33, 34},
	}
	provAreaConf = make(map[uint16]uint16, 50)
	for area, provs := range conf {
		for _, prov := range provs {
			provAreaConf[prov] = area
		}
	}

	return nil
}

func monitorIpCacheMap() {
	var hitRatio float32

	ticker := time.NewTicker(3 * time.Minute)
	for {
		<-ticker.C
		num_miss := atomic.LoadInt32(&numMiss)
		num_hit := atomic.LoadInt32(&numHit)
		if num_miss >= max_capacity {
			cacheMap := &sync.Map{}
			ipCacheMap = cacheMap
			atomic.StoreInt32(&numMiss, 0)
			atomic.StoreInt32(&numHit, 0)

			hitRatio = float32(num_hit) / float32(num_hit+num_miss)

			log.Info("query total num: %d hitCached num: %d missCache num: %d hitRatio: %.3f", num_hit+num_miss, num_hit, num_miss, hitRatio)
		}

	}
}
