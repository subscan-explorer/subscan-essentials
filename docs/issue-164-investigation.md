# Issue #164: Revive mapped-account balance lookup

Public report: https://github.com/subscan-explorer/subscan-issue-tracker/issues/164

## Reproduction (2026-10-08)

Only the report's third-party public sample was used:

- H160: `0x2ee6ace1cdc464a61d3f3910a6853ab5d093a107`
- SS58: `12u9Ha4PxyyQPvJgq3BghnqNXDwLqTnnJFuXV7aZQoiregT2`
- AccountId32: `5400e2f7f5669b26998d8e4d3c1a2c8a2d0a9af827ca54a1cc3509105035c32e`

The loaded public Subscan pages showed **5.37 DOT** for SS58 and **0 DOT** for H160, with reciprocal account links. Repeating the issue's read-only `eth_getBalance` request to `https://eth-rpc.polkadot.io/` at `0x1492abd` returned `0x4a63e5e33a54e000` (5.3603807456 DOT in the RPC's 18-decimal representation).

## Source trace and fix

`plugins/evm/http/etherscan.go` handles `account-balance` and `account-balancemulti` through `API_GetAccounts` in `plugins/evm/dao/api.go`. That method converts H160 addresses with `h160ToAccountIdByNetwork`, then queries indexed balance records by `address`.

Previously the Revive branch always selected `H160 || 0xEE * 12`. For the public sample that is `2ee6ace1cdc464a61d3f3910a6853ab5d093a107eeeeeeeeeeeeeeeeeeeeeeee`, which differs from its native AccountId. A native balance indexed under the original account cannot match that query.

The shared resolver now reads `Revive.OriginalAccount` using the active runtime metadata and existing WebSocket RPC pool. A registered AccountId is normalized to the database representation. Only an absent storage value selects the existing padding fallback. Transport/RPC errors, unavailable metadata storage, and malformed values propagate to callers. Both the balance API and `TouchAccount` use this resolver. The padding helper is unchanged.

This follows the upstream `AccountId32Mapper::to_account_id` behavior in `substrate/frame/revive/src/address.rs`: query `OriginalAccount`, otherwise use `to_fallback_account_id`. The storage declaration in `revive/src/lib.rs` is `StorageMap<_, Identity, H160, AccountId32>`; raw values are 32 bytes, without an Option discriminator.

## Regression evidence

`address_mapping_test.go` uses the public address pair with a stubbed mapping lookup. It covers registered mappings, absent mappings, lookup errors, malformed mappings, H160 normalization, canceled context, and non-Revive/invalid-address branches. A deliberately constructed indexed-balance fixture demonstrates the old key miss and corrected key match; it is not a production database snapshot.

A temporary copy of the resolver with the old padding-only branch fails the public-sample regression. The corrected resolver passes it.

Run the isolated regression suite from the repository root:

```sh
GOPATH=/private/tmp/subscan-go GOCACHE=/private/tmp/subscan-go-cache \
  go test $(rg --files plugins/evm/dao -g '*.go' | rg -v '_test.go$') \
  plugins/evm/dao/address_test.go plugins/evm/dao/address_mapping_test.go
```

Add `-race` to run with the race detector. The explicit file set excludes existing integration-test initializers that connect to configured MySQL, drop tables, and run Redis `FLUSHALL`.

## Limits

This confirms a local backend defect consistent with the observed production discrepancy. Production code/version, account database, caches, and the account-page backend route were not available, so this is not proof that the production page uses this exact resolver. The Ethereum endpoint rejects `state_getStorage`; follow-up verification used the official Substrate endpoint, as described below. Unit tests inject the public mapping; the separate live test independently reads it.

Resolution uses current chain mapping while balance records reflect indexed state; mapping changes can temporarily outpace indexing. Existing `evm_accounts` rows are not backfilled. The websocket client does not accept a request context; cancellation is checked before dispatch, with in-flight timeout behavior supplied by the existing client. The smaller nonzero balance difference remains uninvestigated.

No personal account, chain write, database mutation, commit, push, issue comment, or PR publication was performed. Foundry specialist delegation does not apply to this Go indexer change.

Verification completed: isolated DAO regressions PASS; race-enabled isolated regressions PASS (macOS linker emitted an `LC_DYSYMTAB` warning, but the test executable passed); `go build ./plugins/evm/...` PASS; `go build ./...` PASS; `git diff --check` PASS. The database/Redis integration suite was not run.


## Follow-up verification

The official endpoint `https://polkadot-asset-hub-rpc.polkadot.io` identified the chain as `Polkadot Asset Hub`, runtime `statemint`, spec version `2005000`. A read-only `state_getStorage` request for the sample's `Revive.OriginalAccount` key returned `0x5400e2f7f5669b26998d8e4d3c1a2c8a2d0a9af827ca54a1cc3509105035c32e`.

`TestLiveReviveMapping` then passed through the actual `readReviveOriginalAccount` implementation using the endpoint's fetched V14 metadata and its WebSocket RPC pool. This verifies production lookup code, key encoding, and raw AccountId decoding against live state.

`TestBalanceAPIWithIsolatedDatabase` passed with the race detector. It calls the actual `API_GetAccounts`, uses a SQLite in-memory database through the existing GORM MySQL dialect for the SELECT, and serves mapping responses from a localhost WebSocket server. It verifies mapped and unmapped records in one request and RPC-error propagation. The balance amounts are explicit test fixtures, not live balances. No configured database or Redis service is used. This does not certify MySQL schema or server compatibility.

Run all isolated tests (requires permission to bind localhost in sandboxed environments):

```sh
ETH_RPC=http://127.0.0.1:1 GOPATH=/private/tmp/subscan-go GOCACHE=/private/tmp/subscan-go-cache \
  go test -race $(rg --files plugins/evm/dao -g '*.go' | rg -v '_test.go$') \
  plugins/evm/dao/address_test.go plugins/evm/dao/address_mapping_test.go \
  plugins/evm/dao/address_api_test.go -timeout 60s
```

`ETH_RPC` prevents the imported web3 package's initialization from contacting its default public endpoint during isolated tests. Its expected connection-refused diagnostic does not affect these tests.

For the opt-in live test, save the official endpoint's batch response for `system_chain` (id 1) and `state_getMetadata` (id 2) in a JSON file, then run the same explicit production file set with `address_api_test.go`, `-run '^TestLiveReviveMapping$'`, and `SUBSCAN_MAPPING_METADATA` pointing to that file. The test's live endpoint is fixed to the official Asset Hub endpoint. The default suite skips this network test.

No PR, commit, push, or external publication was performed during follow-up verification. Production backend/version and cache/indexing behavior remain unverified.
