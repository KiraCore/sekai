// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: kira/tokens/tx.proto

package types

import (
	context "context"
	fmt "fmt"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/cosmos/gogoproto/gogoproto"
	grpc1 "github.com/cosmos/gogoproto/grpc"
	proto "github.com/cosmos/gogoproto/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type MsgUpsertTokenInfo struct {
	Denom       string                                        `protobuf:"bytes,1,opt,name=denom,proto3" json:"denom,omitempty"`
	Rate        github_com_cosmos_cosmos_sdk_types.Dec        `protobuf:"bytes,2,opt,name=rate,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"rate"`
	FeePayments bool                                          `protobuf:"varint,3,opt,name=fee_payments,json=feePayments,proto3" json:"fee_payments,omitempty"`
	StakeCap    github_com_cosmos_cosmos_sdk_types.Dec        `protobuf:"bytes,4,opt,name=stake_cap,json=stakeCap,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"stake_cap" yaml:"stake_cap"`
	StakeMin    github_com_cosmos_cosmos_sdk_types.Int        `protobuf:"bytes,5,opt,name=stake_min,json=stakeMin,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"stake_min" yaml:"stake_min"`
	StakeToken  bool                                          `protobuf:"varint,6,opt,name=stake_token,json=stakeToken,proto3" json:"stake_token,omitempty"`
	Invalidated bool                                          `protobuf:"varint,7,opt,name=invalidated,proto3" json:"invalidated,omitempty"`
	Symbol      string                                        `protobuf:"bytes,8,opt,name=symbol,proto3" json:"symbol,omitempty"`
	Name        string                                        `protobuf:"bytes,9,opt,name=name,proto3" json:"name,omitempty"`
	Icon        string                                        `protobuf:"bytes,10,opt,name=icon,proto3" json:"icon,omitempty"`
	Decimals    uint32                                        `protobuf:"varint,11,opt,name=decimals,proto3" json:"decimals,omitempty"`
	Proposer    github_com_cosmos_cosmos_sdk_types.AccAddress `protobuf:"bytes,12,opt,name=proposer,proto3,casttype=github.com/cosmos/cosmos-sdk/types.AccAddress" json:"proposer,omitempty" yaml:"proposer"`
}

func (m *MsgUpsertTokenInfo) Reset()         { *m = MsgUpsertTokenInfo{} }
func (m *MsgUpsertTokenInfo) String() string { return proto.CompactTextString(m) }
func (*MsgUpsertTokenInfo) ProtoMessage()    {}
func (*MsgUpsertTokenInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_9bac5a72e1a3117c, []int{0}
}
func (m *MsgUpsertTokenInfo) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgUpsertTokenInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgUpsertTokenInfo.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgUpsertTokenInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgUpsertTokenInfo.Merge(m, src)
}
func (m *MsgUpsertTokenInfo) XXX_Size() int {
	return m.Size()
}
func (m *MsgUpsertTokenInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgUpsertTokenInfo.DiscardUnknown(m)
}

var xxx_messageInfo_MsgUpsertTokenInfo proto.InternalMessageInfo

func (m *MsgUpsertTokenInfo) GetDenom() string {
	if m != nil {
		return m.Denom
	}
	return ""
}

func (m *MsgUpsertTokenInfo) GetFeePayments() bool {
	if m != nil {
		return m.FeePayments
	}
	return false
}

func (m *MsgUpsertTokenInfo) GetStakeToken() bool {
	if m != nil {
		return m.StakeToken
	}
	return false
}

func (m *MsgUpsertTokenInfo) GetInvalidated() bool {
	if m != nil {
		return m.Invalidated
	}
	return false
}

func (m *MsgUpsertTokenInfo) GetSymbol() string {
	if m != nil {
		return m.Symbol
	}
	return ""
}

func (m *MsgUpsertTokenInfo) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *MsgUpsertTokenInfo) GetIcon() string {
	if m != nil {
		return m.Icon
	}
	return ""
}

func (m *MsgUpsertTokenInfo) GetDecimals() uint32 {
	if m != nil {
		return m.Decimals
	}
	return 0
}

func (m *MsgUpsertTokenInfo) GetProposer() github_com_cosmos_cosmos_sdk_types.AccAddress {
	if m != nil {
		return m.Proposer
	}
	return nil
}

type MsgUpsertTokenInfoResponse struct {
}

func (m *MsgUpsertTokenInfoResponse) Reset()         { *m = MsgUpsertTokenInfoResponse{} }
func (m *MsgUpsertTokenInfoResponse) String() string { return proto.CompactTextString(m) }
func (*MsgUpsertTokenInfoResponse) ProtoMessage()    {}
func (*MsgUpsertTokenInfoResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_9bac5a72e1a3117c, []int{1}
}
func (m *MsgUpsertTokenInfoResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgUpsertTokenInfoResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgUpsertTokenInfoResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgUpsertTokenInfoResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgUpsertTokenInfoResponse.Merge(m, src)
}
func (m *MsgUpsertTokenInfoResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgUpsertTokenInfoResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgUpsertTokenInfoResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgUpsertTokenInfoResponse proto.InternalMessageInfo

type MsgEthereumTx struct {
	TxType string `protobuf:"bytes,1,opt,name=tx_type,json=txType,proto3" json:"tx_type,omitempty"`
	Sender string `protobuf:"bytes,2,opt,name=sender,proto3" json:"sender,omitempty"`
	Hash   string `protobuf:"bytes,3,opt,name=hash,proto3" json:"hash,omitempty"`
	Data   []byte `protobuf:"bytes,4,opt,name=data,proto3" json:"data,omitempty"`
}

func (m *MsgEthereumTx) Reset()         { *m = MsgEthereumTx{} }
func (m *MsgEthereumTx) String() string { return proto.CompactTextString(m) }
func (*MsgEthereumTx) ProtoMessage()    {}
func (*MsgEthereumTx) Descriptor() ([]byte, []int) {
	return fileDescriptor_9bac5a72e1a3117c, []int{2}
}
func (m *MsgEthereumTx) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgEthereumTx) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgEthereumTx.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgEthereumTx) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgEthereumTx.Merge(m, src)
}
func (m *MsgEthereumTx) XXX_Size() int {
	return m.Size()
}
func (m *MsgEthereumTx) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgEthereumTx.DiscardUnknown(m)
}

var xxx_messageInfo_MsgEthereumTx proto.InternalMessageInfo

func (m *MsgEthereumTx) GetTxType() string {
	if m != nil {
		return m.TxType
	}
	return ""
}

func (m *MsgEthereumTx) GetSender() string {
	if m != nil {
		return m.Sender
	}
	return ""
}

func (m *MsgEthereumTx) GetHash() string {
	if m != nil {
		return m.Hash
	}
	return ""
}

func (m *MsgEthereumTx) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

type MsgEthereumTxResponse struct {
}

func (m *MsgEthereumTxResponse) Reset()         { *m = MsgEthereumTxResponse{} }
func (m *MsgEthereumTxResponse) String() string { return proto.CompactTextString(m) }
func (*MsgEthereumTxResponse) ProtoMessage()    {}
func (*MsgEthereumTxResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_9bac5a72e1a3117c, []int{3}
}
func (m *MsgEthereumTxResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgEthereumTxResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgEthereumTxResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgEthereumTxResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgEthereumTxResponse.Merge(m, src)
}
func (m *MsgEthereumTxResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgEthereumTxResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgEthereumTxResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgEthereumTxResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*MsgUpsertTokenInfo)(nil), "kira.tokens.MsgUpsertTokenInfo")
	proto.RegisterType((*MsgUpsertTokenInfoResponse)(nil), "kira.tokens.MsgUpsertTokenInfoResponse")
	proto.RegisterType((*MsgEthereumTx)(nil), "kira.tokens.MsgEthereumTx")
	proto.RegisterType((*MsgEthereumTxResponse)(nil), "kira.tokens.MsgEthereumTxResponse")
}

func init() { proto.RegisterFile("kira/tokens/tx.proto", fileDescriptor_9bac5a72e1a3117c) }

var fileDescriptor_9bac5a72e1a3117c = []byte{
	// 587 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x54, 0x41, 0x6f, 0xd3, 0x30,
	0x14, 0x6e, 0x58, 0xd7, 0xb5, 0x6e, 0xa7, 0x21, 0x6b, 0x30, 0x2b, 0x42, 0x49, 0xc9, 0x01, 0x7a,
	0x59, 0x22, 0xc1, 0x8d, 0xdb, 0xb2, 0x71, 0x98, 0xa0, 0x12, 0x8a, 0xc6, 0x05, 0x0e, 0xc5, 0x4b,
	0x5e, 0xd3, 0xa8, 0x8d, 0x1d, 0xd9, 0x1e, 0x6a, 0xf9, 0x15, 0xfc, 0x14, 0x7e, 0xc6, 0x8e, 0x3b,
	0x70, 0x40, 0x1c, 0x2a, 0xd4, 0xfe, 0x83, 0x1d, 0x39, 0xa1, 0x38, 0x69, 0x97, 0x75, 0x42, 0xb0,
	0x93, 0xdf, 0xfb, 0x3e, 0xfb, 0x7d, 0x79, 0xcf, 0x9f, 0x83, 0xf6, 0xc7, 0x89, 0xa0, 0x9e, 0xe2,
	0x63, 0x60, 0xd2, 0x53, 0x53, 0x37, 0x13, 0x5c, 0x71, 0xdc, 0xce, 0x51, 0xb7, 0x40, 0xcd, 0xfd,
	0x98, 0xc7, 0x5c, 0xe3, 0x5e, 0x1e, 0x15, 0x5b, 0x4c, 0xb3, 0x7a, 0x30, 0x13, 0x3c, 0xe3, 0x92,
	0x4e, 0x4a, 0xee, 0xe0, 0x56, 0xd1, 0x7c, 0x29, 0x09, 0x52, 0x25, 0x86, 0x02, 0xe0, 0x0b, 0x14,
	0x8c, 0xf3, 0xbd, 0x8e, 0x70, 0x5f, 0xc6, 0xef, 0x33, 0x09, 0x42, 0x9d, 0xe5, 0x1b, 0x4e, 0xd9,
	0x90, 0xe3, 0x7d, 0xb4, 0x1d, 0x01, 0xe3, 0x29, 0x31, 0xba, 0x46, 0xaf, 0x15, 0x14, 0x09, 0xf6,
	0x51, 0x5d, 0x50, 0x05, 0xe4, 0x41, 0x0e, 0xfa, 0xee, 0xe5, 0xdc, 0xae, 0xfd, 0x9c, 0xdb, 0xcf,
	0xe2, 0x44, 0x8d, 0x2e, 0xce, 0xdd, 0x90, 0xa7, 0x5e, 0xc8, 0x65, 0xca, 0x65, 0xb9, 0x1c, 0xca,
	0x68, 0xec, 0xa9, 0x59, 0x06, 0xd2, 0x3d, 0x81, 0x30, 0xd0, 0x67, 0xf1, 0x53, 0xd4, 0x19, 0x02,
	0x0c, 0x32, 0x3a, 0x4b, 0x81, 0x29, 0x49, 0xb6, 0xba, 0x46, 0xaf, 0x19, 0xb4, 0x87, 0x00, 0xef,
	0x4a, 0x08, 0x0f, 0x50, 0x4b, 0x2a, 0x3a, 0x86, 0x41, 0x48, 0x33, 0x52, 0xd7, 0x5a, 0xfe, 0xfd,
	0xb4, 0xae, 0xe7, 0xf6, 0xc3, 0x19, 0x4d, 0x27, 0xaf, 0x9c, 0x75, 0x21, 0x27, 0x68, 0xea, 0xf8,
	0x98, 0x66, 0x37, 0x02, 0x69, 0xc2, 0xc8, 0xf6, 0xbd, 0x05, 0x4e, 0x99, 0xda, 0x14, 0x48, 0x13,
	0xb6, 0x12, 0xe8, 0x27, 0x0c, 0xdb, 0xa8, 0x5d, 0xe0, 0x7a, 0xe4, 0xa4, 0xa1, 0x7b, 0x44, 0x1a,
	0xd2, 0x33, 0xc6, 0x5d, 0xd4, 0x4e, 0xd8, 0x67, 0x3a, 0x49, 0x22, 0xaa, 0x20, 0x22, 0x3b, 0xc5,
	0x10, 0x2a, 0x10, 0x7e, 0x8c, 0x1a, 0x72, 0x96, 0x9e, 0xf3, 0x09, 0x69, 0xea, 0x2b, 0x28, 0x33,
	0x8c, 0x51, 0x9d, 0xd1, 0x14, 0x48, 0x4b, 0xa3, 0x3a, 0xce, 0xb1, 0x24, 0xe4, 0x8c, 0xa0, 0x02,
	0xcb, 0x63, 0x6c, 0xa2, 0x66, 0x04, 0x61, 0x92, 0xd2, 0x89, 0x24, 0xed, 0xae, 0xd1, 0xdb, 0x0d,
	0xd6, 0x39, 0xfe, 0x84, 0x9a, 0x85, 0x73, 0x40, 0x90, 0x4e, 0xd7, 0xe8, 0x75, 0xfc, 0x93, 0xeb,
	0xb9, 0xbd, 0x57, 0x34, 0xb4, 0x62, 0x9c, 0xdf, 0x73, 0xfb, 0xf0, 0x3f, 0xa6, 0x71, 0x14, 0x86,
	0x47, 0x51, 0x24, 0x40, 0xca, 0x60, 0x5d, 0xd5, 0x79, 0x82, 0xcc, 0xbb, 0xae, 0x0a, 0x40, 0x66,
	0x9c, 0x49, 0x70, 0x46, 0x68, 0xb7, 0x2f, 0xe3, 0xd7, 0x6a, 0x04, 0x02, 0x2e, 0xd2, 0xb3, 0x29,
	0x3e, 0x40, 0x3b, 0x6a, 0x3a, 0xc8, 0xeb, 0x95, 0x86, 0x6b, 0xa8, 0xe9, 0xd9, 0x2c, 0x03, 0x3d,
	0x05, 0x60, 0x11, 0x88, 0xc2, 0x73, 0x41, 0x99, 0xe5, 0x1d, 0x8f, 0xa8, 0x1c, 0x69, 0xf7, 0xb4,
	0x02, 0x1d, 0xe7, 0x58, 0x44, 0x15, 0xd5, 0x8e, 0xe9, 0x04, 0x3a, 0x76, 0x0e, 0xd0, 0xa3, 0x5b,
	0x4a, 0xab, 0x4f, 0x78, 0xf1, 0xcd, 0x40, 0x5b, 0x7d, 0x19, 0xe3, 0x8f, 0x68, 0x6f, 0xd3, 0xfb,
	0xb6, 0x5b, 0x79, 0x85, 0xee, 0xdd, 0x36, 0xcc, 0xe7, 0xff, 0xd8, 0xb0, 0x12, 0xc1, 0x6f, 0x11,
	0xaa, 0x34, 0x69, 0x6e, 0x1e, 0xbb, 0xe1, 0x4c, 0xe7, 0xef, 0xdc, 0xaa, 0x9a, 0xef, 0x5f, 0x2e,
	0x2c, 0xe3, 0x6a, 0x61, 0x19, 0xbf, 0x16, 0x96, 0xf1, 0x75, 0x69, 0xd5, 0xae, 0x96, 0x56, 0xed,
	0xc7, 0xd2, 0xaa, 0x7d, 0xe8, 0x55, 0xae, 0xe9, 0x4d, 0x22, 0xe8, 0x31, 0x17, 0xe0, 0x49, 0x18,
	0xd3, 0xc4, 0x9b, 0xae, 0x7f, 0x07, 0xf9, 0x65, 0x9d, 0x37, 0xf4, 0xab, 0x7f, 0xf9, 0x27, 0x00,
	0x00, 0xff, 0xff, 0x8b, 0xa1, 0xd0, 0xfc, 0x7f, 0x04, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// MsgClient is the client API for Msg service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type MsgClient interface {
	// UpsertTokenInfo defines a method to upsert token rate
	UpsertTokenInfo(ctx context.Context, in *MsgUpsertTokenInfo, opts ...grpc.CallOption) (*MsgUpsertTokenInfoResponse, error)
	// EthereumTx defines a method to send ethereum transaction
	EthereumTx(ctx context.Context, in *MsgEthereumTx, opts ...grpc.CallOption) (*MsgEthereumTxResponse, error)
}

type msgClient struct {
	cc grpc1.ClientConn
}

func NewMsgClient(cc grpc1.ClientConn) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) UpsertTokenInfo(ctx context.Context, in *MsgUpsertTokenInfo, opts ...grpc.CallOption) (*MsgUpsertTokenInfoResponse, error) {
	out := new(MsgUpsertTokenInfoResponse)
	err := c.cc.Invoke(ctx, "/kira.tokens.Msg/UpsertTokenInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) EthereumTx(ctx context.Context, in *MsgEthereumTx, opts ...grpc.CallOption) (*MsgEthereumTxResponse, error) {
	out := new(MsgEthereumTxResponse)
	err := c.cc.Invoke(ctx, "/kira.tokens.Msg/EthereumTx", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
type MsgServer interface {
	// UpsertTokenInfo defines a method to upsert token rate
	UpsertTokenInfo(context.Context, *MsgUpsertTokenInfo) (*MsgUpsertTokenInfoResponse, error)
	// EthereumTx defines a method to send ethereum transaction
	EthereumTx(context.Context, *MsgEthereumTx) (*MsgEthereumTxResponse, error)
}

// UnimplementedMsgServer can be embedded to have forward compatible implementations.
type UnimplementedMsgServer struct {
}

func (*UnimplementedMsgServer) UpsertTokenInfo(ctx context.Context, req *MsgUpsertTokenInfo) (*MsgUpsertTokenInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpsertTokenInfo not implemented")
}
func (*UnimplementedMsgServer) EthereumTx(ctx context.Context, req *MsgEthereumTx) (*MsgEthereumTxResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EthereumTx not implemented")
}

func RegisterMsgServer(s grpc1.Server, srv MsgServer) {
	s.RegisterService(&_Msg_serviceDesc, srv)
}

func _Msg_UpsertTokenInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgUpsertTokenInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).UpsertTokenInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kira.tokens.Msg/UpsertTokenInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).UpsertTokenInfo(ctx, req.(*MsgUpsertTokenInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_EthereumTx_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgEthereumTx)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).EthereumTx(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kira.tokens.Msg/EthereumTx",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).EthereumTx(ctx, req.(*MsgEthereumTx))
	}
	return interceptor(ctx, in, info, handler)
}

var _Msg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "kira.tokens.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UpsertTokenInfo",
			Handler:    _Msg_UpsertTokenInfo_Handler,
		},
		{
			MethodName: "EthereumTx",
			Handler:    _Msg_EthereumTx_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "kira/tokens/tx.proto",
}

func (m *MsgUpsertTokenInfo) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgUpsertTokenInfo) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgUpsertTokenInfo) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Proposer) > 0 {
		i -= len(m.Proposer)
		copy(dAtA[i:], m.Proposer)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Proposer)))
		i--
		dAtA[i] = 0x62
	}
	if m.Decimals != 0 {
		i = encodeVarintTx(dAtA, i, uint64(m.Decimals))
		i--
		dAtA[i] = 0x58
	}
	if len(m.Icon) > 0 {
		i -= len(m.Icon)
		copy(dAtA[i:], m.Icon)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Icon)))
		i--
		dAtA[i] = 0x52
	}
	if len(m.Name) > 0 {
		i -= len(m.Name)
		copy(dAtA[i:], m.Name)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Name)))
		i--
		dAtA[i] = 0x4a
	}
	if len(m.Symbol) > 0 {
		i -= len(m.Symbol)
		copy(dAtA[i:], m.Symbol)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Symbol)))
		i--
		dAtA[i] = 0x42
	}
	if m.Invalidated {
		i--
		if m.Invalidated {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x38
	}
	if m.StakeToken {
		i--
		if m.StakeToken {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x30
	}
	{
		size := m.StakeMin.Size()
		i -= size
		if _, err := m.StakeMin.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTx(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x2a
	{
		size := m.StakeCap.Size()
		i -= size
		if _, err := m.StakeCap.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTx(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x22
	if m.FeePayments {
		i--
		if m.FeePayments {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x18
	}
	{
		size := m.Rate.Size()
		i -= size
		if _, err := m.Rate.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTx(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if len(m.Denom) > 0 {
		i -= len(m.Denom)
		copy(dAtA[i:], m.Denom)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Denom)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgUpsertTokenInfoResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgUpsertTokenInfoResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgUpsertTokenInfoResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *MsgEthereumTx) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgEthereumTx) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgEthereumTx) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Data) > 0 {
		i -= len(m.Data)
		copy(dAtA[i:], m.Data)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Data)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.Hash) > 0 {
		i -= len(m.Hash)
		copy(dAtA[i:], m.Hash)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Hash)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Sender) > 0 {
		i -= len(m.Sender)
		copy(dAtA[i:], m.Sender)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Sender)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.TxType) > 0 {
		i -= len(m.TxType)
		copy(dAtA[i:], m.TxType)
		i = encodeVarintTx(dAtA, i, uint64(len(m.TxType)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgEthereumTxResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgEthereumTxResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgEthereumTxResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func encodeVarintTx(dAtA []byte, offset int, v uint64) int {
	offset -= sovTx(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *MsgUpsertTokenInfo) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Denom)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = m.Rate.Size()
	n += 1 + l + sovTx(uint64(l))
	if m.FeePayments {
		n += 2
	}
	l = m.StakeCap.Size()
	n += 1 + l + sovTx(uint64(l))
	l = m.StakeMin.Size()
	n += 1 + l + sovTx(uint64(l))
	if m.StakeToken {
		n += 2
	}
	if m.Invalidated {
		n += 2
	}
	l = len(m.Symbol)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(m.Icon)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	if m.Decimals != 0 {
		n += 1 + sovTx(uint64(m.Decimals))
	}
	l = len(m.Proposer)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	return n
}

func (m *MsgUpsertTokenInfoResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *MsgEthereumTx) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.TxType)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(m.Sender)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(m.Hash)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(m.Data)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	return n
}

func (m *MsgEthereumTxResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func sovTx(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTx(x uint64) (n int) {
	return sovTx(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *MsgUpsertTokenInfo) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MsgUpsertTokenInfo: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgUpsertTokenInfo: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Denom", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Denom = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Rate", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Rate.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field FeePayments", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.FeePayments = bool(v != 0)
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field StakeCap", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.StakeCap.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field StakeMin", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.StakeMin.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field StakeToken", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.StakeToken = bool(v != 0)
		case 7:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Invalidated", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.Invalidated = bool(v != 0)
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Symbol", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Symbol = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 10:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Icon", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Icon = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 11:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Decimals", wireType)
			}
			m.Decimals = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Decimals |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 12:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Proposer", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Proposer = append(m.Proposer[:0], dAtA[iNdEx:postIndex]...)
			if m.Proposer == nil {
				m.Proposer = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *MsgUpsertTokenInfoResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MsgUpsertTokenInfoResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgUpsertTokenInfoResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *MsgEthereumTx) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MsgEthereumTx: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgEthereumTx: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TxType", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.TxType = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Sender", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Sender = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Hash", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Hash = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Data", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Data = append(m.Data[:0], dAtA[iNdEx:postIndex]...)
			if m.Data == nil {
				m.Data = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *MsgEthereumTxResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MsgEthereumTxResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgEthereumTxResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipTx(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTx
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowTx
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowTx
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthTx
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTx
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTx
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTx        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTx          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTx = fmt.Errorf("proto: unexpected end of group")
)
