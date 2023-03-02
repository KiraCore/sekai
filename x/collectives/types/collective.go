package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

func (r Collective) GetCollectiveAddress() sdk.AccAddress {
	return authtypes.NewModuleAddress(fmt.Sprintf("collective_account_%s", r.Name))
}

func (r Collective) GetCollectiveDonationAddress() sdk.AccAddress {
	return authtypes.NewModuleAddress(fmt.Sprintf("collective_donation_%s", r.Name))
}
