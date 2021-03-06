/*Package capture contains utilities for MQTT packet-capturing.
 */
package capture

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/KRsully/mqtt-security-analysis/mqttdecode"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/pcapgo"
)

const mqttPort = 1883
const mqttSecPort = 8883

//ListAvailableInterfaces uses pcap to print the name, descriptor, and IP addresses of the machine's network interfaces to the console.
func ListAvailableInterfaces() {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Panic(err)
	}

	for _, device := range devices {
		fmt.Printf("%s (%s):\n", device.Name, device.Description)
		for _, address := range device.Addresses {
			fmt.Printf("%12s %s \n", "-", address.IP)
		}
	}
}

//Create an inactive pcap handle given the user-supplied parameters or the defaults.
//	deviceName is the name of the interface to create the handle on
//	Returns a pointer to the InactiveHandle object
func createInactiveHandle(promiscuity bool, snaplen int, timeout time.Duration, deviceName string) *pcap.InactiveHandle {
	inactive, err := pcap.NewInactiveHandle(deviceName)
	if err != nil {
		log.Fatal(err)
	}

	if err = inactive.SetPromisc(promiscuity); err != nil {
		log.Fatal(err)
	} else if err = inactive.SetSnapLen(snaplen); err != nil {
		log.Fatal(err)
	} else if err = inactive.SetTimeout(timeout); err != nil {
		log.Fatal(err)
	}

	log.Printf("Created inactive pcap handle to capture traffic on device : %s.\n%30s %t\n%30s %d\n%30s %v\n",
		deviceName, "Promiscuous:", promiscuity, "Snapshot Length:", snaplen, "Timeout:", timeout)

	return inactive
}

//Parse the given interface string to determine if pcap can find a valid interface
//	Receives a string which may be the IP address or the name of the network interface
//	Returns a string with the name of one of the machine's interfaces
func parseInterface(deviceIdentifier string) string {

	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Panic(err)
	}

	for _, device := range devices {
		if device.Name == deviceIdentifier {
			return device.Name
		}
		for _, address := range device.Addresses {
			targetIP := net.ParseIP(deviceIdentifier)
			if address.IP.Equal(targetIP) {
				return device.Name
			}
		}
	}

	log.Fatal("Cannot resolve interface name")

	return ""
}

//MQTTPackets uses the pcap and gopacket libraries to capture MQTT packets on the user-specified network interface
func MQTTPackets(ports string, promiscuity bool, snaplen int, timeout time.Duration, device []string) {

	deviceName := parseInterface(device[0])

	inactive := createInactiveHandle(promiscuity, snaplen, timeout, deviceName)

	pcapHandle, err := inactive.Activate()
	if err != nil {
		log.Panic(err)
	}

	inactive.CleanUp()

	defer pcapHandle.Close()

	//Create a BPF filter string that only lets traffic on the desginated ports pass through
	filter := "tcp and port"
	portsSlice := strings.Split(ports, ",")
	for _, port := range portsSlice[:len(portsSlice)-1] {
		filter += fmt.Sprintf(" %s or", port)
	}
	filter += fmt.Sprintf(" %s", portsSlice[len(portsSlice)-1])
	pcapHandle.SetBPFFilter(filter)

	packetSource := gopacket.NewPacketSource(pcapHandle, pcapHandle.LinkType())

	//Create a pcap file with the packets that we capture
	file, err := os.Create("capture.pcap")
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()

	writer := pcapgo.NewWriter(file)
	writer.WriteFileHeader(uint32(snaplen), pcapHandle.LinkType())

	for packet := range packetSource.Packets() {
		writer.WritePacket(packet.Metadata().CaptureInfo, packet.Data())
		//fmt.Println(packet)
		var decodedPacket gopacket.Packet
		if packet.ApplicationLayer() != nil {
			decodedPacket = gopacket.NewPacket(packet.ApplicationLayer().Payload(), layers.LayerTypeTLS, gopacket.DecodeOptions{Lazy: false, DecodeStreamsAsDatagrams: true})
			if decodedPacket.ErrorLayer() != nil {
				decodedPacket = gopacket.NewPacket(packet.ApplicationLayer().Payload(), mqttdecode.MQTTFixedHeaderLayer, gopacket.DecodeOptions{Lazy: false, DecodeStreamsAsDatagrams: true})
			}
			fmt.Println(decodedPacket)
		}

	}
}
