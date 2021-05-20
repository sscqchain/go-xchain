package upgrade

import (
	"gitee.com/xchain/go-xchain/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(&VersionInfo{}, "sscq/upgrade/VersionInfo", nil)
}

var msgCdc = codec.New()

func init() {
	RegisterCodec(msgCdc)
}
