package rest

import (
	"gitee.com/xchain/go-xchain/client/context"
	"gitee.com/xchain/go-xchain/codec"
	"gitee.com/xchain/go-xchain/crypto/keys"

	"github.com/gorilla/mux"
)

// RegisterRoutes registers staking-related REST handlers to a router
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec, kb keys.Keybase) {
	registerQueryRoutes(cliCtx, r, cdc)
	registerTxRoutes(cliCtx, r, cdc, kb)
}
