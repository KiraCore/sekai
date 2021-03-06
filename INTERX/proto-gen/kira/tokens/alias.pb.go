// Code generated by protoc-gen-go. DO NOT EDIT.
// source: kira/tokens/alias.proto

package proto

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type VoteType int32

const (
	VoteType_yes     VoteType = 0
	VoteType_no      VoteType = 1
	VoteType_veto    VoteType = 2
	VoteType_abstain VoteType = 3
)

var VoteType_name = map[int32]string{
	0: "yes",
	1: "no",
	2: "veto",
	3: "abstain",
}

var VoteType_value = map[string]int32{
	"yes":     0,
	"no":      1,
	"veto":    2,
	"abstain": 3,
}

func (x VoteType) String() string {
	return proto.EnumName(VoteType_name, int32(x))
}

func (VoteType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_e5f7b26b3e48e5a6, []int{0}
}

type ProposalStatus int32

const (
	ProposalStatus_undefined ProposalStatus = 0
	ProposalStatus_active    ProposalStatus = 1
	ProposalStatus_rejected  ProposalStatus = 2
	ProposalStatus_passed    ProposalStatus = 3
	ProposalStatus_enacted   ProposalStatus = 4
)

var ProposalStatus_name = map[int32]string{
	0: "undefined",
	1: "active",
	2: "rejected",
	3: "passed",
	4: "enacted",
}

var ProposalStatus_value = map[string]int32{
	"undefined": 0,
	"active":    1,
	"rejected":  2,
	"passed":    3,
	"enacted":   4,
}

func (x ProposalStatus) String() string {
	return proto.EnumName(ProposalStatus_name, int32(x))
}

func (ProposalStatus) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_e5f7b26b3e48e5a6, []int{1}
}

type TokenAlias struct {
	Symbol               string   `protobuf:"bytes,1,opt,name=symbol,proto3" json:"symbol,omitempty"`
	Name                 string   `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Icon                 string   `protobuf:"bytes,3,opt,name=icon,proto3" json:"icon,omitempty"`
	Decimals             uint32   `protobuf:"varint,4,opt,name=decimals,proto3" json:"decimals,omitempty"`
	Denoms               []string `protobuf:"bytes,5,rep,name=denoms,proto3" json:"denoms,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TokenAlias) Reset()         { *m = TokenAlias{} }
func (m *TokenAlias) String() string { return proto.CompactTextString(m) }
func (*TokenAlias) ProtoMessage()    {}
func (*TokenAlias) Descriptor() ([]byte, []int) {
	return fileDescriptor_e5f7b26b3e48e5a6, []int{0}
}

func (m *TokenAlias) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TokenAlias.Unmarshal(m, b)
}
func (m *TokenAlias) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TokenAlias.Marshal(b, m, deterministic)
}
func (m *TokenAlias) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TokenAlias.Merge(m, src)
}
func (m *TokenAlias) XXX_Size() int {
	return xxx_messageInfo_TokenAlias.Size(m)
}
func (m *TokenAlias) XXX_DiscardUnknown() {
	xxx_messageInfo_TokenAlias.DiscardUnknown(m)
}

var xxx_messageInfo_TokenAlias proto.InternalMessageInfo

func (m *TokenAlias) GetSymbol() string {
	if m != nil {
		return m.Symbol
	}
	return ""
}

func (m *TokenAlias) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *TokenAlias) GetIcon() string {
	if m != nil {
		return m.Icon
	}
	return ""
}

func (m *TokenAlias) GetDecimals() uint32 {
	if m != nil {
		return m.Decimals
	}
	return 0
}

func (m *TokenAlias) GetDenoms() []string {
	if m != nil {
		return m.Denoms
	}
	return nil
}

type MsgUpsertTokenAlias struct {
	Symbol               string   `protobuf:"bytes,1,opt,name=symbol,proto3" json:"symbol,omitempty"`
	Name                 string   `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Icon                 string   `protobuf:"bytes,3,opt,name=icon,proto3" json:"icon,omitempty"`
	Decimals             uint32   `protobuf:"varint,4,opt,name=decimals,proto3" json:"decimals,omitempty"`
	Denoms               []string `protobuf:"bytes,5,rep,name=denoms,proto3" json:"denoms,omitempty"`
	Proposer             []byte   `protobuf:"bytes,6,opt,name=proposer,proto3" json:"proposer,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MsgUpsertTokenAlias) Reset()         { *m = MsgUpsertTokenAlias{} }
func (m *MsgUpsertTokenAlias) String() string { return proto.CompactTextString(m) }
func (*MsgUpsertTokenAlias) ProtoMessage()    {}
func (*MsgUpsertTokenAlias) Descriptor() ([]byte, []int) {
	return fileDescriptor_e5f7b26b3e48e5a6, []int{1}
}

func (m *MsgUpsertTokenAlias) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MsgUpsertTokenAlias.Unmarshal(m, b)
}
func (m *MsgUpsertTokenAlias) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MsgUpsertTokenAlias.Marshal(b, m, deterministic)
}
func (m *MsgUpsertTokenAlias) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgUpsertTokenAlias.Merge(m, src)
}
func (m *MsgUpsertTokenAlias) XXX_Size() int {
	return xxx_messageInfo_MsgUpsertTokenAlias.Size(m)
}
func (m *MsgUpsertTokenAlias) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgUpsertTokenAlias.DiscardUnknown(m)
}

var xxx_messageInfo_MsgUpsertTokenAlias proto.InternalMessageInfo

func (m *MsgUpsertTokenAlias) GetSymbol() string {
	if m != nil {
		return m.Symbol
	}
	return ""
}

func (m *MsgUpsertTokenAlias) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *MsgUpsertTokenAlias) GetIcon() string {
	if m != nil {
		return m.Icon
	}
	return ""
}

func (m *MsgUpsertTokenAlias) GetDecimals() uint32 {
	if m != nil {
		return m.Decimals
	}
	return 0
}

func (m *MsgUpsertTokenAlias) GetDenoms() []string {
	if m != nil {
		return m.Denoms
	}
	return nil
}

func (m *MsgUpsertTokenAlias) GetProposer() []byte {
	if m != nil {
		return m.Proposer
	}
	return nil
}

func init() {
	proto.RegisterEnum("kira.tokens.VoteType", VoteType_name, VoteType_value)
	proto.RegisterEnum("kira.tokens.ProposalStatus", ProposalStatus_name, ProposalStatus_value)
	proto.RegisterType((*TokenAlias)(nil), "kira.tokens.TokenAlias")
	proto.RegisterType((*MsgUpsertTokenAlias)(nil), "kira.tokens.MsgUpsertTokenAlias")
}

func init() {
	proto.RegisterFile("kira/tokens/alias.proto", fileDescriptor_e5f7b26b3e48e5a6)
}

var fileDescriptor_e5f7b26b3e48e5a6 = []byte{
	// 386 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xc4, 0x92, 0x41, 0x6f, 0xd3, 0x40,
	0x10, 0x85, 0xeb, 0x38, 0xb8, 0xee, 0xb4, 0x85, 0xd5, 0x82, 0xc0, 0xea, 0xa5, 0x51, 0x0e, 0x28,
	0xaa, 0xd4, 0xac, 0x04, 0x37, 0x6e, 0x29, 0x70, 0x40, 0x08, 0x84, 0x4c, 0x40, 0x88, 0x13, 0x9b,
	0xdd, 0x21, 0x2c, 0xb1, 0x77, 0xac, 0x9d, 0x4d, 0x25, 0xdf, 0xf8, 0xa5, 0xfd, 0x03, 0xdc, 0x38,
	0x72, 0x42, 0xeb, 0x94, 0x88, 0x7f, 0xc0, 0x69, 0xde, 0x7b, 0xdf, 0xcc, 0x93, 0x2d, 0x2d, 0x3c,
	0xda, 0xb8, 0xa0, 0x55, 0xa4, 0x0d, 0x7a, 0x56, 0xba, 0x71, 0x9a, 0xe7, 0x5d, 0xa0, 0x48, 0xf2,
	0x38, 0x81, 0xf9, 0x0e, 0x9c, 0x3d, 0x58, 0xd3, 0x9a, 0x86, 0x5c, 0x25, 0xb5, 0x5b, 0x99, 0xfe,
	0xc8, 0x00, 0x96, 0x69, 0x61, 0x91, 0xee, 0xe4, 0x43, 0x28, 0xb8, 0x6f, 0x57, 0xd4, 0x54, 0xd9,
	0x24, 0x9b, 0x1d, 0xd5, 0xb7, 0x4e, 0x4a, 0x18, 0x7b, 0xdd, 0x62, 0x35, 0x1a, 0xd2, 0x41, 0xa7,
	0xcc, 0x19, 0xf2, 0x55, 0xbe, 0xcb, 0x92, 0x96, 0x67, 0x50, 0x5a, 0x34, 0xae, 0xd5, 0x0d, 0x57,
	0xe3, 0x49, 0x36, 0x3b, 0xad, 0xf7, 0x3e, 0x75, 0x5b, 0xf4, 0xd4, 0x72, 0x75, 0x67, 0x92, 0xa7,
	0xee, 0x9d, 0x9b, 0xfe, 0xcc, 0xe0, 0xfe, 0x1b, 0x5e, 0x7f, 0xe8, 0x18, 0x43, 0xfc, 0xbf, 0xdf,
	0x22, 0xbf, 0x40, 0xd9, 0x05, 0xea, 0x88, 0x31, 0x54, 0xc5, 0x24, 0x9b, 0x9d, 0x5c, 0xbd, 0xf8,
	0x75, 0x73, 0x7e, 0xaf, 0xd7, 0x6d, 0xf3, 0x6c, 0xfa, 0x97, 0x4c, 0x7f, 0xdf, 0x9c, 0x5f, 0xae,
	0x5d, 0xfc, 0xb6, 0x5d, 0xcd, 0x0d, 0xb5, 0xca, 0x10, 0xb7, 0xc4, 0xb7, 0xe3, 0x92, 0xed, 0x46,
	0xc5, 0xbe, 0x43, 0x9e, 0x2f, 0x8c, 0x59, 0x58, 0x1b, 0x90, 0xb9, 0xde, 0xb7, 0x5e, 0x3c, 0x81,
	0xf2, 0x23, 0x45, 0x5c, 0xf6, 0x1d, 0xca, 0x43, 0xc8, 0x7b, 0x64, 0x71, 0x20, 0x0b, 0x18, 0x79,
	0x12, 0x99, 0x2c, 0x61, 0x7c, 0x8d, 0x91, 0xc4, 0x48, 0x1e, 0xc3, 0xa1, 0x5e, 0x71, 0xd4, 0xce,
	0x8b, 0xfc, 0xa2, 0x86, 0xbb, 0xef, 0x86, 0x7b, 0xdd, 0xbc, 0x8f, 0x3a, 0x6e, 0x59, 0x9e, 0xc2,
	0xd1, 0xd6, 0x5b, 0xfc, 0xea, 0x3c, 0x5a, 0x71, 0x20, 0x01, 0x0a, 0x6d, 0xa2, 0xbb, 0x46, 0x91,
	0xc9, 0x13, 0x28, 0x03, 0x7e, 0x47, 0x13, 0xd1, 0x8a, 0x51, 0x22, 0x9d, 0x66, 0x46, 0x2b, 0xf2,
	0xd4, 0x89, 0x5e, 0x0f, 0x60, 0x7c, 0x35, 0xfb, 0xfc, 0xf8, 0x9f, 0x5f, 0x78, 0xed, 0x82, 0x7e,
	0x4e, 0x01, 0x15, 0xe3, 0x46, 0x3b, 0xf5, 0xea, 0xed, 0xf2, 0x65, 0xfd, 0x49, 0x0d, 0x4f, 0x64,
	0x55, 0x0c, 0xe3, 0xe9, 0x9f, 0x00, 0x00, 0x00, 0xff, 0xff, 0xf5, 0x55, 0xdc, 0x97, 0x67, 0x02,
	0x00, 0x00,
}
