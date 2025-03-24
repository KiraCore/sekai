// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: kira/ethereum/ethereum.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/cosmos-sdk/codec/types"
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

type EVMTx struct {
	From     string `protobuf:"bytes,1,opt,name=From,proto3" json:"From,omitempty"`
	To       string `protobuf:"bytes,2,opt,name=To,proto3" json:"To,omitempty"`
	Value    string `protobuf:"bytes,3,opt,name=Value,proto3" json:"Value,omitempty"`
	Gas      string `protobuf:"bytes,4,opt,name=Gas,proto3" json:"Gas,omitempty"`
	GasPrice string `protobuf:"bytes,5,opt,name=GasPrice,proto3" json:"GasPrice,omitempty"`
	Nonce    string `protobuf:"bytes,6,opt,name=Nonce,proto3" json:"Nonce,omitempty"`
	Data     string `protobuf:"bytes,7,opt,name=Data,proto3" json:"Data,omitempty"`
	ChainId  int64  `protobuf:"varint,8,opt,name=ChainId,proto3" json:"ChainId,omitempty"`
	V        string `protobuf:"bytes,9,opt,name=V,proto3" json:"V,omitempty"`
	R        string `protobuf:"bytes,10,opt,name=R,proto3" json:"R,omitempty"`
	S        string `protobuf:"bytes,11,opt,name=S,proto3" json:"S,omitempty"`
}

func (m *EVMTx) Reset()         { *m = EVMTx{} }
func (m *EVMTx) String() string { return proto.CompactTextString(m) }
func (*EVMTx) ProtoMessage()    {}
func (*EVMTx) Descriptor() ([]byte, []int) {
	return fileDescriptor_4ec492796f7b88e4, []int{0}
}
func (m *EVMTx) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *EVMTx) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_EVMTx.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *EVMTx) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EVMTx.Merge(m, src)
}
func (m *EVMTx) XXX_Size() int {
	return m.Size()
}
func (m *EVMTx) XXX_DiscardUnknown() {
	xxx_messageInfo_EVMTx.DiscardUnknown(m)
}

var xxx_messageInfo_EVMTx proto.InternalMessageInfo

func (m *EVMTx) GetFrom() string {
	if m != nil {
		return m.From
	}
	return ""
}

func (m *EVMTx) GetTo() string {
	if m != nil {
		return m.To
	}
	return ""
}

func (m *EVMTx) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

func (m *EVMTx) GetGas() string {
	if m != nil {
		return m.Gas
	}
	return ""
}

func (m *EVMTx) GetGasPrice() string {
	if m != nil {
		return m.GasPrice
	}
	return ""
}

func (m *EVMTx) GetNonce() string {
	if m != nil {
		return m.Nonce
	}
	return ""
}

func (m *EVMTx) GetData() string {
	if m != nil {
		return m.Data
	}
	return ""
}

func (m *EVMTx) GetChainId() int64 {
	if m != nil {
		return m.ChainId
	}
	return 0
}

func (m *EVMTx) GetV() string {
	if m != nil {
		return m.V
	}
	return ""
}

func (m *EVMTx) GetR() string {
	if m != nil {
		return m.R
	}
	return ""
}

func (m *EVMTx) GetS() string {
	if m != nil {
		return m.S
	}
	return ""
}

func init() {
	proto.RegisterType((*EVMTx)(nil), "kira.ethereum.EVMTx")
}

func init() { proto.RegisterFile("kira/ethereum/ethereum.proto", fileDescriptor_4ec492796f7b88e4) }

var fileDescriptor_4ec492796f7b88e4 = []byte{
	// 299 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x44, 0x90, 0xcd, 0x4a, 0x03, 0x31,
	0x10, 0xc7, 0x9b, 0x7e, 0x37, 0x7e, 0x20, 0xa1, 0x87, 0xb1, 0x48, 0x28, 0x9e, 0x8a, 0x87, 0xe6,
	0xe0, 0x1b, 0xd8, 0x6a, 0x11, 0x51, 0x64, 0x5b, 0xf6, 0xe0, 0x2d, 0xad, 0x71, 0x1b, 0xda, 0xdd,
	0x29, 0xd9, 0x5d, 0x68, 0xdf, 0xc2, 0xc7, 0xf2, 0xd8, 0xa3, 0x17, 0x41, 0x76, 0x5f, 0x44, 0x92,
	0xd0, 0x7a, 0xfb, 0xff, 0x7e, 0xff, 0x61, 0x06, 0x86, 0x5e, 0xad, 0xb4, 0x91, 0x42, 0x65, 0x4b,
	0x65, 0x54, 0x1e, 0x1f, 0xc3, 0x70, 0x63, 0x30, 0x43, 0x76, 0x66, 0xdb, 0xe1, 0x41, 0xf6, 0xba,
	0x11, 0x46, 0xe8, 0x1a, 0x61, 0x93, 0x1f, 0xea, 0x5d, 0x46, 0x88, 0xd1, 0x5a, 0x09, 0x47, 0xf3,
	0xfc, 0x43, 0xc8, 0x64, 0xe7, 0xab, 0xeb, 0x1f, 0x42, 0x1b, 0xf7, 0xe1, 0xf3, 0x6c, 0xcb, 0x18,
	0xad, 0x3f, 0x18, 0x8c, 0x81, 0xf4, 0xc9, 0xa0, 0x13, 0xb8, 0xcc, 0xce, 0x69, 0x75, 0x86, 0x50,
	0x75, 0xa6, 0x3a, 0x43, 0xd6, 0xa5, 0x8d, 0x50, 0xae, 0x73, 0x05, 0x35, 0xa7, 0x3c, 0xb0, 0x0b,
	0x5a, 0x9b, 0xc8, 0x14, 0xea, 0xce, 0xd9, 0xc8, 0x7a, 0xb4, 0x3d, 0x91, 0xe9, 0xab, 0xd1, 0x0b,
	0x05, 0x0d, 0xa7, 0x8f, 0x6c, 0x77, 0xbc, 0x60, 0xb2, 0x50, 0xd0, 0xf4, 0x3b, 0x1c, 0xd8, 0xeb,
	0x63, 0x99, 0x49, 0x68, 0xf9, 0xeb, 0x36, 0x33, 0xa0, 0xad, 0xd1, 0x52, 0xea, 0xe4, 0xf1, 0x1d,
	0xda, 0x7d, 0x32, 0xa8, 0x05, 0x07, 0x64, 0xa7, 0x94, 0x84, 0xd0, 0x71, 0xa3, 0x24, 0xb4, 0x14,
	0x00, 0xf5, 0x14, 0x58, 0x9a, 0xc2, 0x89, 0xa7, 0xe9, 0xdd, 0xf8, 0xab, 0xe0, 0x64, 0x5f, 0x70,
	0xf2, 0x5b, 0x70, 0xf2, 0x59, 0xf2, 0xca, 0xbe, 0xe4, 0x95, 0xef, 0x92, 0x57, 0xde, 0x6e, 0x22,
	0x9d, 0x2d, 0xf3, 0xf9, 0x70, 0x81, 0xb1, 0x78, 0xd2, 0x46, 0x8e, 0xd0, 0x28, 0x91, 0xaa, 0x95,
	0xd4, 0x62, 0xfb, 0xff, 0xee, 0x6c, 0xb7, 0x51, 0xe9, 0xbc, 0xe9, 0x9e, 0x75, 0xfb, 0x17, 0x00,
	0x00, 0xff, 0xff, 0x2b, 0x88, 0x2d, 0xe3, 0x8c, 0x01, 0x00, 0x00,
}

func (m *EVMTx) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *EVMTx) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *EVMTx) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.S) > 0 {
		i -= len(m.S)
		copy(dAtA[i:], m.S)
		i = encodeVarintEthereum(dAtA, i, uint64(len(m.S)))
		i--
		dAtA[i] = 0x5a
	}
	if len(m.R) > 0 {
		i -= len(m.R)
		copy(dAtA[i:], m.R)
		i = encodeVarintEthereum(dAtA, i, uint64(len(m.R)))
		i--
		dAtA[i] = 0x52
	}
	if len(m.V) > 0 {
		i -= len(m.V)
		copy(dAtA[i:], m.V)
		i = encodeVarintEthereum(dAtA, i, uint64(len(m.V)))
		i--
		dAtA[i] = 0x4a
	}
	if m.ChainId != 0 {
		i = encodeVarintEthereum(dAtA, i, uint64(m.ChainId))
		i--
		dAtA[i] = 0x40
	}
	if len(m.Data) > 0 {
		i -= len(m.Data)
		copy(dAtA[i:], m.Data)
		i = encodeVarintEthereum(dAtA, i, uint64(len(m.Data)))
		i--
		dAtA[i] = 0x3a
	}
	if len(m.Nonce) > 0 {
		i -= len(m.Nonce)
		copy(dAtA[i:], m.Nonce)
		i = encodeVarintEthereum(dAtA, i, uint64(len(m.Nonce)))
		i--
		dAtA[i] = 0x32
	}
	if len(m.GasPrice) > 0 {
		i -= len(m.GasPrice)
		copy(dAtA[i:], m.GasPrice)
		i = encodeVarintEthereum(dAtA, i, uint64(len(m.GasPrice)))
		i--
		dAtA[i] = 0x2a
	}
	if len(m.Gas) > 0 {
		i -= len(m.Gas)
		copy(dAtA[i:], m.Gas)
		i = encodeVarintEthereum(dAtA, i, uint64(len(m.Gas)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.Value) > 0 {
		i -= len(m.Value)
		copy(dAtA[i:], m.Value)
		i = encodeVarintEthereum(dAtA, i, uint64(len(m.Value)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.To) > 0 {
		i -= len(m.To)
		copy(dAtA[i:], m.To)
		i = encodeVarintEthereum(dAtA, i, uint64(len(m.To)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.From) > 0 {
		i -= len(m.From)
		copy(dAtA[i:], m.From)
		i = encodeVarintEthereum(dAtA, i, uint64(len(m.From)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintEthereum(dAtA []byte, offset int, v uint64) int {
	offset -= sovEthereum(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *EVMTx) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.From)
	if l > 0 {
		n += 1 + l + sovEthereum(uint64(l))
	}
	l = len(m.To)
	if l > 0 {
		n += 1 + l + sovEthereum(uint64(l))
	}
	l = len(m.Value)
	if l > 0 {
		n += 1 + l + sovEthereum(uint64(l))
	}
	l = len(m.Gas)
	if l > 0 {
		n += 1 + l + sovEthereum(uint64(l))
	}
	l = len(m.GasPrice)
	if l > 0 {
		n += 1 + l + sovEthereum(uint64(l))
	}
	l = len(m.Nonce)
	if l > 0 {
		n += 1 + l + sovEthereum(uint64(l))
	}
	l = len(m.Data)
	if l > 0 {
		n += 1 + l + sovEthereum(uint64(l))
	}
	if m.ChainId != 0 {
		n += 1 + sovEthereum(uint64(m.ChainId))
	}
	l = len(m.V)
	if l > 0 {
		n += 1 + l + sovEthereum(uint64(l))
	}
	l = len(m.R)
	if l > 0 {
		n += 1 + l + sovEthereum(uint64(l))
	}
	l = len(m.S)
	if l > 0 {
		n += 1 + l + sovEthereum(uint64(l))
	}
	return n
}

func sovEthereum(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozEthereum(x uint64) (n int) {
	return sovEthereum(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *EVMTx) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowEthereum
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
			return fmt.Errorf("proto: EVMTx: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: EVMTx: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field From", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEthereum
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
				return ErrInvalidLengthEthereum
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEthereum
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.From = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field To", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEthereum
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
				return ErrInvalidLengthEthereum
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEthereum
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.To = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Value", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEthereum
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
				return ErrInvalidLengthEthereum
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEthereum
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Value = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Gas", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEthereum
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
				return ErrInvalidLengthEthereum
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEthereum
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Gas = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field GasPrice", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEthereum
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
				return ErrInvalidLengthEthereum
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEthereum
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.GasPrice = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Nonce", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEthereum
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
				return ErrInvalidLengthEthereum
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEthereum
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Nonce = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Data", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEthereum
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
				return ErrInvalidLengthEthereum
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEthereum
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Data = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 8:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ChainId", wireType)
			}
			m.ChainId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEthereum
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ChainId |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field V", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEthereum
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
				return ErrInvalidLengthEthereum
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEthereum
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.V = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 10:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field R", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEthereum
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
				return ErrInvalidLengthEthereum
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEthereum
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.R = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 11:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field S", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEthereum
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
				return ErrInvalidLengthEthereum
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEthereum
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.S = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipEthereum(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthEthereum
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
func skipEthereum(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowEthereum
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
					return 0, ErrIntOverflowEthereum
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
					return 0, ErrIntOverflowEthereum
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
				return 0, ErrInvalidLengthEthereum
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupEthereum
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthEthereum
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthEthereum        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowEthereum          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupEthereum = fmt.Errorf("proto: unexpected end of group")
)
