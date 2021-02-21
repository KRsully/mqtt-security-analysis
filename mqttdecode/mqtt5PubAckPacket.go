package mqttdecode

import (
	"encoding/binary"

	"github.com/google/gopacket"
)

var MQTT5PubAckPacket = gopacket.RegisterLayerType(
	3898,
	gopacket.LayerTypeMetadata{Name: "MQTT 5.0 PUBACK Packet", Decoder: gopacket.DecodeFunc(DecodeMQTT5PubAckPacket)})

type mqtt5PubAckPacket struct {
	VariableHeader mqtt5PubAckVariableHeader
	Contents       []byte
}

func (layer mqtt5PubAckPacket) LayerType() gopacket.LayerType { return MQTT5PubAckPacket }
func (layer mqtt5PubAckPacket) LayerContents() []byte         { return layer.Contents }
func (layer mqtt5PubAckPacket) LayerPayload() []byte          { return nil }

func DecodeMQTT5PubAckPacket(data []byte, packet gopacket.PacketBuilder) (err error) {
	variableHeader, err := decodeMQTT5PubAckVariableHeader(data)

	packet.AddLayer(&mqtt5PubAckPacket{variableHeader, data})
	return
}

type mqtt5PubAckVariableHeader struct {
	PacketIdentifier uint16
	ReasonCode       byte
	ReasonCodeName   string
	Properties       []MQTT5Property
}

func decodeMQTT5PubAckVariableHeader(data []byte) (header mqtt5PubAckVariableHeader, err error) {
	header.PacketIdentifier = binary.BigEndian.Uint16(data)
	header.ReasonCode = data[2]
	header.ReasonCodeName = resolveReasonCode(header.ReasonCode)
	header.Properties, _ = extractMQTT5Properties(data[3:])
	return
}
