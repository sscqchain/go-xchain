package bank

import (
	sdk "gitee.com/xchain/go-xchain/types"
)

// expected crisis keeper
type CrisisKeeper interface {
	RegisterRoute(moduleName, route string, invar sdk.Invariant)
}
