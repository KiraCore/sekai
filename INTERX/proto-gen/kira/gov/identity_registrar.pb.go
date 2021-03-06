// Code generated by protoc-gen-go. DO NOT EDIT.
// source: kira/gov/identity_registrar.proto

package proto

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
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

type IdentityRecord struct {
	Id                   uint64               `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Address              []byte               `protobuf:"bytes,2,opt,name=address,proto3" json:"address,omitempty"`
	Infos                map[string]string    `protobuf:"bytes,3,rep,name=infos,proto3" json:"infos,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Date                 *timestamp.Timestamp `protobuf:"bytes,4,opt,name=date,proto3" json:"date,omitempty"`
	Verifiers            [][]byte             `protobuf:"bytes,5,rep,name=verifiers,proto3" json:"verifiers,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *IdentityRecord) Reset()         { *m = IdentityRecord{} }
func (m *IdentityRecord) String() string { return proto.CompactTextString(m) }
func (*IdentityRecord) ProtoMessage()    {}
func (*IdentityRecord) Descriptor() ([]byte, []int) {
	return fileDescriptor_cf2821c6a2456955, []int{0}
}

func (m *IdentityRecord) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_IdentityRecord.Unmarshal(m, b)
}
func (m *IdentityRecord) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_IdentityRecord.Marshal(b, m, deterministic)
}
func (m *IdentityRecord) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IdentityRecord.Merge(m, src)
}
func (m *IdentityRecord) XXX_Size() int {
	return xxx_messageInfo_IdentityRecord.Size(m)
}
func (m *IdentityRecord) XXX_DiscardUnknown() {
	xxx_messageInfo_IdentityRecord.DiscardUnknown(m)
}

var xxx_messageInfo_IdentityRecord proto.InternalMessageInfo

func (m *IdentityRecord) GetId() uint64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *IdentityRecord) GetAddress() []byte {
	if m != nil {
		return m.Address
	}
	return nil
}

func (m *IdentityRecord) GetInfos() map[string]string {
	if m != nil {
		return m.Infos
	}
	return nil
}

func (m *IdentityRecord) GetDate() *timestamp.Timestamp {
	if m != nil {
		return m.Date
	}
	return nil
}

func (m *IdentityRecord) GetVerifiers() [][]byte {
	if m != nil {
		return m.Verifiers
	}
	return nil
}

type VerifiedIdRecordsByAddress struct {
	Address              []byte   `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	RecordIds            []uint64 `protobuf:"varint,2,rep,packed,name=recordIds,proto3" json:"recordIds,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *VerifiedIdRecordsByAddress) Reset()         { *m = VerifiedIdRecordsByAddress{} }
func (m *VerifiedIdRecordsByAddress) String() string { return proto.CompactTextString(m) }
func (*VerifiedIdRecordsByAddress) ProtoMessage()    {}
func (*VerifiedIdRecordsByAddress) Descriptor() ([]byte, []int) {
	return fileDescriptor_cf2821c6a2456955, []int{1}
}

func (m *VerifiedIdRecordsByAddress) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_VerifiedIdRecordsByAddress.Unmarshal(m, b)
}
func (m *VerifiedIdRecordsByAddress) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_VerifiedIdRecordsByAddress.Marshal(b, m, deterministic)
}
func (m *VerifiedIdRecordsByAddress) XXX_Merge(src proto.Message) {
	xxx_messageInfo_VerifiedIdRecordsByAddress.Merge(m, src)
}
func (m *VerifiedIdRecordsByAddress) XXX_Size() int {
	return xxx_messageInfo_VerifiedIdRecordsByAddress.Size(m)
}
func (m *VerifiedIdRecordsByAddress) XXX_DiscardUnknown() {
	xxx_messageInfo_VerifiedIdRecordsByAddress.DiscardUnknown(m)
}

var xxx_messageInfo_VerifiedIdRecordsByAddress proto.InternalMessageInfo

func (m *VerifiedIdRecordsByAddress) GetAddress() []byte {
	if m != nil {
		return m.Address
	}
	return nil
}

func (m *VerifiedIdRecordsByAddress) GetRecordIds() []uint64 {
	if m != nil {
		return m.RecordIds
	}
	return nil
}

type IdentityInfoEntry struct {
	Key                  string   `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Info                 string   `protobuf:"bytes,2,opt,name=info,proto3" json:"info,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *IdentityInfoEntry) Reset()         { *m = IdentityInfoEntry{} }
func (m *IdentityInfoEntry) String() string { return proto.CompactTextString(m) }
func (*IdentityInfoEntry) ProtoMessage()    {}
func (*IdentityInfoEntry) Descriptor() ([]byte, []int) {
	return fileDescriptor_cf2821c6a2456955, []int{2}
}

func (m *IdentityInfoEntry) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_IdentityInfoEntry.Unmarshal(m, b)
}
func (m *IdentityInfoEntry) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_IdentityInfoEntry.Marshal(b, m, deterministic)
}
func (m *IdentityInfoEntry) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IdentityInfoEntry.Merge(m, src)
}
func (m *IdentityInfoEntry) XXX_Size() int {
	return xxx_messageInfo_IdentityInfoEntry.Size(m)
}
func (m *IdentityInfoEntry) XXX_DiscardUnknown() {
	xxx_messageInfo_IdentityInfoEntry.DiscardUnknown(m)
}

var xxx_messageInfo_IdentityInfoEntry proto.InternalMessageInfo

func (m *IdentityInfoEntry) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *IdentityInfoEntry) GetInfo() string {
	if m != nil {
		return m.Info
	}
	return ""
}

type MsgCreateIdentityRecord struct {
	Address              []byte               `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	Infos                []*IdentityInfoEntry `protobuf:"bytes,2,rep,name=infos,proto3" json:"infos,omitempty"`
	Date                 *timestamp.Timestamp `protobuf:"bytes,3,opt,name=date,proto3" json:"date,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *MsgCreateIdentityRecord) Reset()         { *m = MsgCreateIdentityRecord{} }
func (m *MsgCreateIdentityRecord) String() string { return proto.CompactTextString(m) }
func (*MsgCreateIdentityRecord) ProtoMessage()    {}
func (*MsgCreateIdentityRecord) Descriptor() ([]byte, []int) {
	return fileDescriptor_cf2821c6a2456955, []int{3}
}

func (m *MsgCreateIdentityRecord) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MsgCreateIdentityRecord.Unmarshal(m, b)
}
func (m *MsgCreateIdentityRecord) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MsgCreateIdentityRecord.Marshal(b, m, deterministic)
}
func (m *MsgCreateIdentityRecord) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgCreateIdentityRecord.Merge(m, src)
}
func (m *MsgCreateIdentityRecord) XXX_Size() int {
	return xxx_messageInfo_MsgCreateIdentityRecord.Size(m)
}
func (m *MsgCreateIdentityRecord) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgCreateIdentityRecord.DiscardUnknown(m)
}

var xxx_messageInfo_MsgCreateIdentityRecord proto.InternalMessageInfo

func (m *MsgCreateIdentityRecord) GetAddress() []byte {
	if m != nil {
		return m.Address
	}
	return nil
}

func (m *MsgCreateIdentityRecord) GetInfos() []*IdentityInfoEntry {
	if m != nil {
		return m.Infos
	}
	return nil
}

func (m *MsgCreateIdentityRecord) GetDate() *timestamp.Timestamp {
	if m != nil {
		return m.Date
	}
	return nil
}

type MsgEditIdentityRecord struct {
	RecordId             uint64               `protobuf:"varint,1,opt,name=record_id,json=recordId,proto3" json:"record_id,omitempty"`
	Address              []byte               `protobuf:"bytes,2,opt,name=address,proto3" json:"address,omitempty"`
	Infos                []*IdentityInfoEntry `protobuf:"bytes,3,rep,name=infos,proto3" json:"infos,omitempty"`
	Date                 *timestamp.Timestamp `protobuf:"bytes,4,opt,name=date,proto3" json:"date,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *MsgEditIdentityRecord) Reset()         { *m = MsgEditIdentityRecord{} }
func (m *MsgEditIdentityRecord) String() string { return proto.CompactTextString(m) }
func (*MsgEditIdentityRecord) ProtoMessage()    {}
func (*MsgEditIdentityRecord) Descriptor() ([]byte, []int) {
	return fileDescriptor_cf2821c6a2456955, []int{4}
}

func (m *MsgEditIdentityRecord) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MsgEditIdentityRecord.Unmarshal(m, b)
}
func (m *MsgEditIdentityRecord) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MsgEditIdentityRecord.Marshal(b, m, deterministic)
}
func (m *MsgEditIdentityRecord) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgEditIdentityRecord.Merge(m, src)
}
func (m *MsgEditIdentityRecord) XXX_Size() int {
	return xxx_messageInfo_MsgEditIdentityRecord.Size(m)
}
func (m *MsgEditIdentityRecord) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgEditIdentityRecord.DiscardUnknown(m)
}

var xxx_messageInfo_MsgEditIdentityRecord proto.InternalMessageInfo

func (m *MsgEditIdentityRecord) GetRecordId() uint64 {
	if m != nil {
		return m.RecordId
	}
	return 0
}

func (m *MsgEditIdentityRecord) GetAddress() []byte {
	if m != nil {
		return m.Address
	}
	return nil
}

func (m *MsgEditIdentityRecord) GetInfos() []*IdentityInfoEntry {
	if m != nil {
		return m.Infos
	}
	return nil
}

func (m *MsgEditIdentityRecord) GetDate() *timestamp.Timestamp {
	if m != nil {
		return m.Date
	}
	return nil
}

type IdentityRecordsVerify struct {
	Id                   uint64   `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Address              []byte   `protobuf:"bytes,2,opt,name=address,proto3" json:"address,omitempty"`
	Verifier             []byte   `protobuf:"bytes,3,opt,name=verifier,proto3" json:"verifier,omitempty"`
	RecordIds            []uint64 `protobuf:"varint,4,rep,packed,name=recordIds,proto3" json:"recordIds,omitempty"`
	Tip                  string   `protobuf:"bytes,5,opt,name=tip,proto3" json:"tip,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *IdentityRecordsVerify) Reset()         { *m = IdentityRecordsVerify{} }
func (m *IdentityRecordsVerify) String() string { return proto.CompactTextString(m) }
func (*IdentityRecordsVerify) ProtoMessage()    {}
func (*IdentityRecordsVerify) Descriptor() ([]byte, []int) {
	return fileDescriptor_cf2821c6a2456955, []int{5}
}

func (m *IdentityRecordsVerify) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_IdentityRecordsVerify.Unmarshal(m, b)
}
func (m *IdentityRecordsVerify) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_IdentityRecordsVerify.Marshal(b, m, deterministic)
}
func (m *IdentityRecordsVerify) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IdentityRecordsVerify.Merge(m, src)
}
func (m *IdentityRecordsVerify) XXX_Size() int {
	return xxx_messageInfo_IdentityRecordsVerify.Size(m)
}
func (m *IdentityRecordsVerify) XXX_DiscardUnknown() {
	xxx_messageInfo_IdentityRecordsVerify.DiscardUnknown(m)
}

var xxx_messageInfo_IdentityRecordsVerify proto.InternalMessageInfo

func (m *IdentityRecordsVerify) GetId() uint64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *IdentityRecordsVerify) GetAddress() []byte {
	if m != nil {
		return m.Address
	}
	return nil
}

func (m *IdentityRecordsVerify) GetVerifier() []byte {
	if m != nil {
		return m.Verifier
	}
	return nil
}

func (m *IdentityRecordsVerify) GetRecordIds() []uint64 {
	if m != nil {
		return m.RecordIds
	}
	return nil
}

func (m *IdentityRecordsVerify) GetTip() string {
	if m != nil {
		return m.Tip
	}
	return ""
}

type MsgRequestIdentityRecordsVerify struct {
	Address              []byte   `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	Verifier             []byte   `protobuf:"bytes,2,opt,name=verifier,proto3" json:"verifier,omitempty"`
	RecordIds            []uint64 `protobuf:"varint,3,rep,packed,name=recordIds,proto3" json:"recordIds,omitempty"`
	Tip                  string   `protobuf:"bytes,4,opt,name=tip,proto3" json:"tip,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MsgRequestIdentityRecordsVerify) Reset()         { *m = MsgRequestIdentityRecordsVerify{} }
func (m *MsgRequestIdentityRecordsVerify) String() string { return proto.CompactTextString(m) }
func (*MsgRequestIdentityRecordsVerify) ProtoMessage()    {}
func (*MsgRequestIdentityRecordsVerify) Descriptor() ([]byte, []int) {
	return fileDescriptor_cf2821c6a2456955, []int{6}
}

func (m *MsgRequestIdentityRecordsVerify) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MsgRequestIdentityRecordsVerify.Unmarshal(m, b)
}
func (m *MsgRequestIdentityRecordsVerify) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MsgRequestIdentityRecordsVerify.Marshal(b, m, deterministic)
}
func (m *MsgRequestIdentityRecordsVerify) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgRequestIdentityRecordsVerify.Merge(m, src)
}
func (m *MsgRequestIdentityRecordsVerify) XXX_Size() int {
	return xxx_messageInfo_MsgRequestIdentityRecordsVerify.Size(m)
}
func (m *MsgRequestIdentityRecordsVerify) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgRequestIdentityRecordsVerify.DiscardUnknown(m)
}

var xxx_messageInfo_MsgRequestIdentityRecordsVerify proto.InternalMessageInfo

func (m *MsgRequestIdentityRecordsVerify) GetAddress() []byte {
	if m != nil {
		return m.Address
	}
	return nil
}

func (m *MsgRequestIdentityRecordsVerify) GetVerifier() []byte {
	if m != nil {
		return m.Verifier
	}
	return nil
}

func (m *MsgRequestIdentityRecordsVerify) GetRecordIds() []uint64 {
	if m != nil {
		return m.RecordIds
	}
	return nil
}

func (m *MsgRequestIdentityRecordsVerify) GetTip() string {
	if m != nil {
		return m.Tip
	}
	return ""
}

type MsgApproveIdentityRecords struct {
	Proposer             []byte   `protobuf:"bytes,1,opt,name=proposer,proto3" json:"proposer,omitempty"`
	Verifier             []byte   `protobuf:"bytes,2,opt,name=verifier,proto3" json:"verifier,omitempty"`
	VerifyRequestId      uint64   `protobuf:"varint,3,opt,name=verifyRequestId,proto3" json:"verifyRequestId,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MsgApproveIdentityRecords) Reset()         { *m = MsgApproveIdentityRecords{} }
func (m *MsgApproveIdentityRecords) String() string { return proto.CompactTextString(m) }
func (*MsgApproveIdentityRecords) ProtoMessage()    {}
func (*MsgApproveIdentityRecords) Descriptor() ([]byte, []int) {
	return fileDescriptor_cf2821c6a2456955, []int{7}
}

func (m *MsgApproveIdentityRecords) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MsgApproveIdentityRecords.Unmarshal(m, b)
}
func (m *MsgApproveIdentityRecords) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MsgApproveIdentityRecords.Marshal(b, m, deterministic)
}
func (m *MsgApproveIdentityRecords) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgApproveIdentityRecords.Merge(m, src)
}
func (m *MsgApproveIdentityRecords) XXX_Size() int {
	return xxx_messageInfo_MsgApproveIdentityRecords.Size(m)
}
func (m *MsgApproveIdentityRecords) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgApproveIdentityRecords.DiscardUnknown(m)
}

var xxx_messageInfo_MsgApproveIdentityRecords proto.InternalMessageInfo

func (m *MsgApproveIdentityRecords) GetProposer() []byte {
	if m != nil {
		return m.Proposer
	}
	return nil
}

func (m *MsgApproveIdentityRecords) GetVerifier() []byte {
	if m != nil {
		return m.Verifier
	}
	return nil
}

func (m *MsgApproveIdentityRecords) GetVerifyRequestId() uint64 {
	if m != nil {
		return m.VerifyRequestId
	}
	return 0
}

func init() {
	proto.RegisterType((*IdentityRecord)(nil), "kira.gov.IdentityRecord")
	proto.RegisterMapType((map[string]string)(nil), "kira.gov.IdentityRecord.InfosEntry")
	proto.RegisterType((*VerifiedIdRecordsByAddress)(nil), "kira.gov.VerifiedIdRecordsByAddress")
	proto.RegisterType((*IdentityInfoEntry)(nil), "kira.gov.IdentityInfoEntry")
	proto.RegisterType((*MsgCreateIdentityRecord)(nil), "kira.gov.MsgCreateIdentityRecord")
	proto.RegisterType((*MsgEditIdentityRecord)(nil), "kira.gov.MsgEditIdentityRecord")
	proto.RegisterType((*IdentityRecordsVerify)(nil), "kira.gov.IdentityRecordsVerify")
	proto.RegisterType((*MsgRequestIdentityRecordsVerify)(nil), "kira.gov.MsgRequestIdentityRecordsVerify")
	proto.RegisterType((*MsgApproveIdentityRecords)(nil), "kira.gov.MsgApproveIdentityRecords")
}

func init() {
	proto.RegisterFile("kira/gov/identity_registrar.proto", fileDescriptor_cf2821c6a2456955)
}

var fileDescriptor_cf2821c6a2456955 = []byte{
	// 652 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xbc, 0x56, 0x4d, 0x6f, 0xd3, 0x30,
	0x18, 0x26, 0x1f, 0x85, 0xd6, 0x9b, 0x36, 0x88, 0x36, 0x51, 0x3a, 0xa4, 0x94, 0x20, 0x41, 0x2e,
	0x4b, 0xc4, 0x38, 0xec, 0xe3, 0xd6, 0x94, 0x1d, 0x2a, 0x54, 0x90, 0xac, 0x09, 0x21, 0x24, 0x34,
	0xb2, 0xda, 0x33, 0x56, 0xd7, 0x3a, 0xd8, 0x6e, 0xa5, 0xfc, 0x01, 0x24, 0x6e, 0x5c, 0xf9, 0x3d,
	0x5c, 0xe0, 0x2f, 0x70, 0x28, 0x3f, 0x80, 0xdb, 0x8e, 0x5c, 0x40, 0x71, 0xe2, 0x7e, 0xc1, 0x24,
	0x3a, 0xd4, 0x9d, 0xe2, 0xf8, 0xf5, 0xfb, 0xbe, 0x79, 0x9e, 0xf7, 0x79, 0xac, 0x80, 0x7b, 0x5d,
	0xca, 0xe3, 0x90, 0xb0, 0x61, 0x48, 0x11, 0xee, 0x4b, 0x2a, 0xd3, 0x63, 0x8e, 0x09, 0x15, 0x92,
	0xc7, 0x3c, 0x48, 0x38, 0x93, 0xcc, 0x29, 0x67, 0x47, 0x02, 0xc2, 0x86, 0xb5, 0x0d, 0xc2, 0x08,
	0x53, 0x9b, 0x61, 0xb6, 0xca, 0xe3, 0x35, 0x97, 0x30, 0x46, 0xce, 0x70, 0xa8, 0xde, 0x4e, 0x06,
	0xa7, 0xa1, 0xa4, 0x3d, 0x2c, 0x64, 0xdc, 0x4b, 0xf2, 0x03, 0xde, 0x2f, 0x13, 0xac, 0xb5, 0x8a,
	0xea, 0x10, 0x77, 0x18, 0x47, 0xce, 0x1a, 0x30, 0x29, 0xaa, 0x1a, 0x75, 0xc3, 0xb7, 0xa1, 0x49,
	0x91, 0xf3, 0x1a, 0xdc, 0x88, 0x11, 0xe2, 0x58, 0x88, 0xaa, 0x59, 0x37, 0xfc, 0xd5, 0xa8, 0x79,
	0x3e, 0x72, 0xd7, 0xd2, 0xb8, 0x77, 0x76, 0xe0, 0x15, 0x01, 0xef, 0xe7, 0xc8, 0xdd, 0x26, 0x54,
	0xbe, 0x1d, 0x9c, 0x04, 0x1d, 0xd6, 0x0b, 0x3b, 0x4c, 0xf4, 0x98, 0x28, 0x1e, 0xdb, 0x02, 0x75,
	0x43, 0x99, 0x26, 0x58, 0x04, 0x8d, 0x4e, 0xa7, 0x91, 0x67, 0x40, 0x5d, 0xd3, 0xd9, 0x07, 0x25,
	0xda, 0x3f, 0x65, 0xa2, 0x6a, 0xd5, 0x2d, 0x7f, 0x65, 0xe7, 0x7e, 0xa0, 0x21, 0x05, 0xb3, 0xdf,
	0x15, 0xb4, 0xb2, 0x53, 0x87, 0x7d, 0xc9, 0x53, 0x98, 0x67, 0x38, 0x7b, 0xc0, 0x46, 0xb1, 0xc4,
	0x55, 0xbb, 0x6e, 0xf8, 0x2b, 0x3b, 0xb5, 0x20, 0x07, 0x1b, 0x68, 0xb0, 0xc1, 0x91, 0x06, 0x1b,
	0x95, 0xbf, 0x8c, 0xdc, 0x6b, 0x1f, 0xbf, 0xbb, 0x06, 0x54, 0x19, 0xce, 0x73, 0x50, 0x19, 0x62,
	0x4e, 0x4f, 0x29, 0xe6, 0xa2, 0x5a, 0xaa, 0x5b, 0xfe, 0x6a, 0xf4, 0x68, 0x71, 0x0c, 0x93, 0x1a,
	0xb5, 0x3d, 0x00, 0x26, 0xdf, 0xe7, 0xdc, 0x04, 0x56, 0x17, 0xa7, 0x8a, 0xc3, 0x0a, 0xcc, 0x96,
	0xce, 0x06, 0x28, 0x0d, 0xe3, 0xb3, 0x01, 0x56, 0x14, 0x56, 0x60, 0xfe, 0x72, 0x60, 0xee, 0x19,
	0xde, 0x27, 0x03, 0xd4, 0x5e, 0xe4, 0x75, 0x50, 0x0b, 0xe5, 0x58, 0x45, 0x94, 0x16, 0x3d, 0xa6,
	0xd9, 0x37, 0x96, 0xc0, 0xfe, 0x5d, 0x50, 0xe1, 0xaa, 0x65, 0x0b, 0x65, 0xe3, 0xb5, 0x7c, 0x1b,
	0x4e, 0x36, 0xbc, 0x7d, 0x70, 0x4b, 0x0f, 0x21, 0x43, 0x77, 0x11, 0x38, 0x07, 0xd8, 0xd9, 0x40,
	0x0a, 0x6c, 0x6a, 0xed, 0xfd, 0x30, 0xc0, 0xed, 0xb6, 0x20, 0x4d, 0x8e, 0x63, 0x89, 0xe7, 0x14,
	0xb6, 0x64, 0x4c, 0xbb, 0x5a, 0x51, 0xa6, 0x52, 0xd4, 0xd6, 0x9f, 0x8a, 0x1a, 0x83, 0x89, 0xec,
	0x4c, 0x18, 0xf3, 0x7a, 0xb2, 0x16, 0xd5, 0x93, 0xf7, 0xde, 0x04, 0x9b, 0x6d, 0x41, 0x0e, 0x11,
	0x95, 0x73, 0x58, 0xb7, 0x34, 0xc1, 0xc7, 0x63, 0x53, 0x95, 0x35, 0xc1, 0xcb, 0xb6, 0xd6, 0xee,
	0xac, 0xb5, 0x16, 0x27, 0x62, 0x61, 0x63, 0x79, 0x5f, 0x4d, 0xb0, 0x39, 0xcb, 0x80, 0x50, 0xe2,
	0x4e, 0xaf, 0xfa, 0x5a, 0x79, 0x03, 0xca, 0xda, 0x9d, 0x6a, 0x9e, 0xab, 0xd1, 0x93, 0xf3, 0x91,
	0xbb, 0x9e, 0xd7, 0xd7, 0x91, 0x4b, 0x34, 0x18, 0x57, 0x9d, 0xb5, 0x8e, 0x3d, 0x67, 0x1d, 0xa7,
	0x01, 0x2c, 0x49, 0x93, 0x6a, 0x29, 0xb3, 0x44, 0x14, 0x66, 0x2c, 0x7d, 0x1b, 0xb9, 0x0f, 0xff,
	0xa1, 0x57, 0x93, 0xd1, 0x3e, 0xcc, 0x72, 0xbd, 0xcf, 0x26, 0x70, 0xdb, 0x82, 0x40, 0xfc, 0x6e,
	0x80, 0x85, 0xfc, 0x3b, 0xab, 0x4b, 0xb6, 0xd2, 0x34, 0x8b, 0xe6, 0xf2, 0x59, 0xb4, 0x2e, 0x60,
	0xd1, 0xfe, 0x0f, 0x16, 0x3f, 0x98, 0xe0, 0x4e, 0x5b, 0x90, 0x46, 0x92, 0x70, 0x36, 0x9c, 0xbb,
	0x89, 0x14, 0xc0, 0x84, 0xb3, 0x84, 0x09, 0xcc, 0x0b, 0x02, 0xa7, 0x00, 0xea, 0xc8, 0x65, 0x00,
	0xea, 0xdc, 0x2b, 0xa0, 0xd0, 0x07, 0xeb, 0x6a, 0x9d, 0x8e, 0x95, 0xa2, 0x14, 0x6f, 0xc3, 0xf9,
	0xed, 0xc8, 0x7f, 0xf5, 0x60, 0xaa, 0xc9, 0x53, 0xca, 0xe3, 0x26, 0xe3, 0x38, 0x14, 0xb8, 0x1b,
	0xd3, 0xb0, 0xf5, 0xec, 0xe8, 0x10, 0xbe, 0x2c, 0x7e, 0x15, 0xae, 0xab, 0xc7, 0xe3, 0xdf, 0x01,
	0x00, 0x00, 0xff, 0xff, 0xa6, 0xf7, 0xf8, 0x7b, 0x84, 0x08, 0x00, 0x00,
}
