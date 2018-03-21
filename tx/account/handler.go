package account

import (
	"fmt"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(am AccountManager) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case FollowMsg:
			return handleFollowMsg(ctx, am, msg)
		case UnfollowMsg:
			return handleUnfollowMsg(ctx, am, msg)
		case TransferMsg:
			return handleTransferMsg(ctx, am, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized account Msg type: %v", reflect.TypeOf(msg).Name())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle FollowMsg
func handleFollowMsg(ctx sdk.Context, am AccountManager, msg FollowMsg) sdk.Result {
	followerList, err := am.GetFollower(ctx, msg.Followee)
	if err != nil {
		return ErrAccountManagerFail("Get follower list failed").Result()
	}

	followingList, err := am.GetFollowing(ctx, msg.Follower)
	if err != nil {
		return ErrAccountManagerFail("Get following list failed").Result()
	}

	// add the "msg.Follower" to the "msg.Followee" 's follower list.
	if findAccountInList(msg.Follower, followerList.Follower) == -1 {
		followerList.Follower = append(followerList.Follower, msg.Follower)
		if err := am.SetFollower(ctx, msg.Followee, followerList); err != nil {
			return ErrAccountManagerFail("Set follower failed").Result()
		}
	}

	// add the "msg.Followee" to the "msg.Follower" 's following list.
	if findAccountInList(msg.Followee, followingList.Following) == -1 {
		followingList.Following = append(followingList.Following, msg.Followee)
		if err := am.SetFollowing(ctx, msg.Follower, followingList); err != nil {
			return ErrAccountManagerFail("Set following failed").Result()
		}
	}

	return sdk.Result{}
}

// Handle UnfollowMsg
func handleUnfollowMsg(ctx sdk.Context, am AccountManager, msg UnfollowMsg) sdk.Result {
	followerList, err := am.GetFollower(ctx, msg.Followee)
	if err != nil {
		return ErrAccountManagerFail("Get follower list failed").Result()
	}

	followingList, err := am.GetFollowing(ctx, msg.Follower)
	if err != nil {
		return ErrAccountManagerFail("Get following list failed").Result()
	}

	// remove the "msg.Follower" from the "msg.Followee" 's follower list.
	if idx := findAccountInList(msg.Follower, followerList.Follower); idx != -1 {
		followerList.Follower = append(followerList.Follower[:idx], followerList.Follower[idx+1:]...)
		if err := am.SetFollower(ctx, msg.Followee, followerList); err != nil {
			return ErrAccountManagerFail("Set follower failed").Result()
		}
	}

	// remove the "msg.Followee" from the "msg.Follower" 's following list.
	if idx := findAccountInList(msg.Followee, followingList.Following); idx != -1 {
		followingList.Following = append(followingList.Following[:idx], followingList.Following[idx+1:]...)
		if err := am.SetFollowing(ctx, msg.Follower, followingList); err != nil {
			return ErrAccountManagerFail("Set following failed").Result()
		}
	}

	return sdk.Result{}
}

// Handle TransferMsg
func handleTransferMsg(ctx sdk.Context, am AccountManager, msg TransferMsg) sdk.Result {
	// check if the sender has enough money
	senderBank, err := am.GetBankFromAccountKey(ctx, msg.Sender)
	if err != nil {
		return ErrAccountManagerFail("Get sender's account bank failed").Result()
	}

	if !senderBank.Balance.IsGTE(msg.Amount) {
		return ErrAccountManagerFail("Sender's coins are not enough").Result()
	}

	// withdraw money from sender's bank
	senderBank.Balance = senderBank.Balance.Minus(msg.Amount)
	if err := am.SetBankFromAccountKey(ctx, msg.Sender, senderBank); err != nil {
		return ErrAccountManagerFail("Set sender's bank failed").Result()
	}

	// send coins using username
	if am.AccountExist(ctx, msg.ReceiverName) {
		if receiverBank, err := am.GetBankFromAccountKey(ctx, msg.ReceiverName); err == nil {
			receiverBank.Balance = receiverBank.Balance.Plus(msg.Amount)
			if setErr := am.SetBankFromAccountKey(ctx, msg.ReceiverName, receiverBank); setErr != nil {
				return ErrAccountManagerFail("Set receiver's bank failed").Result()
			}
			return sdk.Result{}
		}
	}

	// send coins using address
	receiverBank, err := am.GetBankFromAddress(ctx, msg.ReceiverAddr)
	if err == nil {
		// account bank exists
		receiverBank.Balance = receiverBank.Balance.Plus(msg.Amount)
	} else {
		// account bank not found, create a new one for this address
		receiverBank = &AccountBank{
			Address: msg.ReceiverAddr,
			Balance: msg.Amount,
		}
	}

	if setErr := am.SetBankFromAddress(ctx, msg.ReceiverAddr, receiverBank); setErr != nil {
		return ErrAccountManagerFail("Set receiver's bank failed").Result()
	}
	return sdk.Result{}
}

// helper function
func findAccountInList(me AccountKey, lst []AccountKey) int {
	for index, user := range lst {
		if user == me {
			return index
		}
	}
	return -1
}