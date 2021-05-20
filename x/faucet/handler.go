package faucet

import (
	"encoding/json"
	"os"

	sdk "gitee.com/xchain/go-xchain/types"
	"gitee.com/xchain/go-xchain/x/bank"
	log "github.com/sirupsen/logrus"
)

func init() {
	// junying-todo,2020-01-17
	lvl, ok := os.LookupEnv("LOG_LEVEL")
	// LOG_LEVEL not set, let's default to debug
	if !ok {
		lvl = "info" //trace/debug/info/warn/error/parse/fatal/panic
	}
	// parse string, this is built-in feature of logrus
	ll, err := log.ParseLevel(lvl)
	if err != nil {
		ll = log.FatalLevel //TraceLevel/DebugLevel/InfoLevel/WarnLevel/ErrorLevel/ParseLevel/FatalLevel/PanicLevel
	}
	// set global log level
	log.SetLevel(ll)
	log.SetFormatter(&log.TextFormatter{}) //&log.JSONFormatter{})
}

//
type SendTxResp struct {
	ErrCode         sdk.CodeType `json:"code"`
	ErrMsg          string       `json:"message"`
	ContractAddress string       `json:"contract_address"`
	EvmOutput       string       `json:"evm_output"`
}

//
func (rsp SendTxResp) String() string {
	rsp.ErrMsg = sdk.GetErrMsg(rsp.ErrCode)
	data, _ := json.Marshal(&rsp)
	return string(data)
}

// New HTDF Message Handler
// connected to handler.go
// HandleMsgSend, HandleMsgAdd upgraded to EVM version
// commented by junying, 2019-08-21
func NewHandler(bankKeeper bank.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {

		switch msg := msg.(type) {
		case MsgAdd:
			return HandleMsgAdd(ctx, bankKeeper, msg)
		default:
			return HandleUnknownMsg(msg)
		}
	}

}

// junying-todo, 2019-08-26
func HandleUnknownMsg(msg sdk.Msg) sdk.Result {
	var sendTxResp SendTxResp
	log.Debugf("msgType error|mstType=%v\n", msg.Type())
	sendTxResp.ErrCode = sdk.ErrCode_Param
	return sdk.Result{Code: sendTxResp.ErrCode, Log: sendTxResp.String()}
}

// Handle a message to add
func HandleMsgAdd(ctx sdk.Context, keeper bank.Keeper, msg MsgAdd) sdk.Result {
	// CurSystemIssuer, err := GetSystemIssuerFromRoot()
	// if err != nil {
	// 	return sdk.NewError("htdfservice", 101, "system_issuer failed to be found or genesis.json doesn't exists").Result()
	// }

	// if !msg.SystemIssuer.Equals(CurSystemIssuer) {
	// 	return sdk.NewError("htdfservice", 101, "requester is not the system_issuer").Result()
	// }

	_, tags, err := keeper.AddCoins(ctx, msg.ToAddress, msg.Amount)
	if err != nil {
		return sdk.NewError("htdfservice", 101, "keeper failed to add requested amount").Result()
	}
	return sdk.Result{
		Tags: tags,
	}
}
