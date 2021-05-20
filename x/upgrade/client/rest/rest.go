package rest

import (
	"fmt"

	"gitee.com/xchain/go-xchain/client/context"
	"gitee.com/xchain/go-xchain/codec"

	"github.com/gorilla/mux"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	r.HandleFunc(fmt.Sprintf("/upgrade_info"), QueryUpgradeInfoRequestHandlerFn(cliCtx, cdc)).Methods("GET")
}
