# Changelog

## [v0.1.22.6] - 01.11.2021

- sekaid init home config resolve

## [v0.1.22.4] - 12.10.2021

- Add assign claim validator permission and allow claim validator action for the new permission

## [v0.1.22.3] - 12.10.2021

- Upgrade Cosmos SDK to v0.44.2 for chain halt issue fix

## [v0.1.22.1] - 07.10.2021

- Convert identity registrar address proto fields to string
- Resolve pagination issue for all identity record verify requests query

## [v0.1.22] - 04.10.2021

- Upgrade Cosmos SDK to v0.44.1 from v0.42.9

## [v0.1.21.28] - 04.10.2021

- Resolve signing info iterator

## [v0.1.21.27] - 03.10.2021

- Added command for exporting new genesis from old genesis

## [v0.1.21.26] - 03.10.2021

- add README for new genesis file manual generation

## [v0.1.21.25] - 26.09.2021

- Add chain ids for upgrade plan

## [v0.1.21.24] - 26.09.2021

- Resolve tokens module genesis export issue

## [v0.1.21.23] - 26.09.2021

- Add skip handler field for upgrade plan
- Add example script for proposing upgrade with skip handler

## [v0.1.21.22] - 22.09.2021

- Update proto files
- Update node discovery module to search only connected peers

## [v0.1.21.21] - 19.09.2021

- Remove flag height from upgrade module

## [v0.1.21.20] - 19.09.2021

- Add reboot required field for plan
- Restrict the upgrade time to be not less than current block time
- Always halt if InstateUpgrade is set to false

## [v0.1.21.19] - 19.09.2021

- Remove height from upgrade plan
- Add current plan and next plan query and add genesis for the plans
- Resolve export genesis for tokens module
- Add scripts examples for plan modifications

## [v0.1.21.18] - 15.09.2021

- Introduce unique identity required keys into identity registrar
- Changes to allow string network property
- Add unique identity keys property into network properties
- Add basic validation of network properties
- Add example script for testing setting network property for unique identity keys
- Modification into SetNetworkPropertyProposal for new properties added

## [v0.1.21.17] - 14.09.2021

- add addrbook query api on interxd
- add net_info query api on interxd
- fix validators query

## [v0.1.21.16] - 10.09.2021

- Restrict moniker length to be less than 32
- Resolve Identity registrar requests querying by approver and requester
- Test on CLI command for requesters and approver and add examples on scripts
- Remove commission from validator

## [v0.1.21.15] - 08.09.2021

- Add querying records by filtering keys
- add MIN_IDENTITY_APPROVAL_TIP into network properties
- Implement tip checker
- implement auto reject if edited after creating verification request
- add balance check for automatic reject
- Add identity registrar key validation and automatic lowercase converter
- Resolve invalid implementation of SetIdentityRegistrar for genesis initialization process of address key matching to recordId
- Update for moniker field management to use identity registrar
- Add address catcher from identity key record pair
- Resolve error handling of identity registrar to log the errors properly on CLI
- Resolve sample script for identity registrar
- Add CLI command for querying validators and add example cli command for querying
- Add range of tests for the changes

## [v0.1.21.14] - 06.09.2021

- Update pagination limit on sekaid
- Update identity registrar for validators query
- Update snapshot extension
- Add ip_only, connected paramere for node list queries
- Update sort feature of node list queries
- Query all validators using pagination

## [v0.1.21.13] - 27.08.2021

- Identity registrar records structure change
- Identity registrar cli UX changes
- Add scripts example for cli commands
- Resolve tests for identity registrar changes
- Commands description fix

## [v0.1.21.12] - 16.08.2021

- Add description on the script to see the version
- Add reference for identity registrar
- Remove long json script for unjail testing

## [v0.1.21.11] - 11.08.2021

- Upgrade cosmos sdk version to be latest
- Remove all flag in proposals query
- Resolve invalid default home flag
- CLI modifications for missing endpoints
- Remove deprecated use of FlagNode on gov and tokens query
- Resolve CLI test caused by removing unused flag

## [v0.1.21.10] - 11.08.2021

- INTERX: Update node discovery module.
- INTERX: Update proto-gen script.
- INTERX: Update interx metadata query.
- INTERX: Fix voters, DRR, network properties query.
- INTERX: Update configurations/proposal proto types.
- INTERX: Add upgrade plan query. (/api/kira/upgrade/current_plan)

## [v0.1.21.9] - 05.08.2021

- Add pagination for identity registrar grpc_queries
- Modify protogen script to use fixed cosmos sdk branch to prevent autoupgrade
- Ante test for Ante decorators
- Add version info on CLI
- Modify CLI command to be more organized and descriptions to be well-formed along with scripts examples"
- Combine SekaiApp with simap for easier maintenance
- Remove unused codes

## [v0.1.21.8] - 28.07.2021

- Remove resolved TODO in readme for upgrade module
- Add token rates keeper functions test
- Add token alias keeper functions test
- Add test for tokens grpc query functions
- Add test for ProposalTokensWhiteBlackChange
- Add missing CLI test for tokens module
- Remove 'all' flag for staking.ValidatorsRequest and slashing.SigningInfosRequest to prevent performance issue
- Refactor proto-gen script
- Remove duplications in protobuf definition
- Utilize existing sdk protobuf
- Resolve codebase to work with latest protobuf
- Set IdentityRecord to be only one by address

## [v0.1.21.7] - 22.07.2021

- update governance apis
- add interx list query api
- update deposit/withdraw transactions query
- update interx functions query

## [v0.1.21.6] - 19.07.2021

- Update env to be latest
- Remove unused message type and commented codebase
- Check changes to make for codec registration
- Modification for codec meta registration
- Remove unused Proposer field from MsgApproveIdentityRecords
- Add title into the proposal
- For identity record creation/edit, update date to use blockTime
- Shell script modifications

## [v0.1.21.5] - 18.07.2021

- Zero gas implementation

## [v0.1.21.4] - 16.07.2021

- Refactor governance module for easier maintenance

## [v0.1.21.3] - 15.07.2021

## Added

- public p2p node list API
- private p2p node list API
- p2p node_id verification

## Updated

- validators query to use identity registrar
- interx configurations/cli configurations

## [v0.1.21.2] - 14.07.2021

- Remove identity part from claim validator

## [v0.1.21.1] - 13.07.2021

- Add instate upgrade feature

## [v0.1.21] - 08.07.2021

- Implement identity registrar
- Code cleanup for package names

## [v0.1.20] - 10.06.2021

- Add own upgrade module

## [v0.1.19.6] - 12.05.2021

- Add reverse order querying for proposals querying
- Add pagination limit of 512 in grpc execution level

## [v0.1.19.5] - 07.05.2021

- update mischanceconfidence counter logic.
- update protos to latest release.
- add icon to tokens aliases query.

## [v0.1.19.4] - 06.05.2021

### Modified

- Refactor jail / unjail logic
- Modify mischance, MischanceConfidence counter logic
- Implement additional logic for mischance_confidence to count only when active
- Remove Tombstoned status on slashing module
- Add more status checks for validator transition
- Accurate error logging
- Fix tests for uptime counter and validator status transition modification
- Add more tests to handle modified logic

## [v0.1.19.3] - 03.05.2021

### Added

- [Interx] Network Properties query endpoint
- [Interx] Pagination in Proposal query
- [Interx] Network Properties proto description

### Changed

- [Interx] Refactor validator query to use less call for the performance

### Fixed

- [Interx] Multiple lint issues

## [v0.1.19.3] - 03.05.2021

### Added

- [Interx] Network Properties query endpoint
- [Interx] Pagination in Proposal query
- [Interx] Network Properties proto description

### Changed

- [Interx] Refactor validator query to use less call for the performance

### Fixed

- [Interx] Multiple lint issues

## [v0.1.19.2] - 28.04.2021

### Added

- [Interx] MischanceConfidence in ValidatorQuery

### Fixed

- [Interx] MissedBlocksCounter in ValidatorQuery

## [v0.1.19.1] - 26.04.2021

### Added

- Update validator signing info query based on the latest release
- Add configurations for node discovery

## [v0.1.19] - 21.04.2021

### Added

- New uptime counter
- Move properties from slashing to gov for uptime properties
- Proposal to reset whole validators rank

## [v0.1.18.19] - 22.04.2021

### Changed

- Upgrade SDK to version 0.42.4
- Validators cannot use same moniker and is space trimmed.

## [v0.1.18.17] - 17.04.2021

### Added

- Implement infinite gas meter decorator

## [v0.1.18.16] - 09.04.2021

### Fixed

- Add minimum blocks for voting and enactment time on-chain param and implement logic
- Modify error messages for slashing module

## [v0.1.18.14] - 01.04.2021

### Fixed

- Fixed problem with network stopping on pause.
- Validator cannot be claimed if already did.

## [v0.1.18.13] - 02.04.2021

## Added

- Add event manager to all msg handlers
- Fix few permission issues
- Fix tests

## [v0.1.18.11] - 01.04.2021

## Added

- Add missing permissions to sudo role
- Fix vote result unknown state for set poor network messages proposal
- Split huge README into several shell scripts

## [v0.1.18.7] - 26.03.2021

### Added

- Shell script to setup environment variables especially for permissions
- README for upsert token alias, upsert token rates by governance

## [v0.1.18.5] - 25.03.2021

### Fixed

- Some proposals were created in minutes instead of seconds pattern.

## [v0.1.18.4] - 25.03.2021

### Added

- Add a field in all proposals to be able to set some description.

### Fixed

- The actor when it receives a permission becomes active.
- Fixed problem when voting unjail validator proposal.

## [v0.1.18.2] - 25.03.2021

### Changed

- Permission numbers to an organized way
- Cleanup gov codebase function names and vars
- Fix SetPoorNetworkMessagesProposal codec registration
- Add logic for Mischance and ProducedBlocksCounter, MissedBlocksCounter

## [v0.1.18] - 19.03.2021

### Added

- Ante handler to check frozen tokens movement
- Add network properties for ENABLE_TOKEN_WHITELIST / ENABLE_TOKEN_BLACKLIST
- Add permissions for creation and vote on blacklist and whitelist change of tokens
- Added CLI command to submit proposal to change blacklist and whitelist
- Added CLI command to query current blacklist and whitelist
- Network Properties management code modification for boolean properties

## [v0.1.17.6] - 18.03.2021

### Added

- Now when the proposal passes it enter ins status Enactment.
- Add proposal to create a Role.
- Fix GetTxProposalUpsertDataRegistry and make it appear on client.

### Fixed

- Fix and clean some CLI commands (proposal upsert token rates, proposal upsert token alias, proposal upsert data registry).

## [v0.1.17.3] - 03.08.2021

### Fixed

- Validators query to include mischance.

### Added

- Tokens alias/rate query.
- Voters/votes query.

## [v0.1.17.2] - 03.02.2021

### Fixed

- Mischance querying CLI command

### Added

- genutil module to handle validator status
- CLI utility command to get valcons address from account address
- ValidatorJoined hook that's derivated from Cosmos SDK's `ValidatorBonded` hook

## [v0.1.17.1] - 02.25.2021

### Added

- GRPC query for proposals
- GRPC query for validators + validator_signing_infos

## [v0.1.17] - 02.16.2021

### Fixed

- Problem with ClaimValidator and PubKey encoding due to protocol buff Any type.
- Fix bug that made that you can vote when the proposal ended.
- When a proposal does not reach quorum it ends being Quorum not reached.
- Proposal voting time and enactment time now are defined in seconds.
- It shows the votes that a proposal has on client query.

## [v0.1.16.2a] - 02.08.2021

### Added

- Custom evidence module to jail a double signed validator
- CLI command for writing proposal to unjail a validator
- CLI command for setting max jail time network property proposal

## [v0.1.16.1] - 02.04.2021

### Added

- CLI command to set poor network messages
- CLI command to query poor network messages
- Add POOR_NETWORK_MAX_BANK_TX_SEND feature for poor network for restriction (only bond denom is allowed)
- Reject feature for not allowed messages on poor network

## [v0.1.15.2] - 01.21.2021

### Added

- CLI command GetTxProposalUpsertTokenAliasCmd and GetTxProposalUpsertTokenRatesCmd are now exposed.

## [v0.1.15.1] - 01.21.2021

### Added

- CLI command to get ValAddress from AccAddress

## [0.1.15] - 01.15.2021

### Added

- Added custom slashing module for validator's block signing info management, inactivate, activate, pause, unpause
- Added validator performance calculator using `rank` and `streak`
- Upgraded Cosmos SDK to v0.40.0 (stargate)

### Removed

- Old staking, slashing, evidence, distribution module

## [0.1.14.3] - 12.30.2020

### Added

- Added GRPC query for Data Reference Registry.
- Update response caching for data references. (KIP_47.1)
- Added file hosting feature. (KIP_47.1)

## [0.1.14.2] - 11.30.2020

### Added

- Update Cosmos SDK to v0.40.0-rc4.

## [0.1.14.1] - 11.30.2020

### Added

- Proposal to upsert the Token Rates. (CLI too)

## [0.1.14] - 11.26.2020

### Added

- Added a wrapper to register messages with function metadata.
- Added function_id for message types.
- Registered function meta for existing messages.
- Added INTERX api for Kira functions list.
- Added INTERX api for INTERX functions list.

## [0.1.13] - 11.20.2020

### Added

- Proposal to upsert the Token Aliases KIP_24. (CLI too)

## [0.1.12.1] - 11.15.2020

### Added

- Proposal to upsert the Data Registry. (CLI too)
- Proposal to change Network Properties. (CLI too)

### Changed

- Now it is more generic to be able to add new proposals in the complete flow.

## [0.1.12] - 11.13.2020

### Added

KIP_8

- Added grpc gateway
- Added status, balances, transaction hash queries
- Added transaction encode/broadcast
- Added response format

KIP_9

- Added endpoint whitelist

KIP_48

- Added INTERX faucet

KIP_47

- Added response caching

KIP_32

- Added withdraws, deposits

## [0.1.7.3] - 11.12.2020

### Added

- Added CLI for querying proposals / individual proposal
- Added CLI for querying votes / individual vote
- Added CLI for querying whitelisted proposal voters

### Changed

- Updated genesis actor initialization process
- Updated proposal end time and enactment time
- Fixed end blocker concert not registered issue for MsgClaimValidator

## [0.1.7.2] - 11.11.2020

### Changed

- There is a new permission for all role related changes. PermUpsertRole.

## [0.1.7.1] - 11.09.2020

### Changed

- Proposal is now a generic type, the Content part is what changes between different proposal types.

## [0.1.7] - 11.09.2020

- Added CLI command to upsert token rates per denom
- Added CLI commands to query token rates
- Implemented feeprocessing module for new fee processing logic
- Implemented foreign currency fee payment

## [0.1.6.3] - 11.07.2020

### Added

- We can propose to SetPermissions to an actor.
- We can vote a proposal to SetPermissions to an actor.
- Added proposal endtime and proposal enactment time into
  the network properties.

## [0.1.6.2] - 10.23.2020

### Added

- The keeper has method to get all Actors by witelisted permission.
- The keeper has method to get All actors that have specific role.
- The keeper has method to get all roles that have a whitelist permission.

### Changed

- Big refactor on the way Role and Permissions are stored.
- In keeper we don't expose SetPermissionsForRole anymore.

## [0.1.6.1] - 10.19.2020

### Added

- Added CLI command to send a SetPermission proposal.
- Added CLI command to vote a SetPermission proposal.

### Changed

- Now Role and Permissions are persisted differently in order to be able to get
  actors by permission and actors by role.

- Now the commands for all Governance module is simplified in a better hierarchical style.

```
Available Commands:
  councilor   Councilor subcommands
  permission  Permission subcommands
  proposal    Proposal subcommands
  role        Role subcommands
```

## [0.1.6] - 10.16.2020

### Added

- Added CLI command to upsert token alias per symbol
- Added CLI commands to query token aliases per symbol and denom
- Added CLI command to query all token aliases

### Modified

- Modified execution fee to use transaction type as identifier
- Modified min/max fee range to [100 - 1'000'000] in ukex

## [0.1.4] - 10.05.2020

### Added

- Added CLI command to change execution fee per message type
- Added CLI command to change transaction fee range
- Added CLI command to query execution fee
- Added CLI command to query transaction fee range

## [0.1.2.4] - 09.24.2020

### Added

- Added CLI command to remove blacklist permissions into a specific role.
- Added CLI command to create new role.
- Added CLI command to assign new role.
- Added CLI command to remove assignation for a role.

## [0.1.2.3] - 09.17.2020

### Changed

- Updated cosmos SDK to last version of 17th september .

## [0.1.2.2] - 09.14.2020

### Added

- Added CLI command to claim governance seat.
- Added CLI command to set whitelist permissions into a specific role.
- Added CLI command to set blacklist permissions into a specific role.
- Added CLI command to remove whitelist permissions into a specific role.

## [0.1.2.1] - 09.06.2020

### Added

- Added CLI command to Set Blacklist Permissions too.
- Module CustomGov defines in genesis by default Permissions by roles Validator (0x2) and Sudo (0x1).

### Changed

- Now the roles are validated when taking some action. It checks if the user has permissions either in the role or individually.

## [0.1.2] - 09.01.2020

### Added

- Add command to add whitelist permissions to an address, that address is included
  in theNetworkActor registry with the specified permission added.
- Now the user that generates the network has AddPermissions by default, so he is the only one
  that can add permissions into the registry.

### Changed

- Now the ClaimValidator message takes care that the user has ClaimValidator permissions,
  if he the user does not have, it fails.
