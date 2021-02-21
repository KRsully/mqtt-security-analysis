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

func decodeVariableByteInteger(packet []byte) (value int, numberOfBytes int, err error) {
	multiplier := 1

	for _, nextByte := range packet {
		value += int(nextByte&127) * multiplier
		numberOfBytes += 1
		if multiplier > 128*128*128 {
			err = errors.New("Malformed Variable Length Integer")
		}
		multiplier *= 128
		if nextByte&128 == 0 {
			break
		}
	}
	return value, numberOfBytes, err
}

func extractUTF8String(data []byte) (string, int, error) {
	stringLength := binary.BigEndian.Uint16(data)
	return string(data[2 : 2+stringLength]), int(stringLength), nil
}

func resolveReasonCode(rc byte) (name string) {
	switch rc {
	case 0:
		name = "Success"
	case 1:
		name = "Granted QoS 1"
	case 2:
		name = "Granted QoS 2"
	case 4:
		name = "Disconnect with Will Message"
	case 16:
		name = "No matching subscribers"
	case 17:
		name = "No subscription existed"
	case 24:
		name = "Continue authentication"
	case 25:
		name = "Re-authenticate"
	case 128:
		name = "Unspecified error"
	case 129:
		name = "Malformed Packet"
	case 130:
		name = "Protocol Error"
	case 131:
		name = "Implementation specific error"
	case 132:
		name = "Unsupported Protocol Version"
	case 133:
		name = "Client Identifier not valid"
	case 134:
		name = "Bad User Name or Password"
	case 135:
		name = "Not authorized"
	case 136:
		name = "Server unavailable"
	case 137:
		name = "Server busy"
	case 138:
		name = "Banned"
	case 139:
		name = "Server shutting down"
	case 140:
		name = "Bad authentication method"
	case 141:
		name = "Keep Alive timeout"
	case 142:
		name = "Session taken over"
	case 143:
		name = "Topic Filter invalid"
	case 144:
		name = "Topic Name invalid"
	case 145:
		name = "Packet identifier in use"
	case 146:
		name = "Packet Identifier not found"
	case 147:
		name = "Receive Maximum exceeded"
	case 148:
		name = "Topic Alias invalid"
	case 149:
		name = "Packet too large"
	case 150:
		name = "Message rate too high"
	case 151:
		name = "Quota exceeded"
	case 152:
		name = "Administrative action"
	case 153:
		name = "Payload format invalid"
	case 154:
		name = "Retain not supported"
	case 155:
		name = "QoS not supported"
	case 156:
		name = "Use another server"
	case 157:
		name = "Server moved"
	case 158:
		name = "Shared Subscriptions not supported"
	case 159:
		name = "Connection rate exceeded"
	case 160:
		name = "Maximum connect time"
	case 161:
		name = "Subscription Identifiers not supported"
	case 162:
		name = "Wildcard Subscriptions not supported"
	default:
		name = "Unrecognized reason code"
	}
	return name
}
