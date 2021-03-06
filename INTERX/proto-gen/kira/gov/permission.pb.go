// Code generated by protoc-gen-go. DO NOT EDIT.
// source: kira/gov/permission.proto

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

type PermValue int32

const (
	// PERMISSION_ZERO is a no-op permission.
	PermValue_PERMISSION_ZERO PermValue = 0
	// PERMISSION_SET_PERMISSIONS defines the permission that allows to Set Permissions to other actors.
	PermValue_PERMISSION_SET_PERMISSIONS PermValue = 1
	// PERMISSION_CLAIM_VALIDATOR defines the permission that allows to Claim a validator Seat.
	PermValue_PERMISSION_CLAIM_VALIDATOR PermValue = 2
	// PERMISSION_CLAIM_COUNCILOR defines the permission that allows to Claim a Councilor Seat.
	PermValue_PERMISSION_CLAIM_COUNCILOR PermValue = 3
	// PERMISSION_CREATE_SET_PERMISSIONS_PROPOSAL defines the permission needed to create proposals for setting permissions.
	PermValue_PERMISSION_CREATE_SET_PERMISSIONS_PROPOSAL PermValue = 4
	// PERMISSION_VOTE_SET_PERMISSIONS_PROPOSAL defines the permission that an actor must have in order to vote a
	// Proposal to set permissions.
	PermValue_PERMISSION_VOTE_SET_PERMISSIONS_PROPOSAL PermValue = 5
	//  PERMISSION_UPSERT_TOKEN_ALIAS
	PermValue_PERMISSION_UPSERT_TOKEN_ALIAS PermValue = 6
	// PERMISSION_CHANGE_TX_FEE
	PermValue_PERMISSION_CHANGE_TX_FEE PermValue = 7
	// PERMISSION_UPSERT_TOKEN_RATE
	PermValue_PERMISSION_UPSERT_TOKEN_RATE PermValue = 8
	// PERMISSION_UPSERT_ROLE makes possible to add, modify and assign roles.
	PermValue_PERMISSION_UPSERT_ROLE PermValue = 9
	// PERMISSION_CREATE_UPSERT_DATA_REGISTRY_PROPOSAL makes possible to create a proposal to change the Data Registry.
	PermValue_PERMISSION_CREATE_UPSERT_DATA_REGISTRY_PROPOSAL PermValue = 10
	// PERMISSION_VOTE_UPSERT_DATA_REGISTRY_PROPOSAL makes possible to create a proposal to change the Data Registry.
	PermValue_PERMISSION_VOTE_UPSERT_DATA_REGISTRY_PROPOSAL PermValue = 11
	// PERMISSION_CREATE_SET_NETWORK_PROPERTY_PROPOSAL defines the permission needed to create proposals for setting network property.
	PermValue_PERMISSION_CREATE_SET_NETWORK_PROPERTY_PROPOSAL PermValue = 12
	// PERMISSION_VOTE_SET_NETWORK_PROPERTY_PROPOSAL defines the permission that an actor must have in order to vote a
	// Proposal to set network property.
	PermValue_PERMISSION_VOTE_SET_NETWORK_PROPERTY_PROPOSAL PermValue = 13
	// PERMISSION_CREATE_UPSERT_TOKEN_ALIAS_PROPOSAL defines the permission needed to create proposals for upsert token Alias.
	PermValue_PERMISSION_CREATE_UPSERT_TOKEN_ALIAS_PROPOSAL PermValue = 14
	// PERMISSION_VOTE_UPSERT_TOKEN_ALIAS_PROPOSAL defines the permission needed to vote proposals for upsert token.
	PermValue_PERMISSION_VOTE_UPSERT_TOKEN_ALIAS_PROPOSAL PermValue = 15
	// PERMISSION_CREATE_SET_POOR_NETWORK_MESSAGES defines the permission needed to create proposals for setting poor network messages
	PermValue_PERMISSION_CREATE_SET_POOR_NETWORK_MESSAGES PermValue = 16
	// PERMISSION_VOTE_SET_POOR_NETWORK_MESSAGES_PROPOSAL defines the permission needed to vote proposals to set poor network messages
	PermValue_PERMISSION_VOTE_SET_POOR_NETWORK_MESSAGES_PROPOSAL PermValue = 17
	// PERMISSION_CREATE_UPSERT_TOKEN_RATE_PROPOSAL defines the permission needed to create proposals for upsert token rate.
	PermValue_PERMISSION_CREATE_UPSERT_TOKEN_RATE_PROPOSAL PermValue = 18
	// PERMISSION_VOTE_UPSERT_TOKEN_RATE_PROPOSAL defines the permission needed to vote proposals for upsert token rate.
	PermValue_PERMISSION_VOTE_UPSERT_TOKEN_RATE_PROPOSAL PermValue = 19
	// PERMISSION_CREATE_UNJAIL_VALIDATOR_PROPOSAL defines the permission needed to create a proposal to unjail a validator.
	PermValue_PERMISSION_CREATE_UNJAIL_VALIDATOR_PROPOSAL PermValue = 20
	// PERMISSION_VOTE_UNJAIL_VALIDATOR_PROPOSAL defines the permission needed to vote a proposal to unjail a validator.
	PermValue_PERMISSION_VOTE_UNJAIL_VALIDATOR_PROPOSAL PermValue = 21
	// PERMISSION_CREATE_CREATE_ROLE_PROPOSAL defines the permission needed to create a proposal to create a role.
	PermValue_PERMISSION_CREATE_CREATE_ROLE_PROPOSAL PermValue = 22
	// PERMISSION_VOTE_CREATE_ROLE_PROPOSAL defines the permission needed to vote a proposal to create a role.
	PermValue_PERMISSION_VOTE_CREATE_ROLE_PROPOSAL PermValue = 23
	// PERMISSION_CREATE_TOKENS_WHITE_BLACK_CHANGE_PROPOSAL defines the permission needed to create a proposal to blacklist/whitelisted tokens
	PermValue_PERMISSION_CREATE_TOKENS_WHITE_BLACK_CHANGE_PROPOSAL PermValue = 24
	// PERMISSION_VOTE_TOKENS_WHITE_BLACK_CHANGE_PROPOSAL defines the permission needed to vote on blacklist/whitelisted tokens proposal
	PermValue_PERMISSION_VOTE_TOKENS_WHITE_BLACK_CHANGE_PROPOSAL PermValue = 25
	// PERMISSION_CREATE_RESET_WHOLE_VALIDATOR_RANK_PROPOSAL defines the permission needed to create a proposal to reset whole validator rank
	PermValue_PERMISSION_CREATE_RESET_WHOLE_VALIDATOR_RANK_PROPOSAL PermValue = 26
	// PERMISSION_VOTE_RESET_WHOLE_VALIDATOR_RANK_PROPOSAL defines the permission needed to vote on reset whole validator rank proposal
	PermValue_PERMISSION_VOTE_RESET_WHOLE_VALIDATOR_RANK_PROPOSAL PermValue = 27
	// PERMISSION_CREATE_SOFTWARE_UPGRADE_PROPOSAL defines the permission needed to create a proposal for software upgrade
	PermValue_PERMISSION_CREATE_SOFTWARE_UPGRADE_PROPOSAL PermValue = 28
	// PERMISSION_SOFTWARE_UPGRADE_PROPOSAL defines the permission needed to vote on software upgrade proposal
	PermValue_PERMISSION_SOFTWARE_UPGRADE_PROPOSAL PermValue = 29
)

var PermValue_name = map[int32]string{
	0:  "PERMISSION_ZERO",
	1:  "PERMISSION_SET_PERMISSIONS",
	2:  "PERMISSION_CLAIM_VALIDATOR",
	3:  "PERMISSION_CLAIM_COUNCILOR",
	4:  "PERMISSION_CREATE_SET_PERMISSIONS_PROPOSAL",
	5:  "PERMISSION_VOTE_SET_PERMISSIONS_PROPOSAL",
	6:  "PERMISSION_UPSERT_TOKEN_ALIAS",
	7:  "PERMISSION_CHANGE_TX_FEE",
	8:  "PERMISSION_UPSERT_TOKEN_RATE",
	9:  "PERMISSION_UPSERT_ROLE",
	10: "PERMISSION_CREATE_UPSERT_DATA_REGISTRY_PROPOSAL",
	11: "PERMISSION_VOTE_UPSERT_DATA_REGISTRY_PROPOSAL",
	12: "PERMISSION_CREATE_SET_NETWORK_PROPERTY_PROPOSAL",
	13: "PERMISSION_VOTE_SET_NETWORK_PROPERTY_PROPOSAL",
	14: "PERMISSION_CREATE_UPSERT_TOKEN_ALIAS_PROPOSAL",
	15: "PERMISSION_VOTE_UPSERT_TOKEN_ALIAS_PROPOSAL",
	16: "PERMISSION_CREATE_SET_POOR_NETWORK_MESSAGES",
	17: "PERMISSION_VOTE_SET_POOR_NETWORK_MESSAGES_PROPOSAL",
	18: "PERMISSION_CREATE_UPSERT_TOKEN_RATE_PROPOSAL",
	19: "PERMISSION_VOTE_UPSERT_TOKEN_RATE_PROPOSAL",
	20: "PERMISSION_CREATE_UNJAIL_VALIDATOR_PROPOSAL",
	21: "PERMISSION_VOTE_UNJAIL_VALIDATOR_PROPOSAL",
	22: "PERMISSION_CREATE_CREATE_ROLE_PROPOSAL",
	23: "PERMISSION_VOTE_CREATE_ROLE_PROPOSAL",
	24: "PERMISSION_CREATE_TOKENS_WHITE_BLACK_CHANGE_PROPOSAL",
	25: "PERMISSION_VOTE_TOKENS_WHITE_BLACK_CHANGE_PROPOSAL",
	26: "PERMISSION_CREATE_RESET_WHOLE_VALIDATOR_RANK_PROPOSAL",
	27: "PERMISSION_VOTE_RESET_WHOLE_VALIDATOR_RANK_PROPOSAL",
	28: "PERMISSION_CREATE_SOFTWARE_UPGRADE_PROPOSAL",
	29: "PERMISSION_SOFTWARE_UPGRADE_PROPOSAL",
}

var PermValue_value = map[string]int32{
	"PERMISSION_ZERO":                                       0,
	"PERMISSION_SET_PERMISSIONS":                            1,
	"PERMISSION_CLAIM_VALIDATOR":                            2,
	"PERMISSION_CLAIM_COUNCILOR":                            3,
	"PERMISSION_CREATE_SET_PERMISSIONS_PROPOSAL":            4,
	"PERMISSION_VOTE_SET_PERMISSIONS_PROPOSAL":              5,
	"PERMISSION_UPSERT_TOKEN_ALIAS":                         6,
	"PERMISSION_CHANGE_TX_FEE":                              7,
	"PERMISSION_UPSERT_TOKEN_RATE":                          8,
	"PERMISSION_UPSERT_ROLE":                                9,
	"PERMISSION_CREATE_UPSERT_DATA_REGISTRY_PROPOSAL":       10,
	"PERMISSION_VOTE_UPSERT_DATA_REGISTRY_PROPOSAL":         11,
	"PERMISSION_CREATE_SET_NETWORK_PROPERTY_PROPOSAL":       12,
	"PERMISSION_VOTE_SET_NETWORK_PROPERTY_PROPOSAL":         13,
	"PERMISSION_CREATE_UPSERT_TOKEN_ALIAS_PROPOSAL":         14,
	"PERMISSION_VOTE_UPSERT_TOKEN_ALIAS_PROPOSAL":           15,
	"PERMISSION_CREATE_SET_POOR_NETWORK_MESSAGES":           16,
	"PERMISSION_VOTE_SET_POOR_NETWORK_MESSAGES_PROPOSAL":    17,
	"PERMISSION_CREATE_UPSERT_TOKEN_RATE_PROPOSAL":          18,
	"PERMISSION_VOTE_UPSERT_TOKEN_RATE_PROPOSAL":            19,
	"PERMISSION_CREATE_UNJAIL_VALIDATOR_PROPOSAL":           20,
	"PERMISSION_VOTE_UNJAIL_VALIDATOR_PROPOSAL":             21,
	"PERMISSION_CREATE_CREATE_ROLE_PROPOSAL":                22,
	"PERMISSION_VOTE_CREATE_ROLE_PROPOSAL":                  23,
	"PERMISSION_CREATE_TOKENS_WHITE_BLACK_CHANGE_PROPOSAL":  24,
	"PERMISSION_VOTE_TOKENS_WHITE_BLACK_CHANGE_PROPOSAL":    25,
	"PERMISSION_CREATE_RESET_WHOLE_VALIDATOR_RANK_PROPOSAL": 26,
	"PERMISSION_VOTE_RESET_WHOLE_VALIDATOR_RANK_PROPOSAL":   27,
	"PERMISSION_CREATE_SOFTWARE_UPGRADE_PROPOSAL":           28,
	"PERMISSION_SOFTWARE_UPGRADE_PROPOSAL":                  29,
}

func (x PermValue) String() string {
	return proto.EnumName(PermValue_name, int32(x))
}

func (PermValue) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_214168f8815c1062, []int{0}
}

func init() {
	proto.RegisterEnum("kira.gov.PermValue", PermValue_name, PermValue_value)
}

func init() {
	proto.RegisterFile("kira/gov/permission.proto", fileDescriptor_214168f8815c1062)
}

var fileDescriptor_214168f8815c1062 = []byte{
	// 902 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x96, 0x6f, 0x73, 0xdb, 0x34,
	0x1c, 0xc7, 0x07, 0x8c, 0xd1, 0x89, 0xb1, 0x1a, 0xb7, 0x74, 0x9d, 0x58, 0x87, 0x0a, 0xa5, 0x74,
	0xdd, 0xd6, 0x1c, 0x0c, 0xb8, 0xe3, 0x78, 0xa4, 0xa6, 0x6a, 0x6b, 0xf2, 0xc7, 0x39, 0xd9, 0x69,
	0xb6, 0xde, 0xed, 0x7c, 0x6a, 0xab, 0x25, 0x26, 0x69, 0x94, 0x93, 0xdd, 0x0d, 0xde, 0x01, 0xa7,
	0xf7, 0xa0, 0x47, 0xbc, 0x47, 0x1e, 0x73, 0x76, 0x5b, 0xcb, 0x76, 0x9c, 0xb4, 0x8f, 0x9c, 0xc4,
	0xfa, 0x7d, 0x7e, 0x5f, 0x7d, 0x7f, 0x7f, 0x2e, 0xe0, 0xf1, 0x30, 0x94, 0xac, 0xd6, 0x17, 0xef,
	0x6b, 0x13, 0x2e, 0xcf, 0xc3, 0x28, 0x0a, 0xc5, 0x78, 0x67, 0x22, 0x45, 0x2c, 0xec, 0x85, 0xe4,
	0xd5, 0x4e, 0x5f, 0xbc, 0x87, 0xcb, 0x7d, 0xd1, 0x17, 0xe9, 0x8f, 0xb5, 0xe4, 0xd3, 0xe5, 0xfb,
	0xed, 0xff, 0x6c, 0x70, 0xbf, 0xc3, 0xe5, 0xf9, 0x11, 0x1b, 0x5d, 0x70, 0x7b, 0x1d, 0x2c, 0x76,
	0x08, 0x6d, 0x39, 0x9e, 0xe7, 0xb8, 0xed, 0xe0, 0x98, 0x50, 0xd7, 0xba, 0x03, 0x1f, 0x28, 0x8d,
	0x16, 0x92, 0x33, 0xc7, 0x5c, 0x0a, 0xfb, 0x57, 0x00, 0x73, 0x47, 0x3c, 0xe2, 0x07, 0xe6, 0xab,
	0x67, 0x7d, 0x04, 0x57, 0x94, 0x46, 0x76, 0x72, 0xda, 0xe3, 0x71, 0x27, 0x53, 0x13, 0x95, 0xe2,
	0xea, 0x4d, 0xec, 0xb4, 0x82, 0x23, 0xdc, 0x74, 0xf6, 0xb0, 0xef, 0x52, 0xeb, 0x63, 0x13, 0x57,
	0x1f, 0xb1, 0x30, 0x91, 0x13, 0x9e, 0xb1, 0x58, 0xc8, 0xca, 0xb8, 0xba, 0xdb, 0x6d, 0xd7, 0x9d,
	0xa6, 0x4b, 0xad, 0x4f, 0x4a, 0x71, 0x75, 0x71, 0x31, 0x3e, 0x0d, 0x47, 0x42, 0xda, 0x3e, 0xd8,
	0xce, 0xc7, 0x51, 0x82, 0x7d, 0x52, 0x96, 0x1b, 0x74, 0xa8, 0xdb, 0x71, 0x3d, 0xdc, 0xb4, 0xee,
	0xc2, 0x0d, 0xa5, 0x11, 0x4a, 0x39, 0x92, 0xb3, 0x98, 0x17, 0xd5, 0x77, 0xa4, 0x98, 0x88, 0x88,
	0x8d, 0x6c, 0x17, 0x6c, 0xe5, 0xa8, 0x47, 0xee, 0x3c, 0xe6, 0xa7, 0x70, 0x5d, 0x69, 0xb4, 0x96,
	0xba, 0x2b, 0x4a, 0xc4, 0x0c, 0xf8, 0x3b, 0x58, 0xcb, 0x01, 0xbb, 0x1d, 0x8f, 0x50, 0x3f, 0xf0,
	0xdd, 0x06, 0x69, 0x07, 0xb8, 0xe9, 0x60, 0xcf, 0xba, 0x07, 0x57, 0x95, 0x46, 0xcb, 0x49, 0x68,
	0x77, 0x12, 0x71, 0x19, 0xfb, 0x62, 0xc8, 0xc7, 0x78, 0x14, 0xb2, 0xc8, 0xfe, 0x11, 0xac, 0xe6,
	0xef, 0x78, 0x88, 0xdb, 0x07, 0x24, 0xf0, 0x5f, 0x07, 0xfb, 0x84, 0x58, 0x9f, 0xc1, 0x25, 0xa5,
	0xd1, 0x62, 0x7a, 0xa3, 0x01, 0x1b, 0xf7, 0xb9, 0xff, 0xd7, 0x3e, 0xe7, 0xf6, 0x6f, 0xe0, 0xc9,
	0xac, 0x7c, 0x14, 0xfb, 0xc4, 0x5a, 0x80, 0x8f, 0x94, 0x46, 0x4b, 0xa5, 0x74, 0x94, 0xc5, 0xdc,
	0xde, 0x01, 0x2b, 0xd3, 0xa1, 0xd4, 0x6d, 0x12, 0xeb, 0x3e, 0xb4, 0x95, 0x46, 0x0f, 0x4d, 0x10,
	0x15, 0x23, 0x6e, 0xbf, 0x05, 0xb5, 0xe9, 0x0a, 0x5c, 0x85, 0xed, 0x61, 0x1f, 0x07, 0x94, 0x1c,
	0x38, 0x9e, 0x4f, 0xdf, 0x18, 0xcb, 0x00, 0xdc, 0x52, 0x1a, 0x6d, 0x98, 0x32, 0x5c, 0xe2, 0xf6,
	0x58, 0xcc, 0x28, 0xef, 0x87, 0x51, 0x2c, 0xff, 0xce, 0x9c, 0x7b, 0x03, 0x5e, 0x96, 0x4b, 0x31,
	0x1f, 0xfe, 0x39, 0xdc, 0x54, 0x1a, 0x7d, 0x7b, 0x5d, 0x8f, 0x39, 0xe8, 0x4a, 0xe5, 0x49, 0x9d,
	0xdb, 0xc4, 0xef, 0xb9, 0xb4, 0x91, 0x32, 0x09, 0xf5, 0x73, 0xf0, 0x07, 0x65, 0xe5, 0x1e, 0x8f,
	0xdb, 0x3c, 0xfe, 0x20, 0xe4, 0x30, 0xc1, 0x72, 0x19, 0xcf, 0x55, 0x3e, 0x1f, 0xfe, 0x45, 0x51,
	0xf9, 0xad, 0xd1, 0x45, 0xcf, 0x73, 0x5d, 0x65, 0xd0, 0x0f, 0x0d, 0x3a, 0xef, 0xb8, 0x69, 0xb2,
	0x0c, 0xdd, 0x05, 0xcf, 0x67, 0xf8, 0x5d, 0x09, 0x5e, 0x34, 0x13, 0x65, 0xdc, 0xae, 0xc0, 0xbe,
	0x2d, 0x60, 0xf3, 0x73, 0xea, 0xba, 0x34, 0xf3, 0xa4, 0x45, 0x3c, 0x0f, 0x1f, 0x10, 0xcf, 0xb2,
	0xe0, 0x0b, 0xa5, 0xd1, 0x56, 0x71, 0x50, 0x85, 0x90, 0x57, 0x86, 0xb4, 0x78, 0x14, 0xb1, 0x3e,
	0x37, 0xf8, 0x13, 0xf0, 0x53, 0xe5, 0xc0, 0x56, 0xc1, 0x8d, 0xf8, 0x2f, 0xe1, 0xb6, 0xd2, 0x68,
	0x33, 0x3f, 0xba, 0x73, 0x72, 0xf4, 0xc0, 0x8b, 0x1b, 0x4c, 0x4f, 0x46, 0xcb, 0xd0, 0x6d, 0xf8,
	0xbd, 0xd2, 0x68, 0xbd, 0xd2, 0xf3, 0x64, 0xd2, 0x32, 0xb0, 0x57, 0xd8, 0x61, 0xd3, 0x96, 0x17,
	0xb1, 0x4b, 0xf0, 0x3b, 0xa5, 0xd1, 0x37, 0x15, 0x8e, 0x17, 0xa0, 0x47, 0x55, 0x86, 0x77, 0xdb,
	0x7f, 0x60, 0xa7, 0x69, 0x16, 0xb2, 0xa1, 0x2e, 0x4f, 0x89, 0x1d, 0xff, 0xc9, 0xc2, 0x51, 0xb6,
	0xa0, 0x33, 0x2e, 0x05, 0xcf, 0xa6, 0xc4, 0xce, 0xa4, 0x7e, 0x55, 0xd2, 0x3a, 0x83, 0xb9, 0x0f,
	0x36, 0xa7, 0xb5, 0x5e, 0x3d, 0x92, 0xcd, 0x63, 0x80, 0x2b, 0x10, 0x2a, 0x8d, 0x56, 0x8c, 0xcc,
	0x64, 0x05, 0x65, 0x9c, 0x43, 0xb0, 0x51, 0xd6, 0x56, 0x49, 0x79, 0x04, 0x9f, 0x2a, 0x8d, 0xe0,
	0xb5, 0xac, 0x0a, 0xd2, 0x3b, 0xf0, 0xf3, 0xb4, 0xa2, 0xb4, 0x1a, 0x5e, 0xd0, 0x3b, 0x74, 0x7c,
	0x12, 0xec, 0x36, 0x71, 0xbd, 0x71, 0xbd, 0x8c, 0x33, 0xf2, 0x6a, 0xb9, 0x6f, 0xd3, 0xc2, 0x44,
	0xbd, 0x41, 0x18, 0xf3, 0xdd, 0x11, 0x3b, 0x1d, 0x5e, 0x2e, 0xe9, 0x79, 0x7d, 0x7b, 0x8b, 0x2c,
	0x8f, 0x8b, 0x7d, 0x7b, 0x43, 0x8e, 0x01, 0xf8, 0x65, 0xfa, 0x2e, 0x94, 0x24, 0xf3, 0xd1, 0x3b,
	0x4c, 0x7c, 0x31, 0x85, 0xa3, 0xb8, 0xdd, 0x30, 0x69, 0x20, 0x7c, 0xa9, 0x34, 0x7a, 0x96, 0x33,
	0x9b, 0x47, 0x3c, 0xee, 0x0d, 0xc4, 0x88, 0x67, 0x35, 0xa4, 0x6c, 0x3c, 0xcc, 0x32, 0x9d, 0x81,
	0x57, 0xe5, 0xdb, 0xdc, 0x26, 0xcf, 0xd7, 0xf0, 0xb9, 0xd2, 0xe8, 0x87, 0xeb, 0xeb, 0xdc, 0x94,
	0xa5, 0xb2, 0xb3, 0x3d, 0x77, 0xdf, 0xef, 0x61, 0x9a, 0x4c, 0xce, 0x01, 0xc5, 0x7b, 0x39, 0xb3,
	0x9e, 0x94, 0x3b, 0xdb, 0x13, 0xef, 0xe2, 0x0f, 0x4c, 0xf2, 0xee, 0xa4, 0x2f, 0xd9, 0x99, 0xf1,
	0xa9, 0x55, 0xe8, 0x9e, 0xd9, 0xc0, 0xb5, 0x62, 0x53, 0xcf, 0xc0, 0xc1, 0xbb, 0xff, 0xfc, 0xfb,
	0xf4, 0xce, 0xee, 0xd6, 0xf1, 0x66, 0x3f, 0x8c, 0x07, 0x17, 0x27, 0x3b, 0xa7, 0xe2, 0xbc, 0xd6,
	0x08, 0x25, 0xab, 0x0b, 0xc9, 0x6b, 0x11, 0x1f, 0xb2, 0xb0, 0xe6, 0xb4, 0x7d, 0x42, 0x5f, 0xd7,
	0xd2, 0xbf, 0x68, 0x27, 0xf7, 0xd2, 0xc7, 0xab, 0xff, 0x03, 0x00, 0x00, 0xff, 0xff, 0x63, 0x16,
	0x55, 0x36, 0xe6, 0x09, 0x00, 0x00,
}
