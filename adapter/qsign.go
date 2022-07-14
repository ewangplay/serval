package adapter

import (
	"crypto/sha256"
	"encoding/base64"
	"hash"

	"github.com/jerray/qsign"
)

// InitQsign initializes the Qsign instance
func InitQsign() (qs *qsign.Qsign, err error) {
	// New Qsign instance
	qs = qsign.NewQsign(qsign.Options{
		// To use a hash.Hash other than md5
		Hasher: func() hash.Hash {
			return sha256.New()
		},
		// To use a encoding other than hex
		Encoder: func() qsign.Encoding {
			return base64.StdEncoding
		},
	})

	return qs, nil
}
