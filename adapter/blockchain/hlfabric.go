package blockchain

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

// Constants definition
const (
	ConfigPath   = "."
	IdentityName = "appUser"
	ChannelName  = "mychannel"
	ContractID   = "did"
	MSPID        = "Org1MSP"
)

// HLFabricBlockChain represents hyperledger fabric blockchain
type HLFabricBlockChain struct {
	wallet   *gateway.Wallet
	gateway  *gateway.Gateway
	network  *gateway.Network
	contract *gateway.Contract
}

// CreateHLFabricBlockChain creates an instance of ed25519 crypto hub
func CreateHLFabricBlockChain() (*HLFabricBlockChain, error) {

	walletPath := filepath.Join(ConfigPath, "wallet")
	wallet, err := gateway.NewFileSystemWallet(walletPath)
	if err != nil {
		fmt.Printf("Failed to create wallet: %s\n", err)
		return nil, err
	}

	if !wallet.Exists(IdentityName) {
		err = populateWallet(wallet)
		if err != nil {
			fmt.Printf("Failed to populate wallet contents: %s\n", err)
			return nil, err
		}
	}

	ccpPath := filepath.Join(
		ConfigPath,
		"blockchain",
		"connection-org1.yaml",
	)

	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))),
		gateway.WithIdentity(wallet, IdentityName),
	)
	if err != nil {
		fmt.Printf("Failed to connect to gateway: %s\n", err)
		return nil, err
	}

	network, err := gw.GetNetwork(ChannelName)
	if err != nil {
		fmt.Printf("Failed to get network: %s\n", err)
		return nil, err
	}

	contract := network.GetContract(ContractID)

	hlf := &HLFabricBlockChain{
		wallet:   wallet,
		gateway:  gw,
		network:  network,
		contract: contract,
	}

	return hlf, nil
}

func populateWallet(wallet *gateway.Wallet) error {
	credPath := filepath.Join(
		ConfigPath,
		"blockchain",
		"appUser",
		"msp",
	)

	// read the certificate pem
	certPath := filepath.Join(credPath, "signcerts", "cert.pem")
	cert, err := ioutil.ReadFile(filepath.Clean(certPath))
	if err != nil {
		return err
	}

	keyDir := filepath.Join(credPath, "keystore")
	// there's a single file in this dir containing the private key
	files, err := ioutil.ReadDir(keyDir)
	if err != nil {
		return err
	}
	if len(files) != 1 {
		return errors.New("keystore folder should have contain one file")
	}
	keyPath := filepath.Join(keyDir, files[0].Name())
	key, err := ioutil.ReadFile(filepath.Clean(keyPath))
	if err != nil {
		return err
	}

	identity := gateway.NewX509Identity(MSPID, string(cert), string(key))

	err = wallet.Put(IdentityName, identity)
	if err != nil {
		return err
	}
	return nil
}

// Submit will submit a transaction to the ledger
func (hlf *HLFabricBlockChain) Submit(fn string, args ...string) ([]byte, error) {
	return hlf.contract.SubmitTransaction(fn, args...)
}

// Evaluate will evaluate a transaction function and return its results
func (hlf *HLFabricBlockChain) Evaluate(fn string, args ...string) ([]byte, error) {
	return hlf.contract.EvaluateTransaction(fn, args...)
}
