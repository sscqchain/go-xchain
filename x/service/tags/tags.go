package tags

import (
	sdk "gitee.com/xchain/go-xchain/types"
)

var (
	ActionSvcCallTimeOut = "service-call-expiration"

	Action = sdk.TagAction

	Provider   = "provider"
	Consumer   = "consumer"
	RequestID  = "request-id"
	ServiceFee = "service-fee"
	SlashCoins = "service-slash-coins"
)
