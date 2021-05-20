package app

import (
	"os"
	"testing"

	"github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/stretchr/testify/require"

	v0 "gitee.com/xchain/go-xchain/app/v0"
	"gitee.com/xchain/go-xchain/codec"
	sdk "gitee.com/xchain/go-xchain/types"
	"gitee.com/xchain/go-xchain/x/auth"
	"gitee.com/xchain/go-xchain/x/crisis"
	distr "gitee.com/xchain/go-xchain/x/distribution"
	"gitee.com/xchain/go-xchain/x/gov"
	"gitee.com/xchain/go-xchain/x/guardian"
	"gitee.com/xchain/go-xchain/x/mint"
	"gitee.com/xchain/go-xchain/x/service"
	"gitee.com/xchain/go-xchain/x/slashing"
	"gitee.com/xchain/go-xchain/x/staking"
	"gitee.com/xchain/go-xchain/x/upgrade"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

func setGenesis(happ *SscqServiceApp, accs ...*auth.BaseAccount) error {
	genaccs := make([]v0.GenesisAccount, len(accs))
	for i, acc := range accs {
		genaccs[i] = v0.NewGenesisAccount(acc)
	}

	genesisState := v0.NewGenesisState(
		genaccs,
		auth.DefaultGenesisState(),
		staking.DefaultGenesisState(),
		mint.DefaultGenesisState(),
		distr.DefaultGenesisState(),
		gov.DefaultGenesisState(),
		upgrade.DefaultGenesisState(),
		service.DefaultGenesisState(),
		guardian.DefaultGenesisState(),
		slashing.DefaultGenesisState(),
		crisis.DefaultGenesisState(),
	)

	stateBytes, err := codec.MarshalJSONIndent(v0.MakeLatestCodec(), genesisState)
	if err != nil {
		return err
	}

	// Initialize the chain
	vals := []abci.ValidatorUpdate{}
	happ.InitChain(abci.RequestInitChain{Validators: vals, AppStateBytes: stateBytes})
	happ.Commit()

	return nil
}

func TestHsdExport(t *testing.T) {
	db := db.NewMemDB()

	happ := NewSscqServiceApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), config.TestInstrumentationConfig(), db, nil, true, 0)
	// accs added by junying, 2019-11-20
	key := secp256k1.GenPrivKey()
	pub := key.PubKey()
	addr := sdk.AccAddress(pub.Address())
	acc := auth.NewBaseAccountWithAddress(addr)
	setGenesis(happ, &acc)

	// Making a new app object with the db, so that initchain hasn't been called
	// panic: consensus params is empty
	_, _, err := happ.ExportAppStateAndValidators(false)
	require.NoError(t, err, "ExportAppStateAndValidators should not have an error")
}
