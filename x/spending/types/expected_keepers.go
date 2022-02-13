package types

import (
	govtypes "github.com/KiraCore/sekai/x/gov/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// CustomGovKeeper defines the expected interface contract the tokens module requires
type CustomGovKeeper interface {
	GetNetworkActorsByRole(ctx sdk.Context, role uint64) sdk.Iterator
	GetNetworkActorByAddress(ctx sdk.Context, address sdk.AccAddress) (govtypes.NetworkActor, bool)
	CheckIfAllowedPermission(ctx sdk.Context, addr sdk.AccAddress, permValue govtypes.PermValue) bool
	GetNextProposalIDAndIncrement(ctx sdk.Context) uint64
	GetNetworkProperties(ctx sdk.Context) *govtypes.NetworkProperties
	SaveProposal(ctx sdk.Context, proposal govtypes.Proposal)
	AddToActiveProposals(ctx sdk.Context, proposal govtypes.Proposal)
	CreateAndSaveProposalWithContent(ctx sdk.Context, title, description string, content govtypes.Content) (uint64, error)
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
}
