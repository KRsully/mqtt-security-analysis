package mqttdecode

import (
	"encoding/binary"
	"errors"
)

type controlPacketType int

const (
	CONNECT = controlPacketType(iota + 1)
	CONNACK
	PUBLISH
	PUBACK
	PUBREC
	PUBREL
	PUBCOMP
	SUBSCRIBE
	SUBACK
	UNSUBSCRIBE
	UNSUBACK
	PINGREQ
	PINGRESP
	DISCONNECT
	AUTH
)

// type mqttPacket struct {
// 	FixedHeader mqttFixedHeader
// 	VariableHeader mqttVariableHeader
// 	Payload mqttPayload
// }

func (cpt controlPacketType) String() string {
	return [...]string{"", "CONNECT", "CONNACK", "PUBLISH", "PUBACK", "PUBREC", "PUBREL",
		"PUBCOMP", "SUBSCRIBE", "SUBACK", "UNSUBSCRIBE", "UNSUBACK", "PINGREQ", "PINGRESP", "DISCONNECT", "AUTH"}[cpt]
}

func decodeVariableByteInteger(packet []byte) (remainingLength int, err error) {
	multiplier := 1
	value := 0

	for _, nextByte := range packet {
		value += int(nextByte&127) * multiplier
		if multiplier > 128*128*128 {
			err = errors.New("Malformed Variable Length Integer")
		}
		multiplier *= 128
		if nextByte&128 == 0 {
			break
		}
	}
	return value, err
}

func extractUTF8String(data []byte) (string, uint16, error) {
	stringLength := binary.BigEndian.Uint16(data)
	return string(data[2 : 2+stringLength]), stringLength, nil
}
