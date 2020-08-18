# sekai
Kira Hub

## Create order book
```sh
sekaid tx kiraHub createOrderBook base quote mnemonic --from validator --keyring-backend=test --chain-id testing
{}

confirm transaction before signing and broadcasting [y/N]: y
{"height":"2","txhash":"71C1ED0A380EBF22547EEDE4550926D9421E00B250C7FF2D4EE2179E09358AAF","data":"0A110A0F6372656174656F72646572626F6F6B","raw_log":"[{\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"createorderbook\"}]}]}]","logs":[{"events":[{"type":"message","attributes":[{"key":"action","value":"createorderbook"}]}]}],"gas_wanted":"200000","gas_used":"49486"}
```

## upsertSignerKey

- CLI

```sh
# Secp256k1 type
sekaid tx kiraHub upsertSignerKey AnzIM9IcLb07Cvwq3hdMJuuRofAgxfDekkD3nJUPPw0w
sekaid tx kiraHub upsertSignerKey AnzIM9IcLb07Cvwq3hdMJuuRofAgxfDekkD3nJUPPw0w --key-type=Secp256k1

# Ed25519 type
sekaid tx kiraHub upsertSignerKey TXgDkmTYpPRwU/PvDbfbhbwiYA7jXMwQgNffHVey1dC644OBBI4OQdf4Tro6hzimT1dHYzPiGZB0aYWJBC2keQ== --key-type=Ed25519

# enabled false
sekaid tx kiraHub upsertSignerKey AnzIM9IcLb07Cvwq3hdMJuuRofAgxfDekkD3nJUPPw0w --enabled=false

# permissions
sekaid tx kiraHub upsertSignerKey AnzIM9IcLb07Cvwq3hdMJuuRofAgxfDekkD3nJUPPw0w --permissions=1,2

# expiry-time (set when this key expire if does not set it's automatically set to 10 days after current timestamp)
sekaid tx kiraHub upsertSignerKey AnzIM9IcLb07Cvwq3hdMJuuRofAgxfDekkD3nJUPPw0w --expiry-time=1598247750
```

- Rest

```sh
# TODO: should add readme for sample upsertSignerKey with same examples of CLI
```
---
`dev` branch
