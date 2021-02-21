package mqttdecode

import (
	"github.com/google/gopacket"
)

var MQTT5DisconnectPacket = gopacket.RegisterLayerType(
	3906,
	gopacket.LayerTypeMetadata{Name: "MQTT 5.0 DISCONNECT Packet", Decoder: gopacket.DecodeFunc(DecodeMQTT5DisconnectPacket)})

type mqtt5DisconnectPacket struct {
	VariableHeader mqtt5DisconnectVariableHeader
	Contents       []byte
}

func (layer mqtt5DisconnectPacket) LayerType() gopacket.LayerType { return MQTT5DisconnectPacket }
func (layer mqtt5DisconnectPacket) LayerContents() []byte         { return layer.Contents }
func (layer mqtt5DisconnectPacket) LayerPayload() []byte          { return nil }

func DecodeMQTT5DisconnectPacket(data []byte, packet gopacket.PacketBuilder) (err error) {
	variableHeader, err := decodeMQTT5DisconnectVariableHeader(data)

	packet.AddLayer(&mqtt5DisconnectPacket{variableHeader, data})
	return
}

type mqtt5DisconnectVariableHeader struct {
	ReasonCode     byte
	ReasonCodeName string
	Properties     []MQTT5Property
}

func decodeMQTT5DisconnectVariableHeader(data []byte) (header mqtt5DisconnectVariableHeader, err error) {
	header.ReasonCode = data[0]
	header.ReasonCodeName = resolveReasonCode(header.ReasonCode)
	header.Properties, _ = extractMQTT5Properties(data[1:])
	return
}
