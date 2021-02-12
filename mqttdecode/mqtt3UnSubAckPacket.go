package mqttdecode

import (
	"encoding/binary"

	"github.com/google/gopacket"
)

var MQTT3UnsubAckPacket = gopacket.RegisterLayerType(
	3894,
	gopacket.LayerTypeMetadata{Name: "MQTT 3.1.1 UNSUBACK Packet", Decoder: gopacket.DecodeFunc(DecodeMQTT3UnsubAckPacket)})

type mqtt3UnsubAckPacket struct {
	VariableHeader mqtt3UnsubAckVariableHeader
	Contents       []byte
}

func (layer mqtt3UnsubAckPacket) LayerType() gopacket.LayerType { return MQTT3UnsubAckPacket }
func (layer mqtt3UnsubAckPacket) LayerContents() []byte         { return layer.Contents }
func (layer mqtt3UnsubAckPacket) LayerPayload() []byte          { return nil }

func DecodeMQTT3UnsubAckPacket(data []byte, packet gopacket.PacketBuilder) (err error) {
	variableHeader, err := decodeMQTT3UnsubAckVariableHeader(data)

	packet.AddLayer(&mqtt3UnsubAckPacket{variableHeader, data})
	return
}

type mqtt3UnsubAckVariableHeader struct {
	PacketIdentifier uint16
}

func decodeMQTT3UnsubAckVariableHeader(data []byte) (header mqtt3UnsubAckVariableHeader, err error) {
	header.PacketIdentifier = binary.BigEndian.Uint16(data)

	return
}
