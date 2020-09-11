# sekai
Kira Hub

## Set ChangeTxFee permission
```sh
# command to set changeTxFee permission
sekaid tx customgov set-whitelist-permissions --from validator --keyring-backend=test --permission=3 --addr=$(sekaid keys show -a validator --keyring-backend=test) --chain-id=testing --fees=100stake

# response
{"height":"101","txhash":"D584594958BE83482C19B687A17C4A00591C0128308D32D64A211860F6826611","codespace":"","code":0,"data":"0A170A1577686974656C6973742D7065726D697373696F6E73","raw_log":"[{\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"whitelist-permissions\"}]}]}]","logs":[{"msg_index":0,"log":"","events":[{"type":"message","attributes":[{"key":"action","value":"whitelist-permissions"}]}]}],"info":"","gas_wanted":"200000","gas_used":"51921","tx":null,"timestamp":""}
```
## Set network properties
```sh

# command with fee set
sekaid tx customgov set-network-properties --from validator --min_tx_fee="2" --max_tx_fee="2000" --keyring-backend=test --chain-id=testing --fees=100stake

confirm transaction before signing and broadcasting [y/N]: y

# response when all are ok
{"height":"10","txhash":"838448F164CF1C94577B6B4B3810C537F8563AC907B5DDA15E4BF087A12B02AA","codespace":"","code":0,"data":"0A180A167365742D6E6574776F726B2D70726F70657274696573","raw_log":"[{\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"set-network-properties\"}]}]}]","logs":[{"msg_index":0,"log":"","events":[{"type":"message","attributes":[{"key":"action","value":"set-network-properties"}]}]}],"info":"","gas_wanted":"200000","gas_used":"49851","tx":null,"timestamp":""}

# response when not enough permissions to change tx fee
{"height":"3","txhash":"032EF37E996A9D9060A70F74F2C78956FA95F39EDE6A91E1C8BC27EE75C62826","codespace":"customgov","code":5,"data":"","raw_log":"failed to execute message; message index: 0: PermChangeTxFee: not enough permissions","logs":[],"info":"","gas_wanted":"200000","gas_used":"52429","tx":null,"timestamp":""}

# command without fee set
sekaid tx customgov set-network-properties --from validator --min_tx_fee="2" --max_tx_fee="2000" --keyring-backend=test --chain-id=testing

# response
confirm transaction before signing and broadcasting [y/N]: y
{"height":"0","txhash":"9003558A51D7067085FF6F42C28CEB974ACEC357845174FAE3ECA75E0306BED3","codespace":"sdk","code":18,"data":"","raw_log":"fee out of range [1, 10000]: invalid request","logs":[],"info":"","gas_wanted":"200000","gas_used":"13063","tx":null,"timestamp":""}

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
