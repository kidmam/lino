package model

import (
	crypto "github.com/tendermint/tendermint/crypto"

	"github.com/lino-network/lino/types"
)

// PendingCoinDayQueueIR - TotalCoinDay: rat -> string, int -> string
type PendingCoinDayQueueIR struct {
	LastUpdatedAt   string             `json:"last_updated_at"`
	TotalCoinDay    string             `json:"total_coin_day"`
	TotalCoin       types.Coin         `json:"total_coin"`
	PendingCoinDays []PendingCoinDayIR `json:"pending_coin_days"`
}

// PendingCoinDayIR - int to string
type PendingCoinDayIR struct {
	StartTime string     `json:"start_time"`
	EndTime   string     `json:"end_time"`
	Coin      types.Coin `json:"coin"`
}

// AccountInfoIR - IR
type AccountInfoIR struct {
	Username       types.AccountKey `json:"username"`
	CreatedAt      string           `json:"created_at"`
	ResetKey       crypto.PubKey    `json:"reset_key"`
	TransactionKey crypto.PubKey    `json:"transaction_key"`
	AppKey         crypto.PubKey    `json:"app_key"`
}

// ToState - convert IR back to state.
func (a AccountInfoIR) ToState() *AccountInfo {
	return &AccountInfo{
		Username:       a.Username,
		CreatedAt:      types.MustParseInt64(a.CreatedAt, 10, 64),
		ResetKey:       a.ResetKey,
		TransactionKey: a.TransactionKey,
		AppKey:         a.AppKey,
	}
}

// FrozenMoneyIR - IR
type FrozenMoneyIR struct {
	Amount   types.Coin `json:"amount"`
	StartAt  string     `json:"start_at"`
	Times    string     `json:"times"`
	Interval string     `json:"interval"`
}

// AccountBankIR - user balance
type AccountBankIR struct {
	Saving          types.Coin      `json:"saving"`
	CoinDay         types.Coin      `json:"coin_day"`
	FrozenMoneyList []FrozenMoneyIR `json:"frozen_money_list"`
}

// ToState - convert IR back to state.
func (a AccountBankIR) ToState() *AccountBank {
	frozen := make([]FrozenMoney, 0)
	for _, v := range a.FrozenMoneyList {
		frozen = append(frozen, *v.ToState())
	}
	return &AccountBank{
		Saving:          a.Saving,
		CoinDay:         a.CoinDay,
		FrozenMoneyList: frozen,
	}
}

// AccountMetaIR - int to string and internal conversions
type AccountMetaIR struct {
	Sequence             string     `json:"sequence"`
	LastActivityAt       string     `json:"last_activity_at"`
	TransactionCapacity  types.Coin `json:"transaction_capacity"`
	JSONMeta             string     `json:"json_meta"`
	LastReportOrUpvoteAt string     `json:"last_report_or_upvote_at"`
	LastPostAt           string     `json:"last_post_at"`
}

// ToState - convert IR back to state.
func (a AccountMetaIR) ToState() *AccountMeta {
	return &AccountMeta{
		Sequence:             uint64(types.MustParseInt64(a.Sequence, 10, 64)),
		LastActivityAt:       types.MustParseInt64(a.LastActivityAt, 10, 64),
		TransactionCapacity:  a.TransactionCapacity,
		JSONMeta:             a.JSONMeta,
		LastReportOrUpvoteAt: types.MustParseInt64(a.LastReportOrUpvoteAt, 10, 64),
		LastPostAt:           types.MustParseInt64(a.LastPostAt, 10, 64),
	}
}

// AccountRowIR account related information when migrate, pk: Username
type AccountRowIR struct {
	Username            types.AccountKey      `json:"username"`
	Info                AccountInfoIR         `json:"info"`
	Bank                AccountBankIR         `json:"bank"`
	Meta                AccountMetaIR         `json:"meta"`
	PendingCoinDayQueue PendingCoinDayQueueIR `json:"pending_coin_day_queue"`
}

// GrantPubKeyIR - int to string and internal conversions
type GrantPubKeyIR struct {
	Username   types.AccountKey `json:"username"`
	Permission types.Permission `json:"permission"`
	CreatedAt  string           `json:"created_at"`
	ExpiresAt  string           `json:"expires_at"`
	Amount     types.Coin       `json:"amount"`
}

// GrantPubKeyRowIR - int to string and internal conversions
type GrantPubKeyRowIR struct {
	Username    types.AccountKey `json:"username"`
	PubKey      crypto.PubKey    `json:"pub_key"`
	GrantPubKey GrantPubKeyIR    `json:"grant_pub_key"`
}

// AccountTablesIR -
type AccountTablesIR struct {
	Accounts            []AccountRowIR     `json:"accounts"`
	AccountGrantPubKeys []GrantPubKeyRowIR `json:"account_grant_pub_keys"`
}
