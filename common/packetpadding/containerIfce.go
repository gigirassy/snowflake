package packetpadding

// PacketPaddingContainer is an interface that defines methods to pad packets
// with a given number of bytes, and to unpack the padding from a padded packet.
// The packet format is as follows if the desired output length is greater than
// 2 bytes:
// | data | padding | data length |
// The data length is a 16-bit big-endian integer that represents the length of
// the data in bytes.
// If the desired output length is 2 bytes or less, the packet format is as
// follows:
// | padding |
// No payload will be included in the packet.
type PacketPaddingContainer interface {
	// Pack pads the given data with the given number of bytes, and appends the
	// length of the data to the end of the data. The returned byte slice
	// contains the padded data.
	// This generates a packet with a length of
	// len(data_OWNERSHIP_RELINQUISHED) + padding + 2
	// @param data_OWNERSHIP_RELINQUISHED - The payload, this reference is consumed and should not be used after this call.
	// @param padding - The number of padding bytes to add to the data.
	Pack(data_OWNERSHIP_RELINQUISHED []byte, paddingLength int) []byte

	// Unpack extracts the data and padding from the given padded data. It
	// returns the data and the number of padding bytes.
	// the data may be nil.
	// @param wrappedData_OWNERSHIP_RELINQUISHED - The packet, this reference is consumed and should not be used after this call.
	Unpack(wrappedData_OWNERSHIP_RELINQUISHED []byte) ([]byte, int)

	// Pad returns a padding packet of padding length.
	// If the padding length is less than 0, nil is returned.
	// @param padding - The number of padding bytes to add to the data.
	Pad(paddingLength int) []byte
}
