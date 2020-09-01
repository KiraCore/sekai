# Changelog


## [Unreleased]

- Module CustomGov defines in genesis by default Permissions by roles Validator (0x2) and Sudo (0x1).

## [0.1.2] - 09.01.2020
### Added

- Add command to add whitelist permissions to an address, that address is included
in theNetworkActor registry with the specified permission added.
- Now the user that generates the network has AddPermissions by default, so he is the only one
that can add permissions into the registry.

### Changed

- Now the ClaimValidator message takes care that the user has ClaimValidator permissions,
if he the user does not have, it fails.