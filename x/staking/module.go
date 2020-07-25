package staking

import (
	"github.com/KiraCore/cosmos-sdk/codec"
	sdk "github.com/KiraCore/cosmos-sdk/types"
	"github.com/KiraCore/cosmos-sdk/types/module"
	"github.com/KiraCore/cosmos-sdk/x/auth"
	"github.com/KiraCore/cosmos-sdk/x/staking"
	"github.com/KiraCore/cosmos-sdk/x/staking/types"
	"github.com/KiraCore/sekai/x/staking/keeper"
)

var (
	_ module.AppModule = AppModule{}
)

// AppModule extends the cosmos SDK staking.
type AppModule struct {
	staking.AppModule

	stakingKeeper       staking.Keeper
	customStakingKeeper keeper.Keeper
}

// NewHandler returns an sdk.Handler for the staking module.
func (am AppModule) NewHandler() sdk.Handler {
	return NewHandler(am.stakingKeeper, am.customStakingKeeper)
}

// NewAppModule returns a new Custom Staking module.
func NewAppModule(
	cdc codec.Marshaler,
	keeper staking.Keeper,
	ak auth.AccountKeeper,
	bk types.BankKeeper,
) AppModule {
	return AppModule{
		AppModule:     staking.NewAppModule(cdc, keeper, ak, bk),
		stakingKeeper: keeper,
	}
}
