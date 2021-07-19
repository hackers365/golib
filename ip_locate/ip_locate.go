/* ../common/ip_locate.go -   */
/*
modification history
--------------------
2017/04/06, by Lei Hong, Create
*/
/*
DESCRIPTION

*/
package ip_locate

import (
	"fmt"
	"net"
)

type IpLocation struct {
	Country uint16 `json:"country"`
	Area    uint16 `json:"area"`
	Prov    uint16 `json:"prov"`
	City    uint16 `json:"city"`
	Isp     uint16 `json:"isp"`
}

func ipLocate(ipAddr string) (*IpLocation, error) {
	if ipAddr == "" {
		return nil, fmt.Errorf("Ip addr is empty.")
	}
	if ipdb == nil {
		return nil, fmt.Errorf("Ipdb is not inited.")
	}

	ip := net.ParseIP(ipAddr)
	if ip == nil {
		return nil, fmt.Errorf("Ip is nil after net.ParseIP().")
	}

	var record map[string]interface{}
	if err := ipdb.Lookup(ip, &record); err != nil {
		return nil, fmt.Errorf("Lookup for location failed.")
	}

	if record["country"] == nil || record["prov"] == nil || record["isp"] == nil {
		return nil, fmt.Errorf("Ipdb.Lookup failed, record: %+v", record)
	}
	city := uint16(0)
	if record["city"] != nil {
		city = uint16(record["city"].(uint64))
	}

	location := &IpLocation{
		Country: uint16(record["country"].(uint64)),
		Prov:    uint16(record["prov"].(uint64)),
		Isp:     uint16(record["isp"].(uint64)),
		City:    city,
	}

	area, ok := provAreaConf[location.Prov]
	if !ok {
		area = 0
	}
	location.Area = area
	return location, nil
}
