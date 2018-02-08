package state

import (
	"fmt"
	abci "github.com/tendermint/abci/types"
	"github.com/lino-network/lino/types"
	eyes "github.com/tendermint/merkleeyes/client"
	"github.com/tendermint/go-wire"
	"github.com/tendermint/tmlibs/log"
)

// CONTRACT: State should be quick to copy.
// See CacheWrap().
type State struct {
	chainID    string
	store      types.KVStore
	readCache  map[string][]byte // optional, for caching writes to store
	writeCache *types.KVCache    // optional, for caching writes w/o writing to store
	logger     log.Logger
}

func NewState(store types.KVStore) *State {
	return &State{
		chainID:    "",
		store:      store,
		readCache:  make(map[string][]byte),
		writeCache: nil,
		logger:     log.NewNopLogger(),
	}
}

func (s *State) SetLogger(l log.Logger) {
	s.logger = l
}

func (s *State) SetChainID(chainID string) {
	s.chainID = chainID
	s.store.Set([]byte("base/chain_id"), []byte(chainID))
}

func (s *State) GetChainID() string {
	if s.chainID != "" {
		return s.chainID
	}
	s.chainID = string(s.store.Get([]byte("base/chain_id")))
	return s.chainID
}

func (s *State) Get(key []byte) (value []byte) {
	if s.readCache != nil { //if not a cachewrap
		value, ok := s.readCache[string(key)]
		if ok {
			return value
		}
	}
	return s.store.Get(key)
}

func (s *State) Set(key []byte, value []byte) {
	if s.readCache != nil { //if not a cachewrap
		s.readCache[string(key)] = value
	}
	s.store.Set(key, value)
}


// Account
func AccountKey(username types.AccountName) []byte {
	return append([]byte("account/"), username...)
}

func (s *State) GetAccount(username types.AccountName) *types.Account {
	data := s.Get(AccountKey(username))
	if len(data) == 0 {
		return nil
	}
	var acc *types.Account
	err := wire.ReadBinaryBytes(data, &acc)
	if err != nil {
		panic(fmt.Sprintf("Error reading account %X error: %v",
			data, err.Error()))
	}
	return acc
}

func (s *State) SetAccount(username types.AccountName, acc *types.Account) {
	accBytes := wire.BinaryBytes(acc)
	s.Set(AccountKey(username), accBytes)
}

// Post
func (s *State) GetPost(pid []byte) *types.Post {
	return types.GetPost(s, pid)
}

func (s *State) SetPost(pid []byte, post *types.Post) {
	types.SetPost(s, pid, post)
}

func (s *State) UpdateCommentParent(post *types.Post, parent *types.Post) {
	fmt.Println("Not implemented yet.", post, parent)
}

// Like

func (s *State) GetLikesByPostId(post_id []byte) []types.Like {
	return types.GetLikesByPostId(s, post_id);
}

func (s *State) AddLike(like types.Like) {
	types.AddLike(s, like)
}

func (s *State) CacheWrap() *State {
	cache := types.NewKVCache(s)
	return &State{
		chainID:    s.chainID,
		store:      cache,
		readCache:  nil,
		writeCache: cache,
		logger:     s.logger,
	}
}

// Donate
func (s *State) UpdateDonatePost(post *types.Post, acc *types.Account, coin types.Coins) {
	fmt.Println("Not implemented yet.", post, acc, coin)
}

// NOTE: errors if s is not from CacheWrap()
func (s *State) CacheSync() {
	s.writeCache.Sync()
}

func (s *State) Commit() abci.Result {
	switch s.store.(type) {
	case *eyes.Client:
		s.readCache = make(map[string][]byte)
		return s.store.(*eyes.Client).CommitSync()
	default:
		return abci.NewError(abci.CodeType_InternalError, "can only use Commit if store is merkleeyes")
	}

}
