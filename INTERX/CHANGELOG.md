# Changelog

### Updated

- updated validators query to have `top` field - sort validators by `top` field
- update interx metadata - sync to latest

### Added

- add rosetta - network list api
- add rosetta - network options api

## [v0.1.x.x] - xx.xx.xxxx

## [v0.1.17.3] - 03.08.2021
### Fixed
- Validators query to include mischance.
- Updated validators statistics
- Updated validators query to include zero values.

### Added
- Tokens alias/rate query.
- Voters/votes query.
- Added configuration for https/http

## [v0.1.17.1] - 02.25.2021
### Added
- Add proposals query
- Add proposal query by a given proposal_id
- Add query for validators + validator_signing_infos

## [v0.1.16.3]] - 02.09.2021

### Added

- Add query blocks.
- Add query block by height or hash.
- Add query block tractions.
- Add query transaction by hash.
- Add cli configurations: `max_cache_size`, `max_download_size`, `cache_dir`, `caching_duration`, `faucet_time_limit`

### Fix

- Fix caching issues
- Fix mutex issues

## [0.1.16] - 01.29.2021

### Changed

- Update configuration to have mnemonic filename. And read mnemonic from the file.
- Update configuration to group cache related configurations.
- Update Readme.md
- Update Interx Status API
- Changed interx to interxd (cli command are available: `interxd init`, `interxd start`)

### Added

- Add mnemonic validation step.
- Add query validators.
- Add query validator by address, valley, pubkey, moniker and status.
- Add query genesis, genesis-checksum.
