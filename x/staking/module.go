package staking

import (
	"github.com/KiraCore/cosmos-sdk/codec"
	"github.com/KiraCore/cosmos-sdk/types/module"
	"github.com/KiraCore/cosmos-sdk/x/auth"
	"github.com/KiraCore/cosmos-sdk/x/staking"
	"github.com/KiraCore/cosmos-sdk/x/staking/types"
)

var (
	_ module.AppModule = AppModule{}
)

// AppModule extends the cosmos SDK staking.
type AppModule struct {
	staking.AppModule
}

// NewAppModule returns a new Custom Staking module.
func NewAppModule(
	cdc codec.Marshaler,
	keeper staking.Keeper,
	ak auth.AccountKeeper,
	bk types.BankKeeper,
) AppModule {
	return AppModule{
		AppModule: staking.NewAppModule(cdc, keeper, ak, bk),
	}
}
