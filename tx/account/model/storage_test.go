package model

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lino-network/lino/types"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/abci/types"
	"github.com/tendermint/go-crypto"
	dbm "github.com/tendermint/tmlibs/db"
)

var (
	TestKVStoreKey = sdk.NewKVStoreKey("account")
)

func getContext() sdk.Context {
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(TestKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	return sdk.NewContext(ms, abci.Header{}, false, nil)
}

func TestAccountInfo(t *testing.T) {
	as := NewAccountStorage(TestKVStoreKey)
	ctx := getContext()

	priv := crypto.GenPrivKeyEd25519()
	accInfo := AccountInfo{
		Username:       types.AccountKey("test"),
		Created:        0,
		MasterKey:      priv.PubKey(),
		TransactionKey: priv.Generate(1).PubKey(),
		PostKey:        priv.Generate(2).PubKey(),
		Address:        priv.PubKey().Address(),
	}
	err := as.SetInfo(ctx, types.AccountKey("test"), &accInfo)
	assert.Nil(t, err)

	resultPtr, err := as.GetInfo(ctx, types.AccountKey("test"))
	assert.Nil(t, err)
	assert.Equal(t, accInfo, *resultPtr, "Account info should be equal")
}

func TestInvalidAccountInfo(t *testing.T) {
	as := NewAccountStorage(TestKVStoreKey)
	ctx := getContext()

	resultPtr, err := as.GetInfo(ctx, types.AccountKey("test"))
	assert.Nil(t, resultPtr)
	assert.Equal(t, err, ErrAccountInfoDoesntExist())
}

func TestAccountBank(t *testing.T) {
	as := NewAccountStorage(TestKVStoreKey)
	ctx := getContext()

	priv := crypto.GenPrivKeyEd25519()
	accInfo := AccountInfo{
		Username:       types.AccountKey("test"),
		Created:        0,
		MasterKey:      priv.PubKey(),
		TransactionKey: priv.Generate(1).PubKey(),
		PostKey:        priv.Generate(2).PubKey(),
		Address:        priv.PubKey().Address(),
	}
	err := as.SetInfo(ctx, types.AccountKey("test"), &accInfo)
	assert.Nil(t, err)

	accBank := AccountBank{
		Address: priv.PubKey().Address(),
		Balance: types.NewCoin(int64(123)),
	}
	err = as.SetBankFromAddress(ctx, priv.PubKey().Address(), &accBank)
	assert.Nil(t, err)

	resultPtr, err := as.GetBankFromAccountKey(ctx, types.AccountKey("test"))
	assert.Nil(t, err)
	assert.Equal(t, accBank, *resultPtr, "Account bank should be equal")

	resultPtr, err = as.GetBankFromAddress(ctx, priv.PubKey().Address())
	assert.Nil(t, err)
	assert.Equal(t, accBank, *resultPtr, "Account bank should be equal")
}

func TestAccountMeta(t *testing.T) {
	as := NewAccountStorage(TestKVStoreKey)
	ctx := getContext()

	accMeta := AccountMeta{}
	err := as.SetMeta(ctx, types.AccountKey("test"), &accMeta)
	assert.Nil(t, err)

	resultPtr, err := as.GetMeta(ctx, types.AccountKey("test"))
	assert.Nil(t, err)
	assert.Equal(t, accMeta, *resultPtr, "Account meta should be equal")
}

func TestAccountReward(t *testing.T) {
	as := NewAccountStorage(TestKVStoreKey)
	ctx := getContext()

	reward := Reward{}
	err := as.SetReward(ctx, types.AccountKey("test"), &reward)
	assert.Nil(t, err)

	resultPtr, err := as.GetReward(ctx, types.AccountKey("test"))
	assert.Nil(t, err)
	assert.Equal(t, reward, *resultPtr, "Account reward should be equal")
}

func TestAccountRelationShip(t *testing.T) {
	as := NewAccountStorage(TestKVStoreKey)
	ctx := getContext()

	relationship := Relationship{}
	err := as.SetRelationship(
		ctx, types.AccountKey("me"), types.AccountKey("other"), &relationship)
	assert.Nil(t, err)

	resultPtr, err := as.GetRelationship(ctx, types.AccountKey("me"), types.AccountKey("other"))
	assert.Nil(t, err)
	assert.Equal(t, relationship, *resultPtr, "Account relationship should be equal")
}
