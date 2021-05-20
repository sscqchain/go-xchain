package params

import (
	"gitee.com/xchain/go-xchain/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*ParamSet)(nil), nil)
}
