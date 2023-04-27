package types

import (
	fmt "fmt"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	proto "github.com/gogo/protobuf/proto"
)

// PackTxMsgAny marshals the sdk.Msg payload to a protobuf Any type
func PackTxMsgAny(sdkMsg sdk.Msg) (*codectypes.Any, error) {
	msg, ok := sdkMsg.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("can't proto marshal %T", sdkMsg)
	}

	return codectypes.NewAnyWithValue(msg)
}

// UnpackInterfaces implements codectypes.UnpackInterfacesMessage
func (session DappSession) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var sdkMsg sdk.Msg

	for _, any := range session.OnchainMessages {
		err := unpacker.UnpackAny(any, &sdkMsg)
		if err != nil {
			return err
		}
	}
	return nil
}
