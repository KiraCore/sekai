# Bugs you can meet

## clientCtx.GetFromAddress() is returning empty address, and more in general can't read command flags

This can be related to you didn't run ReadTxCommandFlags before the command.
```go
    clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
    if err != nil {
        return err
    }
```

## Error: UnmarshalBinaryBare expected to read prefix bytes 282816A9 (since it is registered concrete) but got 0A530A51...

```sh
confirm transaction before signing and broadcasting [y/N]: y
{"height":"148","txhash":"ED178AD70244D721495143898488850EF1369CDA9C8DE83A77B3B3FBC85D633D","data":"0A110A0F6372656174656F72646572626F6F6B","raw_log":"[{\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"createorderbook\"}]}]}]","logs":[{"events":[{"type":"message","attributes":[{"key":"action","value":"createorderbook"}]}]}],"gas_wanted":"200000","gas_used":"47127"}
sekaid query tx ED178AD70244D721495143898488850EF1369CDA9C8DE83A77B3B3FBC85D633D
```
I tried to debug but the only result I got is like below. It seems it's a cosmos-sdk internal bug.
```log
BBBB:: 
S
Q
 /kira.kiraHub.MsgCreateOrderBook-
?b?3quotmnemonic"i?9Ø?Y=~Qag???G
+
#
!?
  u^??Z??Toh
            ??i?G????x?9?%
??
  @?z??????32????D`??T?lw?A??T?L?Y??k?I\??.v?Aé???x? /
goroutine 1 [running]:
runtime/debug.Stack(0xd4, 0x0, 0x0)
	/usr/local/Cellar/go/1.14.7/libexec/src/runtime/debug/stack.go:24 +0x9d
runtime/debug.PrintStack()
	/usr/local/Cellar/go/1.14.7/libexec/src/runtime/debug/stack.go:16 +0x22
github.com/tendermint/go-amino.(*Codec).UnmarshalBinaryBare(0xc001019260, 0xc0001361a0, 0xcc, 0xcc, 0x52d5440, 0xc000e62000, 0xc0010f4101, 0xc000e62000)
	/Users/admin/go/pkg/mod/github.com/tendermint/go-amino@v0.15.1/amino.go:343 +0x6cf
github.com/cosmos/cosmos-sdk/codec.(*Codec).UnmarshalBinaryBare(0xc0001af378, 0xc0001361a0, 0xcc, 0xcc, 0x52d5440, 0xc000e62000, 0xc000fdbd20, 0x0)
	/Users/admin/go/pkg/mod/github.com/!kira!core/cosmos-sdk@v1.0.1-0.20200811231814-95fd54f999e1/codec/amino.go:114 +0x64
github.com/cosmos/cosmos-sdk/x/auth/client.parseTx(0xc0001af378, 0xc0001361a0, 0xcc, 0xcc, 0xc000fe2780, 0x5855220, 0xc001034660, 0x53b5146)
	/Users/admin/go/pkg/mod/github.com/!kira!core/cosmos-sdk@v1.0.1-0.20200811231814-95fd54f999e1/x/auth/client/query.go:147 +0x92
github.com/cosmos/cosmos-sdk/x/auth/client.formatTxResult(0xc0001af378, 0xc000fc9a20, 0xc0001b7700, 0xc000224ee8, 0xc0010f4180, 0x0)
	/Users/admin/go/pkg/mod/github.com/!kira!core/cosmos-sdk@v1.0.1-0.20200811231814-95fd54f999e1/x/auth/client/query.go:136 +0x58
github.com/cosmos/cosmos-sdk/x/auth/client.QueryTx(0x0, 0x0, 0x0, 0x589c880, 0xc0010f4180, 0x0, 0x0, 0x804a958, 0xc000fdbd20, 0x0, ...)
	/Users/admin/go/pkg/mod/github.com/!kira!core/cosmos-sdk@v1.0.1-0.20200811231814-95fd54f999e1/x/auth/client/query.go:91 +0x217
github.com/cosmos/cosmos-sdk/x/auth/client/cli.QueryTxCmd.func1(0xc001063080, 0xc0010b3a70, 0x1, 0x1, 0x0, 0x0)
	/Users/admin/go/pkg/mod/github.com/!kira!core/cosmos-sdk@v1.0.1-0.20200811231814-95fd54f999e1/x/auth/client/cli/query.go:209 +0x133
github.com/spf13/cobra.(*Command).execute(0xc001063080, 0xc0010b3a50, 0x1, 0x1, 0xc001063080, 0xc0010b3a50)
	/Users/admin/go/pkg/mod/github.com/spf13/cobra@v1.0.0/command.go:842 +0x453
github.com/spf13/cobra.(*Command).ExecuteC(0x64b8ae0, 0x0, 0x0, 0xc001011340)
	/Users/admin/go/pkg/mod/github.com/spf13/cobra@v1.0.0/command.go:950 +0x349
github.com/spf13/cobra.(*Command).Execute(...)
	/Users/admin/go/pkg/mod/github.com/spf13/cobra@v1.0.0/command.go:887
github.com/spf13/cobra.(*Command).ExecuteContext(...)
	/Users/admin/go/pkg/mod/github.com/spf13/cobra@v1.0.0/command.go:880
main.main()
	/Users/admin/sekai/cmd/sekaid/main.go:97 +0x158
```

## cosmos-sdk simd testing

start node
```sh
rm -rf $HOME/.simapp/
./build/simd init --chain-id=testing testing
./build/simd keys add validator --keyring-backend=test
./build/simd add-genesis-account $(./build/simd keys show validator -a --keyring-backend=test) 1000000000stake,1000000000validatortoken --keyring-backend=test
./build/simd gentx validator --keyring-backend=test --chain-id=testing
./build/simd collect-gentxs
./build/simd start
```

Test bank module.
```
./build/simd tx bank send validator cosmos1f339fylyn4czjz06l2s0laxa7377yr72qp880u 1000stake --keyring-backend=test --chain-id=testing
```

This bank module is common and it's working for sekaid too.
```
sekaid tx bank send validator cosmos1f339fylyn4czjz06l2s0laxa7377yr72qp880u 1000stake --keyring-backend=test --chain-id=testing
```

## guidance project upgrade for Cosmos SDK upgrade

- Check changes for cosmos-sdk/simapp and make following changes for the project.
- Check module changes for cosmos-sdk/x/bank or cosmos-sdk/x/auth and make following changes for all the modules of the project.

## How to check cosmos-sdk bug on local and report
[Ref](https://github.com/cosmos/cosmos-sdk/issues/7133)

```sh
# clone cosmos-sdk
git clone https://github.com/cosmos/cosmos-sdk

# build and run daemon
make build-simd
rm -rf $HOME/.simapp/
./build/simd init --chain-id=testing testing
./build/simd keys add validator --keyring-backend=test
./build/simd add-genesis-account $(./build/simd keys show validator -a --keyring-backend=test) 1000000000stake,1000000000validatortoken --keyring-backend=test
./build/simd gentx validator --keyring-backend=test --chain-id=testing
./build/simd collect-gentxs
./build/simd start

# test cli and check if it work
./build/simd tx bank send validator cosmos1f339fylyn4czjz06l2s0laxa7377yr72qp880u 1000stake --keyring-backend=test --chain-id=testing
```