package blockchain

import (
	"errors"
)

// BlockChain represents the block chain interface
type BlockChain interface {
	Submit(string, ...string) ([]byte, error)
	Evaluate(string, ...string) ([]byte, error)
}

// Global block chain instance
var gBC BlockChain

// Config defines BlockChain config
type Config struct {
	HLFabric *HLFabricConfig
}

// InitBlockChain initializes the block chain instance with singleton mode
// This initialization function MUST be called before any other method can be called
func InitBlockChain(cfg *Config) error {

	if gBC == nil {
		var err error

		if cfg == nil || cfg.HLFabric == nil {
			err = errors.New("missing blockchain config")
			return err
		}

		gBC, err = CreateHLFabric(cfg.HLFabric)
		if err != nil {
			return err
		}
	}

	return nil
}

// Submit will submit a transaction to the ledger
func Submit(fn string, args ...string) ([]byte, error) {
	if gBC == nil {
		panic("blockchain not initialized")
	}
	return gBC.Submit(fn, args...)
}

// Evaluate will evaluate a transaction function and return its results
func Evaluate(fn string, args ...string) ([]byte, error) {
	if gBC == nil {
		panic("blockchain not initialized")
	}
	return gBC.Evaluate(fn, args...)
}
