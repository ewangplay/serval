package blockchain

import (
	"errors"
)

// BlockChain represents the blockchain interface
type BlockChain interface {
	Submit(string, ...string) ([]byte, error)
	Evaluate(string, ...string) ([]byte, error)
}

// Global blockchain instance
var gBC BlockChain

// Config defines BlockChain config
type Config struct {
	HLFabric *HLFabricConfig
}

// InitBlockChain initializes the blockchain instance with singleton mode
// This method MUST be called before any other method can be called.
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
// MUST call the InitBlockChain method before you can call this method,
// otherwise it will cause a panic of "blockchain not initialized".
func Submit(fn string, args ...string) ([]byte, error) {
	checkInitState()
	return gBC.Submit(fn, args...)
}

// Evaluate will evaluate a transaction function and return its results
// MUST call the InitBlockChain method before you can call this method,
// otherwise it will cause a panic of "blockchain not initialized".
func Evaluate(fn string, args ...string) ([]byte, error) {
	checkInitState()
	return gBC.Evaluate(fn, args...)
}

func checkInitState() {
	if gBC == nil {
		panic("blockchain not initialized")
	}
}
