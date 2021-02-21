package mqttdecode

import (
	"encoding/binary"

	"github.com/google/gopacket"
)

var MQTT5UnsubAckPacket = gopacket.RegisterLayerType(
	3905,
	gopacket.LayerTypeMetadata{Name: "MQTT 5.0 UNSUBACK Packet", Decoder: gopacket.DecodeFunc(DecodeMQTT5UnsubAckPacket)})

type mqtt5UnsubAckPacket struct {
	VariableHeader  mqtt5UnsubAckVariableHeader
	ReturnCodes     []byte
	ReturnCodeNames []string
	Contents        []byte
}

func (layer mqtt5UnsubAckPacket) LayerType() gopacket.LayerType { return MQTT5UnsubAckPacket }
func (layer mqtt5UnsubAckPacket) LayerContents() []byte         { return layer.Contents }
func (layer mqtt5UnsubAckPacket) LayerPayload() []byte          { return nil }

func DecodeMQTT5UnsubAckPacket(data []byte, packet gopacket.PacketBuilder) (err error) {
	variableHeader, err := decodeMQTT5UnsubAckVariableHeader(data)
	rcNames := decodeMQTT5UnsubAckPayload(data[variableHeader.Length:])
	packet.AddLayer(&mqtt5UnsubAckPacket{variableHeader, data[variableHeader.Length:], rcNames, data})
	return
}

type mqtt5UnsubAckVariableHeader struct {
	PacketIdentifier uint16
	Properties       []MQTT5Property
	Length           int
}

func decodeMQTT5UnsubAckVariableHeader(data []byte) (header mqtt5UnsubAckVariableHeader, err error) {
	header.PacketIdentifier = binary.BigEndian.Uint16(data)
	var propertiesLength int
	header.Properties, propertiesLength = extractMQTT5Properties(data[2:])
	header.Length = 2 + propertiesLength
	return
}

func decodeMQTT5UnsubAckPayload(data []byte) (names []string) {
	for _, rc := range data {
		names = append(names, resolveReasonCode(rc))
	}
	return names
}
