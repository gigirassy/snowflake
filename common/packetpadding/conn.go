package packetpadding

import (
	"io"
)

type ReadWriteCloserPreservesBoundary interface {
	io.ReadWriteCloser
	MessageBoundaryPreserved()
}

type PaddableConnection interface {
	ReadWriteCloserPreservesBoundary
}

func NewPaddableConnection(rwc ReadWriteCloserPreservesBoundary, padding PacketPaddingContainer) PaddableConnection {
	return &paddableConnection{
		ReadWriteCloserPreservesBoundary: rwc,
		padding:                          padding,
	}
}

type paddableConnection struct {
	ReadWriteCloserPreservesBoundary
	padding PacketPaddingContainer
}

func (c *paddableConnection) Write(p []byte) (n int, err error) {
	dataLen := len(p)
	if _, err = c.ReadWriteCloserPreservesBoundary.Write(c.padding.Pack(p, 0)); err != nil {
		return 0, err
	}
	return dataLen, nil
}

func (c *paddableConnection) Read(p []byte) (n int, err error) {
	if n, err = c.ReadWriteCloserPreservesBoundary.Read(p); err != nil {
		return 0, err
	}

	payload, _ := c.padding.Unpack(p[:n])
	if payload != nil {
		copy(p, payload)
	}
	return len(payload), nil
}
