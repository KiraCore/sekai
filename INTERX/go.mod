module github.com/KiraCore/sekai/INTERX

go 1.14

require (
	github.com/FactomProject/basen v0.0.0-20150613233007-fe3947df716e // indirect
	github.com/FactomProject/btcutilecc v0.0.0-20130527213604-d3a63a5752ec // indirect
	github.com/KiraCore/sekai v0.0.0-20200823002648-c9c157f71380
	github.com/cosmos/cosmos-sdk v0.40.0-rc4
	github.com/gofrs/uuid v3.2.0+incompatible
	github.com/golang/protobuf v1.4.3
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/grpc-gateway v1.15.2
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.0.1
	github.com/iancoleman/strcase v0.1.2
	github.com/igorsobreira/kvstore v0.0.0-20131025205959-a8574822a4b3
	github.com/inhies/go-bytesize v0.0.0-20200716184324-4fe85e9b81b2
	github.com/juju/fslock v0.0.0-20160525022230-4d5c94c67b4b
	github.com/nightlyone/lockfile v1.0.0
	github.com/rakyll/statik v0.1.7
	github.com/rs/cors v1.7.0
	github.com/sonyarouje/simdb v0.0.0-20181202125413-c2488dfc374a
	github.com/tendermint/tendermint v0.34.0-rc6
	github.com/tyler-smith/go-bip32 v0.0.0-20170922074101-2c9cfd177564
	github.com/tyler-smith/go-bip39 v1.0.2
	golang.org/x/crypto v0.0.0-20201012173705-84dcc777aaee
	google.golang.org/genproto v0.0.0-20201019141844-1ed22bb0c154
	google.golang.org/grpc v1.33.1
	google.golang.org/protobuf v1.25.0
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.4

replace github.com/KiraCore/sekai => ../
