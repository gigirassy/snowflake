package media

import (
	"crypto/rand"
	"log"
	"math/big"
	"time"

	"github.com/pion/interceptor"
	"github.com/pion/webrtc/v4"
	"github.com/pion/webrtc/v4/pkg/media"
)

func randomInt(min, max int) int {
	nBig, err := rand.Int(rand.Reader, big.NewInt(int64(max-min)))
	if err != nil {
		panic(err)
	}
	return int(nBig.Int64()) + min
}

type RTPReader interface {
	Read(b []byte) (n int, a interceptor.Attributes, err error)
}

// MediaChannel handles media track simulation for WebRTC connections
type MediaChannel struct {
	stopCh chan struct{}
}

// NewMediaChannel creates a new media channel
func NewMediaChannel() *MediaChannel {
	return &MediaChannel{
		stopCh: make(chan struct{}),
	}
}

// StartVideoTrack starts video track simulation on the given peer connection
func (mc *MediaChannel) StartVideoTrack(pc *webrtc.PeerConnection) error {
	videoTrack, err := webrtc.NewTrackLocalStaticSample(
		webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeAV1}, "video", "pion",
	)
	if err != nil {
		log.Printf("webrtc.NewTrackLocalStaticSample ERROR: %s", err)
		return err
	}

	rtpSender, err := pc.AddTrack(videoTrack)
	if err != nil {
		log.Printf("webrtc.AddTrack ERROR: %s", err)
		return err
	}

	go mc.handleRTCP(rtpSender)
	go mc.simulateVideoFrames(videoTrack)

	log.Println("WebRTC: Media track opened")
	return nil
}

// Stop stops the media simulation
func (mc *MediaChannel) Stop() {
	close(mc.stopCh)
}

func (mc *MediaChannel) handleRTCP(reader RTPReader) {
	rtcpBuf := make([]byte, 1500)
	for {
		select {
		case <-mc.stopCh:
			return
		default:
			if _, _, err := reader.Read(rtcpBuf); err != nil {
				return
			}
		}
	}
}

func (mc *MediaChannel) simulateVideoFrames(track *webrtc.TrackLocalStaticSample) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-mc.stopCh:
			return
		case <-ticker.C:
			// Add jitter to simulate "realistic" media patterns
			jitterDelay := time.Duration(randomInt(0, 200)) * time.Millisecond
			time.Sleep(jitterDelay)

			// Vary packet sizes for specific frames types
			var bufSize int
			frameType := randomInt(1, 100)
			switch {
			case frameType <= 5: // I-frames: 5% chance, larger
				bufSize = randomInt(8000, 15000)
			case frameType <= 35: // P-frames: 30% chance, medium
				bufSize = randomInt(2000, 5000)
			default: // B-frames: 65% chance, smaller
				bufSize = randomInt(500, 2000)
			}

			buf := make([]byte, bufSize)

			// Add some timing variation
			frameDuration := time.Duration(randomInt(900, 1100)) * time.Millisecond

			err := track.WriteSample(media.Sample{Data: buf, Duration: frameDuration})
			if err != nil {
				log.Printf("webrtc.WriteSample ERROR: %s", err)
			}

			// Simulate some burst of smaller packets
			if randomInt(1, 10) == 1 { // 10% chance
				burstCount := randomInt(2, 5)
				for i := 0; i < burstCount; i++ {
					smallBuf := make([]byte, randomInt(100, 400))
					time.Sleep(time.Duration(randomInt(10, 50)) * time.Millisecond)

					frameDuration = time.Duration(randomInt(16, 33)) * time.Millisecond

					err = track.WriteSample(media.Sample{Data: smallBuf, Duration: frameDuration})
					if err != nil {
						log.Printf("webrtc.WriteSample burst ERROR: %s", err)
						break
					}
				}
			}
		}
	}
}

// Start sets up duplex media handling (both incoming and outgoing tracks)
func (mc *MediaChannel) Start(pc *webrtc.PeerConnection) error {
	// Set up handler for incoming tracks
	pc.OnTrack(func(remote *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		log.Printf("Media Track received: streamId(%s) id(%s) rid(%s)", remote.StreamID(), remote.ID(), remote.RID())
		go mc.handleRTCP(receiver)
	})

	// Set up outgoing media track
	return mc.StartVideoTrack(pc)
}
