package rest

import (
	"github.com/gorilla/mux"

	"gitee.com/xchain/go-xchain/client/context"
	"gitee.com/xchain/go-xchain/codec"
)

// RegisterRoutes registers minting module REST handlers on the provided router.
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	registerQueryRoutes(cliCtx, r, cdc)
}
