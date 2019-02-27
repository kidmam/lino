package app

import (
	// "context"
	"reflect"
	// "time"
	"unsafe"
	// "fmt"

	"github.com/cosmos/cosmos-sdk/store/cachekv"
	"github.com/cosmos/cosmos-sdk/store/gaskv"
	"github.com/cosmos/cosmos-sdk/store/iavl"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// commit internal store by reflecting to the inner store.
// panic if failed.
func commitCacheByRefl(ctx sdk.Context, key *sdk.KVStoreKey) sdk.CommitStore {
	// gas to cache
	gasKvVal := reflect.ValueOf(ctx.KVStore(key).(*gaskv.Store)).Elem()
	gasParentVal := gasKvVal.FieldByName("parent")
	// convert from reflect.Value, which is an exported field, to a struct pointer.
	// though using unsafe, it's actually safe.
	cacheKv := reflect.NewAt(gasParentVal.Type(),
		unsafe.Pointer(gasParentVal.UnsafeAddr())).Elem().Interface().(*cachekv.Store)
	cacheKv.Write()

	// cache to iavl
	cacheKvVal := reflect.ValueOf(cacheKv).Elem()
	iavlVal := cacheKvVal.FieldByName("parent")
	iavlKv := reflect.NewAt(iavlVal.Type(),
		unsafe.Pointer(iavlVal.UnsafeAddr())).Elem().Interface().(*iavl.Store)
	iavlKv.Commit()
	return iavlKv
}

// This won't help to minimize the memory footprint.
// func makeStoreCleaner(ctx context.Context, cacheStore *cachekv.Store) {
// 	ticker := time.NewTicker(500 * time.Millisecond)
// 	go func() {
// 		for range ticker.C {
// 			cacheStore.Write()
// 		}
// 	}()
// 	<-ctx.Done()
// 	ticker.Stop()
// }
