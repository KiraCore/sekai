// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: kira/ethereum/query.proto

package types

import (
	context "context"
	fmt "fmt"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/cosmos/gogoproto/gogoproto"
	grpc1 "github.com/cosmos/gogoproto/grpc"
	proto "github.com/cosmos/gogoproto/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

type RelayByAddressRequest struct {
	Addr github_com_cosmos_cosmos_sdk_types.AccAddress `protobuf:"bytes,1,opt,name=addr,proto3,casttype=github.com/cosmos/cosmos-sdk/types.AccAddress" json:"addr,omitempty" yaml:"addr"`
}

func (m *RelayByAddressRequest) Reset()         { *m = RelayByAddressRequest{} }
func (m *RelayByAddressRequest) String() string { return proto.CompactTextString(m) }
func (*RelayByAddressRequest) ProtoMessage()    {}
func (*RelayByAddressRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_cfdbf4c942f97f7e, []int{0}
}
func (m *RelayByAddressRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *RelayByAddressRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_RelayByAddressRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *RelayByAddressRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RelayByAddressRequest.Merge(m, src)
}
func (m *RelayByAddressRequest) XXX_Size() int {
	return m.Size()
}
func (m *RelayByAddressRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RelayByAddressRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RelayByAddressRequest proto.InternalMessageInfo

func (m *RelayByAddressRequest) GetAddr() github_com_cosmos_cosmos_sdk_types.AccAddress {
	if m != nil {
		return m.Addr
	}
	return nil
}

type RelayByAddressResponse struct {
	MsgRelay *MsgRelay `protobuf:"bytes,1,opt,name=msg_relay,json=msgRelay,proto3" json:"msg_relay,omitempty"`
}

func (m *RelayByAddressResponse) Reset()         { *m = RelayByAddressResponse{} }
func (m *RelayByAddressResponse) String() string { return proto.CompactTextString(m) }
func (*RelayByAddressResponse) ProtoMessage()    {}
func (*RelayByAddressResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_cfdbf4c942f97f7e, []int{1}
}
func (m *RelayByAddressResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *RelayByAddressResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_RelayByAddressResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *RelayByAddressResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RelayByAddressResponse.Merge(m, src)
}
func (m *RelayByAddressResponse) XXX_Size() int {
	return m.Size()
}
func (m *RelayByAddressResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_RelayByAddressResponse.DiscardUnknown(m)
}

var xxx_messageInfo_RelayByAddressResponse proto.InternalMessageInfo

func (m *RelayByAddressResponse) GetMsgRelay() *MsgRelay {
	if m != nil {
		return m.MsgRelay
	}
	return nil
}

func init() {
	proto.RegisterType((*RelayByAddressRequest)(nil), "kira.ethereum.RelayByAddressRequest")
	proto.RegisterType((*RelayByAddressResponse)(nil), "kira.ethereum.RelayByAddressResponse")
}

func init() { proto.RegisterFile("kira/ethereum/query.proto", fileDescriptor_cfdbf4c942f97f7e) }

var fileDescriptor_cfdbf4c942f97f7e = []byte{
	// 355 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0xcc, 0xce, 0x2c, 0x4a,
	0xd4, 0x4f, 0x2d, 0xc9, 0x48, 0x2d, 0x4a, 0x2d, 0xcd, 0xd5, 0x2f, 0x2c, 0x4d, 0x2d, 0xaa, 0xd4,
	0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x05, 0x49, 0xe9, 0xc1, 0xa4, 0xa4, 0x44, 0xd2, 0xf3,
	0xd3, 0xf3, 0xc1, 0x32, 0xfa, 0x20, 0x16, 0x44, 0x91, 0x94, 0x4c, 0x7a, 0x7e, 0x7e, 0x7a, 0x4e,
	0xaa, 0x7e, 0x62, 0x41, 0xa6, 0x7e, 0x62, 0x5e, 0x5e, 0x7e, 0x49, 0x62, 0x49, 0x66, 0x7e, 0x5e,
	0x31, 0x4c, 0x16, 0xd5, 0x74, 0x18, 0x03, 0x2a, 0x2b, 0x86, 0x2a, 0x5b, 0x52, 0x01, 0x11, 0x57,
	0xca, 0xe5, 0x12, 0x0d, 0x4a, 0xcd, 0x49, 0xac, 0x74, 0xaa, 0x74, 0x4c, 0x49, 0x29, 0x4a, 0x2d,
	0x2e, 0x0e, 0x4a, 0x2d, 0x2c, 0x4d, 0x2d, 0x2e, 0x11, 0x0a, 0xe1, 0x62, 0x49, 0x4c, 0x49, 0x29,
	0x92, 0x60, 0x54, 0x60, 0xd4, 0xe0, 0x71, 0x72, 0xf8, 0x74, 0x4f, 0x9e, 0xbb, 0x32, 0x31, 0x37,
	0xc7, 0x4a, 0x09, 0x24, 0xaa, 0xf4, 0xeb, 0x9e, 0xbc, 0x6e, 0x7a, 0x66, 0x49, 0x46, 0x69, 0x92,
	0x5e, 0x72, 0x7e, 0xae, 0x7e, 0x72, 0x7e, 0x71, 0x6e, 0x7e, 0x31, 0x94, 0xd2, 0x2d, 0x4e, 0xc9,
	0xd6, 0x2f, 0xa9, 0x2c, 0x48, 0x2d, 0xd6, 0x73, 0x4c, 0x4e, 0x86, 0x19, 0x0b, 0x36, 0x4d, 0xc9,
	0x8f, 0x4b, 0x0c, 0xdd, 0xba, 0xe2, 0x82, 0xfc, 0xbc, 0xe2, 0x54, 0x21, 0x13, 0x2e, 0xce, 0xdc,
	0xe2, 0xf4, 0xf8, 0x22, 0x90, 0x2c, 0xd8, 0x52, 0x6e, 0x23, 0x71, 0x3d, 0x94, 0x50, 0xd1, 0xf3,
	0x2d, 0x4e, 0x07, 0x6b, 0x0e, 0xe2, 0xc8, 0x85, 0xb2, 0x8c, 0x7a, 0x18, 0xb9, 0x58, 0x03, 0x41,
	0xe1, 0x28, 0xd4, 0xc4, 0xc8, 0xc5, 0x87, 0x6a, 0xb4, 0x90, 0x0a, 0x9a, 0x7e, 0xac, 0x1e, 0x95,
	0x52, 0x25, 0xa0, 0x0a, 0xe2, 0x3e, 0x25, 0xe5, 0xa6, 0xcb, 0x4f, 0x26, 0x33, 0xc9, 0x0a, 0x49,
	0xeb, 0xa3, 0x86, 0x24, 0xd8, 0xc1, 0xfa, 0xd5, 0x20, 0xdf, 0xd5, 0x3a, 0xb9, 0x9c, 0x78, 0x24,
	0xc7, 0x78, 0xe1, 0x91, 0x1c, 0xe3, 0x83, 0x47, 0x72, 0x8c, 0x13, 0x1e, 0xcb, 0x31, 0x5c, 0x78,
	0x2c, 0xc7, 0x70, 0xe3, 0xb1, 0x1c, 0x43, 0x94, 0x16, 0x52, 0x68, 0x79, 0x67, 0x16, 0x25, 0x3a,
	0xe7, 0x17, 0xa5, 0xea, 0x17, 0xa7, 0x66, 0x27, 0x66, 0xea, 0x57, 0x20, 0x45, 0x0b, 0x28, 0xd4,
	0x92, 0xd8, 0xc0, 0x51, 0x63, 0x0c, 0x08, 0x00, 0x00, 0xff, 0xff, 0x37, 0x16, 0xba, 0xf5, 0x30,
	0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// QueryClient is the client API for Query service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type QueryClient interface {
	RelayByAddress(ctx context.Context, in *RelayByAddressRequest, opts ...grpc.CallOption) (*RelayByAddressResponse, error)
}

type queryClient struct {
	cc grpc1.ClientConn
}

func NewQueryClient(cc grpc1.ClientConn) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) RelayByAddress(ctx context.Context, in *RelayByAddressRequest, opts ...grpc.CallOption) (*RelayByAddressResponse, error) {
	out := new(RelayByAddressResponse)
	err := c.cc.Invoke(ctx, "/kira.ethereum.Query/RelayByAddress", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
type QueryServer interface {
	RelayByAddress(context.Context, *RelayByAddressRequest) (*RelayByAddressResponse, error)
}

// UnimplementedQueryServer can be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func (*UnimplementedQueryServer) RelayByAddress(ctx context.Context, req *RelayByAddressRequest) (*RelayByAddressResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RelayByAddress not implemented")
}

func RegisterQueryServer(s grpc1.Server, srv QueryServer) {
	s.RegisterService(&_Query_serviceDesc, srv)
}

func _Query_RelayByAddress_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RelayByAddressRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).RelayByAddress(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kira.ethereum.Query/RelayByAddress",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).RelayByAddress(ctx, req.(*RelayByAddressRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Query_serviceDesc = grpc.ServiceDesc{
	ServiceName: "kira.ethereum.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RelayByAddress",
			Handler:    _Query_RelayByAddress_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "kira/ethereum/query.proto",
}

func (m *RelayByAddressRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *RelayByAddressRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *RelayByAddressRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Addr) > 0 {
		i -= len(m.Addr)
		copy(dAtA[i:], m.Addr)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.Addr)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *RelayByAddressResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *RelayByAddressResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *RelayByAddressResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.MsgRelay != nil {
		{
			size, err := m.MsgRelay.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintQuery(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintQuery(dAtA []byte, offset int, v uint64) int {
	offset -= sovQuery(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *RelayByAddressRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Addr)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *RelayByAddressResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.MsgRelay != nil {
		l = m.MsgRelay.Size()
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func sovQuery(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozQuery(x uint64) (n int) {
	return sovQuery(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *RelayByAddressRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: RelayByAddressRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: RelayByAddressRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Addr", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Addr = append(m.Addr[:0], dAtA[iNdEx:postIndex]...)
			if m.Addr == nil {
				m.Addr = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *RelayByAddressResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: RelayByAddressResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: RelayByAddressResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MsgRelay", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.MsgRelay == nil {
				m.MsgRelay = &MsgRelay{}
			}
			if err := m.MsgRelay.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func skipQuery(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowQuery
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
					return 0, ErrIntOverflowQuery
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
					return 0, ErrIntOverflowQuery
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
				return 0, ErrInvalidLengthQuery
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupQuery
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthQuery
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthQuery        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowQuery          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupQuery = fmt.Errorf("proto: unexpected end of group")
)
