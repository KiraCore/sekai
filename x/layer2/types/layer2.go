package types

import (
	fmt "fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (dapp Dapp) LpToken() string {
	return fmt.Sprintf("lp/%s", dapp.Name)
}
func (dapp Dapp) GetSpendingPoolLpDeposit() sdk.Int {
	return dapp.TotalBond.Amount.ToDec().Mul(dapp.Pool.Ratio).RoundInt()
}

func (dapp Dapp) GetLpTokenSupply() sdk.Int {
	spendingPoolDeposit := dapp.GetSpendingPoolLpDeposit()
	totalSupply := spendingPoolDeposit.Add(dapp.Issurance.Postmint).Add(dapp.Issurance.Premint)
	return totalSupply
}

func (dapp Dapp) GetAccount() sdk.AccAddress {
	return sdk.AccAddress(dapp.Name)
}

func (dapp Dapp) Version() string {
	if len(dapp.Bin) > 0 {
		return dapp.Bin[0].Hash
	}
	return ""
}
