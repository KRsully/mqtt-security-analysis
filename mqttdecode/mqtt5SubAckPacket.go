package mqttdecode

import (
	"encoding/binary"

	"github.com/google/gopacket"
)

var MQTT5SubAckPacket = gopacket.RegisterLayerType(
	3903,
	gopacket.LayerTypeMetadata{Name: "MQTT 5.0 SUBACK Packet", Decoder: gopacket.DecodeFunc(DecodeMQTT5SubAckPacket)})

type mqtt5SubAckPacket struct {
	VariableHeader mqtt5SubAckVariableHeader
	ReturnCodes    []byte
	Contents       []byte
}

func (layer mqtt5SubAckPacket) LayerType() gopacket.LayerType { return MQTT5SubAckPacket }
func (layer mqtt5SubAckPacket) LayerContents() []byte         { return layer.Contents }
func (layer mqtt5SubAckPacket) LayerPayload() []byte          { return nil }

func DecodeMQTT5SubAckPacket(data []byte, packet gopacket.PacketBuilder) (err error) {
	variableHeader, err := decodeMQTT5SubAckVariableHeader(data)

	packet.AddLayer(&mqtt5SubAckPacket{variableHeader, data[variableHeader.Length:], data})
	return
}

type mqtt5SubAckVariableHeader struct {
	PacketIdentifier uint16
	Properties       []MQTT5Property
	Length           int
}

func decodeMQTT5SubAckVariableHeader(data []byte) (header mqtt5SubAckVariableHeader, err error) {
	header.PacketIdentifier = binary.BigEndian.Uint16(data)
	var propertiesLength int
	header.Properties, propertiesLength = extractMQTT5Properties(data[2:])
	header.Length = 2 + propertiesLength
	return
}
