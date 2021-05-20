package main

import (
	"fmt"
	"os"
	"path"

	"gitee.com/xchain/go-xchain/client/bech32"

	"gitee.com/xchain/go-xchain/params"
	svrConfig "gitee.com/xchain/go-xchain/server/config"

	"gitee.com/xchain/go-xchain/client"
	"gitee.com/xchain/go-xchain/client/lcd"
	"gitee.com/xchain/go-xchain/client/rpc"
	"gitee.com/xchain/go-xchain/client/tx"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/libs/cli"

	sdk "gitee.com/xchain/go-xchain/types"
	authcmd "gitee.com/xchain/go-xchain/x/auth/client/cli"
	sscqservicecmd "gitee.com/xchain/go-xchain/x/core/client/cli"
	faucetcmd "gitee.com/xchain/go-xchain/x/faucet/client/cli"

	accounts "gitee.com/xchain/go-xchain/accounts/cli"
	accrest "gitee.com/xchain/go-xchain/accounts/rest"
	"gitee.com/xchain/go-xchain/app"
	ssrest "gitee.com/xchain/go-xchain/x/core/client/rest"

	dist "gitee.com/xchain/go-xchain/x/distribution/client/rest"
	gv "gitee.com/xchain/go-xchain/x/gov"
	gov "gitee.com/xchain/go-xchain/x/gov/client/rest"
	mint "gitee.com/xchain/go-xchain/x/mint/client/rest"
	sl "gitee.com/xchain/go-xchain/x/slashing"
	slashing "gitee.com/xchain/go-xchain/x/slashing/client/rest"
	st "gitee.com/xchain/go-xchain/x/staking"
	staking "gitee.com/xchain/go-xchain/x/staking/client/rest"

	sscliversion "gitee.com/xchain/go-xchain/server"
	distcmd "gitee.com/xchain/go-xchain/x/distribution"
	ssdistClient "gitee.com/xchain/go-xchain/x/distribution/client"
	ssgovClient "gitee.com/xchain/go-xchain/x/gov/client"
	ssmintClient "gitee.com/xchain/go-xchain/x/mint/client/cli"
	sslashingClient "gitee.com/xchain/go-xchain/x/slashing/client"
	sstakingClient "gitee.com/xchain/go-xchain/x/staking/client"
	upgradecmd "gitee.com/xchain/go-xchain/x/upgrade/client/cli"
	upgraderest "gitee.com/xchain/go-xchain/x/upgrade/client/rest"
)

const (
	storeAcc = "acc"
	storeHS  = "ss"
)

var (
	DEBUGAPI  = "OFF"
	GitCommit = ""
	GitBranch = ""
)

func main() {
	cobra.EnableCommandSorting = false

	if DEBUGAPI == svrConfig.ValueDebugApi_On {
		svrConfig.ApiSecurityLevel = svrConfig.ValueSecurityLevel_Low
	}

	cdc := app.MakeLatestCodec()

	// set address prefix
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(sdk.Bech32PrefixAccAddr, sdk.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(sdk.Bech32PrefixValAddr, sdk.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(sdk.Bech32PrefixConsAddr, sdk.Bech32PrefixConsPub)
	config.Seal()

	mc := []sdk.ModuleClients{
		ssgovClient.NewModuleClient(gv.StoreKey, cdc),
		ssdistClient.NewModuleClient(distcmd.StoreKey, cdc),
		sstakingClient.NewModuleClient(st.StoreKey, cdc),
		sslashingClient.NewModuleClient(sl.StoreKey, cdc),
	}

	rootCmd := &cobra.Command{
		Use:   "sscli",
		Short: "sscqservice Client",
	}

	// Add --chain-id to persistent flags and mark it required
	rootCmd.PersistentFlags().String(client.FlagChainID, "", "Chain ID of tendermint node")
	rootCmd.PersistentPreRunE = func(_ *cobra.Command, _ []string) error {
		return initConfig(rootCmd)
	}

	// Construct Root Command
	rootCmd.AddCommand(
		rpc.StatusCommand(),
		client.ConfigCmd(app.DefaultCLIHome),
		queryCmd(cdc, mc), // check the below
		txCmd(cdc, mc),    // check the below
		versionCmd(cdc, mc),
		client.LineBreak,
		lcd.ServeCommand(cdc, registerRoutes),
		client.LineBreak,
		accounts.Commands(),
		client.LineBreak,
		sscliversion.VersionHscliCmd,
		bech32.Bech32Commands(),
	)

	executor := cli.PrepareMainCmd(rootCmd, "SS", app.DefaultCLIHome)
	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}

func registerRoutes(rs *lcd.RestServer) {
	rs.CliCtx = rs.CliCtx.WithAccountDecoder(rs.Cdc)
	rpc.RegisterRoutes(rs.CliCtx, rs.Mux, rs.Cdc)
	tx.RegisterRoutes(rs.CliCtx, rs.Mux, rs.Cdc)
	ssrest.RegisterRoutes(rs.CliCtx, rs.Mux, rs.Cdc, storeHS)
	accrest.RegisterRoutes(rs.CliCtx, rs.Mux, rs.Cdc)
	accrest.RegisterRoute(rs.CliCtx, rs.Mux, rs.Cdc, storeAcc)
	dist.RegisterRoutes(rs.CliCtx, rs.Mux, rs.Cdc, distcmd.StoreKey)
	staking.RegisterRoutes(rs.CliCtx, rs.Mux, rs.Cdc, rs.KeyBase)
	slashing.RegisterRoutes(rs.CliCtx, rs.Mux, rs.Cdc, rs.KeyBase)
	gov.RegisterRoutes(rs.CliCtx, rs.Mux, rs.Cdc)
	mint.RegisterRoutes(rs.CliCtx, rs.Mux, rs.Cdc)
	upgraderest.RegisterRoutes(rs.CliCtx, rs.Mux, rs.Cdc)
}

func versionCmd(cdc *amino.Codec, mc []sdk.ModuleClients) *cobra.Command {
	cbCmd := &cobra.Command{
		Use:   "version",
		Short: "print version, api security level",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("GitCommit=%s|version=%s|GitBranch=%s|DEBUGAPI=%s|ApiSecurityLevel=%s\n", GitCommit, params.Version, GitBranch, DEBUGAPI, svrConfig.ApiSecurityLevel)
		},
	}

	return cbCmd
}

func queryCmd(cdc *amino.Codec, mc []sdk.ModuleClients) *cobra.Command {
	queryCmd := &cobra.Command{
		Use:     "query",
		Aliases: []string{"q"},
		Short:   "Querying subcommands",
	}

	queryCmd.AddCommand(
		rpc.ValidatorCommand(cdc),
		rpc.BlockCommand(),
		tx.SearchTxCmd(cdc),
		tx.QueryTxCmd(cdc),
		client.LineBreak,
		authcmd.GetAccountCmd(storeAcc, cdc),
		sscqservicecmd.GetCmdCall(cdc),
		ssmintClient.GetCmdQueryBlockRewards(cdc),
		ssmintClient.GetCmdQueryTotalProvisions(cdc),
		upgradecmd.GetInfoCmd("upgrade", cdc),
		upgradecmd.GetCmdQuerySignals("upgrade", cdc),
	)

	for _, m := range mc {
		queryCmd.AddCommand(m.GetQueryCmd())
	}

	return queryCmd
}

func txCmd(cdc *amino.Codec, mc []sdk.ModuleClients) *cobra.Command {
	txCmd := &cobra.Command{
		Use:   "tx",
		Short: "Transactions subcommands",
	}

	if svrConfig.ApiSecurityLevel == svrConfig.ValueSecurityLevel_Low {
		txCmd.AddCommand(
			sscqservicecmd.GetCmdBurn(cdc),
			sscqservicecmd.GetCmdCreate(cdc),
			sscqservicecmd.GetCmdSend(cdc),
			sscqservicecmd.GetCmdSign(cdc),
			faucetcmd.GetCmdAdd(cdc),
		)
	}

	txCmd.AddCommand(
		sscqservicecmd.GetCmdBroadCast(cdc),
		client.LineBreak,
	)

	for _, m := range mc {
		txCmd.AddCommand(m.GetTxCmd())
	}

	return txCmd
}

func initConfig(cmd *cobra.Command) error {
	home, err := cmd.PersistentFlags().GetString(cli.HomeFlag)
	if err != nil {
		return err
	}

	cfgFile := path.Join(home, "config", "config.toml")
	if _, err := os.Stat(cfgFile); err == nil {
		viper.SetConfigFile(cfgFile)

		if err := viper.ReadInConfig(); err != nil {
			return err
		}
	}
	if err := viper.BindPFlag(client.FlagChainID, cmd.PersistentFlags().Lookup(client.FlagChainID)); err != nil {
		return err
	}
	if err := viper.BindPFlag(cli.EncodingFlag, cmd.PersistentFlags().Lookup(cli.EncodingFlag)); err != nil {
		return err
	}
	return viper.BindPFlag(cli.OutputFlag, cmd.PersistentFlags().Lookup(cli.OutputFlag))
}
