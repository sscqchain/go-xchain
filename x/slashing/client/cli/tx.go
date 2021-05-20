package cli

import (
	"gitee.com/xchain/go-xchain/client/context"
	"gitee.com/xchain/go-xchain/client/utils"
	"gitee.com/xchain/go-xchain/codec"
	sdk "gitee.com/xchain/go-xchain/types"
	authtxb "gitee.com/xchain/go-xchain/x/auth/client/txbuilder"
	"gitee.com/xchain/go-xchain/x/slashing"
	sscorecli "gitee.com/xchain/go-xchain/x/core/client/cli"
	"github.com/spf13/cobra"
)

// GetCmdUnjail implements the create unjail validator command.
func GetCmdUnjail(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "unjail [keyaddr]",
		Short: "unjail validator previously jailed for downtime",
		Long: `unjail a jailed validator:

$ sscli tx slashing unjail [keyaddr]
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			validatorAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			msg := slashing.NewMsgUnjail(sdk.ValAddress(validatorAddr))
			return sscorecli.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg}, validatorAddr)
		},
	}
}
