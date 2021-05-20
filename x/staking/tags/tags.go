// nolint
package tags

import (
	sdk "gitee.com/xchain/go-xchain/types"
)

var (
	ActionCompleteUnbonding     = "complete-unbonding"
	ActionCompleteRedelegation  = "complete-redelegation"
	ActionCompleteAuthorization = "complete-authorization"

	Action       = sdk.TagAction
	SrcValidator = sdk.TagSrcValidator
	DstValidator = sdk.TagDstValidator
	Delegator    = sdk.TagDelegator
	Moniker      = "moniker"
	Identity     = "identity"
	EndTime      = "end-time"
)
