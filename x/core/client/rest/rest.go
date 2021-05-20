package rest

import (
	"fmt"
	"net/http"

	"gitee.com/xchain/go-xchain/client/context"
	"gitee.com/xchain/go-xchain/codec"
	sdk "gitee.com/xchain/go-xchain/types"
	"gitee.com/xchain/go-xchain/types/rest"

	svrConfig "gitee.com/xchain/go-xchain/server/config"

	"github.com/gorilla/mux"
	"gitee.com/xchain/go-xchain/x/auth"
	sscqservice "gitee.com/xchain/go-xchain/x/core"
)

const (
	restName = "custom"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec, storeName string) {

	if svrConfig.ApiSecurityLevel == svrConfig.ValueSecurityLevel_Low {
		r.HandleFunc(fmt.Sprintf("/%s/send", storeName), SendTxRequestHandlerFn(cdc, cliCtx)).Methods("POST")
		r.HandleFunc(fmt.Sprintf("/%s/create", storeName), CreateTxRequestHandlerFn(cdc, cliCtx)).Methods("POST")
		r.HandleFunc(fmt.Sprintf("/%s/sign", storeName), SignTxRawRequestHandlerFn(cdc, cliCtx)).Methods("POST")
	}
	r.HandleFunc(fmt.Sprintf("/%s/broadcast", storeName), BroadcastTxRawRequestHandlerFn(cdc, cliCtx)).Methods("POST")
	//
	r.HandleFunc(
		fmt.Sprintf("/%s/contract/{address}/{code}", storeName),
		QueryContractRequestHandlerFn(storeName, cdc, context.GetAccountDecoder(cdc), cliCtx),
	).Methods("GET")
}

// query contractREST Handler
func QueryContractRequestHandlerFn(
	storeName string, cdc *codec.Codec,
	decoder auth.AccountDecoder, cliCtx context.CLIContext,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bech32addr := vars["address"]
		code := vars["code"]

		contractaddr, err := sdk.AccAddressFromBech32(bech32addr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		bz, err := cliCtx.Codec.MarshalJSON(sscqservice.NewQueryContractParams(contractaddr, code))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound,
				"ERROR: height must be integer.")
			return
		}
		//
		route := fmt.Sprintf("custom/%s/%s", sscqservice.QuerierRoute, sscqservice.QueryContract)
		res, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}
