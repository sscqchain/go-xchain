package upgrade

import (
	sdk "gitee.com/xchain/go-xchain/types"
	"gitee.com/xchain/go-xchain/version"
)
const defaultProtocolVersion = version.ProtocolVersion


// GenesisState - all upgrade state that must be provided at genesis
type GenesisState struct {
	GenesisVersion VersionInfo `json:genesis_version`
}

// InitGenesis - build the genesis version For first Version
func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) {
	genesisVersion := data.GenesisVersion

	k.AddNewVersionInfo(ctx, genesisVersion)
	k.protocolKeeper.ClearUpgradeConfig(ctx)
	k.protocolKeeper.SetCurrentVersion(ctx, genesisVersion.UpgradeInfo.Protocol.Version)
}

// WriteGenesis - output genesis parameters
func ExportGenesis(ctx sdk.Context) GenesisState {
	return GenesisState{
		NewVersionInfo(sdk.DefaultUpgradeConfig(defaultProtocolVersion, "https://gitee.com/xchain/go-xchain/releases/tag/v"+version.Version), true),
	}
}

// get raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	return GenesisState{
		NewVersionInfo(sdk.DefaultUpgradeConfig(defaultProtocolVersion, "https://gitee.com/xchain/go-xchain/releases/tag/v"+version.Version), true),
	}
}

// get raw genesis raw message for testing
func DefaultGenesisStateForTest() GenesisState {
	return GenesisState{
		NewVersionInfo(sdk.DefaultUpgradeConfig(defaultProtocolVersion, "https://gitee.com/xchain/go-xchain/releases/tag/v"+version.Version), true),
	}
}
