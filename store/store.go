package store

import (
	dbm "github.com/tendermint/tendermint/libs/db"

	"gitee.com/xchain/go-xchain/store/rootmulti"
	"gitee.com/xchain/go-xchain/store/types"
)

func NewCommitMultiStore(db dbm.DB) types.CommitMultiStore {
	return rootmulti.NewStore(db)
}

func NewPruningOptionsFromString(strategy string) (opt PruningOptions) {
	switch strategy {
	case "nothing":
		opt = PruneNothing
	case "everything":
		opt = PruneEverything
	case "syncable":
		opt = PruneSyncable
	default:
		opt = PruneSyncable
	}
	return
}
