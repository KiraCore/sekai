// Code generated by protoc-gen-go. DO NOT EDIT.
// source: query.proto

package proto

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

type Endpoint struct {
	Url                  string   `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
	Method               string   `protobuf:"bytes,2,opt,name=method,proto3" json:"method,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Endpoint) Reset()         { *m = Endpoint{} }
func (m *Endpoint) String() string { return proto.CompactTextString(m) }
func (*Endpoint) ProtoMessage()    {}
func (*Endpoint) Descriptor() ([]byte, []int) {
	return fileDescriptor_5c6ac9b241082464, []int{0}
}

func (m *Endpoint) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Endpoint.Unmarshal(m, b)
}
func (m *Endpoint) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Endpoint.Marshal(b, m, deterministic)
}
func (m *Endpoint) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Endpoint.Merge(m, src)
}
func (m *Endpoint) XXX_Size() int {
	return xxx_messageInfo_Endpoint.Size(m)
}
func (m *Endpoint) XXX_DiscardUnknown() {
	xxx_messageInfo_Endpoint.DiscardUnknown(m)
}

var xxx_messageInfo_Endpoint proto.InternalMessageInfo

func (m *Endpoint) GetUrl() string {
	if m != nil {
		return m.Url
	}
	return ""
}

func (m *Endpoint) GetMethod() string {
	if m != nil {
		return m.Method
	}
	return ""
}

type RPCMethod struct {
	Description          string   `protobuf:"bytes,1,opt,name=description,proto3" json:"description,omitempty"`
	Enabled              bool     `protobuf:"varint,2,opt,name=enabled,proto3" json:"enabled,omitempty"`
	RateLimit            float64  `protobuf:"fixed64,3,opt,name=rate_limit,json=rateLimit,proto3" json:"rate_limit,omitempty"`
	AuthRateLimit        float64  `protobuf:"fixed64,4,opt,name=auth_rate_limit,json=authRateLimit,proto3" json:"auth_rate_limit,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RPCMethod) Reset()         { *m = RPCMethod{} }
func (m *RPCMethod) String() string { return proto.CompactTextString(m) }
func (*RPCMethod) ProtoMessage()    {}
func (*RPCMethod) Descriptor() ([]byte, []int) {
	return fileDescriptor_5c6ac9b241082464, []int{1}
}

func (m *RPCMethod) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RPCMethod.Unmarshal(m, b)
}
func (m *RPCMethod) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RPCMethod.Marshal(b, m, deterministic)
}
func (m *RPCMethod) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RPCMethod.Merge(m, src)
}
func (m *RPCMethod) XXX_Size() int {
	return xxx_messageInfo_RPCMethod.Size(m)
}
func (m *RPCMethod) XXX_DiscardUnknown() {
	xxx_messageInfo_RPCMethod.DiscardUnknown(m)
}

var xxx_messageInfo_RPCMethod proto.InternalMessageInfo

func (m *RPCMethod) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *RPCMethod) GetEnabled() bool {
	if m != nil {
		return m.Enabled
	}
	return false
}

func (m *RPCMethod) GetRateLimit() float64 {
	if m != nil {
		return m.RateLimit
	}
	return 0
}

func (m *RPCMethod) GetAuthRateLimit() float64 {
	if m != nil {
		return m.AuthRateLimit
	}
	return 0
}

type RPCMethods struct {
	GET                  map[string]*RPCMethod `protobuf:"bytes,1,rep,name=GET,proto3" json:"GET,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	POST                 map[string]*RPCMethod `protobuf:"bytes,2,rep,name=POST,proto3" json:"POST,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *RPCMethods) Reset()         { *m = RPCMethods{} }
func (m *RPCMethods) String() string { return proto.CompactTextString(m) }
func (*RPCMethods) ProtoMessage()    {}
func (*RPCMethods) Descriptor() ([]byte, []int) {
	return fileDescriptor_5c6ac9b241082464, []int{2}
}

func (m *RPCMethods) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RPCMethods.Unmarshal(m, b)
}
func (m *RPCMethods) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RPCMethods.Marshal(b, m, deterministic)
}
func (m *RPCMethods) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RPCMethods.Merge(m, src)
}
func (m *RPCMethods) XXX_Size() int {
	return xxx_messageInfo_RPCMethods.Size(m)
}
func (m *RPCMethods) XXX_DiscardUnknown() {
	xxx_messageInfo_RPCMethods.DiscardUnknown(m)
}

var xxx_messageInfo_RPCMethods proto.InternalMessageInfo

func (m *RPCMethods) GetGET() map[string]*RPCMethod {
	if m != nil {
		return m.GET
	}
	return nil
}

func (m *RPCMethods) GetPOST() map[string]*RPCMethod {
	if m != nil {
		return m.POST
	}
	return nil
}

// RPCMethodsRequest is the request type for the query RPC methods.
type RPCMethodsRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RPCMethodsRequest) Reset()         { *m = RPCMethodsRequest{} }
func (m *RPCMethodsRequest) String() string { return proto.CompactTextString(m) }
func (*RPCMethodsRequest) ProtoMessage()    {}
func (*RPCMethodsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_5c6ac9b241082464, []int{3}
}

func (m *RPCMethodsRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RPCMethodsRequest.Unmarshal(m, b)
}
func (m *RPCMethodsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RPCMethodsRequest.Marshal(b, m, deterministic)
}
func (m *RPCMethodsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RPCMethodsRequest.Merge(m, src)
}
func (m *RPCMethodsRequest) XXX_Size() int {
	return xxx_messageInfo_RPCMethodsRequest.Size(m)
}
func (m *RPCMethodsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RPCMethodsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RPCMethodsRequest proto.InternalMessageInfo

// RPCMethodsResponse is the response type for the query RPC methods.
type RPCMethodsResponse struct {
	ChainId              string      `protobuf:"bytes,1,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty"`
	Block                uint64      `protobuf:"varint,2,opt,name=block,proto3" json:"block,omitempty"`
	BlockTime            string      `protobuf:"bytes,3,opt,name=block_time,json=blockTime,proto3" json:"block_time,omitempty"`
	Timestamp            uint64      `protobuf:"varint,4,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Response             *RPCMethods `protobuf:"bytes,5,opt,name=response,proto3" json:"response,omitempty"`
	Error                *Error      `protobuf:"bytes,6,opt,name=error,proto3" json:"error,omitempty"`
	Signature            string      `protobuf:"bytes,7,opt,name=signature,proto3" json:"signature,omitempty"`
	Hash                 string      `protobuf:"bytes,8,opt,name=hash,proto3" json:"hash,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *RPCMethodsResponse) Reset()         { *m = RPCMethodsResponse{} }
func (m *RPCMethodsResponse) String() string { return proto.CompactTextString(m) }
func (*RPCMethodsResponse) ProtoMessage()    {}
func (*RPCMethodsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_5c6ac9b241082464, []int{4}
}

func (m *RPCMethodsResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RPCMethodsResponse.Unmarshal(m, b)
}
func (m *RPCMethodsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RPCMethodsResponse.Marshal(b, m, deterministic)
}
func (m *RPCMethodsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RPCMethodsResponse.Merge(m, src)
}
func (m *RPCMethodsResponse) XXX_Size() int {
	return xxx_messageInfo_RPCMethodsResponse.Size(m)
}
func (m *RPCMethodsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_RPCMethodsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_RPCMethodsResponse proto.InternalMessageInfo

func (m *RPCMethodsResponse) GetChainId() string {
	if m != nil {
		return m.ChainId
	}
	return ""
}

func (m *RPCMethodsResponse) GetBlock() uint64 {
	if m != nil {
		return m.Block
	}
	return 0
}

func (m *RPCMethodsResponse) GetBlockTime() string {
	if m != nil {
		return m.BlockTime
	}
	return ""
}

func (m *RPCMethodsResponse) GetTimestamp() uint64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

func (m *RPCMethodsResponse) GetResponse() *RPCMethods {
	if m != nil {
		return m.Response
	}
	return nil
}

func (m *RPCMethodsResponse) GetError() *Error {
	if m != nil {
		return m.Error
	}
	return nil
}

func (m *RPCMethodsResponse) GetSignature() string {
	if m != nil {
		return m.Signature
	}
	return ""
}

func (m *RPCMethodsResponse) GetHash() string {
	if m != nil {
		return m.Hash
	}
	return ""
}

func init() {
	proto.RegisterType((*Endpoint)(nil), "query.Endpoint")
	proto.RegisterType((*RPCMethod)(nil), "query.RPCMethod")
	proto.RegisterType((*RPCMethods)(nil), "query.RPCMethods")
	proto.RegisterMapType((map[string]*RPCMethod)(nil), "query.RPCMethods.GETEntry")
	proto.RegisterMapType((map[string]*RPCMethod)(nil), "query.RPCMethods.POSTEntry")
	proto.RegisterType((*RPCMethodsRequest)(nil), "query.RPCMethodsRequest")
	proto.RegisterType((*RPCMethodsResponse)(nil), "query.RPCMethodsResponse")
}

func init() {
	proto.RegisterFile("query.proto", fileDescriptor_5c6ac9b241082464)
}

var fileDescriptor_5c6ac9b241082464 = []byte{
	// 629 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x94, 0x4f, 0x6f, 0xd3, 0x4c,
	0x10, 0xc6, 0xe5, 0x24, 0x6e, 0xe3, 0xc9, 0xfb, 0x8a, 0x76, 0x41, 0xc8, 0x35, 0x45, 0x32, 0x39,
	0x54, 0x11, 0x22, 0x31, 0x04, 0x0e, 0x55, 0x0f, 0x48, 0x25, 0x44, 0xa5, 0xe2, 0x5f, 0x58, 0x72,
	0x00, 0x2e, 0xd1, 0x26, 0x59, 0xec, 0x55, 0x6c, 0xaf, 0xbb, 0x5e, 0x17, 0xf9, 0x86, 0x38, 0x70,
	0xe4, 0x50, 0xae, 0x9c, 0xf9, 0x42, 0x7c, 0x05, 0x3e, 0x08, 0xda, 0x5d, 0x27, 0x2d, 0xa4, 0x12,
	0x07, 0x4e, 0xde, 0x79, 0x9e, 0x5f, 0x66, 0xc6, 0xb3, 0xe3, 0x40, 0xeb, 0xa4, 0xa0, 0xa2, 0xec,
	0x65, 0x82, 0x4b, 0x8e, 0x6c, 0x1d, 0x78, 0xbb, 0x21, 0xe7, 0x61, 0x4c, 0x03, 0x92, 0xb1, 0x80,
	0xa4, 0x29, 0x97, 0x44, 0x32, 0x9e, 0xe6, 0x06, 0xf2, 0xcc, 0x63, 0xd6, 0x0d, 0x69, 0xda, 0xe5,
	0x19, 0x4d, 0x49, 0xc6, 0x4e, 0xfb, 0x01, 0xcf, 0x34, 0x73, 0x09, 0x0f, 0xb2, 0xcc, 0xa8, 0x39,
	0xb7, 0x1f, 0x40, 0x73, 0x98, 0xce, 0x33, 0xce, 0x52, 0x89, 0xb6, 0xa0, 0x5e, 0x88, 0xd8, 0xb5,
	0x7c, 0xab, 0xe3, 0x60, 0x75, 0x44, 0xd7, 0x61, 0x23, 0xa1, 0x32, 0xe2, 0x73, 0xb7, 0xa6, 0xc5,
	0x2a, 0x6a, 0x7f, 0xb1, 0xc0, 0xc1, 0xa3, 0xc1, 0x73, 0x1d, 0x21, 0x1f, 0x5a, 0x73, 0x9a, 0xcf,
	0x04, 0xd3, 0x15, 0xab, 0xdf, 0x5f, 0x94, 0x90, 0x0b, 0x9b, 0x34, 0x25, 0xd3, 0x98, 0x9a, 0x44,
	0x4d, 0xbc, 0x0c, 0xd1, 0x4d, 0x00, 0x41, 0x24, 0x9d, 0xc4, 0x2c, 0x61, 0xd2, 0xad, 0xfb, 0x56,
	0xc7, 0xc2, 0x8e, 0x52, 0x9e, 0x29, 0x01, 0xed, 0xc1, 0x15, 0x52, 0xc8, 0x68, 0x72, 0x81, 0x69,
	0x68, 0xe6, 0x7f, 0x25, 0xe3, 0x25, 0xd7, 0xfe, 0x58, 0x03, 0x58, 0x35, 0x94, 0xa3, 0x3b, 0x50,
	0x3f, 0x1a, 0x8e, 0x5d, 0xcb, 0xaf, 0x77, 0x5a, 0x7d, 0xaf, 0x67, 0x26, 0x7a, 0xee, 0xf7, 0x8e,
	0x86, 0xe3, 0x61, 0x2a, 0x45, 0x89, 0x15, 0x86, 0x02, 0x68, 0x8c, 0x5e, 0xbe, 0x1e, 0xbb, 0x35,
	0x8d, 0xdf, 0x58, 0xc7, 0x95, 0x6b, 0x78, 0x0d, 0x7a, 0x4f, 0xa0, 0xb9, 0xcc, 0xa0, 0x86, 0xb6,
	0xa0, 0xe5, 0x72, 0x68, 0x0b, 0x5a, 0xa2, 0x3d, 0xb0, 0x4f, 0x49, 0x5c, 0x50, 0xfd, 0xaa, 0xad,
	0xfe, 0xd6, 0x9f, 0xf9, 0xb0, 0xb1, 0x0f, 0x6a, 0xfb, 0x96, 0x77, 0x0c, 0xce, 0x2a, 0xf9, 0xbf,
	0xa5, 0x6a, 0x5f, 0x85, 0xed, 0xf3, 0x96, 0x31, 0x3d, 0x29, 0x68, 0x2e, 0xdb, 0x9f, 0x6b, 0x80,
	0x2e, 0xaa, 0x79, 0xc6, 0xd3, 0x9c, 0xa2, 0x1d, 0x68, 0xce, 0x22, 0xc2, 0xd2, 0x09, 0x9b, 0x57,
	0xe5, 0x36, 0x75, 0x7c, 0x3c, 0x47, 0xd7, 0xc0, 0x9e, 0xc6, 0x7c, 0xb6, 0xd0, 0x25, 0x1b, 0xd8,
	0x04, 0xea, 0x9a, 0xf4, 0x61, 0x22, 0x59, 0x42, 0xf5, 0x35, 0x39, 0xd8, 0xd1, 0xca, 0x98, 0x25,
	0x14, 0xed, 0x82, 0xa3, 0x8c, 0x5c, 0x92, 0x24, 0xd3, 0x17, 0xd4, 0xc0, 0xe7, 0x02, 0xea, 0x42,
	0x53, 0x54, 0x95, 0x5d, 0x5b, 0xbf, 0xc8, 0xf6, 0xda, 0x8c, 0xf1, 0x0a, 0x41, 0xb7, 0xc0, 0xa6,
	0x42, 0x70, 0xe1, 0x6e, 0x68, 0xb6, 0xd5, 0xd3, 0xeb, 0x3a, 0x54, 0x12, 0x36, 0x8e, 0xaa, 0x97,
	0xb3, 0x30, 0x25, 0xb2, 0x10, 0xd4, 0xdd, 0x34, 0xdd, 0xac, 0x04, 0x84, 0xa0, 0x11, 0x91, 0x3c,
	0x72, 0x9b, 0xda, 0xd0, 0xe7, 0xfe, 0x77, 0x0b, 0xec, 0x57, 0xaa, 0x26, 0xfa, 0x66, 0xfd, 0xb6,
	0x2a, 0xee, 0x7a, 0x2b, 0x66, 0x76, 0xde, 0xce, 0x25, 0x8e, 0x69, 0xb1, 0xfd, 0xf6, 0xec, 0xf0,
	0x21, 0x98, 0x6f, 0x13, 0xb5, 0xf0, 0x68, 0xe0, 0x57, 0x80, 0xd7, 0xd1, 0x35, 0x7c, 0x19, 0x51,
	0x3f, 0x66, 0xb9, 0xf4, 0xf9, 0x7b, 0x9f, 0x9c, 0x12, 0x16, 0xab, 0x65, 0xf7, 0x15, 0x98, 0x54,
	0x2b, 0xf5, 0xe9, 0xc7, 0xcf, 0xaf, 0x35, 0x84, 0xb6, 0xf4, 0x27, 0x2d, 0xb2, 0xd9, 0xa4, 0x32,
	0x1e, 0xf1, 0xb3, 0xc3, 0xc7, 0xc8, 0xee, 0xd7, 0xef, 0xf5, 0xee, 0xde, 0xb6, 0x2c, 0xb1, 0x0f,
	0xff, 0x85, 0x78, 0x34, 0xe8, 0x86, 0x44, 0xd2, 0x0f, 0xa4, 0x44, 0x9d, 0x48, 0xca, 0x2c, 0x3f,
	0x08, 0x82, 0x90, 0xc9, 0xa8, 0x98, 0xf6, 0x66, 0x3c, 0x09, 0x9e, 0x32, 0x41, 0x06, 0x5c, 0xd0,
	0x20, 0xa7, 0x0b, 0xc2, 0x82, 0xe3, 0x17, 0xe3, 0x21, 0x7e, 0xf3, 0x6e, 0xef, 0x6f, 0x44, 0xa0,
	0xff, 0x01, 0xa6, 0x1b, 0xfa, 0x71, 0xff, 0x57, 0x00, 0x00, 0x00, 0xff, 0xff, 0x97, 0x82, 0x41,
	0x72, 0x78, 0x04, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// QueryClient is the client API for Query service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type QueryClient interface {
	RPCMethods(ctx context.Context, in *RPCMethodsRequest, opts ...grpc.CallOption) (*RPCMethodsResponse, error)
}

type queryClient struct {
	cc grpc.ClientConnInterface
}

func NewQueryClient(cc grpc.ClientConnInterface) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) RPCMethods(ctx context.Context, in *RPCMethodsRequest, opts ...grpc.CallOption) (*RPCMethodsResponse, error) {
	out := new(RPCMethodsResponse)
	err := c.cc.Invoke(ctx, "/query.Query/RPCMethods", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
type QueryServer interface {
	RPCMethods(context.Context, *RPCMethodsRequest) (*RPCMethodsResponse, error)
}

// UnimplementedQueryServer can be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func (*UnimplementedQueryServer) RPCMethods(ctx context.Context, req *RPCMethodsRequest) (*RPCMethodsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RPCMethods not implemented")
}

func RegisterQueryServer(s *grpc.Server, srv QueryServer) {
	s.RegisterService(&_Query_serviceDesc, srv)
}

func _Query_RPCMethods_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RPCMethodsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).RPCMethods(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/query.Query/RPCMethods",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).RPCMethods(ctx, req.(*RPCMethodsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Query_serviceDesc = grpc.ServiceDesc{
	ServiceName: "query.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RPCMethods",
			Handler:    _Query_RPCMethods_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "query.proto",
}
