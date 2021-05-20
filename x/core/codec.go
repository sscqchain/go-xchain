package sscqservice

import (
	"gitee.com/xchain/go-xchain/codec"
)

// RegisterCodec registers concrete types on the Amino codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSend{}, "sscqservice/send", nil)
	// cdc.RegisterConcrete(MsgAdd{}, "sscqservice/add", nil)
}
