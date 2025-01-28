package packetPaddingContainer

import "encoding/binary"

func New() PacketPaddingContainer {
	return packetPaddingContainer{}
}

type packetPaddingContainer struct {
}

func (c packetPaddingContainer) Pack(data_OWNERSHIP_RELINQUISHED []byte, padding int) []byte {
	data := append(data_OWNERSHIP_RELINQUISHED, make([]byte, padding)...)
	dataLength := len(data_OWNERSHIP_RELINQUISHED)
	data = binary.BigEndian.AppendUint16(data, uint16(dataLength))
	return data
}

func (c packetPaddingContainer) Pad(padding int) []byte {
	if assertPaddingLengthIsNotNegative := padding < 0; assertPaddingLengthIsNotNegative {
		return nil
	}
	switch padding {
	case 0:
		return []byte{}
	case 1:
		return []byte{0}
	case 2:
		return []byte{0, 0}
	default:
		return append(make([]byte, padding-2), byte(padding>>8), byte(padding))
	}

}

func (c packetPaddingContainer) Unpack(wrappedData_OWNERSHIP_RELINQUISHED []byte) ([]byte, int) {
	dataLength := len(wrappedData_OWNERSHIP_RELINQUISHED)
	if dataLength < 2 {
		return nil, dataLength
	}

	dataLen := int(binary.BigEndian.Uint16(wrappedData_OWNERSHIP_RELINQUISHED[dataLength-2:]))
	paddingLength := dataLength - dataLen - 2
	if paddingLength < 0 {
		return nil, paddingLength
	}

	return wrappedData_OWNERSHIP_RELINQUISHED[:dataLen], paddingLength
}
