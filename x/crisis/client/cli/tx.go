// nolint
package cli

import (
	"github.com/spf13/cobra"

	"gitee.com/xchain/go-xchain/client/context"
	"gitee.com/xchain/go-xchain/client/utils"
	"gitee.com/xchain/go-xchain/codec"
	sdk "gitee.com/xchain/go-xchain/types"
	authtxb "gitee.com/xchain/go-xchain/x/auth/client/txbuilder"
	"gitee.com/xchain/go-xchain/x/crisis"
)

// command to replace a delegator's withdrawal address
func GetCmdInvariantBroken(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "invariant-broken [module-name] [invariant-route]",
		Short: "submit proof that an invariant broken to halt the chain",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {

			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			senderAddr := cliCtx.GetFromAddress()
			moduleName, route := args[0], args[1]
			msg := crisis.NewMsgVerifyInvariant(senderAddr, moduleName, route)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg}, false)
		},
	}
	return cmd
}
