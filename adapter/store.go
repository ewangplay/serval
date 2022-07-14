package adapter

import (
	"fmt"

	"github.com/ewangplay/gokv/hlfabric"
	"github.com/philippgille/gokv"
	"github.com/philippgille/gokv/badgerdb"
)

type StoreOptions struct {
	Backend  string
	Badgerdb *badgerdb.Options
	Hlfabric *hlfabric.Options
}

// InitStore initializes the store instance with singleton mode
func InitStore(opts *StoreOptions) (store gokv.Store, err error) {
	switch opts.Backend {
	case "badgerdb":
		if opts.Badgerdb == nil {
			return nil, fmt.Errorf("badgerdb backend options invalid")
		}
		fmt.Println("badger options:", *opts.Badgerdb)
		store, err = badgerdb.NewStore(*opts.Badgerdb)
	case "hlfabric":
		if opts.Hlfabric == nil {
			return nil, fmt.Errorf("hlfabric backend options invalid")
		}
		fmt.Println("hlfabric options:", *opts.Hlfabric)
		store, err = hlfabric.NewClient(*opts.Hlfabric)
	default:
		err = fmt.Errorf("backend not supported: %v", opts.Backend)
	}

	return
}
