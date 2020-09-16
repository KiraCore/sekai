module github.com/KiraCore/sekai/INTERX

go 1.12

require (
	github.com/KiraCore/sekai v0.0.0-20200823002648-c9c157f71380
	github.com/cosmos/cosmos-sdk v0.34.4-0.20200821154312-2e1fbaed9c41
	github.com/gofrs/uuid v3.2.0+incompatible
	github.com/golang/protobuf v1.4.2
	github.com/gorilla/mux v1.7.4
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.0.0-beta.4
	github.com/rakyll/statik v0.1.7
	github.com/tendermint/tendermint v0.34.0-rc3
	golang.org/x/crypto v0.0.0-20200709230013-948cd5f35899
	golang.org/x/text v0.3.2 // indirect
	google.golang.org/genproto v0.0.0-20200808173500-a06252235341
	google.golang.org/grpc v1.31.0
	google.golang.org/protobuf v1.25.0
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.4
