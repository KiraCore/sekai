package layer2

import (
	kiratypes "github.com/KiraCore/sekai/types"
	govtypes "github.com/KiraCore/sekai/x/gov/types"
	"github.com/KiraCore/sekai/x/layer2/keeper"
	"github.com/KiraCore/sekai/x/layer2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ApplyJoinDappProposalHandler struct {
	keeper keeper.Keeper
}

func NewApplyJoinDappProposalHandler(keeper keeper.Keeper) *ApplyJoinDappProposalHandler {
	return &ApplyJoinDappProposalHandler{
		keeper: keeper,
	}
}

func (a ApplyJoinDappProposalHandler) ProposalType() string {
	return kiratypes.ProposalTypeJoinDapp
}

func (a ApplyJoinDappProposalHandler) IsAllowedAddress(ctx sdk.Context, addr sdk.AccAddress, proposal govtypes.Content) bool {
	p := proposal.(*types.ProposalJoinDapp)

	dapp := a.keeper.GetDapp(ctx, p.DappName)
	if dapp.Name == "" {
		return false
	}

	return a.keeper.IsAllowedAddress(ctx, addr, dapp.Controllers)
}

func (a ApplyJoinDappProposalHandler) AllowedAddresses(ctx sdk.Context, proposal govtypes.Content) []string {
	p := proposal.(*types.ProposalJoinDapp)

	dapp := a.keeper.GetDapp(ctx, p.DappName)
	if dapp.Name == "" {
		return []string{}
	}
	return a.keeper.AllowedAddresses(ctx, dapp.Controllers)
}

func (a ApplyJoinDappProposalHandler) Quorum(ctx sdk.Context, proposal govtypes.Content) uint64 {
	p := proposal.(*types.ProposalJoinDapp)

	dapp := a.keeper.GetDapp(ctx, p.DappName)
	if dapp.Name == "" {
		return 0
	}

	return dapp.VoteQuorum
}

func (a ApplyJoinDappProposalHandler) VotePeriod(ctx sdk.Context, proposal govtypes.Content) uint64 {
	p := proposal.(*types.ProposalJoinDapp)

	dapp := a.keeper.GetDapp(ctx, p.DappName)
	if dapp.Name == "" {
		return 0
	}

	return dapp.VotePeriod
}

func (a ApplyJoinDappProposalHandler) VoteEnactment(ctx sdk.Context, proposal govtypes.Content) uint64 {
	p := proposal.(*types.ProposalJoinDapp)

	dapp := a.keeper.GetDapp(ctx, p.DappName)
	if dapp.Name == "" {
		return 0
	}

	return dapp.VoteEnactment
}

func (a ApplyJoinDappProposalHandler) Apply(ctx sdk.Context, proposalID uint64, proposal govtypes.Content, slash sdk.Dec) error {
	p := proposal.(*types.ProposalJoinDapp)

	return a.keeper.ExecuteJoinDappProposal(ctx, p)
}

type ApplyTransitionDappProposalHandler struct {
	keeper keeper.Keeper
}

func NewApplyTransitionDappProposalHandler(keeper keeper.Keeper) *ApplyTransitionDappProposalHandler {
	return &ApplyTransitionDappProposalHandler{
		keeper: keeper,
	}
}

func (a ApplyTransitionDappProposalHandler) ProposalType() string {
	return kiratypes.ProposalTypeTransitionDapp
}

func (a ApplyTransitionDappProposalHandler) IsAllowedAddress(ctx sdk.Context, addr sdk.AccAddress, proposal govtypes.Content) bool {
	p := proposal.(*types.ProposalTransitionDapp)

	dapp := a.keeper.GetDapp(ctx, p.DappName)
	if dapp.Name == "" {
		return false
	}
	// TODO: probably, dapp transition will need to use raw Messages
	return a.keeper.IsAllowedAddress(ctx, addr, dapp.Controllers)
}

func (a ApplyTransitionDappProposalHandler) AllowedAddresses(ctx sdk.Context, proposal govtypes.Content) []string {
	p := proposal.(*types.ProposalTransitionDapp)

	operators := a.keeper.GetDappOperators(ctx, p.DappName)
	addrs := []string{}
	for _, operator := range operators {
		addrs = append(addrs, operator.Operator)
	}
	return addrs
}

func (a ApplyTransitionDappProposalHandler) Quorum(ctx sdk.Context, proposal govtypes.Content) uint64 {
	p := proposal.(*types.ProposalTransitionDapp)

	dapp := a.keeper.GetDapp(ctx, p.DappName)
	if dapp.Name == "" {
		return 0
	}

	// TODO: probably, dapp transition will need to use raw Messages
	return dapp.VerifiersMin
}

func (a ApplyTransitionDappProposalHandler) VotePeriod(ctx sdk.Context, proposal govtypes.Content) uint64 {
	p := proposal.(*types.ProposalTransitionDapp)

	dapp := a.keeper.GetDapp(ctx, p.DappName)
	if dapp.Name == "" {
		return 0
	}

	// TODO: probably, dapp transition will need to use raw Messages
	return dapp.VotePeriod
}

func (a ApplyTransitionDappProposalHandler) VoteEnactment(ctx sdk.Context, proposal govtypes.Content) uint64 {
	p := proposal.(*types.ProposalTransitionDapp)

	dapp := a.keeper.GetDapp(ctx, p.DappName)
	if dapp.Name == "" {
		return 0
	}

	// TODO: probably, dapp transition will need to use raw Messages
	return dapp.VoteEnactment
}

func (a ApplyTransitionDappProposalHandler) Apply(ctx sdk.Context, proposalID uint64, proposal govtypes.Content, slash sdk.Dec) error {
	p := proposal.(*types.ProposalTransitionDapp)

	dapp := a.keeper.GetDapp(ctx, p.DappName)
	if dapp.Name == "" {
		return types.ErrDappDoesNotExist
	}

	dapp.StatusHash = p.StatusHash
	a.keeper.SetDapp(ctx, dapp)
	// TODO: probably, dapp transition will need to use raw Messages

	// TODO: handle operator rank/streak to be similar to validator rank/streak calculation
	// For both verifiers and executors, we will utilize dApp performance counters in a similar manner
	// to which we are utilizing validator performance counters to determine “inactive” operators
	// and ranks (the difference is that there is no need for a mischance confidence counter,
	// only mischance alone since sending a verification tx will not be probabilistic like
	// in the case of proposing a block).

	// **Performance Counters**

	// We will begin by defining the following [Network Properties](https://www.notion.so/de74fe4b731a47df86683f2e9eefa793):
	// `dAppMischanceRankDecreaseAmount`, `dAppMaxMischance`, `dAppInactiveRankDecreasePercent`.
	// The ranking and performance counter system will be based on the dApp session change submission and its verification.
	// If the dApp operator participates in the production of a dApp session (sends session or verification tx) his
	// `rank` and `streak` must be increased by `1` while `mischance` re-set to 0,
	// otherwise in the case of failure to participate the `mischance` counter must be increased by `1`,
	// the `streak` re-set to `0` and `rank` decreased by `dAppMischanceRankDecreaseAmount`.
	// If the dApp operator does not enable maintenance mode by using the `pause-dapp-tx` and the `mischance` counter
	// exceeds `dAppMaxMischance` then his ranks should be slashed by `dAppInactiveRankDecreasePercent`.
	// Alongside `rank`, `streak`, and `mischance` we also need to include `verified_sessions_counter` increased by `1`
	//  every time the verifier or executor submits verification tx,  `created_sessions_counter` every time the executor
	//  proposed a new session, and `missed_sessions_counter` whenever verification was missed by the verifier or executor.
	// All the performance counter values must be positive integers, meaning ranks can’t be negative and the smallest possible rank
	// will be `0`.

	// 	properties := k.gk.GetNetworkProperties(ctx)
	//     // Update uptime counter
	//     missed := !signed
	//     if missed { // increment counter
	//         signInfo.MissedBlocksCounter++
	//         // increment mischance only when missed blocks are bigger than mischance confidence
	//         if signInfo.MischanceConfidence >= int64(properties.MischanceConfidence) {
	//             signInfo.Mischance++
	//         } else {
	//             signInfo.MischanceConfidence++
	//         }
	//     } else { // set counter to 0
	//         signInfo.Mischance = 0
	//         signInfo.MischanceConfidence = 0
	//         signInfo.ProducedBlocksCounter++
	//         signInfo.LastPresentBlock = ctx.BlockHeight()
	//     }

	//     // handle staking module's validator object update actions
	//     k.sk.HandleValidatorSignature(ctx, validator.ValKey, missed, signInfo.Mischance)

	//     // HandleValidatorSignature manage rank and streak by block miss / sign result
	// func (k Keeper) HandleValidatorSignature(ctx sdk.Context, valAddress sdk.ValAddress, missed bool, mischance int64) error {
	//     validator, err := k.GetValidator(ctx, valAddress)
	//     if err != nil {
	//         return err
	//     }
	//     networkProperties := k.govkeeper.GetNetworkProperties(ctx)
	//     if missed {
	//         if mischance > 0 { // it means mischance confidence is set, we update streak and rank properties
	//             // set validator streak by 0 and decrease rank by X
	//             validator.Streak = 0
	//             validator.Rank -= int64(networkProperties.MischanceRankDecreaseAmount)
	//             if validator.Rank < 0 {
	//                 validator.Rank = 0
	//             }
	//         }
	//     } else {
	//         // increase streak and reset rank if streak is higher than rank
	//         validator.Streak++
	//         if validator.Streak > validator.Rank {
	//             validator.Rank = validator.Streak
	//         }
	//     }
	//     k.AddValidator(ctx, validator)
	//     return nil
	// }
	return nil
}

type ApplyUpsertDappProposalHandler struct {
	keeper keeper.Keeper
}

func NewApplyUpsertDappProposalHandler(keeper keeper.Keeper) *ApplyUpsertDappProposalHandler {
	return &ApplyUpsertDappProposalHandler{
		keeper: keeper,
	}
}

func (a ApplyUpsertDappProposalHandler) ProposalType() string {
	return kiratypes.ProposalTypeUpsertDapp
}

func (a ApplyUpsertDappProposalHandler) IsAllowedAddress(ctx sdk.Context, addr sdk.AccAddress, proposal govtypes.Content) bool {
	p := proposal.(*types.ProposalUpsertDapp)

	dapp := a.keeper.GetDapp(ctx, p.Dapp.Name)
	if dapp.Name == "" {
		return false
	}

	return a.keeper.IsAllowedAddress(ctx, addr, dapp.Controllers)
}

func (a ApplyUpsertDappProposalHandler) AllowedAddresses(ctx sdk.Context, proposal govtypes.Content) []string {
	p := proposal.(*types.ProposalUpsertDapp)

	dapp := a.keeper.GetDapp(ctx, p.Dapp.Name)
	if dapp.Name == "" {
		return []string{}
	}

	return a.keeper.AllowedAddresses(ctx, dapp.Controllers)
}

func (a ApplyUpsertDappProposalHandler) Quorum(ctx sdk.Context, proposal govtypes.Content) uint64 {
	p := proposal.(*types.ProposalUpsertDapp)

	dapp := a.keeper.GetDapp(ctx, p.Dapp.Name)
	if dapp.Name == "" {
		return 0
	}

	return dapp.VoteQuorum
}

func (a ApplyUpsertDappProposalHandler) VotePeriod(ctx sdk.Context, proposal govtypes.Content) uint64 {
	p := proposal.(*types.ProposalUpsertDapp)

	dapp := a.keeper.GetDapp(ctx, p.Dapp.Name)
	if dapp.Name == "" {
		return 0
	}

	return dapp.VotePeriod
}

func (a ApplyUpsertDappProposalHandler) VoteEnactment(ctx sdk.Context, proposal govtypes.Content) uint64 {
	p := proposal.(*types.ProposalUpsertDapp)

	dapp := a.keeper.GetDapp(ctx, p.Dapp.Name)
	if dapp.Name == "" {
		return 0
	}

	return dapp.VoteEnactment
}

func (a ApplyUpsertDappProposalHandler) Apply(ctx sdk.Context, proposalID uint64, proposal govtypes.Content, slash sdk.Dec) error {
	p := proposal.(*types.ProposalUpsertDapp)

	dapp := a.keeper.GetDapp(ctx, p.Dapp.Name)
	if dapp.Name == "" {
		return types.ErrDappDoesNotExist
	}

	dapp.Name = p.Dapp.Name
	dapp.Denom = p.Dapp.Denom
	dapp.Description = p.Dapp.Description
	dapp.Website = p.Dapp.Website
	dapp.Logo = p.Dapp.Logo
	dapp.Social = p.Dapp.Social
	dapp.Docs = p.Dapp.Docs
	dapp.Controllers = p.Dapp.Controllers
	dapp.Bin = p.Dapp.Bin
	dapp.Pool = p.Dapp.Pool
	dapp.Issurance = p.Dapp.Issurance
	dapp.UpdateTimeMax = p.Dapp.UpdateTimeMax
	dapp.ExecutorsMin = p.Dapp.ExecutorsMin
	dapp.ExecutorsMax = p.Dapp.ExecutorsMax
	dapp.VerifiersMin = p.Dapp.VerifiersMin
	dapp.TotalBond = p.Dapp.TotalBond
	dapp.CreationTime = p.Dapp.CreationTime
	dapp.StatusHash = p.Dapp.StatusHash
	dapp.Status = p.Dapp.Status
	dapp.VoteQuorum = p.Dapp.VoteQuorum
	dapp.VotePeriod = p.Dapp.VotePeriod
	dapp.VoteEnactment = p.Dapp.VoteEnactment

	a.keeper.SetDapp(ctx, p.Dapp)
	return nil
}
