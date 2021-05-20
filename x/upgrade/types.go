package upgrade

import (
	sdk "gitee.com/xchain/go-xchain/types"
)

type VersionInfo struct {
	UpgradeInfo sdk.UpgradeConfig
	Success     bool
}

func NewVersionInfo(upgradeConfig sdk.UpgradeConfig, success bool) VersionInfo {
	return VersionInfo{
		upgradeConfig,
		success,
	}
}
