package slashing

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "gitee.com/xchain/go-xchain/types"
)

func TestMsgUnjailGetSignBytes(t *testing.T) {
	addr := sdk.AccAddress("abcd")
	msg := NewMsgUnjail(sdk.ValAddress(addr))
	bytes := msg.GetSignBytes()
	require.Equal(t, string(bytes), `{"address":"sscqvaloper1v93xxeqdmnmfc"}`)
}
