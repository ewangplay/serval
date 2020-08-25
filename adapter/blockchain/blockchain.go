package blockchain

// Constants definition
const (
	BlockChainKey = "BlockChain"
)

// BlockChain represents the block chain interface
type BlockChain interface {
	Submit(string, ...string) ([]byte, error)
	Evaluate(string, ...string) ([]byte, error)
}

// Global block chain instance
var gBlockChain BlockChain

// GetBlockChain returns the block chain instance in singleton mode
func GetBlockChain() BlockChain {
	if gBlockChain == nil {
		var err error
		gBlockChain, err = CreateHLFabricBlockChain()
		if err != nil {
			panic(err)
		}
	}
	return gBlockChain
}
