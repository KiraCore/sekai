# Changelog

## [Unreleased]
### Added

- Added CLI command to remove blacklist permissions into a specific role.

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