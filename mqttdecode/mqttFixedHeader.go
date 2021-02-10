package mqttdecode

import (
	"errors"

	"github.com/google/gopacket"
)

var MQTTFixedHeaderLayer = gopacket.RegisterLayerType(
	3883,
	gopacket.LayerTypeMetadata{Name: "MQTT Fixed Header", Decoder: gopacket.DecodeFunc(decodeMQTTFixedHeader)})

type mqttFixedHeader struct {
	ControlPacketType string
	Flags             byte
	RemainingLength   int
	Contents          []byte
	Payload           []byte
}

func (layer mqttFixedHeader) LayerType() gopacket.LayerType { return MQTTFixedHeaderLayer }
func (layer mqttFixedHeader) LayerContents() []byte         { return layer.Contents }
func (layer mqttFixedHeader) LayerPayload() []byte          { return nil }

func decodeControlPacketType(header byte) controlPacketType {
	//MQTT control packet type is determined by the value of the 4 highest bits of the packet's first byte
	return controlPacketType((header & 0xF0) >> 4)
}

func decodeMQTTFixedHeader(data []byte, packet gopacket.PacketBuilder) (err error) {
	remainingLength, err := decodeVariableByteInteger(data[1:])
	ctlPacketType := decodeControlPacketType(data[0])

	if err != nil {
		return err
	}

	if ctlPacketType < 1 || ctlPacketType > 15 {
		return errors.New("Invalid Control Type Error")
	}

	packet.AddLayer(&mqttFixedHeader{decodeControlPacketType(data[0]).String(),
		data[0] & 0xF, remainingLength, data[0 : len(data)-remainingLength], data[len(data)-remainingLength:]})
	//log.Printf("Payload?: %v", data[len(data)-remainingLength:])
	switch ctlPacketType {
	case 1:
		DecodeMQTT3ConnectPacket(data[len(data)-remainingLength:], packet)
	case 2:
		DecodeMQTT3ConnAckPacket(data[len(data)-remainingLength:], packet)
	default:
	}

	return nil

}
