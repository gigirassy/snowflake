package media

import (
	"testing"
)

func TestMediaChannelStop(t *testing.T) {
	mc := NewMediaChannel()

	// This should not panic
	mc.Stop()

	// Verify that the stop channel is closed
	select {
	case <-mc.stopCh:
		// Channel is closed, which is expected
	default:
		t.Fatal("MediaChannel stopCh should be closed after Stop()")
	}
}
