package keeper

import (
	"github.com/KiraCore/sekai/x/layer2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetDappSession(ctx sdk.Context, session types.ExecutionRegistrar) {
	bz := k.cdc.MustMarshal(&session)
	store := ctx.KVStore(k.storeKey)
	store.Set(types.ExecutionRegistrarKey(session.DappName), bz)
}

func (k Keeper) DeleteDappSession(ctx sdk.Context, name string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.ExecutionRegistrarKey(name))
}

func (k Keeper) GetDappSession(ctx sdk.Context, name string) types.ExecutionRegistrar {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ExecutionRegistrarKey(name))
	if bz == nil {
		return types.ExecutionRegistrar{}
	}

	sessionInfo := types.ExecutionRegistrar{}
	k.cdc.MustUnmarshal(bz, &sessionInfo)
	return sessionInfo
}

func (k Keeper) GetAllDappSessions(ctx sdk.Context) []types.ExecutionRegistrar {
	store := ctx.KVStore(k.storeKey)

	sessions := []types.ExecutionRegistrar{}
	it := sdk.KVStorePrefixIterator(store, []byte(types.PrefixDappSessionKey))
	defer it.Close()

	for ; it.Valid(); it.Next() {
		session := types.ExecutionRegistrar{}
		k.cdc.MustUnmarshal(it.Value(), &session)
		sessions = append(sessions, session)
	}
	return sessions
}

// Halt the dapp if verifiers and executors number does not meet on ExitDapp and after Jail operation
func (k Keeper) HaltDappIfNoEnoughActiveOperators(ctx sdk.Context, dappName string) bool {
	dapp := k.GetDapp(ctx, dappName)
	executors := k.GetDappExecutors(ctx, dappName)
	verifiers := k.GetDappVerifiers(ctx, dappName)
	// active executors
	activeExecutors := int(0)
	for _, executor := range executors {
		if executor.Status == types.OperatorActive {
			activeExecutors++
		}
	}
	activeVerifiers := int(0)
	for _, verifier := range verifiers {
		if verifier.Status == types.OperatorActive {
			activeVerifiers++
		}
	}

	if activeExecutors < int(dapp.ExecutorsMin) || activeVerifiers < int(dapp.VerifiersMin) {
		dapp.Status = types.Halted
		k.SetDapp(ctx, dapp)
		return true
	}
	return false
}

func (k Keeper) ResetNewSession(ctx sdk.Context, name string, prevLeader string) {
	operators := k.GetDappOperators(ctx, name)
	for _, operator := range operators {
		if operator.Status == types.OperatorExiting {
			// If the operator leaving the dApp was a verifier then
			// as the result of the exit tx his LP tokens bond should be returned once the record is deleted.
			// The bond can only be claimed if and only if the status didn’t change to jailed in the meantime.
			if operator.BondedLpAmount.IsPositive() {
				dapp := k.GetDapp(ctx, name)
				dappLpToken := dapp.LpToken()
				lpCoins := sdk.NewCoins(sdk.NewCoin(dappLpToken, operator.BondedLpAmount))
				addr := sdk.MustAccAddressFromBech32(operator.Operator)
				err := k.bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, lpCoins)
				if err != nil {
					panic(err)
				}
			}
			// remove operators exiting when session ends
			k.DeleteDappOperator(ctx, name, operator.Operator)
		}
	}

	// halt the dapp if operators condition does not meet on session reset
	halted := k.HaltDappIfNoEnoughActiveOperators(ctx, name)
	if halted {
		return
	}

	leader := ""
	executors := k.GetDappExecutors(ctx, name)
	if len(executors) > 0 {
		leader = executors[0].Operator
		for index, executor := range executors {
			if executor.Operator == prevLeader {
				leader = executors[(index+1)%len(executors)].Interx
			}
		}
	}
	if leader == "" {
		operator := k.GetDappOperator(ctx, name, prevLeader)
		if operator.Status == types.OperatorActive {
			leader = prevLeader
		}
	}

	session := k.GetDappSession(ctx, name)
	session.NextSession = &types.DappSession{
		Leader:     leader,
		Start:      uint64(ctx.BlockTime().Unix()),
		StatusHash: "",
		Status:     types.SessionScheduled,
	}
	k.SetDappSession(ctx, session)

	// halt the dapp if next session leader is not available
	if session.NextSession.Leader == "" {
		dapp := k.GetDapp(ctx, name)
		dapp.Status = types.Halted
		k.SetDapp(ctx, dapp)
	}
}

func (k Keeper) CreateNewSession(ctx sdk.Context, name string, prevLeader string) {
	session := k.GetDappSession(ctx, name)
	session.PrevSession = session.CurrSession
	session.CurrSession = session.NextSession
	k.SetDappSession(ctx, session)

	// handle bridge and mint messages
	msgServer := NewMsgServerImpl(k)
	for _, msg := range session.PrevSession.OnchainMessages {
		cacheCtx, write := ctx.CacheContext()
		var err error
		switch msg := msg.GetCachedValue().(type) {
		case *types.MsgTransferDappTx:
			_, err = msgServer.TransferDappTx(sdk.WrapSDKContext(cacheCtx), msg)
		case *types.MsgAckTransferDappTx:
			_, err = msgServer.AckTransferDappTx(sdk.WrapSDKContext(ctx), msg)
		case *types.MsgMintCreateFtTx:
			_, err = msgServer.MintCreateFtTx(sdk.WrapSDKContext(ctx), msg)
		case *types.MsgMintCreateNftTx:
			_, err = msgServer.MintCreateNftTx(sdk.WrapSDKContext(ctx), msg)
		case *types.MsgMintIssueTx:
			_, err = msgServer.MintIssueTx(sdk.WrapSDKContext(ctx), msg)
		case *types.MsgMintBurnTx:
			_, err = msgServer.MintBurnTx(sdk.WrapSDKContext(ctx), msg)
		}
		if err == nil {
			write()
		}
	}
	k.ResetNewSession(ctx, name, prevLeader)
}

// **Next Session**

// A summary of possible “next session” states:
// - `unscheduled` - The session is NOT ready yet to be picked up by another dApp Leader (e.g. previous dApp Session is ongoing)
// - `scheduled`- The session is ready to be picked up by the next dApp Leader
// - `ongoing` - Session was claimed by the dApp Leader and is currently being executed

// **Current Session**

// Allow every executor to send a `denounce-leader-tx` (true/false) which at any point in time via majority vote can invalidate
// further changes to the current session
// and will cause the next session to change its status to `scheduled`.

// It must be possible for the automatic denouncement of the leader to happen if the current session data is not updated within
// `update_time_max` since `sessions.next.data.time` while the executor
// did NOT change his status to `paused`
// If automatic denouement happens then the same as in the case of manual votes of denouement the current session status
// should change to `denounced` while the next session status should be set to `scheduled`.

// The most important property of the current session is `data`.
// The objective of the leader is to periodically update the `sessions.current.data` property with hashes of the `old` state,
// user `input`, the `new` proposed state, expected changes to account balances,
// and optional `proof` of the correctness of the execution

// The `transition-dapp-tx` can be submitted ONLY by the leader and can be sent at ANY point in time allowing for finality
// the moment a sufficient number of verifiers or executors accept the changes**.

// Whenever leader submits the new dApp state transition the `version` MUST be included and match the application version**.
// The transition tx should fail if the dApp `version` is NOT correct, the leader status is NOT active or
// if the `old` state hash does NOT match the `sessions.previous.data.new` hash.

// There are no limitations to how many times the current session leader can submit `transition-dapp-tx` to replace the content
// of the current session data, however, every time a new dApp state transition is submitted,
// the list of approvals must be wiped, meaning any verifications will be lost.

// Regardless if the approvals are lost or not the performance counters should maintain their count allowing for a fair
// reward distribution to operators later on.
// Additionally, the dApp leader must send the `transition-dapp-tx` before `update_time_max` elapses,
// otherwise his session will expire, and sending of any further `transition-dapp-tx` should fail, and the
// [sessions.currend.data.final](http://sessions.currend.data.final)` flag should be set to true.
// Finally, the session leader must be able to include a boolean flag in the `transition-dapp-tx`
// indicating if the transition tx is final on his own.
// If the final session flag is set then the leader should NOT be allowed to submit another
// `transition-dapp-tx` to the current session.

// If the non-final session is approved before the finality flag is set and the session changes state to `approved`
// then the default next session leader should become the current session leader.
// This way the operator can continue execution with all data already available to him and provide fully uninterrupted service
// to the users.
// In order to give verifiers sufficient time to approve non-final sessions we should prohibit the submission of new
// `transition-dapp-tx` by the leader unless no less than `update_time_max/2` seconds elapsed since the last submission.

// A summary of possible “current session” states:

// - `accepted` - Session was accepted by verifiers and the new state is irreversible (changes will be applied to the blockchain unless status changes to halted or failed)
// - `ongoing` - Session was claimed by the dApp Leader and is currently being executed
// - `denounced` - The session was rejected by either a manual vote of executors to change the leader or due to the leader becoming jailed.
// - `halted` - The session was halted because it was questioned by one of the verifiers submitting evidence against it. Validators will have to assist unhalting and decide slashing penalty to either validator or fisherman
// - `failed` - The session failed to transition from “current” into “previous session” because internal data submitted by the leader was invalid and could NOT result in the modification of the blockchain state.

// In the case where the current dApp Session becomes `failed` or `halted` but the new Leader already started execution based on an unverified state then both sessions (current and next) should be rejected by the network and a consecutive dApp Leader MUST start the execution from the dApp State that was persisted as a result of previous Session. If by any chance resources from the previous dApp state are no longer accessible (e.g. IPFS/URL link does not work and no validator has it saved on their execution node) then users can begin recovery of funds in accordance to the latest known settlement of balances within the Execution Registrar.

// **Previous Session**

// In order for the already approved “current session” to become a “previous session” the changes proposed in the `data.internal`
// section of the current session properties must be applied to the blockchain state.
// We will differentiate two types of internal data:

// - Application Data (`appdata`) - A simple key-value pair dictionary, each application should be able to store on-chain
// up to `256` keys with up to `8192` characters each.
// The purpose of `appdata` is to enable applications communicate informations to the users that can be easily accessed
// and can be trusted, this might include location of IPFS gateways preserving dApp state files,
// deposit addresses of bridges, communicates to other applications and more.
//     - JSON Structure:

//         ```json
//         { "<key-1>": "<value-1>", "<key-2>": "<value-2>", ... }
//         ```

//     - Key names must be unique and adhere to the same rules as IR keys
//     - Values must be a string or set to integer value `0` indicating that the key should be deleted
//     - Compression Algorithm: `zip (max compression, level 9)`
//     - Format: `base64`
// - Cross-Application Message Data (`xamdata`) - A simple json structure containing changes that must be communicated
// to the blockchain and in particular to the Application Bridge Registrar (ABR).
// The ABR is responsible for maintaining balances of users who deposit tokens to various dApps,
// minting new tokens and sending communicates between dApps & internal modules.
//     - JSON Structure:

//         ```json
//         { // key-array list defninig changes that must be applied to the Application Bridge Registrar
//         	"<key-1>": [ "<value-1>", "<value-2>", "<value-3>", ...  ], ... }
//         ```

//     - Compression Algorithm: `zip (max compression, level 9)`
//     - Format: `base64`

// It is NOT guaranteed that even if dApp state was accepted the changes to the blockchain can actually be executed.
// For example application might be requesting to transfer to the user more assets then it actually holds.
// In such cases the State of the “current session” state must change to `failed` and the next session should automatically be wiped
// and set to `scheaduled`.
// The operators must also be explicitly informed what is the exact reason for the failure though `errors`
// property in the current session.
// Any failures should automatically result in the current session leader being denounced,
//  meaning that the next session leader should be changed even if session finalization flag was not set in `transition-dapp-tx`.

// **Reporting Failed Execution**

// Fisherman (verifiers and executors who are not a leader) must be able to report any potential misbehavior of a leader
// and/or otherwise issues with the application to ensure that funds locked within can remain safe.
// In certain cases the finality gadget might not even exists or be possible to be created while fisherman might want
// to create custom systems monitoring behavior of the application far beyond the code logic itself.
// ANY fisherman must be able to send a `reject-dapp-transition-tx` with a message body (up to `8192` characters)
// containing a dApp name, session id and a reason for halting the application.
// As the result of `reject-dapp-transition-tx` the status of the current session should immediately change to
// `halted` and a `dapp-slash` proposal raised, similar to the way that the validator slashing proposals are created.
// In this case however the validators will need to make following judgements:

// - should the dApp session be unhalted and state transition accepted (boolean) - true / false
//     - if `true` then status of session should change to `ongoing` unless sufficient number of approvals is sent then `accepted`
//     - if `false` then status should be set to `failed` and `errors` specify that governance vote decided to reject the session
// - what percentage of the fisherman bond should be slashed (percentage) - from 0 to 1 (100%)
//     - slashed dApp LP tokens should be simply deminted (increasing value of other LP token holders positions)
//     - if fisherman is slashed the validator can NOT be slashed
// - what percentage of the validator stake should be slashed (percentage) - from 0 to 0.1 (10%)
//     - From the amount slashed 1% should be sent to fisherman address as reward and 99% to community pool
//     - if validator is slashed the fisherman can NOT be slashed

// Voting on `dapp-slash` proposal should only be possible by validators **with exception for the dApp leader**.
// If quorum is not reached or result is inconclusive then the default behavior should be no slashing and
// rejection of the dApp session state. If as the result of the proposal the bond of the fisherman was slashed,
// then he should be removed from the fisherman position until sufficient amount is locked by him again.
