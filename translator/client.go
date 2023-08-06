package translator

import (
	"fmt"

	"github.com/elastic/apm-data/model/modelpb"
	"go.opentelemetry.io/collector/pdata/pcommon"
)

func parseClient(client *modelpb.Client, attrs pcommon.Map) {
	if client == nil {
		return
	}

	parseClientIP(client.Ip, attrs)
	PutOptionalStr(attrs, "client.domain", &client.Domain)
	PutOptionalInt(attrs, "client.port", &client.Port)
}

func parseClientIP(ip *modelpb.IP, attrs pcommon.Map) {
	if ip == nil {
		return
	}

	if &ip.V4 != nil {
		attrs.PutStr("client.ip.v4", parseIPV4(ip.V4))
	}

	if &ip.V6 != nil && len(ip.V6) == 16 {
		attrs.PutStr("client.ip.v6", parseIPV6(ip.V6))
	}
}

// Turn a uint32 into an IPv4 address string (e.g. 192.168.1.1)
func parseIPV4(ipv4 uint32) string {
	return fmt.Sprintf("%d.%d.%d.%d", ipv4>>24, (ipv4>>16)&0xff, (ipv4>>8)&0xff, ipv4&0xff)
}

// Turn a []byte into an IPv6 address string (e.g. 2001:0db8:85a3:0000:0000:8a2e:0370:7334)
func parseIPV6(ipv6 []byte) string {
	return fmt.Sprintf("%x:%x:%x:%x:%x:%x:%x:%x", ipv6[0:2], ipv6[2:4], ipv6[4:6], ipv6[6:8], ipv6[8:10], ipv6[10:12], ipv6[12:14], ipv6[14:16])
}
