package types

import (
	fmt "fmt"
	io "io"
	math "math"
	math_bits "math/bits"

	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
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

type NetworkPropertiesV0228 struct {
	MinTxFee                    uint64 `protobuf:"varint,1,opt,name=min_tx_fee,json=minTxFee,proto3" json:"min_tx_fee,omitempty"`
	MaxTxFee                    uint64 `protobuf:"varint,2,opt,name=max_tx_fee,json=maxTxFee,proto3" json:"max_tx_fee,omitempty"`
	VoteQuorum                  uint64 `protobuf:"varint,3,opt,name=vote_quorum,json=voteQuorum,proto3" json:"vote_quorum,omitempty"`
	DefaultProposalEndTime      uint64 `protobuf:"varint,4,opt,name=default_proposal_end_time,json=defaultProposalEndTime,proto3" json:"default_proposal_end_time,omitempty"`
	MinimumProposalEndTime      uint64 `protobuf:"varint,5,opt,name=minimum_proposal_end_time,json=minimumProposalEndTime,proto3" json:"minimum_proposal_end_time,omitempty"`
	ProposalEnactmentTime       uint64 `protobuf:"varint,6,opt,name=proposal_enactment_time,json=proposalEnactmentTime,proto3" json:"proposal_enactment_time,omitempty"`
	MinProposalEndBlocks        uint64 `protobuf:"varint,7,opt,name=min_proposal_end_blocks,json=minProposalEndBlocks,proto3" json:"min_proposal_end_blocks,omitempty"`
	MinProposalEnactmentBlocks  uint64 `protobuf:"varint,8,opt,name=min_proposal_enactment_blocks,json=minProposalEnactmentBlocks,proto3" json:"min_proposal_enactment_blocks,omitempty"`
	EnableForeignFeePayments    bool   `protobuf:"varint,9,opt,name=enable_foreign_fee_payments,json=enableForeignFeePayments,proto3" json:"enable_foreign_fee_payments,omitempty"`
	MischanceRankDecreaseAmount uint64 `protobuf:"varint,10,opt,name=mischance_rank_decrease_amount,json=mischanceRankDecreaseAmount,proto3" json:"mischance_rank_decrease_amount,omitempty"`
	MaxMischance                uint64 `protobuf:"varint,11,opt,name=max_mischance,json=maxMischance,proto3" json:"max_mischance,omitempty"`
	MischanceConfidence         uint64 `protobuf:"varint,12,opt,name=mischance_confidence,json=mischanceConfidence,proto3" json:"mischance_confidence,omitempty"`
	InactiveRankDecreasePercent uint64 `protobuf:"varint,13,opt,name=inactive_rank_decrease_percent,json=inactiveRankDecreasePercent,proto3" json:"inactive_rank_decrease_percent,omitempty"`
	MinValidators               uint64 `protobuf:"varint,14,opt,name=min_validators,json=minValidators,proto3" json:"min_validators,omitempty"`
	PoorNetworkMaxBankSend      uint64 `protobuf:"varint,15,opt,name=poor_network_max_bank_send,json=poorNetworkMaxBankSend,proto3" json:"poor_network_max_bank_send,omitempty"`
	JailMaxTime                 uint64 `protobuf:"varint,16,opt,name=jail_max_time,json=jailMaxTime,proto3" json:"jail_max_time,omitempty"`
	EnableTokenWhitelist        bool   `protobuf:"varint,17,opt,name=enable_token_whitelist,json=enableTokenWhitelist,proto3" json:"enable_token_whitelist,omitempty"`
	EnableTokenBlacklist        bool   `protobuf:"varint,18,opt,name=enable_token_blacklist,json=enableTokenBlacklist,proto3" json:"enable_token_blacklist,omitempty"`
	MinIdentityApprovalTip      uint64 `protobuf:"varint,19,opt,name=min_identity_approval_tip,json=minIdentityApprovalTip,proto3" json:"min_identity_approval_tip,omitempty"`
	UniqueIdentityKeys          string `protobuf:"bytes,20,opt,name=unique_identity_keys,json=uniqueIdentityKeys,proto3" json:"unique_identity_keys,omitempty"`
}

func (m *NetworkPropertiesV0228) Reset()         { *m = NetworkPropertiesV0228{} }
func (m *NetworkPropertiesV0228) String() string { return proto.CompactTextString(m) }
func (*NetworkPropertiesV0228) ProtoMessage()    {}
func (*NetworkPropertiesV0228) Descriptor() ([]byte, []int) {
	return fileDescriptor_98011a6048e5dde3, []int{2}
}
func (m *NetworkPropertiesV0228) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *NetworkPropertiesV0228) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_NetworkProperties.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *NetworkPropertiesV0228) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NetworkProperties.Merge(m, src)
}
func (m *NetworkPropertiesV0228) XXX_Size() int {
	return m.Size()
}
func (m *NetworkPropertiesV0228) XXX_DiscardUnknown() {
	xxx_messageInfo_NetworkProperties.DiscardUnknown(m)
}

var xxx_messageInfo_NetworkProperties proto.InternalMessageInfo

func (m *NetworkPropertiesV0228) GetMinTxFee() uint64 {
	if m != nil {
		return m.MinTxFee
	}
	return 0
}

func (m *NetworkPropertiesV0228) GetMaxTxFee() uint64 {
	if m != nil {
		return m.MaxTxFee
	}
	return 0
}

func (m *NetworkPropertiesV0228) GetVoteQuorum() uint64 {
	if m != nil {
		return m.VoteQuorum
	}
	return 0
}

func (m *NetworkPropertiesV0228) GetDefaultProposalEndTime() uint64 {
	if m != nil {
		return m.DefaultProposalEndTime
	}
	return 0
}

func (m *NetworkPropertiesV0228) GetMinimumProposalEndTime() uint64 {
	if m != nil {
		return m.MinimumProposalEndTime
	}
	return 0
}

func (m *NetworkPropertiesV0228) GetProposalEnactmentTime() uint64 {
	if m != nil {
		return m.ProposalEnactmentTime
	}
	return 0
}

func (m *NetworkPropertiesV0228) GetMinProposalEndBlocks() uint64 {
	if m != nil {
		return m.MinProposalEndBlocks
	}
	return 0
}

func (m *NetworkPropertiesV0228) GetMinProposalEnactmentBlocks() uint64 {
	if m != nil {
		return m.MinProposalEnactmentBlocks
	}
	return 0
}

func (m *NetworkPropertiesV0228) GetEnableForeignFeePayments() bool {
	if m != nil {
		return m.EnableForeignFeePayments
	}
	return false
}

func (m *NetworkPropertiesV0228) GetMischanceRankDecreaseAmount() uint64 {
	if m != nil {
		return m.MischanceRankDecreaseAmount
	}
	return 0
}

func (m *NetworkPropertiesV0228) GetMaxMischance() uint64 {
	if m != nil {
		return m.MaxMischance
	}
	return 0
}

func (m *NetworkPropertiesV0228) GetMischanceConfidence() uint64 {
	if m != nil {
		return m.MischanceConfidence
	}
	return 0
}

func (m *NetworkPropertiesV0228) GetInactiveRankDecreasePercent() uint64 {
	if m != nil {
		return m.InactiveRankDecreasePercent
	}
	return 0
}

func (m *NetworkPropertiesV0228) GetMinValidators() uint64 {
	if m != nil {
		return m.MinValidators
	}
	return 0
}

func (m *NetworkPropertiesV0228) GetPoorNetworkMaxBankSend() uint64 {
	if m != nil {
		return m.PoorNetworkMaxBankSend
	}
	return 0
}

func (m *NetworkPropertiesV0228) GetJailMaxTime() uint64 {
	if m != nil {
		return m.JailMaxTime
	}
	return 0
}

func (m *NetworkPropertiesV0228) GetEnableTokenWhitelist() bool {
	if m != nil {
		return m.EnableTokenWhitelist
	}
	return false
}

func (m *NetworkPropertiesV0228) GetEnableTokenBlacklist() bool {
	if m != nil {
		return m.EnableTokenBlacklist
	}
	return false
}

func (m *NetworkPropertiesV0228) GetMinIdentityApprovalTip() uint64 {
	if m != nil {
		return m.MinIdentityApprovalTip
	}
	return 0
}

func (m *NetworkPropertiesV0228) GetUniqueIdentityKeys() string {
	if m != nil {
		return m.UniqueIdentityKeys
	}
	return ""
}

func init() {
	proto.RegisterType((*NetworkPropertiesV0228)(nil), "kira.gov.NetworkPropertiesV0228")
}

func init() { proto.RegisterFile("kira/gov/network_properties.proto", fileDescriptor_98011a6048e5dde3) }

var fileDescriptor_98011a6048e5dde3 = []byte{
	// 1240 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x96, 0xcb, 0x6e, 0xdb, 0xc6,
	0x17, 0xc6, 0xad, 0xc4, 0x49, 0xe4, 0xb1, 0x65, 0xd3, 0xb4, 0x62, 0x33, 0x74, 0xfe, 0x32, 0xff,
	0x2e, 0x02, 0x04, 0x05, 0x62, 0x35, 0xbd, 0x01, 0x0d, 0x50, 0x14, 0x94, 0x34, 0x6a, 0x18, 0x89,
	0x97, 0x50, 0x94, 0x92, 0x74, 0x33, 0x18, 0x4b, 0x63, 0x65, 0x2a, 0x91, 0x54, 0x48, 0x4a, 0x91,
	0xdf, 0xa0, 0xd0, 0xaa, 0xdb, 0x2e, 0x08, 0x14, 0xe8, 0x2b, 0xf4, 0x21, 0xba, 0xcc, 0xb2, 0xab,
	0xa2, 0x88, 0x37, 0x7d, 0x86, 0xae, 0x8a, 0x19, 0x52, 0xb2, 0x2e, 0xb6, 0x57, 0x96, 0xf9, 0x7d,
	0xbf, 0x73, 0x0e, 0xe7, 0xf2, 0x81, 0xe0, 0xff, 0x3d, 0x1a, 0xe0, 0x62, 0xd7, 0x1f, 0x15, 0x3d,
	0x12, 0xbd, 0xf7, 0x83, 0x1e, 0x1a, 0x04, 0xfe, 0x80, 0x04, 0x11, 0x25, 0xe1, 0xc9, 0x20, 0xf0,
	0x23, 0x5f, 0xcc, 0x32, 0xcb, 0x49, 0xd7, 0x1f, 0xc9, 0xf9, 0xae, 0xdf, 0xf5, 0xf9, 0xc3, 0x22,
	0xfb, 0x95, 0xe8, 0xc7, 0xbf, 0x67, 0xc0, 0x81, 0x1e, 0x76, 0x1b, 0x24, 0x32, 0x92, 0x12, 0xd6,
	0xac, 0x82, 0xf8, 0x02, 0x88, 0xab, 0x75, 0xa5, 0x8c, 0x92, 0x79, 0xbc, 0xf9, 0xf9, 0xe1, 0xc9,
	0xb4, 0xf0, 0xc9, 0x0a, 0x68, 0xef, 0x7a, 0x2b, 0xb5, 0x74, 0x90, 0x65, 0x35, 0xfc, 0x90, 0x04,
	0xd2, 0x2d, 0x25, 0xf3, 0x78, 0xab, 0xf4, 0xf4, 0xdf, 0xbf, 0x8e, 0x9e, 0x74, 0x69, 0xf4, 0x76,
	0x78, 0x7a, 0xd2, 0xf6, 0xdd, 0x62, 0xdb, 0x0f, 0x5d, 0x3f, 0x4c, 0xff, 0x3c, 0x09, 0x3b, 0xbd,
	0x62, 0x74, 0x3e, 0x20, 0xe1, 0x89, 0xda, 0x6e, 0xab, 0x9d, 0x4e, 0x40, 0xc2, 0xd0, 0x9e, 0x95,
	0x38, 0x36, 0x41, 0x7e, 0xb1, 0xed, 0x79, 0x0b, 0xf7, 0x87, 0x44, 0xcc, 0x83, 0x3b, 0x23, 0xf6,
	0x83, 0x4f, 0xb9, 0x6e, 0x27, 0xff, 0x88, 0x87, 0x60, 0x23, 0x8c, 0x02, 0x94, 0x28, 0xac, 0xfb,
	0x86, 0x9d, 0x0d, 0xa3, 0x80, 0x23, 0xcf, 0xd6, 0xff, 0xf9, 0xf5, 0x28, 0x73, 0xfc, 0x4b, 0x16,
	0xec, 0xae, 0xae, 0xc0, 0x43, 0x00, 0x5c, 0xea, 0xa1, 0x68, 0x8c, 0xce, 0xc8, 0xb4, 0x66, 0xd6,
	0xa5, 0x9e, 0x33, 0xae, 0x12, 0xc2, 0x55, 0x3c, 0x9e, 0xaa, 0xb7, 0x52, 0x15, 0x8f, 0x13, 0xf5,
	0x08, 0x6c, 0x8e, 0xfc, 0x88, 0xa0, 0x77, 0x43, 0x3f, 0x18, 0xba, 0xd2, 0x6d, 0x2e, 0x03, 0xf6,
	0xe8, 0x25, 0x7f, 0x22, 0x7e, 0x03, 0x1e, 0x74, 0xc8, 0x19, 0x1e, 0xf6, 0x23, 0x94, 0xbc, 0x17,
	0xee, 0x23, 0xe2, 0x75, 0x50, 0x44, 0x5d, 0x22, 0xad, 0x73, 0xfb, 0x7e, 0x6a, 0xb0, 0x52, 0x1d,
	0x7a, 0x1d, 0x87, 0xba, 0x84, 0xa1, 0x2e, 0xf5, 0xa8, 0x3b, 0x74, 0xaf, 0x40, 0xef, 0x24, 0x68,
	0x6a, 0x58, 0x46, 0xbf, 0x06, 0x07, 0x73, 0x08, 0x6e, 0x47, 0x2e, 0xf1, 0xa2, 0x04, 0xbc, 0xcb,
	0xc1, 0xfb, 0x83, 0x19, 0x91, 0xaa, 0x9c, 0xfb, 0x0a, 0x1c, 0xb0, 0xa5, 0x58, 0x68, 0x77, 0xda,
	0xf7, 0xdb, 0xbd, 0x50, 0xba, 0xc7, 0xb9, 0xbc, 0x4b, 0xbd, 0xb9, 0x66, 0x25, 0xae, 0x89, 0x2a,
	0xf8, 0xdf, 0x12, 0x36, 0x6d, 0x99, 0xc2, 0x59, 0x0e, 0xcb, 0x0b, 0x70, 0x6a, 0x49, 0x4b, 0x7c,
	0x0b, 0x0e, 0x89, 0x87, 0x4f, 0xfb, 0x04, 0x9d, 0xf9, 0x01, 0xa1, 0x5d, 0x8f, 0x2d, 0x37, 0x1a,
	0xe0, 0x73, 0xe6, 0x09, 0xa5, 0x0d, 0x25, 0xf3, 0x38, 0x6b, 0x4b, 0x89, 0xa5, 0x9a, 0x38, 0xaa,
	0x84, 0x58, 0xa9, 0x2e, 0x96, 0x41, 0xc1, 0xa5, 0x61, 0xfb, 0x2d, 0xf6, 0xda, 0x04, 0x05, 0xd8,
	0xeb, 0xa1, 0x0e, 0x69, 0x07, 0x04, 0x87, 0x04, 0x61, 0xd7, 0x1f, 0x7a, 0x91, 0x04, 0xf8, 0x08,
	0x87, 0x33, 0x97, 0x8d, 0xbd, 0x5e, 0x25, 0xf5, 0xa8, 0xdc, 0x22, 0x7e, 0x02, 0x72, 0x6c, 0xab,
	0x67, 0x16, 0x69, 0x93, 0x33, 0x5b, 0x2e, 0x1e, 0xeb, 0xd3, 0x67, 0xe2, 0x53, 0x90, 0xbf, 0xec,
	0xd4, 0xf6, 0xbd, 0x33, 0xda, 0x21, 0xcc, 0xbb, 0xc5, 0xbd, 0x7b, 0x33, 0xad, 0x3c, 0x93, 0xd8,
	0x70, 0x94, 0xbd, 0x2e, 0x1d, 0x2d, 0xcf, 0x36, 0x20, 0x41, 0x9b, 0x78, 0x91, 0x94, 0x4b, 0x86,
	0x9b, 0xba, 0xe6, 0x67, 0xb3, 0x12, 0x8b, 0xf8, 0x08, 0x6c, 0xb3, 0x35, 0x1e, 0xe1, 0x3e, 0xed,
	0xe0, 0xc8, 0x0f, 0x42, 0x69, 0x9b, 0x43, 0x39, 0x97, 0x7a, 0xad, 0xd9, 0x43, 0xf1, 0x19, 0x90,
	0x07, 0xbe, 0x1f, 0xa0, 0xe9, 0x9d, 0x66, 0x2f, 0x74, 0xca, 0x7a, 0x86, 0xc4, 0xeb, 0x48, 0x3b,
	0xc9, 0xa9, 0x61, 0x8e, 0xf4, 0x1e, 0xe8, 0x78, 0x5c, 0xc2, 0x5e, 0xaf, 0x41, 0xbc, 0x8e, 0x78,
	0x0c, 0x72, 0x3f, 0x62, 0xda, 0xe7, 0x0c, 0x3f, 0x2b, 0x02, 0xb7, 0x6f, 0xb2, 0x87, 0x3a, 0x1e,
	0xf3, 0x13, 0xf2, 0x25, 0xd8, 0x4f, 0xf7, 0x29, 0xf2, 0x7b, 0xc4, 0x43, 0xef, 0xdf, 0xd2, 0x88,
	0xf4, 0x69, 0x18, 0x49, 0xbb, 0x7c, 0x8b, 0xf2, 0x89, 0xea, 0x30, 0xf1, 0xd5, 0x54, 0x5b, 0xa1,
	0x4e, 0xfb, 0xb8, 0xdd, 0xe3, 0x94, 0xb8, 0x42, 0x95, 0xa6, 0x5a, 0x7a, 0x01, 0x10, 0x5b, 0xc5,
	0x88, 0x46, 0xe7, 0x08, 0x0f, 0x06, 0x81, 0x3f, 0xc2, 0x7d, 0x14, 0xd1, 0x81, 0xb4, 0x37, 0xbb,
	0x00, 0x5a, 0xaa, 0xab, 0xa9, 0xec, 0xd0, 0x81, 0xf8, 0x19, 0xc8, 0x0f, 0x3d, 0xfa, 0x6e, 0x48,
	0x2e, 0xe9, 0x1e, 0x39, 0x0f, 0xa5, 0x3c, 0xcf, 0x05, 0x31, 0xd1, 0xa6, 0x60, 0x8d, 0x9c, 0x87,
	0x9f, 0xc6, 0x59, 0xb0, 0xb3, 0x94, 0x36, 0xec, 0xee, 0xeb, 0x9a, 0x81, 0x9c, 0xd7, 0xa8, 0x0a,
	0xa1, 0xb0, 0x26, 0x6f, 0x4d, 0x62, 0x25, 0xab, 0xcf, 0x25, 0x83, 0xae, 0xbe, 0x9e, 0xaa, 0x99,
	0x54, 0x9d, 0x4b, 0x86, 0x96, 0xe9, 0x40, 0xf4, 0xb2, 0x69, 0xda, 0x4d, 0x5d, 0xb8, 0x25, 0x6f,
	0x4f, 0x62, 0x05, 0xb4, 0x16, 0x92, 0xa1, 0x02, 0xab, 0x6a, 0xb3, 0xee, 0x20, 0xcb, 0x36, 0x2d,
	0xb3, 0xa1, 0xd6, 0x11, 0x34, 0x2a, 0xc8, 0xd1, 0x74, 0x28, 0xdc, 0x96, 0xe5, 0x49, 0xac, 0xec,
	0x57, 0xae, 0x4d, 0x06, 0x5d, 0x33, 0x34, 0xbd, 0xa9, 0x5f, 0x81, 0xae, 0x27, 0xa8, 0x7e, 0x6d,
	0x32, 0xcc, 0x21, 0x6a, 0xd9, 0xd1, 0xa1, 0xe1, 0x24, 0xe0, 0x1d, 0xf9, 0xc1, 0x24, 0x56, 0xee,
	0x5b, 0xd7, 0x25, 0x03, 0x5b, 0x8a, 0x85, 0x76, 0xa5, 0xba, 0x59, 0xae, 0x35, 0x84, 0xbb, 0xb2,
	0x34, 0x89, 0x95, 0xbc, 0x7e, 0x4d, 0x32, 0x2c, 0x61, 0xd3, 0x96, 0x29, 0x7c, 0x4f, 0x2e, 0x4c,
	0x62, 0x45, 0xd6, 0x6f, 0x4c, 0x06, 0x68, 0xa8, 0xa5, 0x3a, 0x44, 0x55, 0xd3, 0x86, 0xda, 0xf7,
	0x06, 0x5b, 0x6e, 0x64, 0xa9, 0x6f, 0x58, 0x99, 0x86, 0x90, 0x95, 0x1f, 0x4e, 0x62, 0x45, 0x82,
	0x37, 0x24, 0x83, 0xae, 0x35, 0xca, 0xcf, 0x55, 0xa3, 0x0c, 0x91, 0xad, 0x1a, 0x35, 0x54, 0x81,
	0x65, 0x1b, 0xaa, 0x0d, 0x88, 0x54, 0xdd, 0x6c, 0x1a, 0x8e, 0xb0, 0x21, 0x1f, 0x4d, 0x62, 0xe5,
	0x50, 0xbf, 0x39, 0x19, 0xd8, 0x56, 0xcf, 0x0a, 0x09, 0x40, 0x16, 0x26, 0xb1, 0xb2, 0xa5, 0x2f,
	0x25, 0xc3, 0x65, 0xa7, 0xb2, 0x69, 0x54, 0xb5, 0x0a, 0x64, 0xde, 0x4d, 0xf9, 0x60, 0x12, 0x2b,
	0x7b, 0xfa, 0xd5, 0xc9, 0xa0, 0xb1, 0x15, 0xd1, 0x5a, 0xcb, 0xb3, 0x59, 0xd0, 0x2e, 0x43, 0xc3,
	0x11, 0xb6, 0x92, 0xe1, 0xb4, 0x1b, 0x92, 0xe1, 0x19, 0x90, 0x2d, 0xd3, 0xb4, 0x91, 0x01, 0x9d,
	0x57, 0xa6, 0x5d, 0x43, 0x6c, 0xd2, 0x12, 0x2b, 0xd6, 0x80, 0x46, 0x45, 0xc8, 0x25, 0xc7, 0xc1,
	0xba, 0xfa, 0xca, 0x3f, 0x02, 0xdb, 0x6c, 0x7f, 0x5a, 0x6a, 0x5d, 0xab, 0xa8, 0x8e, 0x69, 0x37,
	0x84, 0x6d, 0x79, 0x77, 0x12, 0x2b, 0x39, 0x7d, 0x21, 0x55, 0x8e, 0x41, 0xee, 0x85, 0xaa, 0xd5,
	0x79, 0x69, 0x7e, 0x56, 0x76, 0xe4, 0x9d, 0x49, 0xac, 0x6c, 0xbe, 0x58, 0x4c, 0x86, 0x74, 0x9f,
	0x1c, 0xb3, 0x06, 0x0d, 0xf4, 0xea, 0xb9, 0xe6, 0xc0, 0xba, 0xd6, 0x70, 0x04, 0x21, 0x39, 0x20,
	0xf0, 0x9a, 0x64, 0x58, 0xa0, 0x4a, 0x75, 0xb5, 0x5c, 0xe3, 0xd4, 0xee, 0x0a, 0xb5, 0x90, 0x0c,
	0x6c, 0x6c, 0xb6, 0xc0, 0x8e, 0xe6, 0xbc, 0x41, 0xaa, 0x65, 0xd9, 0x66, 0x4b, 0xad, 0x23, 0x47,
	0xb3, 0x04, 0x71, 0x76, 0x01, 0xae, 0x49, 0x86, 0xa6, 0xa1, 0xbd, 0x6c, 0xc2, 0x4b, 0xba, 0x06,
	0xdf, 0x34, 0x84, 0x3d, 0x79, 0x7f, 0x12, 0x2b, 0x62, 0x73, 0x25, 0x19, 0xe4, 0xf5, 0x9f, 0x7e,
	0x2b, 0xac, 0x95, 0xbe, 0xfb, 0xe3, 0x63, 0x21, 0xf3, 0xe1, 0x63, 0x21, 0xf3, 0xf7, 0xc7, 0x42,
	0xe6, 0xe7, 0x8b, 0xc2, 0xda, 0x87, 0x8b, 0xc2, 0xda, 0x9f, 0x17, 0x85, 0xb5, 0x1f, 0x1e, 0xcd,
	0x7d, 0xdf, 0xd4, 0x68, 0x80, 0xcb, 0x7e, 0x40, 0x8a, 0x21, 0xe9, 0x61, 0x5a, 0x1c, 0xf3, 0xef,
	0x36, 0xfe, 0x89, 0x73, 0x7a, 0x97, 0x7f, 0x8b, 0x7d, 0xf1, 0x5f, 0x00, 0x00, 0x00, 0xff, 0xff,
	0x2f, 0xed, 0x0f, 0xd1, 0xd0, 0x09, 0x00, 0x00,
}

func (m *NetworkPropertiesV0228) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *NetworkPropertiesV0228) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *NetworkPropertiesV0228) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.UniqueIdentityKeys) > 0 {
		i -= len(m.UniqueIdentityKeys)
		copy(dAtA[i:], m.UniqueIdentityKeys)
		i = encodeVarintNetworkProperties(dAtA, i, uint64(len(m.UniqueIdentityKeys)))
		i--
		dAtA[i] = 0x1
		i--
		dAtA[i] = 0xa2
	}
	if m.MinIdentityApprovalTip != 0 {
		i = encodeVarintNetworkProperties(dAtA, i, uint64(m.MinIdentityApprovalTip))
		i--
		dAtA[i] = 0x1
		i--
		dAtA[i] = 0x98
	}
	if m.EnableTokenBlacklist {
		i--
		if m.EnableTokenBlacklist {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x1
		i--
		dAtA[i] = 0x90
	}
	if m.EnableTokenWhitelist {
		i--
		if m.EnableTokenWhitelist {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x1
		i--
		dAtA[i] = 0x88
	}
	if m.JailMaxTime != 0 {
		i = encodeVarintNetworkProperties(dAtA, i, uint64(m.JailMaxTime))
		i--
		dAtA[i] = 0x1
		i--
		dAtA[i] = 0x80
	}
	if m.PoorNetworkMaxBankSend != 0 {
		i = encodeVarintNetworkProperties(dAtA, i, uint64(m.PoorNetworkMaxBankSend))
		i--
		dAtA[i] = 0x78
	}
	if m.MinValidators != 0 {
		i = encodeVarintNetworkProperties(dAtA, i, uint64(m.MinValidators))
		i--
		dAtA[i] = 0x70
	}
	if m.InactiveRankDecreasePercent != 0 {
		i = encodeVarintNetworkProperties(dAtA, i, uint64(m.InactiveRankDecreasePercent))
		i--
		dAtA[i] = 0x68
	}
	if m.MischanceConfidence != 0 {
		i = encodeVarintNetworkProperties(dAtA, i, uint64(m.MischanceConfidence))
		i--
		dAtA[i] = 0x60
	}
	if m.MaxMischance != 0 {
		i = encodeVarintNetworkProperties(dAtA, i, uint64(m.MaxMischance))
		i--
		dAtA[i] = 0x58
	}
	if m.MischanceRankDecreaseAmount != 0 {
		i = encodeVarintNetworkProperties(dAtA, i, uint64(m.MischanceRankDecreaseAmount))
		i--
		dAtA[i] = 0x50
	}
	if m.EnableForeignFeePayments {
		i--
		if m.EnableForeignFeePayments {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x48
	}
	if m.MinProposalEnactmentBlocks != 0 {
		i = encodeVarintNetworkProperties(dAtA, i, uint64(m.MinProposalEnactmentBlocks))
		i--
		dAtA[i] = 0x40
	}
	if m.MinProposalEndBlocks != 0 {
		i = encodeVarintNetworkProperties(dAtA, i, uint64(m.MinProposalEndBlocks))
		i--
		dAtA[i] = 0x38
	}
	if m.ProposalEnactmentTime != 0 {
		i = encodeVarintNetworkProperties(dAtA, i, uint64(m.ProposalEnactmentTime))
		i--
		dAtA[i] = 0x30
	}
	if m.MinimumProposalEndTime != 0 {
		i = encodeVarintNetworkProperties(dAtA, i, uint64(m.MinimumProposalEndTime))
		i--
		dAtA[i] = 0x28
	}
	if m.DefaultProposalEndTime != 0 {
		i = encodeVarintNetworkProperties(dAtA, i, uint64(m.DefaultProposalEndTime))
		i--
		dAtA[i] = 0x20
	}
	if m.VoteQuorum != 0 {
		i = encodeVarintNetworkProperties(dAtA, i, uint64(m.VoteQuorum))
		i--
		dAtA[i] = 0x18
	}
	if m.MaxTxFee != 0 {
		i = encodeVarintNetworkProperties(dAtA, i, uint64(m.MaxTxFee))
		i--
		dAtA[i] = 0x10
	}
	if m.MinTxFee != 0 {
		i = encodeVarintNetworkProperties(dAtA, i, uint64(m.MinTxFee))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintNetworkProperties(dAtA []byte, offset int, v uint64) int {
	offset -= sovNetworkProperties(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}

func (m *NetworkPropertiesV0228) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.MinTxFee != 0 {
		n += 1 + sovNetworkProperties(uint64(m.MinTxFee))
	}
	if m.MaxTxFee != 0 {
		n += 1 + sovNetworkProperties(uint64(m.MaxTxFee))
	}
	if m.VoteQuorum != 0 {
		n += 1 + sovNetworkProperties(uint64(m.VoteQuorum))
	}
	if m.DefaultProposalEndTime != 0 {
		n += 1 + sovNetworkProperties(uint64(m.DefaultProposalEndTime))
	}
	if m.MinimumProposalEndTime != 0 {
		n += 1 + sovNetworkProperties(uint64(m.MinimumProposalEndTime))
	}
	if m.ProposalEnactmentTime != 0 {
		n += 1 + sovNetworkProperties(uint64(m.ProposalEnactmentTime))
	}
	if m.MinProposalEndBlocks != 0 {
		n += 1 + sovNetworkProperties(uint64(m.MinProposalEndBlocks))
	}
	if m.MinProposalEnactmentBlocks != 0 {
		n += 1 + sovNetworkProperties(uint64(m.MinProposalEnactmentBlocks))
	}
	if m.EnableForeignFeePayments {
		n += 2
	}
	if m.MischanceRankDecreaseAmount != 0 {
		n += 1 + sovNetworkProperties(uint64(m.MischanceRankDecreaseAmount))
	}
	if m.MaxMischance != 0 {
		n += 1 + sovNetworkProperties(uint64(m.MaxMischance))
	}
	if m.MischanceConfidence != 0 {
		n += 1 + sovNetworkProperties(uint64(m.MischanceConfidence))
	}
	if m.InactiveRankDecreasePercent != 0 {
		n += 1 + sovNetworkProperties(uint64(m.InactiveRankDecreasePercent))
	}
	if m.MinValidators != 0 {
		n += 1 + sovNetworkProperties(uint64(m.MinValidators))
	}
	if m.PoorNetworkMaxBankSend != 0 {
		n += 1 + sovNetworkProperties(uint64(m.PoorNetworkMaxBankSend))
	}
	if m.JailMaxTime != 0 {
		n += 2 + sovNetworkProperties(uint64(m.JailMaxTime))
	}
	if m.EnableTokenWhitelist {
		n += 3
	}
	if m.EnableTokenBlacklist {
		n += 3
	}
	if m.MinIdentityApprovalTip != 0 {
		n += 2 + sovNetworkProperties(uint64(m.MinIdentityApprovalTip))
	}
	l = len(m.UniqueIdentityKeys)
	if l > 0 {
		n += 2 + l + sovNetworkProperties(uint64(l))
	}
	return n
}

func sovNetworkProperties(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozNetworkProperties(x uint64) (n int) {
	return sovNetworkProperties(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}

func (m *NetworkPropertiesV0228) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowNetworkProperties
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
			return fmt.Errorf("proto: NetworkProperties: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: NetworkProperties: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MinTxFee", wireType)
			}
			m.MinTxFee = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNetworkProperties
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MinTxFee |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaxTxFee", wireType)
			}
			m.MaxTxFee = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNetworkProperties
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MaxTxFee |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field VoteQuorum", wireType)
			}
			m.VoteQuorum = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNetworkProperties
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.VoteQuorum |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field DefaultProposalEndTime", wireType)
			}
			m.DefaultProposalEndTime = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNetworkProperties
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.DefaultProposalEndTime |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MinimumProposalEndTime", wireType)
			}
			m.MinimumProposalEndTime = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNetworkProperties
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MinimumProposalEndTime |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ProposalEnactmentTime", wireType)
			}
			m.ProposalEnactmentTime = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNetworkProperties
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ProposalEnactmentTime |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 7:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MinProposalEndBlocks", wireType)
			}
			m.MinProposalEndBlocks = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNetworkProperties
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MinProposalEndBlocks |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 8:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MinProposalEnactmentBlocks", wireType)
			}
			m.MinProposalEnactmentBlocks = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNetworkProperties
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MinProposalEnactmentBlocks |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 9:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field EnableForeignFeePayments", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNetworkProperties
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
			m.EnableForeignFeePayments = bool(v != 0)
		case 10:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MischanceRankDecreaseAmount", wireType)
			}
			m.MischanceRankDecreaseAmount = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNetworkProperties
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MischanceRankDecreaseAmount |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 11:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaxMischance", wireType)
			}
			m.MaxMischance = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNetworkProperties
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MaxMischance |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 12:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MischanceConfidence", wireType)
			}
			m.MischanceConfidence = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNetworkProperties
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MischanceConfidence |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 13:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field InactiveRankDecreasePercent", wireType)
			}
			m.InactiveRankDecreasePercent = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNetworkProperties
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.InactiveRankDecreasePercent |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 14:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MinValidators", wireType)
			}
			m.MinValidators = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNetworkProperties
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MinValidators |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 15:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field PoorNetworkMaxBankSend", wireType)
			}
			m.PoorNetworkMaxBankSend = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNetworkProperties
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.PoorNetworkMaxBankSend |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 16:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field JailMaxTime", wireType)
			}
			m.JailMaxTime = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNetworkProperties
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.JailMaxTime |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 17:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field EnableTokenWhitelist", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNetworkProperties
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
			m.EnableTokenWhitelist = bool(v != 0)
		case 18:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field EnableTokenBlacklist", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNetworkProperties
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
			m.EnableTokenBlacklist = bool(v != 0)
		case 19:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MinIdentityApprovalTip", wireType)
			}
			m.MinIdentityApprovalTip = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNetworkProperties
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MinIdentityApprovalTip |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 20:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field UniqueIdentityKeys", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNetworkProperties
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
				return ErrInvalidLengthNetworkProperties
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthNetworkProperties
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.UniqueIdentityKeys = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipNetworkProperties(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthNetworkProperties
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
func skipNetworkProperties(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowNetworkProperties
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
					return 0, ErrIntOverflowNetworkProperties
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
					return 0, ErrIntOverflowNetworkProperties
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
				return 0, ErrInvalidLengthNetworkProperties
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupNetworkProperties
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthNetworkProperties
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthNetworkProperties        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowNetworkProperties          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupNetworkProperties = fmt.Errorf("proto: unexpected end of group")
)
