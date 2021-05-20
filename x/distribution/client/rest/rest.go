package rest

import (
	"github.com/gorilla/mux"

	"gitee.com/xchain/go-xchain/client/context"
	"gitee.com/xchain/go-xchain/codec"
)

// RegisterRoutes register distribution REST routes.
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec, queryRoute string) {
	registerQueryRoutes(cliCtx, r, cdc, queryRoute)
	registerTxRoutes(cliCtx, r, cdc, queryRoute)
}
