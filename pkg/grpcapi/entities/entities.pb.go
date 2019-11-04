// Code generated by protoc-gen-go. DO NOT EDIT.
// source: entities/entities.proto

package entities

import (
	fmt "fmt"
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

type AccountType int32

const (
	AccountType_ACCOUNT_TYPE_UNKNOWN AccountType = 0
	AccountType_ACCOUNT_TYPE_COMPANY AccountType = 1
	AccountType_ACCOUNT_TYPE_PERSONA AccountType = 2
)

var AccountType_name = map[int32]string{
	0: "ACCOUNT_TYPE_UNKNOWN",
	1: "ACCOUNT_TYPE_COMPANY",
	2: "ACCOUNT_TYPE_PERSONA",
}

var AccountType_value = map[string]int32{
	"ACCOUNT_TYPE_UNKNOWN": 0,
	"ACCOUNT_TYPE_COMPANY": 1,
	"ACCOUNT_TYPE_PERSONA": 2,
}

func (x AccountType) String() string {
	return proto.EnumName(AccountType_name, int32(x))
}

func (AccountType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_db230341a541e6ba, []int{0}
}

type Empty struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Empty) Reset()         { *m = Empty{} }
func (m *Empty) String() string { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()    {}
func (*Empty) Descriptor() ([]byte, []int) {
	return fileDescriptor_db230341a541e6ba, []int{0}
}

func (m *Empty) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Empty.Unmarshal(m, b)
}
func (m *Empty) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Empty.Marshal(b, m, deterministic)
}
func (m *Empty) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Empty.Merge(m, src)
}
func (m *Empty) XXX_Size() int {
	return xxx_messageInfo_Empty.Size(m)
}
func (m *Empty) XXX_DiscardUnknown() {
	xxx_messageInfo_Empty.DiscardUnknown(m)
}

var xxx_messageInfo_Empty proto.InternalMessageInfo

func init() {
	proto.RegisterEnum("personaappapi.entities.AccountType", AccountType_name, AccountType_value)
	proto.RegisterType((*Empty)(nil), "personaappapi.entities.Empty")
}

func init() { proto.RegisterFile("entities/entities.proto", fileDescriptor_db230341a541e6ba) }

var fileDescriptor_db230341a541e6ba = []byte{
	// 185 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x4f, 0xcd, 0x2b, 0xc9,
	0x2c, 0xc9, 0x4c, 0x2d, 0xd6, 0x87, 0x31, 0xf4, 0x0a, 0x8a, 0xf2, 0x4b, 0xf2, 0x85, 0xc4, 0x0a,
	0x52, 0x8b, 0x8a, 0xf3, 0xf3, 0x12, 0x13, 0x0b, 0x0a, 0x12, 0x0b, 0x32, 0xf5, 0x60, 0xb2, 0x4a,
	0xec, 0x5c, 0xac, 0xae, 0xb9, 0x05, 0x25, 0x95, 0x5a, 0xd1, 0x5c, 0xdc, 0x8e, 0xc9, 0xc9, 0xf9,
	0xa5, 0x79, 0x25, 0x21, 0x95, 0x05, 0xa9, 0x42, 0x12, 0x5c, 0x22, 0x8e, 0xce, 0xce, 0xfe, 0xa1,
	0x7e, 0x21, 0xf1, 0x21, 0x91, 0x01, 0xae, 0xf1, 0xa1, 0x7e, 0xde, 0x7e, 0xfe, 0xe1, 0x7e, 0x02,
	0x0c, 0x18, 0x32, 0xce, 0xfe, 0xbe, 0x01, 0x8e, 0x7e, 0x91, 0x02, 0x8c, 0x18, 0x32, 0x01, 0xae,
	0x41, 0xc1, 0xfe, 0x7e, 0x8e, 0x02, 0x4c, 0x4e, 0xae, 0x5c, 0x62, 0xf9, 0x79, 0x39, 0x99, 0x79,
	0xa9, 0x7a, 0x08, 0x67, 0xe8, 0x25, 0xe7, 0x17, 0xa5, 0x3a, 0x71, 0xb9, 0x17, 0x15, 0x24, 0xbb,
	0x82, 0x5c, 0x53, 0x19, 0x25, 0x8f, 0x90, 0xd4, 0x2f, 0xc8, 0x4e, 0xd7, 0x4f, 0x2f, 0x2a, 0x48,
	0x4e, 0x2c, 0xc8, 0x84, 0x7b, 0x25, 0x89, 0x0d, 0xec, 0x17, 0x63, 0x40, 0x00, 0x00, 0x00, 0xff,
	0xff, 0xe4, 0x6d, 0x32, 0x0f, 0xe6, 0x00, 0x00, 0x00,
}
