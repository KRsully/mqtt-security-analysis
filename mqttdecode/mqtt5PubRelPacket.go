package mqttdecode

import (
	"encoding/binary"

	"github.com/google/gopacket"
)

var MQTT5PubRelPacket = gopacket.RegisterLayerType(
	3900,
	gopacket.LayerTypeMetadata{Name: "MQTT 5.0 PUBREL Packet", Decoder: gopacket.DecodeFunc(DecodeMQTT5PubRelPacket)})

type mqtt5PubRelPacket struct {
	VariableHeader mqtt5PubRelVariableHeader
	Contents       []byte
}

func (layer mqtt5PubRelPacket) LayerType() gopacket.LayerType { return MQTT5PubRelPacket }
func (layer mqtt5PubRelPacket) LayerContents() []byte         { return layer.Contents }
func (layer mqtt5PubRelPacket) LayerPayload() []byte          { return nil }

func DecodeMQTT5PubRelPacket(data []byte, packet gopacket.PacketBuilder) (err error) {
	variableHeader, err := decodeMQTT5PubRelVariableHeader(data)

	packet.AddLayer(&mqtt5PubRelPacket{variableHeader, data})
	return
}

type mqtt5PubRelVariableHeader struct {
	PacketIdentifier uint16
	ReasonCode       byte
	ReasonCodeName   string
	Properties       []MQTT5Property
}

func decodeMQTT5PubRelVariableHeader(data []byte) (header mqtt5PubRelVariableHeader, err error) {
	header.PacketIdentifier = binary.BigEndian.Uint16(data)
	header.ReasonCode = data[2]
	header.ReasonCodeName = resolveReasonCode(header.ReasonCode)
	header.Properties, _ = extractMQTT5Properties(data[3:])
	return
}
