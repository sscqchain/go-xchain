package faucet

import (
	"gitee.com/xchain/go-xchain/codec"
)

// RegisterCodec registers concrete types on the Amino codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgAdd{}, "sscq/add", nil)
}
