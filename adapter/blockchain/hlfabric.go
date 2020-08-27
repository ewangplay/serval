package blockchain

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"github.com/spf13/viper"
)

// package-level variables definition
var (
	ChannelName    = "mychannel"
	ContractID     = "did"
	MSPID          = "Org1MSP"
	CCPPath        = "connection-org1.yaml"
	WalletPath     = "wallet"
	AppUserName    = "appUser"
	AppUserMSPPath = "msp"
)

func init() {
	ChannelName = viper.GetString("blockchain.hlfabric.channelName")
	ContractID = viper.GetString("blockchain.hlfabric.contractID")
	MSPID = viper.GetString("blockchain.hlfabric.mspID")
	WalletPath = viper.GetString("blockchain.hlfabric.walletPath")
	CCPPath = viper.GetString("blockchain.hlfabric.ccpPath")
	AppUserName = viper.GetString("blockchain.hlfabric.appUser.Name")
	AppUserMSPPath = viper.GetString("blockchain.hlfabric.appUser.mspPath")
}

// HLFabricBlockChain represents hyperledger fabric blockchain
type HLFabricBlockChain struct {
	wallet   *gateway.Wallet
	gateway  *gateway.Gateway
	network  *gateway.Network
	contract *gateway.Contract
}

// CreateHLFabricBlockChain creates an instance of ed25519 crypto hub
func CreateHLFabricBlockChain() (*HLFabricBlockChain, error) {

	wallet, err := gateway.NewFileSystemWallet(WalletPath)
	if err != nil {
		fmt.Printf("Failed to create wallet: %s\n", err)
		return nil, err
	}

	if !wallet.Exists(AppUserName) {
		err = populateWallet(wallet)
		if err != nil {
			fmt.Printf("Failed to populate wallet contents: %s\n", err)
			return nil, err
		}
	}

	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(CCPPath))),
		gateway.WithIdentity(wallet, AppUserName),
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
	// read the certificate pem
	certPath := filepath.Join(AppUserMSPPath, "signcerts", "cert.pem")
	cert, err := ioutil.ReadFile(filepath.Clean(certPath))
	if err != nil {
		return err
	}

	keyDir := filepath.Join(AppUserMSPPath, "keystore")
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

	err = wallet.Put(AppUserName, identity)
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
