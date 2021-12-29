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
	ProposalEndTime             uint64 `protobuf:"varint,4,opt,name=proposal_end_time,json=proposalEndTime,proto3" json:"proposal_end_time,omitempty"`
	ProposalEnactmentTime       uint64 `protobuf:"varint,5,opt,name=proposal_enactment_time,json=proposalEnactmentTime,proto3" json:"proposal_enactment_time,omitempty"`
	MinProposalEndBlocks        uint64 `protobuf:"varint,6,opt,name=min_proposal_end_blocks,json=minProposalEndBlocks,proto3" json:"min_proposal_end_blocks,omitempty"`
	MinProposalEnactmentBlocks  uint64 `protobuf:"varint,7,opt,name=min_proposal_enactment_blocks,json=minProposalEnactmentBlocks,proto3" json:"min_proposal_enactment_blocks,omitempty"`
	EnableForeignFeePayments    bool   `protobuf:"varint,8,opt,name=enable_foreign_fee_payments,json=enableForeignFeePayments,proto3" json:"enable_foreign_fee_payments,omitempty"`
	MischanceRankDecreaseAmount uint64 `protobuf:"varint,9,opt,name=mischance_rank_decrease_amount,json=mischanceRankDecreaseAmount,proto3" json:"mischance_rank_decrease_amount,omitempty"`
	MaxMischance                uint64 `protobuf:"varint,10,opt,name=max_mischance,json=maxMischance,proto3" json:"max_mischance,omitempty"`
	MischanceConfidence         uint64 `protobuf:"varint,11,opt,name=mischance_confidence,json=mischanceConfidence,proto3" json:"mischance_confidence,omitempty"`
	InactiveRankDecreasePercent uint64 `protobuf:"varint,12,opt,name=inactive_rank_decrease_percent,json=inactiveRankDecreasePercent,proto3" json:"inactive_rank_decrease_percent,omitempty"`
	MinValidators               uint64 `protobuf:"varint,13,opt,name=min_validators,json=minValidators,proto3" json:"min_validators,omitempty"`
	PoorNetworkMaxBankSend      uint64 `protobuf:"varint,14,opt,name=poor_network_max_bank_send,json=poorNetworkMaxBankSend,proto3" json:"poor_network_max_bank_send,omitempty"`
	JailMaxTime                 uint64 `protobuf:"varint,15,opt,name=jail_max_time,json=jailMaxTime,proto3" json:"jail_max_time,omitempty"`
	EnableTokenWhitelist        bool   `protobuf:"varint,16,opt,name=enable_token_whitelist,json=enableTokenWhitelist,proto3" json:"enable_token_whitelist,omitempty"`
	EnableTokenBlacklist        bool   `protobuf:"varint,17,opt,name=enable_token_blacklist,json=enableTokenBlacklist,proto3" json:"enable_token_blacklist,omitempty"`
	MinIdentityApprovalTip      uint64 `protobuf:"varint,18,opt,name=min_identity_approval_tip,json=minIdentityApprovalTip,proto3" json:"min_identity_approval_tip,omitempty"`
	UniqueIdentityKeys          string `protobuf:"bytes,19,opt,name=unique_identity_keys,json=uniqueIdentityKeys,proto3" json:"unique_identity_keys,omitempty"`
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

func (m *NetworkPropertiesV0228) GetProposalEndTime() uint64 {
	if m != nil {
		return m.ProposalEndTime
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
	// 1202 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x56, 0xcb, 0x6e, 0xdb, 0x46,
	0x17, 0xb6, 0x12, 0x27, 0x96, 0xc7, 0x96, 0x4d, 0xd3, 0x8a, 0xcd, 0xd0, 0xf9, 0x65, 0xfe, 0x2e,
	0x02, 0x04, 0x01, 0x62, 0x35, 0xbd, 0x01, 0x0d, 0x50, 0x14, 0x94, 0x3c, 0x6e, 0x18, 0x89, 0x97,
	0x50, 0xb4, 0x92, 0x74, 0x33, 0x18, 0x4b, 0x13, 0x65, 0x2a, 0x91, 0x54, 0x48, 0x4a, 0x91, 0xdf,
	0xa0, 0xe0, 0xaa, 0x2f, 0x40, 0xa0, 0x40, 0x5f, 0xa1, 0x0f, 0xd1, 0x5d, 0x03, 0x74, 0xd3, 0x55,
	0x51, 0xc4, 0x9b, 0x3e, 0x43, 0x57, 0xc5, 0x0c, 0x49, 0xd9, 0x92, 0x62, 0xaf, 0x2c, 0xcf, 0x77,
	0x39, 0xe7, 0xcc, 0xcc, 0x47, 0x0c, 0xf8, 0x7f, 0x9f, 0x06, 0xb8, 0xda, 0xf3, 0xc7, 0x55, 0x8f,
	0x44, 0xef, 0xfc, 0xa0, 0x8f, 0x86, 0x81, 0x3f, 0x24, 0x41, 0x44, 0x49, 0x78, 0x38, 0x0c, 0xfc,
	0xc8, 0x17, 0x8b, 0x8c, 0x72, 0xd8, 0xf3, 0xc7, 0x72, 0xb9, 0xe7, 0xf7, 0x7c, 0xbe, 0x58, 0x65,
	0xbf, 0x52, 0xfc, 0xe0, 0xd7, 0x02, 0xd8, 0xd5, 0xc3, 0x5e, 0x8b, 0x44, 0x46, 0x6a, 0x61, 0x4d,
	0x1d, 0xc4, 0x67, 0x40, 0x5c, 0xf4, 0x95, 0x0a, 0x4a, 0xe1, 0xc1, 0xda, 0x67, 0x7b, 0x87, 0xb9,
	0xf1, 0xe1, 0x82, 0xd0, 0xde, 0xf2, 0x16, 0xbc, 0x74, 0x50, 0x64, 0x1e, 0x7e, 0x48, 0x02, 0xe9,
	0x86, 0x52, 0x78, 0xb0, 0x5e, 0x7b, 0xfc, 0xef, 0x5f, 0xfb, 0x8f, 0x7a, 0x34, 0x7a, 0x33, 0x3a,
	0x3d, 0xec, 0xf8, 0x6e, 0xb5, 0xe3, 0x87, 0xae, 0x1f, 0x66, 0x7f, 0x1e, 0x85, 0xdd, 0x7e, 0x35,
	0x3a, 0x1b, 0x92, 0xf0, 0x50, 0xed, 0x74, 0xd4, 0x6e, 0x37, 0x20, 0x61, 0x68, 0x4f, 0x2d, 0x0e,
	0x4c, 0x50, 0x9e, 0x2d, 0x7b, 0xd6, 0xc6, 0x83, 0x11, 0x11, 0xcb, 0xe0, 0xd6, 0x98, 0xfd, 0xe0,
	0x5d, 0x2e, 0xdb, 0xe9, 0x3f, 0xe2, 0x1e, 0x58, 0x0d, 0xa3, 0x00, 0xa5, 0x08, 0xab, 0xbe, 0x6a,
	0x17, 0xc3, 0x28, 0xe0, 0x92, 0x27, 0xcb, 0xff, 0xfc, 0xbc, 0x5f, 0x38, 0xf8, 0x7d, 0x05, 0x6c,
	0x2d, 0xee, 0xc0, 0x3d, 0x00, 0x5c, 0xea, 0xa1, 0x68, 0x82, 0x5e, 0x93, 0xdc, 0xb3, 0xe8, 0x52,
	0xcf, 0x99, 0x1c, 0x13, 0xc2, 0x51, 0x3c, 0xc9, 0xd1, 0x1b, 0x19, 0x8a, 0x27, 0x29, 0xba, 0x0f,
	0xd6, 0xc6, 0x7e, 0x44, 0xd0, 0xdb, 0x91, 0x1f, 0x8c, 0x5c, 0xe9, 0x26, 0x87, 0x01, 0x5b, 0x7a,
	0xce, 0x57, 0xc4, 0x87, 0x60, 0x2b, 0x9d, 0x07, 0x0f, 0x10, 0xf1, 0xba, 0x28, 0xa2, 0x2e, 0x91,
	0x96, 0x39, 0x6d, 0x33, 0x07, 0xa0, 0xd7, 0x75, 0xa8, 0x4b, 0xc4, 0xaf, 0xc0, 0xee, 0x25, 0x2e,
	0xee, 0x44, 0x2e, 0xf1, 0xa2, 0x54, 0x71, 0x8b, 0x2b, 0xee, 0x5c, 0x28, 0x32, 0x94, 0xeb, 0xbe,
	0x04, 0xbb, 0x6c, 0x80, 0x99, 0x3a, 0xa7, 0x03, 0xbf, 0xd3, 0x0f, 0xa5, 0xdb, 0x5c, 0x57, 0x76,
	0xa9, 0x67, 0x5d, 0x14, 0xab, 0x71, 0x4c, 0x54, 0xc1, 0xff, 0xe6, 0x64, 0x79, 0xc9, 0x4c, 0xbc,
	0xc2, 0xc5, 0xf2, 0x8c, 0x38, 0xa3, 0x64, 0x16, 0xdf, 0x80, 0x3d, 0xe2, 0xe1, 0xd3, 0x01, 0x41,
	0xaf, 0xfd, 0x80, 0xd0, 0x9e, 0xc7, 0x36, 0x09, 0x0d, 0xf1, 0x19, 0xe3, 0x84, 0x52, 0x51, 0x29,
	0x3c, 0x28, 0xda, 0x52, 0x4a, 0x39, 0x4e, 0x19, 0xc7, 0x84, 0x58, 0x19, 0x2e, 0xd6, 0x41, 0xc5,
	0xa5, 0x61, 0xe7, 0x0d, 0xf6, 0x3a, 0x04, 0x05, 0xd8, 0xeb, 0xa3, 0x2e, 0xe9, 0x04, 0x04, 0x87,
	0x04, 0x61, 0xd7, 0x1f, 0x79, 0x91, 0xb4, 0xca, 0x5b, 0xd8, 0x9b, 0xb2, 0x6c, 0xec, 0xf5, 0x8f,
	0x32, 0x8e, 0xca, 0x29, 0xe2, 0x27, 0xa0, 0xc4, 0x0e, 0x68, 0x4a, 0x91, 0x00, 0xd7, 0xac, 0xbb,
	0x78, 0xa2, 0xe7, 0x6b, 0xe2, 0x63, 0x50, 0xbe, 0xa8, 0xd4, 0xf1, 0xbd, 0xd7, 0xb4, 0x4b, 0x18,
	0x77, 0x8d, 0x73, 0xb7, 0xa7, 0x58, 0x7d, 0x0a, 0xb1, 0xe6, 0x28, 0x1b, 0x97, 0x8e, 0xe7, 0x7b,
	0x1b, 0x92, 0xa0, 0x43, 0xbc, 0x48, 0x5a, 0x4f, 0x9b, 0xcb, 0x59, 0x97, 0x7b, 0xb3, 0x52, 0x8a,
	0x78, 0x1f, 0x6c, 0xb0, 0x3d, 0x1e, 0xe3, 0x01, 0xed, 0xe2, 0xc8, 0x0f, 0x42, 0xa9, 0xc4, 0x45,
	0x25, 0x97, 0x7a, 0xed, 0xe9, 0xa2, 0xf8, 0x04, 0xc8, 0x43, 0xdf, 0x0f, 0x50, 0x9e, 0x44, 0x36,
	0xd0, 0x29, 0xab, 0x19, 0x12, 0xaf, 0x2b, 0x6d, 0x70, 0xc9, 0x0e, 0x63, 0x64, 0xb7, 0x57, 0xc7,
	0x93, 0x1a, 0xf6, 0xfa, 0x2d, 0xe2, 0x75, 0xc5, 0x03, 0x50, 0xfa, 0x01, 0xd3, 0x01, 0xd7, 0xf0,
	0xbb, 0xb2, 0xc9, 0xe9, 0x6b, 0x6c, 0x51, 0xc7, 0x13, 0x7e, 0x43, 0xbe, 0x00, 0x3b, 0xd9, 0x39,
	0x45, 0x7e, 0x9f, 0x78, 0xe8, 0xdd, 0x1b, 0x1a, 0x91, 0x01, 0x0d, 0x23, 0x49, 0xe0, 0x47, 0x54,
	0x4e, 0x51, 0x87, 0x81, 0x2f, 0x72, 0x6c, 0x41, 0x75, 0x3a, 0xc0, 0x9d, 0x3e, 0x57, 0x6d, 0x2d,
	0xa8, 0x6a, 0x39, 0x26, 0x7e, 0x0d, 0xee, 0xb2, 0x91, 0xd9, 0x2e, 0x46, 0x34, 0x3a, 0x43, 0x78,
	0x38, 0x0c, 0xfc, 0x31, 0x1e, 0xa0, 0x88, 0x0e, 0x25, 0x31, 0x1d, 0xc5, 0xa5, 0x9e, 0x96, 0xe1,
	0x6a, 0x06, 0x3b, 0x74, 0x28, 0x7e, 0x0a, 0xca, 0x23, 0x8f, 0xbe, 0x1d, 0x91, 0x0b, 0x75, 0x9f,
	0x9c, 0x85, 0xd2, 0x36, 0x4f, 0xb3, 0x98, 0x62, 0xb9, 0xb0, 0x41, 0xce, 0xc2, 0x87, 0x7f, 0xac,
	0x80, 0xcd, 0xb9, 0x6f, 0x04, 0x4b, 0xac, 0xae, 0x19, 0xc8, 0x79, 0x89, 0x8e, 0x21, 0x14, 0x96,
	0xe4, 0xf5, 0x38, 0x51, 0x8a, 0xfa, 0xa5, 0x3c, 0xeb, 0xea, 0xcb, 0x1c, 0x2d, 0x64, 0xe8, 0xa5,
	0x3c, 0xb7, 0x4d, 0x07, 0xa2, 0xe7, 0x27, 0xa6, 0x7d, 0xa2, 0x0b, 0x37, 0xe4, 0x8d, 0x38, 0x51,
	0x40, 0x7b, 0x26, 0xcf, 0x96, 0x6d, 0x5a, 0x66, 0x4b, 0x6d, 0x22, 0x68, 0x1c, 0x21, 0x47, 0xd3,
	0xa1, 0x70, 0x53, 0xde, 0x8e, 0x13, 0x65, 0xd3, 0x5a, 0xcc, 0xf3, 0x25, 0xae, 0x5a, 0x77, 0x74,
	0x68, 0x38, 0xa9, 0x62, 0x59, 0xbe, 0x1b, 0x27, 0xca, 0x1d, 0xeb, 0xaa, 0x3c, 0xb3, 0x01, 0x66,
	0xea, 0xd4, 0x9a, 0x66, 0xbd, 0xd1, 0x12, 0x6e, 0xc9, 0x52, 0x9c, 0x28, 0x65, 0xfd, 0x8a, 0x3c,
	0xcf, 0xc9, 0xf2, 0x92, 0x99, 0xf8, 0xb6, 0x5c, 0x89, 0x13, 0x45, 0xd6, 0xaf, 0xcd, 0x33, 0x34,
	0xd4, 0x5a, 0x13, 0xa2, 0x63, 0xd3, 0x86, 0xda, 0x77, 0x06, 0xdb, 0x24, 0x64, 0xa9, 0xaf, 0x98,
	0x4d, 0x4b, 0x58, 0x91, 0xef, 0xc5, 0x89, 0x22, 0xc1, 0x6b, 0xf2, 0xac, 0x6b, 0xad, 0xfa, 0x53,
	0xd5, 0xa8, 0x43, 0x64, 0xab, 0x46, 0x03, 0x1d, 0xc1, 0xba, 0x0d, 0xd5, 0x16, 0x44, 0xaa, 0x6e,
	0x9e, 0x18, 0x8e, 0x50, 0x94, 0xf7, 0xe3, 0x44, 0xd9, 0xd3, 0xaf, 0xcf, 0x33, 0x3b, 0xa0, 0xa9,
	0x91, 0xb0, 0x2a, 0x0b, 0x71, 0xa2, 0xac, 0xeb, 0x73, 0x79, 0xbe, 0xa8, 0x54, 0x37, 0x8d, 0x63,
	0xed, 0x08, 0x32, 0x2e, 0x90, 0x77, 0xe3, 0x44, 0xd9, 0xd6, 0x3f, 0x9e, 0x67, 0x8d, 0xed, 0x88,
	0xd6, 0x9e, 0xef, 0xcd, 0x82, 0x76, 0x1d, 0x1a, 0x8e, 0xb0, 0x96, 0x36, 0xa7, 0x5d, 0x93, 0xe7,
	0x27, 0x40, 0xb6, 0x4c, 0xd3, 0x46, 0x06, 0x74, 0x5e, 0x98, 0x76, 0x03, 0xb1, 0x4e, 0x6b, 0xcc,
	0xac, 0x05, 0x8d, 0x23, 0x61, 0x5d, 0x96, 0xe3, 0x44, 0xd9, 0xb1, 0x3e, 0x1e, 0xd4, 0xfb, 0x60,
	0x83, 0x9d, 0x4f, 0x5b, 0x6d, 0x6a, 0x47, 0xaa, 0x63, 0xda, 0x2d, 0xa1, 0x24, 0x6f, 0xc5, 0x89,
	0x52, 0xd2, 0x67, 0xbe, 0x05, 0x07, 0xa0, 0xf4, 0x4c, 0xd5, 0x9a, 0xdc, 0x9a, 0xdf, 0x95, 0x0d,
	0x79, 0x33, 0x4e, 0x94, 0xb5, 0x67, 0xb3, 0x79, 0xce, 0xce, 0xc9, 0x31, 0x1b, 0xd0, 0x40, 0x2f,
	0x9e, 0x6a, 0x0e, 0x6c, 0x6a, 0x2d, 0x47, 0xd8, 0x4c, 0x2f, 0x08, 0xbc, 0x22, 0xcf, 0x33, 0xaa,
	0x5a, 0x53, 0xad, 0x37, 0xb8, 0x4a, 0x58, 0x50, 0xcd, 0xe4, 0x99, 0xb5, 0xcd, 0x36, 0xd8, 0xd1,
	0x9c, 0x57, 0x48, 0xb5, 0x2c, 0xdb, 0x6c, 0xab, 0x4d, 0xe4, 0x68, 0x96, 0xb0, 0x95, 0x4e, 0xac,
	0x5f, 0x99, 0xe7, 0x13, 0x43, 0x7b, 0x7e, 0x02, 0x2f, 0xd4, 0x0d, 0xf8, 0xaa, 0x25, 0x88, 0xf2,
	0x4e, 0x9c, 0x28, 0xe2, 0xc9, 0x42, 0x9e, 0xe5, 0xe5, 0x1f, 0x7f, 0xa9, 0x2c, 0xd5, 0xbe, 0xfd,
	0xed, 0x43, 0xa5, 0xf0, 0xfe, 0x43, 0xa5, 0xf0, 0xf7, 0x87, 0x4a, 0xe1, 0xa7, 0xf3, 0xca, 0xd2,
	0xfb, 0xf3, 0xca, 0xd2, 0x9f, 0xe7, 0x95, 0xa5, 0xef, 0xef, 0x5f, 0x7a, 0x4b, 0x34, 0x68, 0x80,
	0xeb, 0x7e, 0x40, 0xaa, 0x21, 0xe9, 0x63, 0x5a, 0x9d, 0xf0, 0x37, 0x12, 0x7f, 0x4e, 0x9c, 0xde,
	0xe6, 0xef, 0x9e, 0xcf, 0xff, 0x0b, 0x00, 0x00, 0xff, 0xff, 0x47, 0x3a, 0x99, 0xc8, 0x3c, 0x09,
	0x00, 0x00,
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
		dAtA[i] = 0x9a
	}
	if m.MinIdentityApprovalTip != 0 {
		i = encodeVarintNetworkProperties(dAtA, i, uint64(m.MinIdentityApprovalTip))
		i--
		dAtA[i] = 0x1
		i--
		dAtA[i] = 0x90
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
		dAtA[i] = 0x88
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
		dAtA[i] = 0x80
	}
	if m.JailMaxTime != 0 {
		i = encodeVarintNetworkProperties(dAtA, i, uint64(m.JailMaxTime))
		i--
		dAtA[i] = 0x78
	}
	if m.PoorNetworkMaxBankSend != 0 {
		i = encodeVarintNetworkProperties(dAtA, i, uint64(m.PoorNetworkMaxBankSend))
		i--
		dAtA[i] = 0x70
	}
	if m.MinValidators != 0 {
		i = encodeVarintNetworkProperties(dAtA, i, uint64(m.MinValidators))
		i--
		dAtA[i] = 0x68
	}
	if m.InactiveRankDecreasePercent != 0 {
		i = encodeVarintNetworkProperties(dAtA, i, uint64(m.InactiveRankDecreasePercent))
		i--
		dAtA[i] = 0x60
	}
	if m.MischanceConfidence != 0 {
		i = encodeVarintNetworkProperties(dAtA, i, uint64(m.MischanceConfidence))
		i--
		dAtA[i] = 0x58
	}
	if m.MaxMischance != 0 {
		i = encodeVarintNetworkProperties(dAtA, i, uint64(m.MaxMischance))
		i--
		dAtA[i] = 0x50
	}
	if m.MischanceRankDecreaseAmount != 0 {
		i = encodeVarintNetworkProperties(dAtA, i, uint64(m.MischanceRankDecreaseAmount))
		i--
		dAtA[i] = 0x48
	}
	if m.EnableForeignFeePayments {
		i--
		if m.EnableForeignFeePayments {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x40
	}
	if m.MinProposalEnactmentBlocks != 0 {
		i = encodeVarintNetworkProperties(dAtA, i, uint64(m.MinProposalEnactmentBlocks))
		i--
		dAtA[i] = 0x38
	}
	if m.MinProposalEndBlocks != 0 {
		i = encodeVarintNetworkProperties(dAtA, i, uint64(m.MinProposalEndBlocks))
		i--
		dAtA[i] = 0x30
	}
	if m.ProposalEnactmentTime != 0 {
		i = encodeVarintNetworkProperties(dAtA, i, uint64(m.ProposalEnactmentTime))
		i--
		dAtA[i] = 0x28
	}
	if m.ProposalEndTime != 0 {
		i = encodeVarintNetworkProperties(dAtA, i, uint64(m.ProposalEndTime))
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
	if m.ProposalEndTime != 0 {
		n += 1 + sovNetworkProperties(uint64(m.ProposalEndTime))
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
		n += 1 + sovNetworkProperties(uint64(m.JailMaxTime))
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
				return fmt.Errorf("proto: wrong wireType = %d for field ProposalEndTime", wireType)
			}
			m.ProposalEndTime = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNetworkProperties
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ProposalEndTime |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
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
		case 6:
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
		case 7:
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
		case 8:
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
		case 9:
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
		case 10:
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
		case 11:
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
		case 12:
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
		case 13:
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
		case 14:
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
		case 15:
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
		case 16:
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
		case 17:
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
		case 18:
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
		case 19:
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
