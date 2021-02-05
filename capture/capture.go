/*
*Package capture contains utilities for MQTT packet-capturing.
 */
package capture

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/KRsully/mqtt-security-analysis/mqttdecode"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/pcapgo"
)

const mqttPort = 1883
const mqttSecPort = 8883

// type controlPacketType int

// const (
// 	CONNECT = controlPacketType(iota + 1)
// 	CONNACK
// 	PUBLISH
// 	PUBACK
// 	PUBREC
// 	PUBREL
// 	PUBCOMP
// 	SUBSCRIBE
// 	SUBACK
// 	UNSUBSCRIBE
// 	UNSUBACK
// 	PINGREQ
// 	PINGRESP
// 	DISCONNECT
// )

// var MQTT3LayerType = gopacket.RegisterLayerType(
// 	3883,
// 	gopacket.LayerTypeMetadata{Name: "MQTT3.1.1", Decoder: gopacket.DecodeFunc(decodeMQTT3)})

// type MQTT3Layer struct {
// 	ControlPacketType controlPacketType
// 	Flags             byte
// 	RemainingLength   int
// 	Contents          []byte
// }

// func (layer MQTT3Layer) LayerType() gopacket.LayerType { return MQTT3LayerType }
// func (layer MQTT3Layer) LayerContents() []byte         { return layer.Contents }
// func (layer MQTT3Layer) LayerPayload() []byte          { return nil }

// func calculateRemainingLength(packet []byte) (remainingLength int) {

// 	multiplier := 1
// 	remainingLength = 0

// 	for _, nextByte := range packet {
// 		remainingLength += int(nextByte&127) * multiplier
// 		multiplier *= 128
// 		if multiplier > 128*128*128 {
// 			log.Println(errors.New("Malformed Remaining Length Header"))
// 		}
// 		if nextByte&128 == 0 {
// 			break
// 		}
// 	}

// 	return
// }

// func decodeControlPacketType(header byte) controlPacketType {
// 	//MQTT control packet type is determined by the value of the 4 highest bits of the packet's first byte
// 	return controlPacketType((header & 0xF0) >> 4)
// }

// func decodeMQTT3(data []byte, packet gopacket.PacketBuilder) error {
// 	remainingLength := calculateRemainingLength(data[1:])
// 	packet.AddLayer(&MQTT3Layer{decodeControlPacketType(data[0]), data[0] & 0xF, remainingLength, data[(remainingLength/128 + 2):]})

// 	return nil
// }

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

//Parse the flags provided through the command line
//	Returns a pointer to a struct with the values of the flags
// func parseFlags() *pcapFlags {
// 	flags := pcapFlags{
// 		port:        flag.Int("p", mqttPort, "port to capture MQTT packets on"),
// 		promiscuity: flag.Bool("pro", false, "promiscuous mode on"),
// 		snaplen:     flag.Int("sl", 65535, "maximum bytes per packet capture"), //Is MQTT packet max 512 bytes?
// 		timeout:     flag.Duration("t", pcap.BlockForever, "time to wait for additional packets after packet capture before returning (min. nanoseconds)"),
// 	}

// 	flag.Parse()

// 	flags.tail = flag.Args()

// 	return &flags
// }

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
func MQTTPackets(port int, promiscuity bool, snaplen int, timeout time.Duration, device []string) {

	deviceName := parseInterface(device[0])

	inactive := createInactiveHandle(promiscuity, snaplen, timeout, deviceName)

	pcapHandle, err := inactive.Activate()
	if err != nil {
		log.Panic(err)
	}

	inactive.CleanUp()

	defer pcapHandle.Close()

	//TODO: Allow for user specification of ports to watch
	//pcapHandle.SetBPFFilter("tcp and port " + strconv.Itoa(port))
	pcapHandle.SetBPFFilter("tcp and port 1883 or port 8883")

	packetSource := gopacket.NewPacketSource(pcapHandle, pcapHandle.LinkType())

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
		if packet.ApplicationLayer() != nil {
			fmt.Println(gopacket.NewPacket(packet.ApplicationLayer().Payload(), mqttdecode.MQTT3LayerType, gopacket.Lazy))
		}
		// for _, layer := range packet.Layers() {
		// 	fmt.Println("LAYER: ", layer.LayerType())
		// }
	}
}
