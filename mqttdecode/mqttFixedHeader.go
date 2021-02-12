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

	if ctlPacketType < 1 {
		return errors.New("Invalid Control Type Error")
	}

	packet.AddLayer(&mqttFixedHeader{decodeControlPacketType(data[0]).String(),
		data[0] & 0xF, remainingLength, data[0 : len(data)-remainingLength], data[len(data)-remainingLength:]})
	switch ctlPacketType {
	case 1:
		err = DecodeMQTT3ConnectPacket(data[len(data)-remainingLength:], packet)
		if err != nil {
			//DecodeMQTT5
		}
	case 2:
		err = DecodeMQTT3ConnAckPacket(data[len(data)-remainingLength:], packet)
		if err != nil {
			//DecodeMQTT5
		}
	case 3:
		err = DecodeMQTT3PublishPacket(data[len(data)-remainingLength:], packet)
		if err != nil {
			//DecodeMQTT5
		}
	case 4:
		err = DecodeMQTT3PubAckPacket(data[len(data)-remainingLength:], packet)
		if err != nil {
			//DecodeMQTT5
		}
	case 5:
		err = DecodeMQTT3PubRecPacket(data[len(data)-remainingLength:], packet)
		if err != nil {
			//DecodeMQTT5
		}
	case 6:
		err = DecodeMQTT3PubRelPacket(data[len(data)-remainingLength:], packet)
		if err != nil {
			//DecodeMQTT5
		}
	case 7:
		err = DecodeMQTT3PubCompPacket(data[len(data)-remainingLength:], packet)
		if err != nil {
			//DecodeMQTT5
		}
	case 8:
		err = DecodeMQTT3SubscribePacket(data[len(data)-remainingLength:], packet)
		if err != nil {
			//DecodeMQTT5
		}
	case 9:
		err = DecodeMQTT3SubAckPacket(data[len(data)-remainingLength:], packet)
		if err != nil {
			//DecodeMQTT5
		}
	case 10:
		err = DecodeMQTT3UnsubscribePacket(data[len(data)-remainingLength:], packet)
		if err != nil {
			//DecodeMQTT5
		}
	case 11:
		err = DecodeMQTT3UnsubAckPacket(data[len(data)-remainingLength:], packet)
		if err != nil {
			//DecodeMQTT5
		}
	case 12:
		//PINGREQ Packet
	case 13:
		//PINGRESP Packet
	case 14:
		//DISCONNECT Packet
	case 15:
		//AUTH Packet (MQTT 5.0 only)
	default:
	}

	return nil

}
