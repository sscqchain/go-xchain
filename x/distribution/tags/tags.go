// nolint
package tags

import (
	sdk "gitee.com/xchain/go-xchain/types"
)

// Distribution tx tags
var (
	Rewards    = "rewards"
	Commission = "commission"

	Validator = sdk.TagSrcValidator
	Delegator = sdk.TagDelegator
)
