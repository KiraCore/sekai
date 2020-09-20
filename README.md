# sekai
Kira Hub

## Create order book
```sh
# command
sekaid tx ixp createOrderBook base quote mnemonic --from validator --keyring-backend=test --chain-id testing --home=$HOME/.sekaid
{"body":{"messages":[{"@type":"/kira.ixp.MsgCreateOrderBook","Base":"base","Quote":"quote","Mnemonic":"mnemonic","Curator":"kira17pzan0q4d5acyykygwqass8z0crjvflhjq3qvm"}],"memo":"","timeout_height":"0","extension_options":[],"non_critical_extension_options":[]},"auth_info":{"signer_infos":[],"fee":{"amount":[],"gas_limit":"200000"}},"signatures":[]}

confirm transaction before signing and broadcasting [y/N]: y

# response
{"height":"37","txhash":"86F154FAD24E330906DC8E983324FDC962B5DBEDC86CECF7A7C93CA9ED474D00","codespace":"","code":0,"data":"0A3C0A0F6372656174656F72646572626F6F6B12297B224944223A226635323533383535663932653135376639663033353830323931623665356462227D","raw_log":"[{\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"createorderbook\"}]}]}]","logs":[{"msg_index":0,"log":"","events":[{"type":"message","attributes":[{"key":"action","value":"createorderbook"}]}]}],"info":"","gas_wanted":"200000","gas_used":"47467","tx":null,"timestamp":""}
```

If you parse returned data field, you can get orderbook ID.
```
0A3C0A0F6372656174656F72646572626F6F6B12297B224944223A226635323533383535663932653135376639663033353830323931623665356462227D
```
is equivalent to
```
createorderbook {"ID":"f5253855f92e157f9f03580291b6e5db"}
```
Here `"f5253855f92e157f9f03580291b6e5db"` is orderbook ID.

## Query order book
- By ID
Ex1.
```sh
# command
sekaid query ixp listorderbooks ID e6a8fc6cf92e157f9f03580291b6e5db --chain-id testing

# response
{"orderbooks":[{"ID":"e6a8fc6cf92e157f9f03580291b6e5db","Index":0,"Base":"base","Quote":"quote","Mnemonic":"mnemonic","Curator":"kira17pzan0q4d5acyykygwqass8z0crjvflhjq3qvm"}]}
```
Ex2.
```sh
# command
sekaid query ixp listorderbooks ID e6a8 --chain-id testing

# response
{"orderbooks":[{"ID":"","Index":0,"Base":"","Quote":"","Mnemonic":"","Curator":""}]}
```
- By curator
```sh
# command
sekaid query ixp listorderbooks Curator $(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) --chain-id testing

# response
{"orderbooks":[{"ID":"e6a8fc6cf92e157f9f03580291b6e5db","Index":0,"Base":"base","Quote":"quote","Mnemonic":"mnemonic","Curator":"kira17pzan0q4d5acyykygwqass8z0crjvflhjq3qvm"}]}
```
- By Index

Ex1.
```sh
# command
sekaid query ixp listorderbooks Index 0 --chain-id testing

# response
{"orderbooks":[{"ID":"e6a8fc6cf92e157f9f03580291b6e5db","Index":0,"Base":"base","Quote":"quote","Mnemonic":"mnemonic","Curator":"kira17pzan0q4d5acyykygwqass8z0crjvflhjq3qvm"}]}
```

Ex2.
```sh
# command
sekaid query ixp listorderbooks Index 10 --chain-id testing

# response
{"orderbooks":[]}
```

- By Quote
Ex1.
```sh
# command
sekaid query ixp listorderbooks Quote quote --chain-id testing

# response
{"orderbooks":[{"ID":"e6a8fc6cf92e157f9f03580291b6e5db","Index":0,"Base":"base","Quote":"quote","Mnemonic":"mnemonic","Curator":"kira17pzan0q4d5acyykygwqass8z0crjvflhjq3qvm"}]}
```
Ex2.
```sh
# command
sekaid query ixp listorderbooks Quote q --chain-id testing

# response
{"orderbooks":[]}
```
- By Base
Ex1.
```sh
# command
sekaid query ixp listorderbooks Base base --chain-id testing

# response
{"orderbooks":[{"ID":"e6a8fc6cf92e157f9f03580291b6e5db","Index":0,"Base":"base","Quote":"quote","Mnemonic":"mnemonic","Curator":"kira17pzan0q4d5acyykygwqass8z0crjvflhjq3qvm"}]}
```
Ex2.
```sh
# command
sekaid query ixp listorderbooks Base b --chain-id testing

# response
{"orderbooks":[]}
```

- By trading pair
```sh
# command
sekaid query ixp listorderbooks_tradingpair base quote --chain-id testing

# response
{"orderbooks":[{"ID":"36e0cd99f92e157f9f03580291b6e5db","Index":0,"Base":"base","Quote":"quote","Mnemonic":"mnemonic","Curator":"kira1hqvm6nup0kkntgfq9xe0hk7gxnm4wjakmmpae0"}]}
```
```sh
# command
sekaid query ixp listorderbooks_tradingpair base quote1 --chain-id testing

# response
{"orderbooks":[]}
```
## Create order
```sh
# command
sekaid tx ixp createOrder 36e0cd99f92e157f9f03580291b6e5db 0 0 0 --from validator --keyring-backend=test --chain-id=testing --home=$HOME/.sekaid
{"body":{"messages":[{"@type":"/kira.ixp.MsgCreateOrder","OrderBookID":"36e0cd99f92e157f9f03580291b6e5db","OrderType":"limitBuy","Amount":"0","LimitPrice":"0","ExpiryTime":"0","Curator":"kira13n3el7pzynrd3d4mn92k49djktf75kngl2xlc7"}],"memo":"","timeout_height":"0","extension_options":[],"non_critical_extension_options":[]},"auth_info":{"signer_infos":[],"fee":{"amount":[],"gas_limit":"200000"}},"signatures":[]}

confirm transaction before signing and broadcasting [y/N]: y

# response
{"height":"8","txhash":"97F57E6C73053F3FF0F6838D556F91DCCC133110F30D4B8A72AB87A96B5FFAF0","codespace":"","code":0,"data":"0A0D0A0B6372656174656F72646572","raw_log":"[{\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"createorder\"}]}]}]","logs":[{"msg_index":0,"log":"","events":[{"type":"message","attributes":[{"key":"action","value":"createorder"}]}]}],"info":"","gas_wanted":"200000","gas_used":"45716","tx":null,"timestamp":""}
```
## Cancel order
```sh
# command
sekaid tx ixp cancelOrder a991645f855efbb8855efbb891b6e5db --from validator --keyring-backend=test --chain-id=testing --home=$HOME/.sekaid
{"body":{"messages":[{"@type":"/kira.ixp.MsgCancelOrder","OrderID":"a991645f855efbb8855efbb891b6e5db","Curator":""}],"memo":"","timeout_height":"0","extension_options":[],"non_critical_extension_options":[]},"auth_info":{"signer_infos":[],"fee":{"amount":[],"gas_limit":"200000"}},"signatures":[]}

confirm transaction before signing and broadcasting [y/N]: y

# response
{"height":"63","txhash":"C8B36D6B608C92A8F0BBDEB5280580DD8D4339B152039C0A8D0B608A48FF8B03","codespace":"","code":0,"data":"0A380A0B6372656174656F7264657212297B224944223A226139393136343566383535656662623838353565666262383931623665356462227D","raw_log":"[{\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"createorder\"}]}]}]","logs":[{"msg_index":0,"log":"","events":[{"type":"message","attributes":[{"key":"action","value":"createorder"}]}]}],"info":"","gas_wanted":"200000","gas_used":"46696","tx":null,"timestamp":""}
```

If you parse returned data field, you can get order ID.
```
0A3C0A0F6372656174656F72646572626F6F6B12297B224944223A226635323533383535663932653135376639663033353830323931623665356462227D
```
is equivalent to
```
createorder {"ID":"a991645f855efbb8855efbb891b6e5db"}
```
Here `"a991645f855efbb8855efbb891b6e5db"` is order ID.

## Query order
```sh
# command
sekaid query ixp listorders 36e0cd99f92e157f9f03580291b6e5db 0 0

# response
{"orders":[{"ID":"","Index":0,"OrderBookID":"36e0cd99f92e157f9f03580291b6e5db","OrderType":"limitBuy","Amount":"0","LimitPrice":"0","ExpiryTime":"0","IsCancelled":false,"Curator":"kira13n3el7pzynrd3d4mn92k49djktf75kngl2xlc7"}]}
```
```sh
# command
sekaid query ixp listorders 36e0cd99f92e157f9f03580291b6e5b 0 0

# response
{"orders":[]}
```
## upsertSignerKey

- CLI

```sh
# Secp256k1 type
sekaid tx ixp upsertSignerKey AnzIM9IcLb07Cvwq3hdMJuuRofAgxfDekkD3nJUPPw0w
sekaid tx ixp upsertSignerKey AnzIM9IcLb07Cvwq3hdMJuuRofAgxfDekkD3nJUPPw0w --key-type=Secp256k1

# Ed25519 type
sekaid tx ixp upsertSignerKey TXgDkmTYpPRwU/PvDbfbhbwiYA7jXMwQgNffHVey1dC644OBBI4OQdf4Tro6hzimT1dHYzPiGZB0aYWJBC2keQ== --key-type=Ed25519

# enabled false
sekaid tx ixp upsertSignerKey AnzIM9IcLb07Cvwq3hdMJuuRofAgxfDekkD3nJUPPw0w --enabled=false

# permissions
sekaid tx ixp upsertSignerKey AnzIM9IcLb07Cvwq3hdMJuuRofAgxfDekkD3nJUPPw0w --permissions=1,2

# expiry-time (set when this key expire if does not set it's automatically set to 10 days after current timestamp)
sekaid tx ixp upsertSignerKey AnzIM9IcLb07Cvwq3hdMJuuRofAgxfDekkD3nJUPPw0w --expiry-time=1598247750
# data
sekaid tx ixp upsertSignerKey AnzIM9IcLb07Cvwq3hdMJuuRofAgxfDekkD3nJUPPw0w --data="192.168.1.1" --keyring-backend=test --chain-id=testing --home=$HOME/.sekaid
```
Ex1.
```sh
# command
sekaid tx ixp upsertSignerKey AnzIM9IcLb07Cvwq3hdMJuuRofAgxfDekkD3nJUPPw0w --from validator --keyring-backend=test --chain-id=testing --home=$HOME/.sekaid
{"body":{"messages":[{"@type":"/kira.ixp.MsgUpsertSignerKey","PubKey":"AnzIM9IcLb07Cvwq3hdMJuuRofAgxfDekkD3nJUPPw0w","KeyType":"Secp256k1","ExpiryTime":"1599186595","Enabled":true,"Permissions":[],"Curator":"kira13n3el7pzynrd3d4mn92k49djktf75kngl2xlc7"}],"memo":"","timeout_height":"0","extension_options":[],"non_critical_extension_options":[]},"auth_info":{"signer_infos":[],"fee":{"amount":[],"gas_limit":"200000"}},"signatures":[]}

confirm transaction before signing and broadcasting [y/N]: y

# response
{"height":"11543","txhash":"34B0D49958A881D280392644D3CBC358100DCF8353F1DCC26899D3B039CAAF02","codespace":"","code":0,"data":"0A110A0F7570736572747369676E65726B6579","raw_log":"[{\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"upsertsignerkey\"}]}]}]","logs":[{"msg_index":0,"log":"","events":[{"type":"message","attributes":[{"key":"action","value":"upsertsignerkey"}]}]}],"info":"","gas_wanted":"200000","gas_used":"44086","tx":null,"timestamp":""}
```

## Query signer keys

```sh
# command
sekaid query ixp getsignerkeys --curator $(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid)

# response
{"signerkeys":[{"PubKey":"AnzIM9IcLb07Cvwq3hdMJuuRofAgxfDekkD3nJUPPw0w","KeyType":"Secp256k1","ExpiryTime":"1599187090","Enabled":true,"Permissions":[],"Curator":"kira1xsq0wapm5t975k3hn2rj4y2zhnm5up9d59uhpy"}]}
```
- Rest

```sh
# TODO: should add readme for sample upsertSignerKey with same examples of CLI
```
---
`dev` branch
