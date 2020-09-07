package main

import (
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing DIDs
type SmartContract struct {
	contractapi.Contract
}

// InitLedger ...
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	return nil
}

// CreateDID adds a new did / ddo record to the world state
func (s *SmartContract) CreateDID(ctx contractapi.TransactionContextInterface, did string, ddo string) error {
	if did == "" {
		return fmt.Errorf("input param did is missing")
	}
	if len(ddo) == 0 {
		return fmt.Errorf("input param ddo is missing")
	}
	return ctx.GetStub().PutState(did, []byte(ddo))
}

// QueryDID returns the ddo stored in the world state with given did
func (s *SmartContract) QueryDID(ctx contractapi.TransactionContextInterface, did string) ([]byte, error) {
	if did == "" {
		return nil, fmt.Errorf("input param did is missing")
	}
	ddo, err := ctx.GetStub().GetState(did)
	if err != nil {
		return nil, err
	}
	if ddo == nil {
		return nil, fmt.Errorf("%s does not exist", did)
	}
	return ddo, nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(new(SmartContract))
	if err != nil {
		fmt.Printf("Error create chaincode: %v\n", err)
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting chaincode: %v\n", err)
	}
}
