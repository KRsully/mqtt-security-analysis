/*
* Little Go program to run the packet-capture utilities in capture\capture.go
 */
package main

import (
	"flag"
	"os"

	"github.com/KRsully/mqtt-security-analysis/capture"
	"github.com/google/gopacket/pcap"
)

func main() {
	ports := flag.String("p", "1883,8883", "comma-separated port(s) to capture MQTT packets on")
	promiscuity := flag.Bool("pro", false, "promiscuous mode on")
	snaplen := flag.Int("sl", 65535, "maximum bytes per packet capture") //Is MQTT packet max 512 bytes?
	timeout := flag.Duration("t", pcap.BlockForever, "time to wait for additional packets after packet capture before returning (min. nanoseconds)")

	flag.Parse()

	flag.Args()

	if len(os.Args) < 2 {
		capture.ListAvailableInterfaces()
	} else {
		capture.MQTTPackets(*ports, *promiscuity, *snaplen, *timeout, flag.Args())
	}
}
