package mqttdecode

import (
	"encoding/binary"

	"github.com/google/gopacket"
)

var MQTT5PubCompPacket = gopacket.RegisterLayerType(
	3901,
	gopacket.LayerTypeMetadata{Name: "MQTT 5.0 PUBCOMP Packet", Decoder: gopacket.DecodeFunc(DecodeMQTT5PubCompPacket)})

type mqtt5PubCompPacket struct {
	VariableHeader mqtt5PubCompVariableHeader
	Contents       []byte
}

func (layer mqtt5PubCompPacket) LayerType() gopacket.LayerType { return MQTT5PubCompPacket }
func (layer mqtt5PubCompPacket) LayerContents() []byte         { return layer.Contents }
func (layer mqtt5PubCompPacket) LayerPayload() []byte          { return nil }

func DecodeMQTT5PubCompPacket(data []byte, packet gopacket.PacketBuilder) (err error) {
	variableHeader, err := decodeMQTT5PubCompVariableHeader(data)

	packet.AddLayer(&mqtt5PubCompPacket{variableHeader, data})
	return
}

type mqtt5PubCompVariableHeader struct {
	PacketIdentifier uint16
	ReasonCode       byte
	ReasonCodeName   string
	Properties       []MQTT5Property
}

func decodeMQTT5PubCompVariableHeader(data []byte) (header mqtt5PubCompVariableHeader, err error) {
	header.PacketIdentifier = binary.BigEndian.Uint16(data)
	header.ReasonCode = data[2]
	header.ReasonCodeName = resolveReasonCode(header.ReasonCode)
	header.Properties, _ = extractMQTT5Properties(data[3:])
	return
}
