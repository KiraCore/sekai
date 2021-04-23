# Changelog

- merge token aliases query and token supply.
- update response type for uncomfirmed_txs
- update postman API collection

## [v0.1.18.15] - 04.06.2021

- update consensus stopped after 6 * average block time.
- add unconfirmed_txs query (wip).

## [v0.1.18.12] - 04.02.2021

- update interx status api
- add nodes cli configurations
- fix consensus_stopped logic

## [v0.1.18.10] - 03.31.2021

- add `addrbook` cli option
- fix faucet amount handling (e.g. 18 decimal)

## [v0.1.18.9] - 03.29.2021

- add `faucet_amounts`, `faucet_minimum_amounts`, `fee_amounts` for interxd init cli configuration.
- fix validators query (calculate consensus_stopped after sync finished).
- update consensus query responses.
- fix configuration init for mnemonics.
- fix privkey generation.

## [v0.1.18.6] - 03.25.2021

- fix proposals query to have description field.

## [v0.1.18.3] - 03.25.2021

- add consensus api
- fix proposal query
- fix votes/voters query
- update protos to latest release

## [v0.1.x.x] - xx.xx.xxxx

### Updated

- updated postman api collection to latest
- updated consensus stop logic with block times.

### Added

- add rosetta - account balance api

## [v0.1.17.5] - 03.17.2021

### Updated

- update cors allow headers

### Added

- add rosetta - network status api

## [v0.1.17.4] - 03.16.2021

### Updated

- updated validators query to have `top` field - sort validators by `top` field
- update interx metadata - sync to latest

### Added

- add rosetta - network list api
- add rosetta - network options api

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
