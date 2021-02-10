package mqttdecode

import (
	"encoding/binary"

	"github.com/google/gopacket"
)

var MQTT3PubRecPacket = gopacket.RegisterLayerType(
	3888,
	gopacket.LayerTypeMetadata{Name: "MQTT 3.1.1 PUBREC Packet", Decoder: gopacket.DecodeFunc(decodeMQTT3PubRecPacket)})

type mqtt3PubRecPacket struct {
	VariableHeader mqtt3PubAckVariableHeader
	Contents       []byte
}

func (layer mqtt3PubRecPacket) LayerType() gopacket.LayerType { return MQTT3PubRecPacket }
func (layer mqtt3PubRecPacket) LayerContents() []byte         { return layer.Contents }
func (layer mqtt3PubRecPacket) LayerPayload() []byte          { return nil }

func decodeMQTT3PubRecPacket(data []byte, packet gopacket.PacketBuilder) (err error) {
	variableHeader, err := decodeMQTT3PubRecVariableHeader(data)

	packet.AddLayer(&mqtt3PubAckPacket{variableHeader, data})
	return
}

type mqtt3PubRecVariableHeader struct {
	PacketIdentifier uint16
}

func decodeMQTT3PubRecVariableHeader(data []byte) (header mqtt3PubAckVariableHeader, err error) {
	header.PacketIdentifier = binary.BigEndian.Uint16(data[:1])

	return
}
