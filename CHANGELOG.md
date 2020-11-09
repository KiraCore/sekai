# Changelog

## [0.1.7] - 11.09.2020
- Added CLI command to upsert token rates per denom
- Added CLI commands to query token rates
- Implemented feeprocessing module for new fee processing logic
- Implemented foreign currency fee payment

## [0.1.6.4] - 11.09.2020

### Changed
- Proposal is now a generic type, the Content part is what changes between different proposal types.
  
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
