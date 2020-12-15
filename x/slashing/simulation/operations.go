package simulation

import (
	"errors"
	"math/rand"

	"github.com/KiraCore/sekai/x/slashing/keeper"
	"github.com/KiraCore/sekai/x/slashing/types"
	stakingkeeper "github.com/KiraCore/sekai/x/staking/keeper"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp/helpers"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
)

// Simulation operation weights constants
const (
	OpWeightMsgActivate = "op_weight_msg_activate"
)

// WeightedOperations returns all the operations from the module with their respective weights
func WeightedOperations(
	appParams simtypes.AppParams, cdc codec.JSONMarshaler, ak types.AccountKeeper,
	bk types.BankKeeper, k keeper.Keeper, sk stakingkeeper.Keeper,
) simulation.WeightedOperations {

	var weightMsgActivate int
	appParams.GetOrGenerate(cdc, OpWeightMsgActivate, &weightMsgActivate, nil,
		func(_ *rand.Rand) {
			weightMsgActivate = simappparams.DefaultWeightMsgUnjail
		},
	)

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(
			weightMsgActivate,
			SimulateMsgActivate(ak, bk, k, sk),
		),
	}
}

// SimulateMsgActivate generates a MsgActivate with random values
// nolint: interfacer
func SimulateMsgActivate(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper, sk stakingkeeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		validator, ok := stakingkeeper.RandomValidator(r, sk, ctx)
		if !ok {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgActivate, "validator is not ok"), nil, nil // skip
		}

		simAccount, found := simtypes.FindAccount(accs, sdk.AccAddress(validator.GetOperator()))
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgActivate, "unable to find account"), nil, nil // skip
		}

		if !validator.IsJailed() {
			// TODO: due to this condition this message is almost, if not always, skipped !
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgActivate, "validator is not jailed"), nil, nil
		}

		cons, err := validator.TmConsPubKey()
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgActivate, "unable to get validator consensus key"), nil, err
		}
		consAddr := sdk.ConsAddress(cons.Address())
		info, found := k.GetValidatorSigningInfo(ctx, consAddr)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgActivate, "unable to find validator signing info"), nil, nil // skip
		}

		account := ak.GetAccount(ctx, sdk.AccAddress(validator.GetOperator()))
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgActivate, "unable to generate fees"), nil, err
		}

		msg := types.NewMsgActivate(validator.GetOperator())

		txGen := simappparams.MakeTestEncodingConfig().TxConfig
		tx, err := helpers.GenTx(
			txGen,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			simAccount.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate mock tx"), nil, err
		}

		_, res, err := app.Deliver(txGen.TxEncoder(), tx)

		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to deliver tx"), nil, errors.New(res.Log)
		}

		return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}
