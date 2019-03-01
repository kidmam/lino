package model

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
	"github.com/tendermint/tendermint/crypto"
)

// AccountInfo - user information
type AccountInfo struct {
	Username       types.AccountKey `json:"username"`
	CreatedAt      int64            `json:"created_at"`
	ResetKey       crypto.PubKey    `json:"reset_key"`
	TransactionKey crypto.PubKey    `json:"transaction_key"`
	AppKey         crypto.PubKey    `json:"app_key"`
}

// ToIR -
func (a AccountInfo) ToIR() AccountInfoIR {
	return AccountInfoIR{
		Username:       a.Username,
		CreatedAt:      strconv.FormatInt(a.CreatedAt, 10),
		ResetKey:       a.ResetKey,
		TransactionKey: a.TransactionKey,
		AppKey:         a.AppKey,
	}
}

// AccountBank - user balance
type AccountBank struct {
	Saving          types.Coin    `json:"saving"`
	CoinDay         types.Coin    `json:"coin_day"`
	FrozenMoneyList []FrozenMoney `json:"frozen_money_list"`
}

// ToIR -
func (a AccountBank) ToIR() AccountBankIR {
	return AccountBankIR{
		Saving:          a.Saving,
		CoinDay:         a.CoinDay,
		FrozenMoneyList: FrozenMoneySliceToIR(a.FrozenMoneyList),
	}
}

// FrozenMoney - frozen money
type FrozenMoney struct {
	Amount   types.Coin `json:"amount"`
	StartAt  int64      `json:"start_at"`
	Times    int64      `json:"times"`
	Interval int64      `json:"interval"`
}

// ToState - convert IR back to state.
func (f FrozenMoneyIR) ToState() *FrozenMoney {
	return &FrozenMoney{
		Amount:   f.Amount,
		StartAt:  types.MustParseInt64(f.StartAt, 10, 64),
		Times:    types.MustParseInt64(f.Times, 10, 64),
		Interval: types.MustParseInt64(f.Interval, 10, 64),
	}
}

// ToIR -
func (F FrozenMoney) ToIR() FrozenMoneyIR {
	return FrozenMoneyIR{
		Amount:   F.Amount,
		StartAt:  strconv.FormatInt(F.StartAt, 10),
		Times:    strconv.FormatInt(F.Times, 10),
		Interval: strconv.FormatInt(F.Interval, 10),
	}
}

// FrozenMoneySliceToIR -
func FrozenMoneySliceToIR(origin []FrozenMoney) (ir []FrozenMoneyIR) {
	for _, v := range origin {
		ir = append(ir, v.ToIR())
	}
	return
}

// PendingCoinDayQueue - stores a list of pending coin day and total number of coin waiting in list
type PendingCoinDayQueue struct {
	LastUpdatedAt   int64            `json:"last_updated_at"`
	TotalCoinDay    sdk.Dec          `json:"total_coin_day"`
	TotalCoin       types.Coin       `json:"total_coin"`
	PendingCoinDays []PendingCoinDay `json:"pending_coin_days"`
}

// ToIR coin.
func (p PendingCoinDayQueue) ToIR() PendingCoinDayQueueIR {
	return PendingCoinDayQueueIR{
		LastUpdatedAt:   strconv.FormatInt(p.LastUpdatedAt, 10),
		TotalCoinDay:    p.TotalCoinDay.String(),
		TotalCoin:       p.TotalCoin,
		PendingCoinDays: PendingCoinDaySliceToIR(p.PendingCoinDays),
	}
}

// PendingCoinDay - pending coin day in the list
type PendingCoinDay struct {
	StartTime int64      `json:"start_time"`
	EndTime   int64      `json:"end_time"`
	Coin      types.Coin `json:"coin"`
}

// PendingCoinDaySliceToIR -
func PendingCoinDaySliceToIR(origin []PendingCoinDay) (ir []PendingCoinDayIR) {
	for _, v := range origin {
		ir = append(ir, v.ToIR())
	}
	return
}

// ToIR -
func (s PendingCoinDay) ToIR() PendingCoinDayIR {
	return PendingCoinDayIR{
		StartTime: strconv.FormatInt(s.StartTime, 10),
		EndTime:   strconv.FormatInt(s.EndTime, 10),
		Coin:      s.Coin,
	}
}

// GrantPubKey - user grant permission to a public key with a certain permission
type GrantPubKey struct {
	Username   types.AccountKey `json:"username"`
	Permission types.Permission `json:"permission"`
	CreatedAt  int64            `json:"created_at"`
	ExpiresAt  int64            `json:"expires_at"`
	Amount     types.Coin       `json:"amount"`
}

// ToIR - int to string and internal conversions
func (g GrantPubKey) ToIR() GrantPubKeyIR {
	return GrantPubKeyIR{
		Username:   g.Username,
		Permission: g.Permission,
		CreatedAt:  strconv.FormatInt(g.CreatedAt, 10),
		ExpiresAt:  strconv.FormatInt(g.ExpiresAt, 10),
		Amount:     g.Amount,
	}
}

// AccountMeta - stores tiny and frequently updated fields.
type AccountMeta struct {
	Sequence             uint64     `json:"sequence"`
	LastActivityAt       int64      `json:"last_activity_at"`
	TransactionCapacity  types.Coin `json:"transaction_capacity"`
	JSONMeta             string     `json:"json_meta"`
	LastReportOrUpvoteAt int64      `json:"last_report_or_upvote_at"`
	LastPostAt           int64      `json:"last_post_at"`
}

// ToIR - int to string and internal conversions
func (a AccountMeta) ToIR() AccountMetaIR {
	return AccountMetaIR{
		Sequence:             strconv.FormatInt(int64(a.Sequence), 10),
		LastActivityAt:       strconv.FormatInt(a.LastActivityAt, 10),
		TransactionCapacity:  a.TransactionCapacity,
		JSONMeta:             a.JSONMeta,
		LastReportOrUpvoteAt: strconv.FormatInt(a.LastReportOrUpvoteAt, 10),
		LastPostAt:           strconv.FormatInt(a.LastPostAt, 10),
	}
}

// AccountInfraConsumption records infra utility consumption
// type AccountInfraConsumption struct {
// 	Storage   int64 `json:"storage"`
// 	Bandwidth int64 `json:"bandwidth"`
// }

// Reward - get from the inflation pool
type Reward struct {
	TotalIncome     types.Coin `json:"total_income"`
	OriginalIncome  types.Coin `json:"original_income"`
	FrictionIncome  types.Coin `json:"friction_income"`
	InflationIncome types.Coin `json:"inflation_income"`
	UnclaimReward   types.Coin `json:"unclaim_reward"`
}
