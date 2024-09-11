package snowflake_client

import (
	"errors"
	"io"
	"net"
	"time"
)

type ReadWriteCloserPreservesBoundary interface {
	io.ReadWriteCloser
	MessageBoundaryPreserved()
}

func ConfirmsReadWriteCloserPreservesMessageBoundary(rwc io.ReadWriteCloser) ReadWriteCloserPreservesBoundary {
	return &messageBoundaryPreservedReadWriteCloser{rwc}
}

type messageBoundaryPreservedReadWriteCloser struct {
	io.ReadWriteCloser
}

func (m *messageBoundaryPreservedReadWriteCloser) MessageBoundaryPreserved() {}

var errENOSYS = errors.New("not implemented")

func newPacketConnWrapper(localAddr, remoteAddr net.Addr, rwc ReadWriteCloserPreservesBoundary) net.PacketConn {
	return &packetConnWrapper{
		ReadWriteCloserPreservesBoundary: rwc,
		remoteAddr:                       remoteAddr,
		localAddr:                        localAddr,
	}
}

type packetConnWrapper struct {
	ReadWriteCloserPreservesBoundary
	remoteAddr net.Addr
	localAddr  net.Addr
}

func (pcw *packetConnWrapper) ReadFrom(p []byte) (n int, addr net.Addr, err error) {
	n, err = pcw.Read(p)
	if err != nil {
		return 0, nil, err
	}
	return n, pcw.remoteAddr, nil
}

func (pcw *packetConnWrapper) WriteTo(p []byte, addr net.Addr) (n int, err error) {
	return pcw.Write(p)
}

func (pcw *packetConnWrapper) Close() error {
	return pcw.ReadWriteCloser.Close()
}

func (pcw *packetConnWrapper) LocalAddr() net.Addr {
	return pcw.localAddr
}

func (pcw *packetConnWrapper) SetDeadline(t time.Time) error {
	return errENOSYS
}

func (pcw *packetConnWrapper) SetReadDeadline(t time.Time) error {
	return errENOSYS
}

func (pcw *packetConnWrapper) SetWriteDeadline(t time.Time) error {
	return errENOSYS
}
