package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/crypto"

	"github.com/KiraCore/sekai/x/slashing/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AfterValidatorCreated adds the address-pubkey relation when a validator is created.
func (k Keeper) AfterValidatorCreated(ctx sdk.Context, valAddr sdk.ValAddress) error {
	fmt.Println("AfterValidatorCreated.hooks1")
	validator, err := k.sk.GetValidator(ctx, valAddr)
	fmt.Println("AfterValidatorCreated.hooks2", validator.ValKey.String(), err.Error())
	validators := k.sk.GetValidatorSet(ctx)
	fmt.Println("registered validators count", len(validators))
	for i, val := range validators {
		fmt.Println("registered validators", i, val.ValKey.String())
	}
	if err != nil {
		return err
	}

	consPk, err := validator.TmConsPubKey()
	fmt.Println("AfterValidatorCreated.hooks3", consPk.Address().String())
	if err != nil {
		return err
	}
	fmt.Println("AfterValidatorCreated.hooks4")
	k.AddPubkey(ctx, consPk)
	return nil
}

// AfterValidatorRemoved deletes the address-pubkey relation when a validator is removed,
func (k Keeper) AfterValidatorRemoved(ctx sdk.Context, address sdk.ConsAddress) {
	k.deleteAddrPubkeyRelation(ctx, crypto.Address(address))
}

//_________________________________________________________________________________________

// Hooks wrapper struct for slashing keeper
type Hooks struct {
	k Keeper
}

var _ types.StakingHooks = Hooks{}

// Return the wrapper struct
func (k Keeper) Hooks() Hooks {
	return Hooks{k}
}

// Implements sdk.ValidatorHooks
func (h Hooks) AfterValidatorRemoved(ctx sdk.Context, consAddr sdk.ConsAddress, _ sdk.ValAddress) {
	h.k.AfterValidatorRemoved(ctx, consAddr)
}

// Implements sdk.ValidatorHooks
func (h Hooks) AfterValidatorCreated(ctx sdk.Context, valAddr sdk.ValAddress) {
	h.k.AfterValidatorCreated(ctx, valAddr)
}

func (h Hooks) BeforeValidatorModified(_ sdk.Context, _ sdk.ValAddress) {}
