package keeper

import (
	"github.com/KiraCore/sekai/x/layer2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetDappOperator(ctx sdk.Context, operator types.DappOperator) {
	bz := k.cdc.MustMarshal(&operator)
	store := ctx.KVStore(k.storeKey)
	store.Set(types.DappOperatorKey(operator.DappName, operator.Operator), bz)
}

func (k Keeper) DeleteDappOperator(ctx sdk.Context, name, address string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.DappOperatorKey(name, address))
}

func (k Keeper) GetDappOperator(ctx sdk.Context, name string, user string) types.DappOperator {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.DappOperatorKey(name, user))
	if bz == nil {
		return types.DappOperator{}
	}

	operator := types.DappOperator{}
	k.cdc.MustUnmarshal(bz, &operator)
	return operator
}

func (k Keeper) GetDappOperators(ctx sdk.Context, name string) []types.DappOperator {
	store := ctx.KVStore(k.storeKey)

	operators := []types.DappOperator{}
	it := sdk.KVStorePrefixIterator(store, append([]byte(types.PrefixDappOperatorKey), name...))
	defer it.Close()

	for ; it.Valid(); it.Next() {
		operator := types.DappOperator{}
		k.cdc.MustUnmarshal(it.Value(), &operator)
		operators = append(operators, operator)
	}
	return operators
}

func (k Keeper) GetDappExecutors(ctx sdk.Context, name string) []types.DappOperator {
	operators := k.GetDappOperators(ctx, name)
	executors := []types.DappOperator{}
	for _, operator := range operators {
		if operator.Executor {
			executors = append(executors, operator)
		}
	}
	return executors
}

func (k Keeper) GetAllDappOperators(ctx sdk.Context) []types.DappOperator {
	store := ctx.KVStore(k.storeKey)

	operators := []types.DappOperator{}
	it := sdk.KVStorePrefixIterator(store, []byte(types.PrefixDappOperatorKey))
	defer it.Close()

	for ; it.Valid(); it.Next() {
		operator := types.DappOperator{}
		k.cdc.MustUnmarshal(it.Value(), &operator)
		operators = append(operators, operator)
	}
	return operators
}

func (k Keeper) ExecuteJoinDappProposal(ctx sdk.Context, p *types.ProposalJoinDapp) error {
	dapp := k.GetDapp(ctx, p.DappName)
	if dapp.Name == "" {
		return types.ErrDappDoesNotExist
	}

	if p.Verifier {
		properties := k.gk.GetNetworkProperties(ctx)
		verifierBond := properties.DappVerifierBond
		totalSupply := dapp.GetLpTokenSupply()
		dappBondLpToken := dapp.LpToken()
		verifierBondCoins := sdk.NewCoins(sdk.NewCoin(dappBondLpToken, totalSupply.ToDec().Mul(verifierBond).RoundInt()))
		addr := sdk.MustAccAddressFromBech32(p.Interx)
		if verifierBondCoins.IsAllPositive() {
			err := k.bk.SendCoinsFromAccountToModule(ctx, addr, types.ModuleName, verifierBondCoins)
			if err != nil {
				return err
			}
		}
	}

	if p.Executor {
		// TODO: ensure executor is a validator

		executors := k.GetDappExecutors(ctx, p.DappName)
		if len(executors) >= int(dapp.ExecutorsMax) {
			return types.ErrNumberOfOperatorsExceedsExecutorsMax
		}

		k.SetDappOperator(ctx, types.DappOperator{
			DappName: p.DappName,
			Operator: p.Sender,
			Executor: p.Executor,
			Verifier: p.Verifier,
			Interx:   p.Interx,
			Status:   types.OperatorActive,
		})

		// when executors_min reaches, session is created
		executors = k.GetDappExecutors(ctx, p.DappName)
		if len(executors) >= int(dapp.ExecutorsMin) && len(executors) >= 1 {
			session := k.GetDappSession(ctx, p.DappName)
			if session.DappName == "" {
				k.CreateNewSession(ctx, p.DappName, "")
			}
		}
		return nil
	}

	k.SetDappOperator(ctx, types.DappOperator{
		DappName: p.DappName,
		Operator: p.Sender,
		Executor: p.Executor,
		Verifier: p.Verifier,
		Interx:   p.Interx,
		Status:   types.OperatorActive,
	})
	return nil
}
