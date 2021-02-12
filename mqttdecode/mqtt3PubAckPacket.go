package mqttdecode

import (
	"encoding/binary"

	"github.com/google/gopacket"
)

var MQTT3PubAckPacket = gopacket.RegisterLayerType(
	3887,
	gopacket.LayerTypeMetadata{Name: "MQTT 3.1.1 PUBACK Packet", Decoder: gopacket.DecodeFunc(DecodeMQTT3PubAckPacket)})

type mqtt3PubAckPacket struct {
	VariableHeader mqtt3PubAckVariableHeader
	Contents       []byte
}

func (layer mqtt3PubAckPacket) LayerType() gopacket.LayerType { return MQTT3PubAckPacket }
func (layer mqtt3PubAckPacket) LayerContents() []byte         { return layer.Contents }
func (layer mqtt3PubAckPacket) LayerPayload() []byte          { return nil }

func DecodeMQTT3PubAckPacket(data []byte, packet gopacket.PacketBuilder) (err error) {
	variableHeader, err := decodeMQTT3PubAckVariableHeader(data)

	packet.AddLayer(&mqtt3PubAckPacket{variableHeader, data})
	return
}

type mqtt3PubAckVariableHeader struct {
	PacketIdentifier uint16
}

func decodeMQTT3PubAckVariableHeader(data []byte) (header mqtt3PubAckVariableHeader, err error) {
	header.PacketIdentifier = binary.BigEndian.Uint16(data)

	return
}
