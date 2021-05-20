package rest

import (
	"net/http"

	"gitee.com/xchain/go-xchain/accounts/keystore"
	"gitee.com/xchain/go-xchain/client/context"
	"gitee.com/xchain/go-xchain/client/utils"
	"gitee.com/xchain/go-xchain/codec"
	"gitee.com/xchain/go-xchain/crypto/keys/keyerror"
	sdk "gitee.com/xchain/go-xchain/types"
	"gitee.com/xchain/go-xchain/types/rest"
	"gitee.com/xchain/go-xchain/x/auth"
	authtxb "gitee.com/xchain/go-xchain/x/auth/client/txbuilder"
	sscqservice "gitee.com/xchain/go-xchain/x/core"
	sscorecli "gitee.com/xchain/go-xchain/x/core/client/cli"
)

// SignBody defines the properties of a sign request's body.
type SignBody struct {
	Tx         auth.StdTx   `json:"tx"`
	BaseReq    rest.BaseReq `json:"base_req"`
	Passphrase string       `json:"passphrase"`
}

// nolint: unparam
// sign tx REST handler
func SignTxRequestHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req SignBody

		if !rest.ReadRESTReq(w, r, cdc, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		// Santize
		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		// validate tx
		// discard error if it's CodeNoSignatures as the tx comes with no signatures
		if err := req.Tx.ValidateBasic(); err != nil && err.Code() != sdk.CodeNoSignatures {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// derive the from account address and name from the Keybase

		txBldr := authtxb.NewTxBuilder(
			utils.GetTxEncoder(cdc),
			req.BaseReq.AccountNumber,
			req.BaseReq.Sequence,
			req.Tx.Fee.GasWanted,
			1.0,
			false,
			req.BaseReq.ChainID,
			req.Tx.GetMemo(),
			req.Tx.Fee.GasPrice,
		)
		
		var signedTx auth.StdTx
		addr := req.BaseReq.From
		ksw := keystore.NewKeyStoreWallet(keystore.DefaultKeyStoreHome())
		signedTx, err := ksw.SignStdTx(txBldr, signedTx, addr, req.Passphrase)
		if keyerror.IsErrKeyNotFound(err) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		} else if keyerror.IsErrWrongPassword(err) {
			rest.WriteErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		} else if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cdc, signedTx, cliCtx.Indent)
	}
}

// SignBody defines the properties of a sign request's body.
type SignRawBody struct {
	Tx         string `json:"tx"`
	Passphrase string `json:"passphrase"`
	Offline    bool   `json:"offline"`
	Encode     bool   `json:"encode"`
}

// sign tx REST handler
func SignTxRawRequestHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req SignRawBody

		if !rest.ReadRESTReq(w, r, cdc, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		// load sign tx from string
		stdTx, err := sscqservice.ReadStdTxFromRawData(cliCtx.Codec, req.Tx)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "transaction decode failed")
			return
		}

		// derive the from account address and name from the Keybase
		if len(stdTx.GetSigners()) == 0 {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "signer not found")
			return
		}

		// sign
		res, err := sscorecli.SignTransaction(authtxb.NewTxBuilderFromCLI(), cliCtx, stdTx, req.Passphrase, req.Offline)
		if err != nil {
			return
		}

		// response
		if !req.Encode {
			w.Write(res)
		} else {
			w.Write([]byte(sscqservice.Encode_Hex(res)))
		}
	}
}
