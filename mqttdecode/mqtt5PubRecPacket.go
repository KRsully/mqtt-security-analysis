package mqttdecode

import (
	"encoding/binary"

	"github.com/google/gopacket"
)

var MQTT5PubRecPacket = gopacket.RegisterLayerType(
	3899,
	gopacket.LayerTypeMetadata{Name: "MQTT 5.0 PUBREC Packet", Decoder: gopacket.DecodeFunc(DecodeMQTT5PubRecPacket)})

type mqtt5PubRecPacket struct {
	VariableHeader mqtt5PubRecVariableHeader
	Contents       []byte
}

func (layer mqtt5PubRecPacket) LayerType() gopacket.LayerType { return MQTT5PubRecPacket }
func (layer mqtt5PubRecPacket) LayerContents() []byte         { return layer.Contents }
func (layer mqtt5PubRecPacket) LayerPayload() []byte          { return nil }

func DecodeMQTT5PubRecPacket(data []byte, packet gopacket.PacketBuilder) (err error) {
	variableHeader, err := decodeMQTT5PubRecVariableHeader(data)

	packet.AddLayer(&mqtt5PubRecPacket{variableHeader, data})
	return
}

type mqtt5PubRecVariableHeader struct {
	PacketIdentifier uint16
	ReasonCode       byte
	ReasonCodeName   string
	Properties       []MQTT5Property
}

func decodeMQTT5PubRecVariableHeader(data []byte) (header mqtt5PubRecVariableHeader, err error) {
	header.PacketIdentifier = binary.BigEndian.Uint16(data)
	header.ReasonCode = data[2]
	header.ReasonCodeName = resolveReasonCode(header.ReasonCode)
	header.Properties, _ = extractMQTT5Properties(data[3:])
	return
}
