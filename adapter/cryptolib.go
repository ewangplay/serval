package adapter

import (
	cl "github.com/ewangplay/cryptolib"
)

// InitCryptolib initializes the cryptolib instance
func InitCryptolib() (csp cl.CSP, err error) {
	cfg := &cl.Config{
		ProviderName: "SW",
	}
	return cl.GetCSP(cfg)
}
