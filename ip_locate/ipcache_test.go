package ip_locate

import (
	"bufio"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
	"unsafe"
)

func TestMonitorIpCacheMap(t *testing.T) {
	InitTest()
	sumIp := 0
	sumTests := 0

	ipCacheMap.Range(func(key, value interface{}) bool {
		sumIp++
		return true
	})
	t.Logf("sum:%d\n", sumIp)
	start := time.Now()
	for i := 0; i < 1000; i++ {
		monitorIpCacheMap()
		sumTests++
		if sumTests%10 == 0 {
			elapsed := int64(time.Since(start) / time.Millisecond)
			t.Logf("sumTests:%d, elapsed:%d\n", sumTests, elapsed)
			start = time.Now()
		}
	}
}

func TestMonitor(t *testing.T) {
	InitTest()
	sumIp := 0
	sumTests := 0

	ipCacheMap.Range(func(key, value interface{}) bool {
		sumIp++
		return true
	})
	t.Logf("sum:%d\n", sumIp)

	start := time.Now()
	for i := 0; i < 1000; i++ {
		monitorIpCacheMap()
		sumTests++
		if sumTests%10 == 0 {
			elapsed := int64(time.Since(start) / time.Millisecond)
			t.Logf("sumTests:%d, elapsed:%d\n", sumTests, elapsed)
			sumIp = 0
			ipCacheMap.Range(func(key, value interface{}) bool {
				sumIp++
				return true
			})
			t.Logf("sum:%d\n", sumIp)
			start = time.Now()
		}
		InitTest()
	}
}

func BenchmarkMonitor(b *testing.B) {
	InitTest()
	sum := 0
	ipCacheMap.Range(func(key, value interface{}) bool {
		sum++
		return true
	})
	b.Logf("sum:%d\n", sum)
	for i := 0; i < b.N; i++ {
		monitorIpCacheMap()
	}
}

func bytes2string(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func InitTest() {
	err := Init("../../../sched_execute/data/ipdb_prov.mmdb")
	if err != nil {
		log.Fatalf("maxminddb.Open(): %s", err.Error())
		return
	}

	f, err := os.Open("./select_sid.txt")
	if err != nil {
		log.Fatalf("failed to open file,err:%s\n", err)
	}
	reader := bufio.NewReader(f)
	var lineByte []byte
	for {
		lineByte, err = reader.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		line := bytes2string(lineByte)
		for i := 0; i < rand.Intn(30); i++ {
			_, err = IpLocate(line)
		}
	}
}

func BenchmarkTest(b *testing.B) {
	b.StopTimer()
	InitTest()
	b.StartTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := IpLocate("36.110.59.146")
			if err != nil {
				b.Errorf("%+v", err)
			}
		}
	})
}

func InitMinitor_manual() {
	ipCacheMap = &sync.Map{}
	var ipAddr strings.Builder
	var ipstat *ipStat
	for a := 0; a < 128; a++ {
		for b := 0; b < 128; b++ {
			for c := 0; c < 128; c++ {
				for d := 0; d < 1; d++ {
					ipAddr.WriteString(strconv.Itoa(a))
					ipAddr.WriteString(".")
					ipAddr.WriteString(strconv.Itoa(b))
					ipAddr.WriteString(".")
					ipAddr.WriteString(strconv.Itoa(c))
					ipAddr.WriteString(".")
					ipAddr.WriteString(strconv.Itoa(d))
					ipstat = &ipStat{
						location: &IpLocation{
							Country: 1,
							Prov:    2,
							City:    3,
							Isp:     4,
						},
						sumQuery: uint32(rand.Intn(20)),
					}
					ipCacheMap.Store(ipAddr.String(), ipstat)
					ipAddr.Reset()
				}
			}
		}
	}
}

func BenchmarkIpCache(b *testing.B) {
	b.StopTimer()
	InitMinitor_manual()
	var sumIp = 0
	ipCacheMap.Range(func(key, value interface{}) bool {
		sumIp++
		return true
	})
	b.Logf("sum:%d\n", sumIp)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		begin := time.Now()
		monitorIpCacheMap()
		elapsed := int64(time.Since(begin) / time.Millisecond)
		b.Logf("elapsed:%dms\n", elapsed)
	}
	sumIp = 0
	ipCacheMap.Range(func(key, value interface{}) bool {
		sumIp++
		return true
	})
	b.Logf("sum:%d\n", sumIp)
	b.ReportAllocs()

}

func TestIpCacheMap(t *testing.T) {

	InitMinitor_manual()
	var sumIp = 0
	ipCacheMap.Range(func(key, value interface{}) bool {
		sumIp++
		return true
	})
	t.Logf("sum:%d\n", sumIp)

	begin := time.Now()
	monitorIpCacheMap()
	elapsed := int64(time.Since(begin) / time.Millisecond)
	t.Logf("elapsed:%dms\n", elapsed)

	sumIp = 0
	ipCacheMap.Range(func(key, value interface{}) bool {
		sumIp++
		return true
	})
	t.Logf("sum:%d\n", sumIp)
}
