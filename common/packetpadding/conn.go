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

func NewPaddableConnection(rwc ReadWriteCloserPreservesBoundary) PaddableConnection {
	return &paddableConnection{
		ReadWriteCloserPreservesBoundary: rwc,
	}
}

type paddableConnection struct {
	ReadWriteCloserPreservesBoundary
}

func (c *paddableConnection) Write(p []byte) (n int, err error) {
	dataLen := len(p)
	if _, err = c.ReadWriteCloserPreservesBoundary.Write(Pack(p, 0)); err != nil {
		return 0, err
	}
	return dataLen, nil
}

func (c *paddableConnection) Read(p []byte) (n int, err error) {
	if n, err = c.ReadWriteCloserPreservesBoundary.Read(p); err != nil {
		return 0, err
	}

	payload, _ := Unpack(p[:n])
	if payload != nil {
		copy(p, payload)
	}
	return len(payload), nil
}
