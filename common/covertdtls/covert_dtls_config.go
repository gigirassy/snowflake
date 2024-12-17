package covertdtls

import (
	"github.com/theodorsm/covert-dtls/pkg/fingerprints"
	"strings"
)

type CovertDTLSConfig struct {
	Randomize   bool
	Mimic       bool
	Fingerprint fingerprints.ClientHelloFingerprint
}

func ParseConfigString(str string) CovertDTLSConfig {
	config := CovertDTLSConfig{}
	str = strings.ToLower(str)
	if strings.Contains(str, "random") {
		config.Randomize = true
	}
	if strings.Contains(str, "mimic") {
		config.Mimic = true
	}
	return config
}
