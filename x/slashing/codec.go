package slashing

import (
	"gitee.com/xchain/go-xchain/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgUnjail{}, "sscq/MsgUnjail", nil)
}

var cdcEmpty = codec.New()
