package simulation

import (
	"fmt"
	"math/rand"

	baseapp "gitee.com/xchain/go-xchain/app"
	sdk "gitee.com/xchain/go-xchain/types"
	"gitee.com/xchain/go-xchain/x/simulation"
	"gitee.com/xchain/go-xchain/x/slashing"
)

// SimulateMsgUnjail
func SimulateMsgUnjail(k slashing.Keeper) simulation.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simulation.Account) (opMsg simulation.OperationMsg, fOps []simulation.FutureOperation, err error) {

		acc := simulation.RandomAcc(r, accs)
		address := sdk.ValAddress(acc.Address)
		msg := slashing.NewMsgUnjail(address)
		if msg.ValidateBasic() != nil {
			return simulation.NoOpMsg(), nil, fmt.Errorf("expected msg to pass ValidateBasic: %s", msg.GetSignBytes())
		}
		ctx, write := ctx.CacheContext()
		ok := slashing.NewHandler(k)(ctx, msg).IsOK()
		if ok {
			write()
		}
		opMsg = simulation.NewOperationMsg(msg, ok, "")
		return opMsg, nil, nil
	}
}
