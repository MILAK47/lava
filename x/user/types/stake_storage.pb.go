// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: user/stake_storage.proto

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

type StakeStorage struct {
	StakedUsers []UserStake `protobuf:"bytes,1,rep,name=stakedUsers,proto3" json:"stakedUsers"`
}

func (m *StakeStorage) Reset()         { *m = StakeStorage{} }
func (m *StakeStorage) String() string { return proto.CompactTextString(m) }
func (*StakeStorage) ProtoMessage()    {}
func (*StakeStorage) Descriptor() ([]byte, []int) {
	return fileDescriptor_67a553e13cbfdb8d, []int{0}
}
func (m *StakeStorage) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *StakeStorage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_StakeStorage.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *StakeStorage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StakeStorage.Merge(m, src)
}
func (m *StakeStorage) XXX_Size() int {
	return m.Size()
}
func (m *StakeStorage) XXX_DiscardUnknown() {
	xxx_messageInfo_StakeStorage.DiscardUnknown(m)
}

var xxx_messageInfo_StakeStorage proto.InternalMessageInfo

func (m *StakeStorage) GetStakedUsers() []UserStake {
	if m != nil {
		return m.StakedUsers
	}
	return nil
}

func init() {
	proto.RegisterType((*StakeStorage)(nil), "lavanet.lava.user.StakeStorage")
}

func init() { proto.RegisterFile("user/stake_storage.proto", fileDescriptor_67a553e13cbfdb8d) }

var fileDescriptor_67a553e13cbfdb8d = []byte{
	// 195 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x28, 0x2d, 0x4e, 0x2d,
	0xd2, 0x2f, 0x2e, 0x49, 0xcc, 0x4e, 0x8d, 0x2f, 0x2e, 0xc9, 0x2f, 0x4a, 0x4c, 0x4f, 0xd5, 0x2b,
	0x28, 0xca, 0x2f, 0xc9, 0x17, 0x12, 0xcc, 0x49, 0x2c, 0x4b, 0xcc, 0x4b, 0x2d, 0xd1, 0x03, 0xd1,
	0x7a, 0x20, 0x65, 0x52, 0xa2, 0x60, 0xc5, 0x20, 0x22, 0x1e, 0xac, 0x03, 0xa2, 0x52, 0x4a, 0x24,
	0x3d, 0x3f, 0x3d, 0x1f, 0xcc, 0xd4, 0x07, 0xb1, 0x20, 0xa2, 0x4a, 0x21, 0x5c, 0x3c, 0xc1, 0x20,
	0x45, 0xc1, 0x10, 0x53, 0x85, 0x5c, 0xb8, 0xb8, 0xc1, 0x9a, 0x52, 0x42, 0x8b, 0x53, 0x8b, 0x8a,
	0x25, 0x18, 0x15, 0x98, 0x35, 0xb8, 0x8d, 0x64, 0xf4, 0x30, 0x6c, 0xd1, 0x03, 0xc9, 0x83, 0x75,
	0x3a, 0xb1, 0x9c, 0xb8, 0x27, 0xcf, 0x10, 0x84, 0xac, 0xcd, 0xc9, 0xee, 0xc4, 0x23, 0x39, 0xc6,
	0x0b, 0x8f, 0xe4, 0x18, 0x1f, 0x3c, 0x92, 0x63, 0x9c, 0xf0, 0x58, 0x8e, 0xe1, 0xc2, 0x63, 0x39,
	0x86, 0x1b, 0x8f, 0xe5, 0x18, 0xa2, 0x54, 0xd2, 0x33, 0x4b, 0x32, 0x4a, 0x93, 0xf4, 0x92, 0xf3,
	0x73, 0xf5, 0xa1, 0x86, 0x82, 0x69, 0xfd, 0x0a, 0xb0, 0x8b, 0xf5, 0x4b, 0x2a, 0x0b, 0x52, 0x8b,
	0x93, 0xd8, 0xc0, 0x8e, 0x33, 0x06, 0x04, 0x00, 0x00, 0xff, 0xff, 0xc4, 0x4d, 0xe2, 0x77, 0xf8,
	0x00, 0x00, 0x00,
}

func (m *StakeStorage) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *StakeStorage) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *StakeStorage) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.StakedUsers) > 0 {
		for iNdEx := len(m.StakedUsers) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.StakedUsers[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintStakeStorage(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func encodeVarintStakeStorage(dAtA []byte, offset int, v uint64) int {
	offset -= sovStakeStorage(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *StakeStorage) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.StakedUsers) > 0 {
		for _, e := range m.StakedUsers {
			l = e.Size()
			n += 1 + l + sovStakeStorage(uint64(l))
		}
	}
	return n
}

func sovStakeStorage(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozStakeStorage(x uint64) (n int) {
	return sovStakeStorage(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *StakeStorage) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowStakeStorage
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
			return fmt.Errorf("proto: StakeStorage: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: StakeStorage: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field StakedUsers", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStakeStorage
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
				return ErrInvalidLengthStakeStorage
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthStakeStorage
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.StakedUsers = append(m.StakedUsers, UserStake{})
			if err := m.StakedUsers[len(m.StakedUsers)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipStakeStorage(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthStakeStorage
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
func skipStakeStorage(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowStakeStorage
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
					return 0, ErrIntOverflowStakeStorage
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
					return 0, ErrIntOverflowStakeStorage
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
				return 0, ErrInvalidLengthStakeStorage
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupStakeStorage
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthStakeStorage
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthStakeStorage        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowStakeStorage          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupStakeStorage = fmt.Errorf("proto: unexpected end of group")
)