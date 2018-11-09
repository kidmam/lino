package commands

import (
	"fmt"

	"github.com/lino-network/lino/client"

	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	sdk "github.com/cosmos/cosmos-sdk/types"
	acc "github.com/lino-network/lino/x/account"
)

// updateAccountTxCmd will create a follow tx and sign it with the given key
func UpdateAccountTxCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-account",
		Short: "update account meta data",
		RunE:  updateAccountCmd(cdc),
	}
	cmd.Flags().String(client.FlagUser, "", "user who change the account meta")
	cmd.Flags().String(client.FlagAccMeta, "", "account meta")
	return cmd
}

// send follow transaction to the blockchain
func updateAccountCmd(cdc *wire.Codec) client.CommandTxCallback {
	return func(cmd *cobra.Command, args []string) error {
		ctx := client.NewCoreContextFromViper()
		user := viper.GetString(client.FlagUser)
		accountMeta := viper.GetString(client.FlagAccMeta)

		var msg sdk.Msg
		msg = acc.NewUpdateAccountMsg(user, accountMeta)

		// build and sign the transaction, then broadcast to Tendermint
		res, err := ctx.SignBuildBroadcast([]sdk.Msg{msg}, cdc)

		if err != nil {
			return err
		}

		fmt.Printf("Committed at block %d. Hash: %s\n", res.Height, res.Hash.String())
		return nil
	}
}
