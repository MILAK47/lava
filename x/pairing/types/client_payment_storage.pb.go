// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: pairing/client_payment_storage.proto

package types

import (
	fmt "fmt"
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

type ClientPaymentStorage struct {
	Index                              string                                `protobuf:"bytes,1,opt,name=index,proto3" json:"index,omitempty"`
	UniquePaymentStorageClientProvider []*UniquePaymentStorageClientProvider `protobuf:"bytes,2,rep,name=uniquePaymentStorageClientProvider,proto3" json:"uniquePaymentStorageClientProvider,omitempty"`
	TotalCU                            uint64                                `protobuf:"varint,3,opt,name=totalCU,proto3" json:"totalCU,omitempty"`
	Epoch                              uint64                                `protobuf:"varint,4,opt,name=epoch,proto3" json:"epoch,omitempty"`
}

func (m *ClientPaymentStorage) Reset()         { *m = ClientPaymentStorage{} }
func (m *ClientPaymentStorage) String() string { return proto.CompactTextString(m) }
func (*ClientPaymentStorage) ProtoMessage()    {}
func (*ClientPaymentStorage) Descriptor() ([]byte, []int) {
	return fileDescriptor_da8671b7794bb7f3, []int{0}
}
func (m *ClientPaymentStorage) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ClientPaymentStorage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ClientPaymentStorage.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ClientPaymentStorage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ClientPaymentStorage.Merge(m, src)
}
func (m *ClientPaymentStorage) XXX_Size() int {
	return m.Size()
}
func (m *ClientPaymentStorage) XXX_DiscardUnknown() {
	xxx_messageInfo_ClientPaymentStorage.DiscardUnknown(m)
}

var xxx_messageInfo_ClientPaymentStorage proto.InternalMessageInfo

func (m *ClientPaymentStorage) GetIndex() string {
	if m != nil {
		return m.Index
	}
	return ""
}

func (m *ClientPaymentStorage) GetUniquePaymentStorageClientProvider() []*UniquePaymentStorageClientProvider {
	if m != nil {
		return m.UniquePaymentStorageClientProvider
	}
	return nil
}

func (m *ClientPaymentStorage) GetTotalCU() uint64 {
	if m != nil {
		return m.TotalCU
	}
	return 0
}

func (m *ClientPaymentStorage) GetEpoch() uint64 {
	if m != nil {
		return m.Epoch
	}
	return 0
}

func init() {
	proto.RegisterType((*ClientPaymentStorage)(nil), "lavanet.lava.pairing.ClientPaymentStorage")
}

func init() {
	proto.RegisterFile("pairing/client_payment_storage.proto", fileDescriptor_da8671b7794bb7f3)
}

var fileDescriptor_da8671b7794bb7f3 = []byte{
	// 256 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x29, 0x48, 0xcc, 0x2c,
	0xca, 0xcc, 0x4b, 0xd7, 0x4f, 0xce, 0xc9, 0x4c, 0xcd, 0x2b, 0x89, 0x2f, 0x48, 0xac, 0xcc, 0x05,
	0xd1, 0xc5, 0x25, 0xf9, 0x45, 0x89, 0xe9, 0xa9, 0x7a, 0x05, 0x45, 0xf9, 0x25, 0xf9, 0x42, 0x22,
	0x39, 0x89, 0x65, 0x89, 0x79, 0xa9, 0x25, 0x7a, 0x20, 0x5a, 0x0f, 0xaa, 0x45, 0xca, 0x04, 0xa6,
	0xb7, 0x34, 0x2f, 0xb3, 0xb0, 0x34, 0x15, 0x5d, 0x6f, 0x3c, 0xcc, 0xc8, 0xa2, 0xfc, 0xb2, 0xcc,
	0x94, 0xd4, 0x22, 0x88, 0x59, 0x4a, 0xcf, 0x19, 0xb9, 0x44, 0x9c, 0xc1, 0x32, 0x01, 0x10, 0xf5,
	0xc1, 0x10, 0xe5, 0x42, 0x22, 0x5c, 0xac, 0x99, 0x79, 0x29, 0xa9, 0x15, 0x12, 0x8c, 0x0a, 0x8c,
	0x1a, 0x9c, 0x41, 0x10, 0x8e, 0x50, 0x07, 0x23, 0x97, 0x12, 0xc4, 0x7c, 0x54, 0xe5, 0x50, 0x23,
	0xa0, 0x66, 0x4b, 0x30, 0x29, 0x30, 0x6b, 0x70, 0x1b, 0x59, 0xe8, 0x61, 0x73, 0xa8, 0x5e, 0x28,
	0x41, 0xfd, 0x41, 0x44, 0xd8, 0x21, 0x24, 0xc1, 0xc5, 0x5e, 0x92, 0x5f, 0x92, 0x98, 0xe3, 0x1c,
	0x2a, 0xc1, 0xac, 0xc0, 0xa8, 0xc1, 0x12, 0x04, 0xe3, 0x82, 0x9c, 0x9e, 0x5a, 0x90, 0x9f, 0x9c,
	0x21, 0xc1, 0x02, 0x16, 0x87, 0x70, 0x9c, 0x1c, 0x4f, 0x3c, 0x92, 0x63, 0xbc, 0xf0, 0x48, 0x8e,
	0xf1, 0xc1, 0x23, 0x39, 0xc6, 0x09, 0x8f, 0xe5, 0x18, 0x2e, 0x3c, 0x96, 0x63, 0xb8, 0xf1, 0x58,
	0x8e, 0x21, 0x4a, 0x3d, 0x3d, 0xb3, 0x24, 0xa3, 0x34, 0x49, 0x2f, 0x39, 0x3f, 0x57, 0x1f, 0xea,
	0x62, 0x30, 0xad, 0x5f, 0xa1, 0x0f, 0x0b, 0xd3, 0x92, 0xca, 0x82, 0xd4, 0xe2, 0x24, 0x36, 0x70,
	0x98, 0x19, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff, 0xe4, 0x1f, 0x34, 0xce, 0xa7, 0x01, 0x00, 0x00,
}

func (m *ClientPaymentStorage) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ClientPaymentStorage) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ClientPaymentStorage) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Epoch != 0 {
		i = encodeVarintClientPaymentStorage(dAtA, i, uint64(m.Epoch))
		i--
		dAtA[i] = 0x20
	}
	if m.TotalCU != 0 {
		i = encodeVarintClientPaymentStorage(dAtA, i, uint64(m.TotalCU))
		i--
		dAtA[i] = 0x18
	}
	if len(m.UniquePaymentStorageClientProvider) > 0 {
		for iNdEx := len(m.UniquePaymentStorageClientProvider) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.UniquePaymentStorageClientProvider[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintClientPaymentStorage(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	if len(m.Index) > 0 {
		i -= len(m.Index)
		copy(dAtA[i:], m.Index)
		i = encodeVarintClientPaymentStorage(dAtA, i, uint64(len(m.Index)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintClientPaymentStorage(dAtA []byte, offset int, v uint64) int {
	offset -= sovClientPaymentStorage(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *ClientPaymentStorage) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Index)
	if l > 0 {
		n += 1 + l + sovClientPaymentStorage(uint64(l))
	}
	if len(m.UniquePaymentStorageClientProvider) > 0 {
		for _, e := range m.UniquePaymentStorageClientProvider {
			l = e.Size()
			n += 1 + l + sovClientPaymentStorage(uint64(l))
		}
	}
	if m.TotalCU != 0 {
		n += 1 + sovClientPaymentStorage(uint64(m.TotalCU))
	}
	if m.Epoch != 0 {
		n += 1 + sovClientPaymentStorage(uint64(m.Epoch))
	}
	return n
}

func sovClientPaymentStorage(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozClientPaymentStorage(x uint64) (n int) {
	return sovClientPaymentStorage(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *ClientPaymentStorage) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowClientPaymentStorage
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
			return fmt.Errorf("proto: ClientPaymentStorage: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ClientPaymentStorage: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Index", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowClientPaymentStorage
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
				return ErrInvalidLengthClientPaymentStorage
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthClientPaymentStorage
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Index = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field UniquePaymentStorageClientProvider", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowClientPaymentStorage
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
				return ErrInvalidLengthClientPaymentStorage
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthClientPaymentStorage
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.UniquePaymentStorageClientProvider = append(m.UniquePaymentStorageClientProvider, &UniquePaymentStorageClientProvider{})
			if err := m.UniquePaymentStorageClientProvider[len(m.UniquePaymentStorageClientProvider)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field TotalCU", wireType)
			}
			m.TotalCU = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowClientPaymentStorage
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.TotalCU |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Epoch", wireType)
			}
			m.Epoch = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowClientPaymentStorage
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Epoch |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipClientPaymentStorage(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthClientPaymentStorage
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
func skipClientPaymentStorage(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowClientPaymentStorage
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
					return 0, ErrIntOverflowClientPaymentStorage
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
					return 0, ErrIntOverflowClientPaymentStorage
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
				return 0, ErrInvalidLengthClientPaymentStorage
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupClientPaymentStorage
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthClientPaymentStorage
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthClientPaymentStorage        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowClientPaymentStorage          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupClientPaymentStorage = fmt.Errorf("proto: unexpected end of group")
)
