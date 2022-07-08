package store

import (
	"fmt"
	"sync"

	"github.com/ewangplay/gokv/hlfabric"
	"github.com/philippgille/gokv"
	"github.com/philippgille/gokv/badgerdb"
)

// Global instance definition
var (
	initOnce sync.Once
	gStore   gokv.Store
)

type Options struct {
	Backend  string
	Badgerdb *badgerdb.Options
	Hlfabric *hlfabric.Options
}

// InitStore initializes the store instance with singleton mode
func InitStore(opts *Options) error {
	var err error

	initOnce.Do(func() {
		err = initStore(opts)
	})

	return err
}

func initStore(opts *Options) (err error) {
	switch opts.Backend {
	case "badgerdb":
		if opts.Badgerdb == nil {
			return fmt.Errorf("badgerdb backend options invalid")
		}
		fmt.Println("badger options:", *opts.Badgerdb)
		gStore, err = badgerdb.NewStore(*opts.Badgerdb)
	case "hlfabric":
		if opts.Hlfabric == nil {
			return fmt.Errorf("hlfabric backend options invalid")
		}
		fmt.Println("hlfabric options:", *opts.Hlfabric)
		gStore, err = hlfabric.NewClient(*opts.Hlfabric)
	default:
		err = fmt.Errorf("backend not supported: %v", opts.Backend)
	}

	if err != nil {
		return err
	}
	return nil
}

func ReleaseStore() {
	if gStore != nil {
		gStore.Close()
	}
}

func Set(k string, v any) error {
	assertStoreValid()
	return gStore.Set(k, v)

}

func Get(k string, v any) (found bool, err error) {
	assertStoreValid()
	return gStore.Get(k, v)

}

func Delete(k string) error {
	assertStoreValid()
	return gStore.Delete(k)
}

func assertStoreValid() {
	if gStore == nil {
		panic("Store not be initialized")
	}
}
