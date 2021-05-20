package cli

import (
	"fmt"
	"os"
	"strings"

	"gitee.com/xchain/go-xchain/accounts/keystore"
	"gitee.com/xchain/go-xchain/client"
	"gitee.com/xchain/go-xchain/client/context"
	"gitee.com/xchain/go-xchain/client/keys"
	"gitee.com/xchain/go-xchain/client/utils"
	"gitee.com/xchain/go-xchain/codec"
	sdk "gitee.com/xchain/go-xchain/types"
	authtxb "gitee.com/xchain/go-xchain/x/auth/client/txbuilder"
	faucet "gitee.com/xchain/go-xchain/x/faucet"
	"github.com/spf13/cobra"
)

// junying-todo-20190409
func GetCmdAdd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add [amount]",
		Short: "publish new coin or add existing coin to system issuer except stake",
		Long:  "hscli tx add sscq1yaaf9egv34zgua4lkanakktmv36ch8cl0lzkn3 5satoshi --gas-price=100",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {

			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			systemissuer, err := faucet.GetSystemIssuerFromValidators(cdc)
			if err != nil {
				return err
			}

			if systemissuer == nil {
				fmt.Print("failed to find systemissuer.\n")
				return nil
			}

			toaddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			if strings.Contains(args[1], "stake") {
				fmt.Print("stake can't be added. Or, system will panic. \n")
				return nil
			}

			coins, err := sdk.ParseCoins(args[1])
			if err != nil {
				return err
			}

			msg := faucet.NewMsgAdd(systemissuer, toaddr, coins)
			cliCtx.PrintResponse = true

			return CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg}, systemissuer) //not completed yet, need account name
		},
	}

	return client.PostCommands(cmd)[0]
}

func PrepareTxBuilder(txBldr authtxb.TxBuilder, cliCtx context.CLIContext, fromaddr sdk.AccAddress) (authtxb.TxBuilder, error) {

	// TODO: (ref #1903) Allow for user supplied account number without
	// automatically doing a manual lookup.
	if txBldr.AccountNumber() == 0 {
		accNum, err := cliCtx.GetAccountNumber(fromaddr)
		if err != nil {
			return txBldr, err
		}
		txBldr = txBldr.WithAccountNumber(accNum)
	}

	// TODO: (ref #1903) Allow for user supplied account sequence without
	// automatically doing a manual lookup.
	if txBldr.Sequence() == 0 {
		accSeq, err := cliCtx.GetAccountSequence(fromaddr)
		if err != nil {
			return txBldr, err
		}
		txBldr = txBldr.WithSequence(accSeq)
	}
	return txBldr, nil
}

func CompleteAndBroadcastTxCLI(txBldr authtxb.TxBuilder, cliCtx context.CLIContext, msgs []sdk.Msg, fromaddr sdk.AccAddress) error {
	//
	txBldr, err := PrepareTxBuilder(txBldr, cliCtx, fromaddr)
	if err != nil {
		return err
	}

	if txBldr.SimulateAndExecute() || cliCtx.Simulate {
		txBldr, err := utils.EnrichWithGas(txBldr, cliCtx, msgs)
		if err != nil {
			return err
		}

		gasEst := utils.GasEstimateResponse{GasEstimate: txBldr.GasWanted()}
		fmt.Fprintf(os.Stderr, "%s\n", gasEst.String())
	}

	passphrase, err := keys.ReadShortPassphraseFromStdin(sdk.AccAddress.String(fromaddr))
	if err != nil {
		return err
	}
	addr := sdk.AccAddress.String(fromaddr)
	ksw := keystore.NewKeyStoreWallet(keystore.DefaultKeyStoreHome())
	txBytes, err := ksw.BuildAndSign(txBldr, addr, passphrase, msgs)
	if err != nil {
		return err
	}
	// broadcast to a Tendermint node
	res, err := cliCtx.BroadcastTx(txBytes)
	cliCtx.PrintOutput(res)
	return err
}
