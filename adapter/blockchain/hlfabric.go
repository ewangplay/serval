package blockchain

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

// AppUser defines hlfabric app user
type AppUser struct {
	Name    string
	MspPath string
}

// HLFabricConfig defines hlfabric config
type HLFabricConfig struct {
	ChannelName    string
	ContractID     string
	MspID          string
	WalletPath     string
	CcpPath        string
	AppUser        AppUser
	EndorsingPeers []string
}

// HLFabric represents hyperledger fabric blockchain
type HLFabric struct {
	config   *HLFabricConfig
	wallet   *gateway.Wallet
	gateway  *gateway.Gateway
	network  *gateway.Network
	contract *gateway.Contract
}

// CreateHLFabric creates an instance of ed25519 crypto hub
func CreateHLFabric(cfg *HLFabricConfig) (*HLFabric, error) {

	wallet, err := gateway.NewFileSystemWallet(cfg.WalletPath)
	if err != nil {
		fmt.Printf("Failed to create wallet: %s\n", err)
		return nil, err
	}

	if !wallet.Exists(cfg.AppUser.Name) {
		err = populateWallet(wallet, cfg.MspID, &cfg.AppUser)
		if err != nil {
			fmt.Printf("Failed to populate wallet contents: %s\n", err)
			return nil, err
		}
	}

	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(cfg.CcpPath))),
		gateway.WithIdentity(wallet, cfg.AppUser.Name),
	)
	if err != nil {
		fmt.Printf("Failed to connect to gateway: %s\n", err)
		return nil, err
	}

	network, err := gw.GetNetwork(cfg.ChannelName)
	if err != nil {
		fmt.Printf("Failed to get network: %s\n", err)
		return nil, err
	}

	contract := network.GetContract(cfg.ContractID)

	hlf := &HLFabric{
		config:   cfg,
		wallet:   wallet,
		gateway:  gw,
		network:  network,
		contract: contract,
	}

	return hlf, nil
}

func populateWallet(wallet *gateway.Wallet, mspID string, appUser *AppUser) error {
	// read the certificate pem
	certPath := filepath.Join(appUser.MspPath, "signcerts", "cert.pem")
	cert, err := ioutil.ReadFile(filepath.Clean(certPath))
	if err != nil {
		return err
	}

	keyDir := filepath.Join(appUser.MspPath, "keystore")
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

	identity := gateway.NewX509Identity(mspID, string(cert), string(key))

	err = wallet.Put(appUser.Name, identity)
	if err != nil {
		return err
	}
	return nil
}

// Submit will submit a transaction to the ledger
func (hlf *HLFabric) Submit(fn string, args ...string) ([]byte, error) {
	txn, err := hlf.contract.CreateTransaction(fn,
		gateway.WithEndorsingPeers(hlf.config.EndorsingPeers...))
	if err != nil {
		return nil, err
	}
	return txn.Submit(args...)
}

// Evaluate will evaluate a transaction function and return its results
func (hlf *HLFabric) Evaluate(fn string, args ...string) ([]byte, error) {
	return hlf.contract.EvaluateTransaction(fn, args...)
}
