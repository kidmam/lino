package model

import (
	"github.com/lino-network/lino/param"
	types "github.com/lino-network/lino/types"
)

type Proposal interface {
	GetProposalInfo() ProposalInfo
	SetProposalInfo(ProposalInfo)
}

type ProposalInfo struct {
	Creator       types.AccountKey     `json:"creator"`
	ProposalID    types.ProposalKey    `json:"proposal_id"`
	AgreeVotes    types.Coin           `json:"agree_vote"`
	DisagreeVotes types.Coin           `json:"disagree_vote"`
	Result        types.ProposalResult `json:"result"`
	CreatedAt     int64                `json:"created_at"`
	ExpiredAt     int64                `json:"expired_at"`
	Reason        string               `json:"reason"`
}

type ChangeParamProposal struct {
	ProposalInfo
	Param  param.Parameter `json:"param"`
	Reason string          `json:"reason"`
}

func (p *ChangeParamProposal) GetProposalInfo() ProposalInfo     { return p.ProposalInfo }
func (p *ChangeParamProposal) SetProposalInfo(info ProposalInfo) { p.ProposalInfo = info }

type ContentCensorshipProposal struct {
	ProposalInfo
	Permlink types.Permlink `json:"permlink"`
	Reason   string         `json:"reason"`
}

func (p *ContentCensorshipProposal) GetProposalInfo() ProposalInfo     { return p.ProposalInfo }
func (p *ContentCensorshipProposal) SetProposalInfo(info ProposalInfo) { p.ProposalInfo = info }

type ProtocolUpgradeProposal struct {
	ProposalInfo
	Link   string `json:"link"`
	Reason string `json:"reason"`
}

func (p *ProtocolUpgradeProposal) GetProposalInfo() ProposalInfo     { return p.ProposalInfo }
func (p *ProtocolUpgradeProposal) SetProposalInfo(info ProposalInfo) { p.ProposalInfo = info }

type NextProposalID struct {
	NextProposalID int64 `json:"next_proposal_id"`
}
