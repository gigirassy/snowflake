package packetpadding

import "encoding/binary"

func New() PacketPaddingContainer {
	return packetPaddingContainer{}
}

type packetPaddingContainer struct {
}

func (c packetPaddingContainer) Pack(data_OWNERSHIP_RELINQUISHED []byte, paddingLength int) []byte {
	data := append(data_OWNERSHIP_RELINQUISHED, make([]byte, paddingLength)...)
	dataLength := len(data_OWNERSHIP_RELINQUISHED)
	data = binary.BigEndian.AppendUint16(data, uint16(dataLength))
	return data
}

func (c packetPaddingContainer) Pad(paddingLength int) []byte {
	if assertPaddingLengthIsNotNegative := paddingLength < 0; assertPaddingLengthIsNotNegative {
		return nil
	}
	switch paddingLength {
	case 0:
		return []byte{}
	case 1:
		return []byte{0}
	case 2:
		return []byte{0, 0}
	default:
		return append(make([]byte, paddingLength-2), byte(paddingLength>>8), byte(paddingLength))
	}

}

func (c packetPaddingContainer) Unpack(wrappedData_OWNERSHIP_RELINQUISHED []byte) ([]byte, int) {
	dataLength := len(wrappedData_OWNERSHIP_RELINQUISHED)
	if dataLength < 2 {
		return nil, dataLength
	}

	dataLen := int(binary.BigEndian.Uint16(wrappedData_OWNERSHIP_RELINQUISHED[dataLength-2:]))
	if dataLen > 2047 {
		return nil, 0
	}
	paddingLength := dataLength - dataLen - 2
	if paddingLength < 0 {
		return nil, paddingLength
	}

	return wrappedData_OWNERSHIP_RELINQUISHED[:dataLen], paddingLength
}
