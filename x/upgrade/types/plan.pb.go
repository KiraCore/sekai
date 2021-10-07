// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: kira/upgrade/plan.proto

package types

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
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

type Plan struct {
	Name                 string     `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Resources            []Resource `protobuf:"bytes,2,rep,name=resources,proto3" json:"resources"`
	UpgradeTime          int64      `protobuf:"varint,3,opt,name=upgrade_time,json=upgradeTime,proto3" json:"upgrade_time,omitempty"`
	OldChainId           string     `protobuf:"bytes,4,opt,name=old_chain_id,json=oldChainId,proto3" json:"old_chain_id,omitempty"`
	NewChainId           string     `protobuf:"bytes,5,opt,name=new_chain_id,json=newChainId,proto3" json:"new_chain_id,omitempty"`
	RollbackChecksum     string     `protobuf:"bytes,6,opt,name=rollback_checksum,json=rollbackChecksum,proto3" json:"rollback_checksum,omitempty"`
	MaxEnrolmentDuration int64      `protobuf:"varint,7,opt,name=max_enrolment_duration,json=maxEnrolmentDuration,proto3" json:"max_enrolment_duration,omitempty"`
	InstateUpgrade       bool       `protobuf:"varint,8,opt,name=instate_upgrade,json=instateUpgrade,proto3" json:"instate_upgrade,omitempty"`
	RebootRequired       bool       `protobuf:"varint,9,opt,name=reboot_required,json=rebootRequired,proto3" json:"reboot_required,omitempty"`
	SkipHandler          bool       `protobuf:"varint,10,opt,name=skip_handler,json=skipHandler,proto3" json:"skip_handler,omitempty"`
}

func (m *Plan) Reset()         { *m = Plan{} }
func (m *Plan) String() string { return proto.CompactTextString(m) }
func (*Plan) ProtoMessage()    {}
func (*Plan) Descriptor() ([]byte, []int) {
	return fileDescriptor_96c9a10c2f596dea, []int{0}
}
func (m *Plan) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Plan) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Plan.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Plan) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Plan.Merge(m, src)
}
func (m *Plan) XXX_Size() int {
	return m.Size()
}
func (m *Plan) XXX_DiscardUnknown() {
	xxx_messageInfo_Plan.DiscardUnknown(m)
}

var xxx_messageInfo_Plan proto.InternalMessageInfo

func (m *Plan) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Plan) GetResources() []Resource {
	if m != nil {
		return m.Resources
	}
	return nil
}

func (m *Plan) GetUpgradeTime() int64 {
	if m != nil {
		return m.UpgradeTime
	}
	return 0
}

func (m *Plan) GetOldChainId() string {
	if m != nil {
		return m.OldChainId
	}
	return ""
}

func (m *Plan) GetNewChainId() string {
	if m != nil {
		return m.NewChainId
	}
	return ""
}

func (m *Plan) GetRollbackChecksum() string {
	if m != nil {
		return m.RollbackChecksum
	}
	return ""
}

func (m *Plan) GetMaxEnrolmentDuration() int64 {
	if m != nil {
		return m.MaxEnrolmentDuration
	}
	return 0
}

func (m *Plan) GetInstateUpgrade() bool {
	if m != nil {
		return m.InstateUpgrade
	}
	return false
}

func (m *Plan) GetRebootRequired() bool {
	if m != nil {
		return m.RebootRequired
	}
	return false
}

func (m *Plan) GetSkipHandler() bool {
	if m != nil {
		return m.SkipHandler
	}
	return false
}

func init() {
	proto.RegisterType((*Plan)(nil), "kira.upgrade.Plan")
}

func init() { proto.RegisterFile("kira/upgrade/plan.proto", fileDescriptor_96c9a10c2f596dea) }

var fileDescriptor_96c9a10c2f596dea = []byte{
	// 393 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0x91, 0x4d, 0x8b, 0xd4, 0x30,
	0x18, 0x80, 0xa7, 0x4e, 0x5d, 0x77, 0x32, 0x83, 0x1f, 0x61, 0x59, 0xc3, 0x1c, 0x6a, 0xf5, 0x62,
	0x45, 0x68, 0x41, 0x3d, 0x79, 0xdc, 0x2a, 0x28, 0x5e, 0xa4, 0xe8, 0xc5, 0x4b, 0x48, 0x9b, 0x97,
	0x36, 0x34, 0x1f, 0x35, 0x4d, 0xd9, 0xf1, 0x5f, 0x78, 0xf1, 0x3f, 0xed, 0x71, 0x8f, 0x9e, 0x44,
	0x66, 0xfe, 0x88, 0xb4, 0xcd, 0xb8, 0x7a, 0x6a, 0x78, 0x9e, 0xa7, 0xe4, 0x25, 0x2f, 0x7a, 0xd8,
	0x0a, 0xcb, 0xb2, 0xa1, 0xab, 0x2d, 0xe3, 0x90, 0x75, 0x92, 0xe9, 0xb4, 0xb3, 0xc6, 0x19, 0xbc,
	0x19, 0x45, 0xea, 0xc5, 0xf6, 0xac, 0x36, 0xb5, 0x99, 0x44, 0x36, 0x9e, 0xe6, 0x66, 0xbb, 0xfd,
	0xef, 0x67, 0xff, 0x9d, 0xdd, 0x93, 0x1f, 0x4b, 0x14, 0x7e, 0x94, 0x4c, 0x63, 0x8c, 0x42, 0xcd,
	0x14, 0x90, 0x20, 0x0e, 0x92, 0x55, 0x31, 0x9d, 0xf1, 0x6b, 0xb4, 0xb2, 0xd0, 0x9b, 0xc1, 0x56,
	0xd0, 0x93, 0x5b, 0xf1, 0x32, 0x59, 0xbf, 0x38, 0x4f, 0xff, 0xbd, 0x30, 0x2d, 0xbc, 0xbe, 0x08,
	0xaf, 0x7e, 0x3d, 0x5a, 0x14, 0x37, 0x39, 0x7e, 0x8c, 0x36, 0x3e, 0xa2, 0x4e, 0x28, 0x20, 0xcb,
	0x38, 0x48, 0x96, 0xc5, 0xda, 0xb3, 0x4f, 0x42, 0x01, 0x8e, 0xd1, 0xc6, 0x48, 0x4e, 0xab, 0x86,
	0x09, 0x4d, 0x05, 0x27, 0xe1, 0x74, 0x35, 0x32, 0x92, 0xe7, 0x23, 0x7a, 0xcf, 0xc7, 0x42, 0xc3,
	0xe5, 0x4d, 0x71, 0x7b, 0x2e, 0x34, 0x5c, 0x1e, 0x8b, 0xe7, 0xe8, 0x81, 0x35, 0x52, 0x96, 0xac,
	0x6a, 0x69, 0xd5, 0x40, 0xd5, 0xf6, 0x83, 0x22, 0x27, 0x53, 0x76, 0xff, 0x28, 0x72, 0xcf, 0xf1,
	0x2b, 0x74, 0xae, 0xd8, 0x8e, 0x82, 0xb6, 0x46, 0x2a, 0xd0, 0x8e, 0xf2, 0xc1, 0x32, 0x27, 0x8c,
	0x26, 0x77, 0xa6, 0xe9, 0xce, 0x14, 0xdb, 0xbd, 0x3d, 0xca, 0x37, 0xde, 0xe1, 0xa7, 0xe8, 0x9e,
	0xd0, 0xbd, 0x63, 0x0e, 0xa8, 0x9f, 0x9e, 0x9c, 0xc6, 0x41, 0x72, 0x5a, 0xdc, 0xf5, 0xf8, 0xf3,
	0x4c, 0xc7, 0xd0, 0x42, 0x69, 0x8c, 0xa3, 0x16, 0xbe, 0x0e, 0xc2, 0x02, 0x27, 0xab, 0x39, 0x9c,
	0x71, 0xe1, 0xe9, 0xf8, 0x36, 0x7d, 0x2b, 0x3a, 0xda, 0x30, 0xcd, 0x25, 0x58, 0x82, 0xa6, 0x6a,
	0x3d, 0xb2, 0x77, 0x33, 0xba, 0xc8, 0xaf, 0xf6, 0x51, 0x70, 0xbd, 0x8f, 0x82, 0xdf, 0xfb, 0x28,
	0xf8, 0x7e, 0x88, 0x16, 0xd7, 0x87, 0x68, 0xf1, 0xf3, 0x10, 0x2d, 0xbe, 0x3c, 0xab, 0x85, 0x6b,
	0x86, 0x32, 0xad, 0x8c, 0xca, 0x3e, 0x08, 0xcb, 0x72, 0x63, 0x21, 0xeb, 0xa1, 0x65, 0x22, 0xdb,
	0xfd, 0x5d, 0xb2, 0xfb, 0xd6, 0x41, 0x5f, 0x9e, 0x4c, 0x3b, 0x7e, 0xf9, 0x27, 0x00, 0x00, 0xff,
	0xff, 0x45, 0x87, 0x1e, 0xc6, 0x3e, 0x02, 0x00, 0x00,
}

func (m *Plan) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Plan) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Plan) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.SkipHandler {
		i--
		if m.SkipHandler {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x50
	}
	if m.RebootRequired {
		i--
		if m.RebootRequired {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x48
	}
	if m.InstateUpgrade {
		i--
		if m.InstateUpgrade {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x40
	}
	if m.MaxEnrolmentDuration != 0 {
		i = encodeVarintPlan(dAtA, i, uint64(m.MaxEnrolmentDuration))
		i--
		dAtA[i] = 0x38
	}
	if len(m.RollbackChecksum) > 0 {
		i -= len(m.RollbackChecksum)
		copy(dAtA[i:], m.RollbackChecksum)
		i = encodeVarintPlan(dAtA, i, uint64(len(m.RollbackChecksum)))
		i--
		dAtA[i] = 0x32
	}
	if len(m.NewChainId) > 0 {
		i -= len(m.NewChainId)
		copy(dAtA[i:], m.NewChainId)
		i = encodeVarintPlan(dAtA, i, uint64(len(m.NewChainId)))
		i--
		dAtA[i] = 0x2a
	}
	if len(m.OldChainId) > 0 {
		i -= len(m.OldChainId)
		copy(dAtA[i:], m.OldChainId)
		i = encodeVarintPlan(dAtA, i, uint64(len(m.OldChainId)))
		i--
		dAtA[i] = 0x22
	}
	if m.UpgradeTime != 0 {
		i = encodeVarintPlan(dAtA, i, uint64(m.UpgradeTime))
		i--
		dAtA[i] = 0x18
	}
	if len(m.Resources) > 0 {
		for iNdEx := len(m.Resources) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Resources[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintPlan(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	if len(m.Name) > 0 {
		i -= len(m.Name)
		copy(dAtA[i:], m.Name)
		i = encodeVarintPlan(dAtA, i, uint64(len(m.Name)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintPlan(dAtA []byte, offset int, v uint64) int {
	offset -= sovPlan(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Plan) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovPlan(uint64(l))
	}
	if len(m.Resources) > 0 {
		for _, e := range m.Resources {
			l = e.Size()
			n += 1 + l + sovPlan(uint64(l))
		}
	}
	if m.UpgradeTime != 0 {
		n += 1 + sovPlan(uint64(m.UpgradeTime))
	}
	l = len(m.OldChainId)
	if l > 0 {
		n += 1 + l + sovPlan(uint64(l))
	}
	l = len(m.NewChainId)
	if l > 0 {
		n += 1 + l + sovPlan(uint64(l))
	}
	l = len(m.RollbackChecksum)
	if l > 0 {
		n += 1 + l + sovPlan(uint64(l))
	}
	if m.MaxEnrolmentDuration != 0 {
		n += 1 + sovPlan(uint64(m.MaxEnrolmentDuration))
	}
	if m.InstateUpgrade {
		n += 2
	}
	if m.RebootRequired {
		n += 2
	}
	if m.SkipHandler {
		n += 2
	}
	return n
}

func sovPlan(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozPlan(x uint64) (n int) {
	return sovPlan(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Plan) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPlan
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
			return fmt.Errorf("proto: Plan: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Plan: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPlan
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
				return ErrInvalidLengthPlan
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPlan
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Resources", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPlan
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
				return ErrInvalidLengthPlan
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthPlan
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Resources = append(m.Resources, Resource{})
			if err := m.Resources[len(m.Resources)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field UpgradeTime", wireType)
			}
			m.UpgradeTime = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPlan
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.UpgradeTime |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field OldChainId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPlan
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
				return ErrInvalidLengthPlan
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPlan
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.OldChainId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field NewChainId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPlan
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
				return ErrInvalidLengthPlan
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPlan
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.NewChainId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RollbackChecksum", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPlan
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
				return ErrInvalidLengthPlan
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPlan
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.RollbackChecksum = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 7:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaxEnrolmentDuration", wireType)
			}
			m.MaxEnrolmentDuration = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPlan
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MaxEnrolmentDuration |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 8:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field InstateUpgrade", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPlan
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
			m.InstateUpgrade = bool(v != 0)
		case 9:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field RebootRequired", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPlan
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
			m.RebootRequired = bool(v != 0)
		case 10:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SkipHandler", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPlan
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
			m.SkipHandler = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipPlan(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthPlan
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
func skipPlan(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowPlan
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
					return 0, ErrIntOverflowPlan
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
					return 0, ErrIntOverflowPlan
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
				return 0, ErrInvalidLengthPlan
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupPlan
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthPlan
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthPlan        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowPlan          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupPlan = fmt.Errorf("proto: unexpected end of group")
)
