package mqttdecode

import (
	"encoding/binary"

	"github.com/google/gopacket"
)

var MQTT3PubCompPacket = gopacket.RegisterLayerType(
	3890,
	gopacket.LayerTypeMetadata{Name: "MQTT 3.1.1 PUBCOMP Packet", Decoder: gopacket.DecodeFunc(DecodeMQTT3PubCompPacket)})

type mqtt3PubCompPacket struct {
	VariableHeader mqtt3PubAckVariableHeader
	Contents       []byte
}

func (layer mqtt3PubCompPacket) LayerType() gopacket.LayerType { return MQTT3PubCompPacket }
func (layer mqtt3PubCompPacket) LayerContents() []byte         { return layer.Contents }
func (layer mqtt3PubCompPacket) LayerPayload() []byte          { return nil }

func DecodeMQTT3PubCompPacket(data []byte, packet gopacket.PacketBuilder) (err error) {
	variableHeader, err := decodeMQTT3PubRecVariableHeader(data)

	packet.AddLayer(&mqtt3PubAckPacket{variableHeader, data})
	return
}

type mqtt3PubCompVariableHeader struct {
	PacketIdentifier uint16
}

func decodeMQTT3PubCompVariableHeader(data []byte) (header mqtt3PubAckVariableHeader, err error) {
	header.PacketIdentifier = binary.BigEndian.Uint16(data[:1])

	return
}
