# sekai
Kira Hub

## upsertSignerKey

- CLI

```sh
# Secp256k1 type
sekaicli tx kiraHub upsertSignerKey AnzIM9IcLb07Cvwq3hdMJuuRofAgxfDekkD3nJUPPw0w
sekaicli tx kiraHub upsertSignerKey AnzIM9IcLb07Cvwq3hdMJuuRofAgxfDekkD3nJUPPw0w --key-type=Secp256k1

# Ed25519 type
sekaicli tx kiraHub upsertSignerKey TXgDkmTYpPRwU/PvDbfbhbwiYA7jXMwQgNffHVey1dC644OBBI4OQdf4Tro6hzimT1dHYzPiGZB0aYWJBC2keQ== --key-type=Ed25519

# enabled false
sekaicli tx kiraHub upsertSignerKey AnzIM9IcLb07Cvwq3hdMJuuRofAgxfDekkD3nJUPPw0w --enabled=false

# permissions
sekaicli tx kiraHub upsertSignerKey AnzIM9IcLb07Cvwq3hdMJuuRofAgxfDekkD3nJUPPw0w --permissions=1,2

# expiry-time (set when this key expire if does not set it's automatically set to 10 days after current timestamp)
sekaicli tx kiraHub upsertSignerKey AnzIM9IcLb07Cvwq3hdMJuuRofAgxfDekkD3nJUPPw0w --expiry-time=1598247750
```

- Rest

```sh
# TODO: should add readme for sample upsertSignerKey with same examples of CLI
```
---
`dev` branch
