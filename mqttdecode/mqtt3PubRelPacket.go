package mqttdecode

import (
	"encoding/binary"

	"github.com/google/gopacket"
)

var MQTT3PubRelPacket = gopacket.RegisterLayerType(
	3889,
	gopacket.LayerTypeMetadata{Name: "MQTT 3.1.1 PUBREL Packet", Decoder: gopacket.DecodeFunc(DecodeMQTT3PubRelPacket)})

type mqtt3PubRelPacket struct {
	VariableHeader mqtt3PubAckVariableHeader
	Contents       []byte
}

func (layer mqtt3PubRelPacket) LayerType() gopacket.LayerType { return MQTT3PubRelPacket }
func (layer mqtt3PubRelPacket) LayerContents() []byte         { return layer.Contents }
func (layer mqtt3PubRelPacket) LayerPayload() []byte          { return nil }

func DecodeMQTT3PubRelPacket(data []byte, packet gopacket.PacketBuilder) (err error) {
	variableHeader, err := decodeMQTT3PubRecVariableHeader(data)

	packet.AddLayer(&mqtt3PubAckPacket{variableHeader, data})
	return
}

type mqtt3PubRelVariableHeader struct {
	PacketIdentifier uint16
}

func decodeMQTT3PubRelVariableHeader(data []byte) (header mqtt3PubAckVariableHeader, err error) {
	header.PacketIdentifier = binary.BigEndian.Uint16(data)

	return
}
