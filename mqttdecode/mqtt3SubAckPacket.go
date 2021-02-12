package mqttdecode

import (
	"encoding/binary"

	"github.com/google/gopacket"
)

var MQTT3SubAckPacket = gopacket.RegisterLayerType(
	3892,
	gopacket.LayerTypeMetadata{Name: "MQTT 3.1.1 SUBACK Packet", Decoder: gopacket.DecodeFunc(DecodeMQTT3SubAckPacket)})

type mqtt3SubAckPacket struct {
	VariableHeader mqtt3SubAckVariableHeader
	ReturnCodes    []byte
	Contents       []byte
}

func (layer mqtt3SubAckPacket) LayerType() gopacket.LayerType { return MQTT3SubAckPacket }
func (layer mqtt3SubAckPacket) LayerContents() []byte         { return layer.Contents }
func (layer mqtt3SubAckPacket) LayerPayload() []byte          { return nil }

func DecodeMQTT3SubAckPacket(data []byte, packet gopacket.PacketBuilder) (err error) {
	variableHeader, err := decodeMQTT3SubAckVariableHeader(data)

	packet.AddLayer(&mqtt3SubAckPacket{variableHeader, data[2:], data})
	return
}

type mqtt3SubAckVariableHeader struct {
	PacketIdentifier uint16
}

func decodeMQTT3SubAckVariableHeader(data []byte) (header mqtt3SubAckVariableHeader, err error) {
	header.PacketIdentifier = binary.BigEndian.Uint16(data)

	return
}
