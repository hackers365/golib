/* ip_locate_test.go -   */
/*
modification history
--------------------
2017/05/25, by Lei Hong, Create
*/
/*
DESCRIPTION

*/
package ip_locate

import (
	"fmt"
	"testing"
	"time"
)

var num = 100

func TestIpLocation(t *testing.T) {
	// ip db
	err := Init("../../../sched_execute/data/ipdb_prov.mmdb")
	if err != nil {
		t.Errorf("maxminddb.Open(): %s", err.Error())
		return
	}

	var location *IpLocation
	begin := time.Now().UnixNano()
	for i := 0; i < num; i++ {
		location, err = IpLocate("36.110.59.146")
		if err != nil {
			t.Errorf("%+v", err)
		}
	}

	fmt.Println(location)

	end := time.Now().UnixNano()
	elapsed := end - begin
	fmt.Printf("elapasd:%dns\n", elapsed)
}

func BenchmarkIpLocate(b *testing.B) {
	// ip db
	err := Init("../../../sched_execute/data/ipdb_prov.mmdb")
	if err != nil {
		b.Errorf("maxminddb.Open(): %s", err.Error())
		return
	}
	var location *IpLocation
	begin := time.Now().UnixNano()
	for i := 0; i < b.N; i++ {
		location, err = IpLocate("36.110.59.146")
		if err != nil {
			b.Errorf("%+v", err)
		}
	}

	fmt.Println(location)

	end := time.Now().UnixNano()
	elapsed := end - begin
	fmt.Printf("elapasd:%dns\n", elapsed)
}
