module github.com/KiraCore/sekai

go 1.14

require (
	github.com/asaskevich/govalidator v0.0.0-20200428143746-21a406dcc535
	github.com/coreos/etcd v3.3.13+incompatible
	github.com/cosmos/cosmos-sdk v0.34.4-0.20200917134338-9cb27fb171aa
	github.com/gogo/protobuf v1.3.1
	github.com/golang/protobuf v1.4.2
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/grpc-gateway v1.14.8
	github.com/otiai10/copy v1.2.0
	github.com/pkg/errors v0.9.1
	github.com/regen-network/cosmos-proto v0.3.0
	github.com/spf13/cast v1.3.1
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.6.1
	github.com/tendermint/tendermint v0.34.0-rc3.0.20200907055413-3359e0bf2f84
	github.com/tendermint/tm-db v0.6.2
	golang.org/x/crypto v0.0.0-20200820211705-5c72a883971a
	google.golang.org/grpc v1.32.0
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.4
