package types

import sdk "github.com/KiraCore/cosmos-sdk/types"

type Validator struct {
	Moniker  string
	Website  string
	Social   string
	Identity string

	Comission sdk.Dec

	ValKey sdk.ValAddress
	PubKey sdk.AccAddress
}
