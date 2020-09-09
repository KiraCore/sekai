# sekai
Kira Hub

## Set network properties
```sh
# command
sekaid tx customgov set-network-properties --from validator --min_tx_fee="2" --max_tx_fee="2000" --keyring-backend=test --chain-id=testing
{"body":{"messages":[{"@type":"/kira.gov.MsgSetNetworkProperties","network_properties":{"min_tx_fee":"2","max_tx_fee":"2000"},"proposer":"kira1jnpaqdqkjvpvkgxpjz57vzjv4r5zsuudmjltkh"}],"memo":"","timeout_height":"0","extension_options":[],"non_critical_extension_options":[]},"auth_info":{"signer_infos":[],"fee":{"amount":[],"gas_limit":"200000"}},"signatures":[]}

confirm transaction before signing and broadcasting [y/N]: y

# response
{"height":"0","txhash":"DDB61A46D581D9813CF993B3037B069EF9820FFF9C0BBAF162FBA7B31920B238","codespace":"sdk","code":2,"data":"","raw_log":"no concrete type registered for type URL /kira.gov.MsgSetNetworkProperties against interface *types.Msg: tx parse error","logs":[],"info":"","gas_wanted":"0","gas_used":"0","tx":null,"timestamp":""}
```
## Query network properties
```sh
# command
sekaid query customgov network-properties

# response
properties:
  max_tx_fee: "10000"
  min_tx_fee: "1"
```
---
`dev` branch
