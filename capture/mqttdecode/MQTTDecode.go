package mqttdecode

import (
	"errors"
	"log"

	"github.com/google/gopacket"
)

type controlPacketType int

const (
	CONNECT = controlPacketType(iota + 1)
	CONNACK
	PUBLISH
	PUBACK
	PUBREC
	PUBREL
	PUBCOMP
	SUBSCRIBE
	SUBACK
	UNSUBSCRIBE
	UNSUBACK
	PINGREQ
	PINGRESP
	DISCONNECT
	AUTH
)

func (cpt controlPacketType) String() string {
	return [...]string{"", "CONNECT", "CONNACK", "PUBLISH", "PUBACK", "PUBREC", "PUBREL",
		"PUBCOMP", "SUBSCRIBE", "SUBACK", "UNSUBSCRIBE", "UNSUBACK", "PINGREQ", "PINGRESP", "DISCONNECT", "AUTH"}[cpt]
}

var MQTT3LayerType = gopacket.RegisterLayerType(
	3883,
	gopacket.LayerTypeMetadata{Name: "MQTT3.1.1", Decoder: gopacket.DecodeFunc(decodeMQTT3)})

type MQTT3Layer struct {
	ControlPacketType string
	Flags             byte
	RemainingLength   int
	Contents          []byte
}

func (layer MQTT3Layer) LayerType() gopacket.LayerType { return MQTT3LayerType }
func (layer MQTT3Layer) LayerContents() []byte         { return layer.Contents }
func (layer MQTT3Layer) LayerPayload() []byte          { return nil }

func calculateRemainingLength(packet []byte) (remainingLength int) {

	multiplier := 1
	remainingLength = 0

	for _, nextByte := range packet {
		remainingLength += int(nextByte&127) * multiplier
		multiplier *= 128
		if multiplier > 128*128*128 {
			log.Println(errors.New("Malformed Remaining Length Header"))
		}
		if nextByte&128 == 0 {
			break
		}
	}

	return
}

func decodeControlPacketType(header byte) controlPacketType {
	//MQTT control packet type is determined by the value of the 4 highest bits of the packet's first byte
	return controlPacketType((header & 0xF0) >> 4)
}

func decodeMQTT3(data []byte, packet gopacket.PacketBuilder) error {
	remainingLength := calculateRemainingLength(data[1:])
	packet.AddLayer(&MQTT3Layer{decodeControlPacketType(data[0]).String(), data[0] & 0xF, remainingLength, data[(remainingLength/128 + 2):]})

	return nil
}
