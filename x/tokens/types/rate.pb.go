// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: kira/tokens/rate.proto

package types

import (
	fmt "fmt"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
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

type TokenRate struct {
	Denom       string                                 `protobuf:"bytes,1,opt,name=denom,proto3" json:"denom,omitempty"`
	FeeRate     github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,2,opt,name=fee_rate,json=feeRate,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"fee_rate" yaml:"fee_rate"`
	FeePayments bool                                   `protobuf:"varint,3,opt,name=fee_payments,json=feePayments,proto3" json:"fee_payments,omitempty"`
	StakeCap    github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,4,opt,name=stake_cap,json=stakeCap,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"stake_cap" yaml:"stake_cap"`
	StakeMin    github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,5,opt,name=stake_min,json=stakeMin,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"stake_min" yaml:"stake_min"`
	StakeToken  bool                                   `protobuf:"varint,6,opt,name=stake_token,json=stakeToken,proto3" json:"stake_token,omitempty"`
	Invalidated bool                                   `protobuf:"varint,7,opt,name=invalidated,proto3" json:"invalidated,omitempty"`
}

func (m *TokenRate) Reset()         { *m = TokenRate{} }
func (m *TokenRate) String() string { return proto.CompactTextString(m) }
func (*TokenRate) ProtoMessage()    {}
func (*TokenRate) Descriptor() ([]byte, []int) {
	return fileDescriptor_d415c64b17c96dda, []int{0}
}
func (m *TokenRate) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TokenRate) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TokenRate.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TokenRate) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TokenRate.Merge(m, src)
}
func (m *TokenRate) XXX_Size() int {
	return m.Size()
}
func (m *TokenRate) XXX_DiscardUnknown() {
	xxx_messageInfo_TokenRate.DiscardUnknown(m)
}

var xxx_messageInfo_TokenRate proto.InternalMessageInfo

func (m *TokenRate) GetDenom() string {
	if m != nil {
		return m.Denom
	}
	return ""
}

func (m *TokenRate) GetFeePayments() bool {
	if m != nil {
		return m.FeePayments
	}
	return false
}

func (m *TokenRate) GetStakeToken() bool {
	if m != nil {
		return m.StakeToken
	}
	return false
}

func (m *TokenRate) GetInvalidated() bool {
	if m != nil {
		return m.Invalidated
	}
	return false
}

type MsgUpsertTokenRate struct {
	Denom       string                                        `protobuf:"bytes,1,opt,name=denom,proto3" json:"denom,omitempty"`
	Rate        github_com_cosmos_cosmos_sdk_types.Dec        `protobuf:"bytes,2,opt,name=rate,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"rate" yaml:"rate"`
	FeePayments bool                                          `protobuf:"varint,3,opt,name=fee_payments,json=feePayments,proto3" json:"fee_payments,omitempty"`
	StakeCap    github_com_cosmos_cosmos_sdk_types.Dec        `protobuf:"bytes,4,opt,name=stake_cap,json=stakeCap,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"stake_cap" yaml:"stake_cap"`
	StakeMin    github_com_cosmos_cosmos_sdk_types.Int        `protobuf:"bytes,5,opt,name=stake_min,json=stakeMin,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"stake_min" yaml:"stake_min"`
	StakeToken  bool                                          `protobuf:"varint,6,opt,name=stake_token,json=stakeToken,proto3" json:"stake_token,omitempty"`
	Invalidated bool                                          `protobuf:"varint,7,opt,name=invalidated,proto3" json:"invalidated,omitempty"`
	Proposer    github_com_cosmos_cosmos_sdk_types.AccAddress `protobuf:"bytes,8,opt,name=proposer,proto3,casttype=github.com/cosmos/cosmos-sdk/types.AccAddress" json:"proposer,omitempty" yaml:"proposer"`
}

func (m *MsgUpsertTokenRate) Reset()         { *m = MsgUpsertTokenRate{} }
func (m *MsgUpsertTokenRate) String() string { return proto.CompactTextString(m) }
func (*MsgUpsertTokenRate) ProtoMessage()    {}
func (*MsgUpsertTokenRate) Descriptor() ([]byte, []int) {
	return fileDescriptor_d415c64b17c96dda, []int{1}
}
func (m *MsgUpsertTokenRate) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgUpsertTokenRate) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgUpsertTokenRate.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgUpsertTokenRate) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgUpsertTokenRate.Merge(m, src)
}
func (m *MsgUpsertTokenRate) XXX_Size() int {
	return m.Size()
}
func (m *MsgUpsertTokenRate) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgUpsertTokenRate.DiscardUnknown(m)
}

var xxx_messageInfo_MsgUpsertTokenRate proto.InternalMessageInfo

func (m *MsgUpsertTokenRate) GetDenom() string {
	if m != nil {
		return m.Denom
	}
	return ""
}

func (m *MsgUpsertTokenRate) GetFeePayments() bool {
	if m != nil {
		return m.FeePayments
	}
	return false
}

func (m *MsgUpsertTokenRate) GetStakeToken() bool {
	if m != nil {
		return m.StakeToken
	}
	return false
}

func (m *MsgUpsertTokenRate) GetInvalidated() bool {
	if m != nil {
		return m.Invalidated
	}
	return false
}

func (m *MsgUpsertTokenRate) GetProposer() github_com_cosmos_cosmos_sdk_types.AccAddress {
	if m != nil {
		return m.Proposer
	}
	return nil
}

func init() {
	proto.RegisterType((*TokenRate)(nil), "kira.tokens.TokenRate")
	proto.RegisterType((*MsgUpsertTokenRate)(nil), "kira.tokens.MsgUpsertTokenRate")
}

func init() { proto.RegisterFile("kira/tokens/rate.proto", fileDescriptor_d415c64b17c96dda) }

var fileDescriptor_d415c64b17c96dda = []byte{
	// 430 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xec, 0x94, 0xb1, 0x6e, 0xd4, 0x30,
	0x18, 0xc7, 0x2f, 0xf4, 0xda, 0xe6, 0x9c, 0x4a, 0x20, 0xab, 0x42, 0x16, 0x43, 0x72, 0x64, 0x40,
	0x59, 0x9a, 0x0c, 0x6c, 0x48, 0x0c, 0x97, 0x76, 0x41, 0xa8, 0x12, 0x44, 0xb0, 0x20, 0xa4, 0xc3,
	0x4d, 0xbe, 0x0b, 0x56, 0x1a, 0xdb, 0xb2, 0x0d, 0xe2, 0x9e, 0x80, 0x95, 0xb7, 0xe0, 0x55, 0x3a,
	0x76, 0x44, 0x0c, 0x11, 0xba, 0x7b, 0x03, 0x46, 0x26, 0x14, 0xbb, 0x69, 0x4f, 0x2c, 0xf4, 0xf6,
	0x4e, 0x71, 0x7e, 0xb6, 0xfe, 0xff, 0xef, 0xfb, 0xfe, 0x89, 0xd1, 0xc3, 0x86, 0x29, 0x9a, 0x19,
	0xd1, 0x00, 0xd7, 0x99, 0xa2, 0x06, 0x52, 0xa9, 0x84, 0x11, 0x38, 0xe8, 0x79, 0xea, 0xf8, 0xa3,
	0xc3, 0x5a, 0xd4, 0xc2, 0xf2, 0xac, 0x5f, 0xb9, 0x23, 0xf1, 0xf7, 0x1d, 0x34, 0x79, 0xd3, 0x1f,
	0x28, 0xa8, 0x01, 0x7c, 0x88, 0x76, 0x2b, 0xe0, 0xa2, 0x25, 0xde, 0xd4, 0x4b, 0x26, 0x85, 0x7b,
	0xc1, 0xef, 0x91, 0xbf, 0x00, 0x98, 0xf7, 0xc2, 0xe4, 0x5e, 0xbf, 0x91, 0xcf, 0x2e, 0xba, 0x68,
	0xf4, 0xb3, 0x8b, 0x9e, 0xd4, 0xcc, 0x7c, 0xfc, 0x74, 0x96, 0x96, 0xa2, 0xcd, 0x4a, 0xa1, 0x5b,
	0xa1, 0xaf, 0x1e, 0x47, 0xba, 0x6a, 0x32, 0xb3, 0x94, 0xa0, 0xd3, 0x13, 0x28, 0x7f, 0x77, 0xd1,
	0xfd, 0x25, 0x6d, 0xcf, 0x9f, 0xc5, 0x83, 0x4e, 0x5c, 0xec, 0x2f, 0x00, 0xac, 0xe7, 0x63, 0x74,
	0xd0, 0x53, 0x49, 0x97, 0x2d, 0x70, 0xa3, 0xc9, 0xce, 0xd4, 0x4b, 0xfc, 0x22, 0x58, 0x00, 0xbc,
	0xba, 0x42, 0x78, 0x8e, 0x26, 0xda, 0xd0, 0x06, 0xe6, 0x25, 0x95, 0x64, 0x6c, 0x2b, 0xc8, 0xb7,
	0xae, 0xe0, 0x81, 0xab, 0xe0, 0x5a, 0x28, 0x2e, 0x7c, 0xbb, 0x3e, 0xa6, 0xf2, 0xc6, 0xa0, 0x65,
	0x9c, 0xec, 0x6e, 0x6d, 0xf0, 0x82, 0x9b, 0x7f, 0x0d, 0x5a, 0xc6, 0x07, 0x83, 0x53, 0xc6, 0x71,
	0x84, 0x02, 0xc7, 0x6d, 0x18, 0x64, 0xcf, 0xf6, 0x88, 0x2c, 0xb2, 0xd3, 0xc7, 0x53, 0x14, 0x30,
	0xfe, 0x99, 0x9e, 0xb3, 0x8a, 0x1a, 0xa8, 0xc8, 0xbe, 0x1b, 0xc2, 0x06, 0x8a, 0xbf, 0x8e, 0x11,
	0x3e, 0xd5, 0xf5, 0x5b, 0xa9, 0x41, 0x99, 0xff, 0x45, 0xf6, 0x1a, 0x8d, 0x37, 0xe2, 0x7a, 0xbe,
	0xf5, 0xb0, 0x02, 0xd7, 0x8b, 0x8b, 0xca, 0x4a, 0xdd, 0xe5, 0x74, 0xbb, 0x9c, 0xf0, 0x07, 0xe4,
	0x4b, 0x25, 0xa4, 0xd0, 0xa0, 0x88, 0x3f, 0xf5, 0x92, 0x83, 0xfc, 0xe4, 0xe6, 0xfb, 0x1f, 0x76,
	0xe2, 0x3f, 0x5d, 0x74, 0x74, 0x8b, 0x8a, 0x67, 0x65, 0x39, 0xab, 0x2a, 0x05, 0x5a, 0x17, 0xd7,
	0xaa, 0x79, 0x7e, 0xb1, 0x0a, 0xbd, 0xcb, 0x55, 0xe8, 0xfd, 0x5a, 0x85, 0xde, 0xb7, 0x75, 0x38,
	0xba, 0x5c, 0x87, 0xa3, 0x1f, 0xeb, 0x70, 0xf4, 0x2e, 0xd9, 0x90, 0x7c, 0xc9, 0x14, 0x3d, 0x16,
	0x0a, 0x32, 0x0d, 0x0d, 0x65, 0xd9, 0x97, 0xe1, 0x7e, 0xb0, 0xc2, 0x67, 0x7b, 0xf6, 0xf7, 0x7f,
	0xfa, 0x37, 0x00, 0x00, 0xff, 0xff, 0xf6, 0xd1, 0xe1, 0x89, 0x3b, 0x04, 0x00, 0x00,
}

func (m *TokenRate) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TokenRate) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TokenRate) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
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
		i = encodeVarintRate(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x2a
	{
		size := m.StakeCap.Size()
		i -= size
		if _, err := m.StakeCap.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintRate(dAtA, i, uint64(size))
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
		size := m.FeeRate.Size()
		i -= size
		if _, err := m.FeeRate.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintRate(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if len(m.Denom) > 0 {
		i -= len(m.Denom)
		copy(dAtA[i:], m.Denom)
		i = encodeVarintRate(dAtA, i, uint64(len(m.Denom)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgUpsertTokenRate) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgUpsertTokenRate) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgUpsertTokenRate) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Proposer) > 0 {
		i -= len(m.Proposer)
		copy(dAtA[i:], m.Proposer)
		i = encodeVarintRate(dAtA, i, uint64(len(m.Proposer)))
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
		i = encodeVarintRate(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x2a
	{
		size := m.StakeCap.Size()
		i -= size
		if _, err := m.StakeCap.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintRate(dAtA, i, uint64(size))
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
		i = encodeVarintRate(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if len(m.Denom) > 0 {
		i -= len(m.Denom)
		copy(dAtA[i:], m.Denom)
		i = encodeVarintRate(dAtA, i, uint64(len(m.Denom)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintRate(dAtA []byte, offset int, v uint64) int {
	offset -= sovRate(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *TokenRate) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Denom)
	if l > 0 {
		n += 1 + l + sovRate(uint64(l))
	}
	l = m.FeeRate.Size()
	n += 1 + l + sovRate(uint64(l))
	if m.FeePayments {
		n += 2
	}
	l = m.StakeCap.Size()
	n += 1 + l + sovRate(uint64(l))
	l = m.StakeMin.Size()
	n += 1 + l + sovRate(uint64(l))
	if m.StakeToken {
		n += 2
	}
	if m.Invalidated {
		n += 2
	}
	return n
}

func (m *MsgUpsertTokenRate) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Denom)
	if l > 0 {
		n += 1 + l + sovRate(uint64(l))
	}
	l = m.Rate.Size()
	n += 1 + l + sovRate(uint64(l))
	if m.FeePayments {
		n += 2
	}
	l = m.StakeCap.Size()
	n += 1 + l + sovRate(uint64(l))
	l = m.StakeMin.Size()
	n += 1 + l + sovRate(uint64(l))
	if m.StakeToken {
		n += 2
	}
	if m.Invalidated {
		n += 2
	}
	l = len(m.Proposer)
	if l > 0 {
		n += 1 + l + sovRate(uint64(l))
	}
	return n
}

func sovRate(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozRate(x uint64) (n int) {
	return sovRate(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *TokenRate) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowRate
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
			return fmt.Errorf("proto: TokenRate: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TokenRate: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Denom", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRate
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
				return ErrInvalidLengthRate
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthRate
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Denom = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field FeeRate", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRate
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
				return ErrInvalidLengthRate
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthRate
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.FeeRate.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
					return ErrIntOverflowRate
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
					return ErrIntOverflowRate
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
				return ErrInvalidLengthRate
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthRate
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
					return ErrIntOverflowRate
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
				return ErrInvalidLengthRate
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthRate
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
					return ErrIntOverflowRate
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
					return ErrIntOverflowRate
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
		default:
			iNdEx = preIndex
			skippy, err := skipRate(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthRate
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
func (m *MsgUpsertTokenRate) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowRate
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
			return fmt.Errorf("proto: MsgUpsertTokenRate: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgUpsertTokenRate: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Denom", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRate
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
				return ErrInvalidLengthRate
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthRate
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
					return ErrIntOverflowRate
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
				return ErrInvalidLengthRate
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthRate
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
					return ErrIntOverflowRate
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
					return ErrIntOverflowRate
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
				return ErrInvalidLengthRate
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthRate
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
					return ErrIntOverflowRate
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
				return ErrInvalidLengthRate
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthRate
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
					return ErrIntOverflowRate
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
					return ErrIntOverflowRate
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
				return fmt.Errorf("proto: wrong wireType = %d for field Proposer", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRate
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
				return ErrInvalidLengthRate
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthRate
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
			skippy, err := skipRate(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthRate
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
func skipRate(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowRate
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
					return 0, ErrIntOverflowRate
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
					return 0, ErrIntOverflowRate
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
				return 0, ErrInvalidLengthRate
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupRate
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthRate
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthRate        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowRate          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupRate = fmt.Errorf("proto: unexpected end of group")
)
