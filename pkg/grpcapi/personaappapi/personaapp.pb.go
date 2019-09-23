// Code generated by protoc-gen-go. DO NOT EDIT.
// source: personaapp.proto

package personaappapi

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
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

type Ping struct {
	Key                  string               `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value                string               `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	CreatedAt            *timestamp.Timestamp `protobuf:"bytes,3,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt            *timestamp.Timestamp `protobuf:"bytes,4,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *Ping) Reset()         { *m = Ping{} }
func (m *Ping) String() string { return proto.CompactTextString(m) }
func (*Ping) ProtoMessage()    {}
func (*Ping) Descriptor() ([]byte, []int) {
	return fileDescriptor_a6083bb8799e537c, []int{0}
}

func (m *Ping) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Ping.Unmarshal(m, b)
}
func (m *Ping) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Ping.Marshal(b, m, deterministic)
}
func (m *Ping) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Ping.Merge(m, src)
}
func (m *Ping) XXX_Size() int {
	return xxx_messageInfo_Ping.Size(m)
}
func (m *Ping) XXX_DiscardUnknown() {
	xxx_messageInfo_Ping.DiscardUnknown(m)
}

var xxx_messageInfo_Ping proto.InternalMessageInfo

func (m *Ping) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *Ping) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

func (m *Ping) GetCreatedAt() *timestamp.Timestamp {
	if m != nil {
		return m.CreatedAt
	}
	return nil
}

func (m *Ping) GetUpdatedAt() *timestamp.Timestamp {
	if m != nil {
		return m.UpdatedAt
	}
	return nil
}

type SetPingRequest struct {
	Key                  string   `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value                string   `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SetPingRequest) Reset()         { *m = SetPingRequest{} }
func (m *SetPingRequest) String() string { return proto.CompactTextString(m) }
func (*SetPingRequest) ProtoMessage()    {}
func (*SetPingRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_a6083bb8799e537c, []int{1}
}

func (m *SetPingRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SetPingRequest.Unmarshal(m, b)
}
func (m *SetPingRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SetPingRequest.Marshal(b, m, deterministic)
}
func (m *SetPingRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SetPingRequest.Merge(m, src)
}
func (m *SetPingRequest) XXX_Size() int {
	return xxx_messageInfo_SetPingRequest.Size(m)
}
func (m *SetPingRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SetPingRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SetPingRequest proto.InternalMessageInfo

func (m *SetPingRequest) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *SetPingRequest) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

type SetPingResponse struct {
	Ping                 *Ping    `protobuf:"bytes,1,opt,name=ping,proto3" json:"ping,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SetPingResponse) Reset()         { *m = SetPingResponse{} }
func (m *SetPingResponse) String() string { return proto.CompactTextString(m) }
func (*SetPingResponse) ProtoMessage()    {}
func (*SetPingResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_a6083bb8799e537c, []int{2}
}

func (m *SetPingResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SetPingResponse.Unmarshal(m, b)
}
func (m *SetPingResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SetPingResponse.Marshal(b, m, deterministic)
}
func (m *SetPingResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SetPingResponse.Merge(m, src)
}
func (m *SetPingResponse) XXX_Size() int {
	return xxx_messageInfo_SetPingResponse.Size(m)
}
func (m *SetPingResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_SetPingResponse.DiscardUnknown(m)
}

var xxx_messageInfo_SetPingResponse proto.InternalMessageInfo

func (m *SetPingResponse) GetPing() *Ping {
	if m != nil {
		return m.Ping
	}
	return nil
}

type GetPingRequest struct {
	Key                  string   `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetPingRequest) Reset()         { *m = GetPingRequest{} }
func (m *GetPingRequest) String() string { return proto.CompactTextString(m) }
func (*GetPingRequest) ProtoMessage()    {}
func (*GetPingRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_a6083bb8799e537c, []int{3}
}

func (m *GetPingRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetPingRequest.Unmarshal(m, b)
}
func (m *GetPingRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetPingRequest.Marshal(b, m, deterministic)
}
func (m *GetPingRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetPingRequest.Merge(m, src)
}
func (m *GetPingRequest) XXX_Size() int {
	return xxx_messageInfo_GetPingRequest.Size(m)
}
func (m *GetPingRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetPingRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetPingRequest proto.InternalMessageInfo

func (m *GetPingRequest) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

type GetPingResponse struct {
	Ping                 *Ping    `protobuf:"bytes,1,opt,name=ping,proto3" json:"ping,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetPingResponse) Reset()         { *m = GetPingResponse{} }
func (m *GetPingResponse) String() string { return proto.CompactTextString(m) }
func (*GetPingResponse) ProtoMessage()    {}
func (*GetPingResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_a6083bb8799e537c, []int{4}
}

func (m *GetPingResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetPingResponse.Unmarshal(m, b)
}
func (m *GetPingResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetPingResponse.Marshal(b, m, deterministic)
}
func (m *GetPingResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetPingResponse.Merge(m, src)
}
func (m *GetPingResponse) XXX_Size() int {
	return xxx_messageInfo_GetPingResponse.Size(m)
}
func (m *GetPingResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetPingResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetPingResponse proto.InternalMessageInfo

func (m *GetPingResponse) GetPing() *Ping {
	if m != nil {
		return m.Ping
	}
	return nil
}

type RegisterCompanyRequest struct {
	CompanyName          string   `protobuf:"bytes,1,opt,name=company_name,json=companyName,proto3" json:"company_name,omitempty"`
	Email                string   `protobuf:"bytes,2,opt,name=email,proto3" json:"email,omitempty"`
	Phone                string   `protobuf:"bytes,3,opt,name=phone,proto3" json:"phone,omitempty"`
	Password             string   `protobuf:"bytes,4,opt,name=password,proto3" json:"password,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RegisterCompanyRequest) Reset()         { *m = RegisterCompanyRequest{} }
func (m *RegisterCompanyRequest) String() string { return proto.CompactTextString(m) }
func (*RegisterCompanyRequest) ProtoMessage()    {}
func (*RegisterCompanyRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_a6083bb8799e537c, []int{5}
}

func (m *RegisterCompanyRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RegisterCompanyRequest.Unmarshal(m, b)
}
func (m *RegisterCompanyRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RegisterCompanyRequest.Marshal(b, m, deterministic)
}
func (m *RegisterCompanyRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RegisterCompanyRequest.Merge(m, src)
}
func (m *RegisterCompanyRequest) XXX_Size() int {
	return xxx_messageInfo_RegisterCompanyRequest.Size(m)
}
func (m *RegisterCompanyRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RegisterCompanyRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RegisterCompanyRequest proto.InternalMessageInfo

func (m *RegisterCompanyRequest) GetCompanyName() string {
	if m != nil {
		return m.CompanyName
	}
	return ""
}

func (m *RegisterCompanyRequest) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *RegisterCompanyRequest) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

func (m *RegisterCompanyRequest) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

type RegisterCompanyResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RegisterCompanyResponse) Reset()         { *m = RegisterCompanyResponse{} }
func (m *RegisterCompanyResponse) String() string { return proto.CompactTextString(m) }
func (*RegisterCompanyResponse) ProtoMessage()    {}
func (*RegisterCompanyResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_a6083bb8799e537c, []int{6}
}

func (m *RegisterCompanyResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RegisterCompanyResponse.Unmarshal(m, b)
}
func (m *RegisterCompanyResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RegisterCompanyResponse.Marshal(b, m, deterministic)
}
func (m *RegisterCompanyResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RegisterCompanyResponse.Merge(m, src)
}
func (m *RegisterCompanyResponse) XXX_Size() int {
	return xxx_messageInfo_RegisterCompanyResponse.Size(m)
}
func (m *RegisterCompanyResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_RegisterCompanyResponse.DiscardUnknown(m)
}

var xxx_messageInfo_RegisterCompanyResponse proto.InternalMessageInfo

type RegisterPersonaRequest struct {
	FirstName            string   `protobuf:"bytes,1,opt,name=firstName,proto3" json:"firstName,omitempty"`
	LastName             string   `protobuf:"bytes,2,opt,name=lastName,proto3" json:"lastName,omitempty"`
	Email                string   `protobuf:"bytes,3,opt,name=email,proto3" json:"email,omitempty"`
	Phone                string   `protobuf:"bytes,4,opt,name=phone,proto3" json:"phone,omitempty"`
	Password             string   `protobuf:"bytes,5,opt,name=password,proto3" json:"password,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RegisterPersonaRequest) Reset()         { *m = RegisterPersonaRequest{} }
func (m *RegisterPersonaRequest) String() string { return proto.CompactTextString(m) }
func (*RegisterPersonaRequest) ProtoMessage()    {}
func (*RegisterPersonaRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_a6083bb8799e537c, []int{7}
}

func (m *RegisterPersonaRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RegisterPersonaRequest.Unmarshal(m, b)
}
func (m *RegisterPersonaRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RegisterPersonaRequest.Marshal(b, m, deterministic)
}
func (m *RegisterPersonaRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RegisterPersonaRequest.Merge(m, src)
}
func (m *RegisterPersonaRequest) XXX_Size() int {
	return xxx_messageInfo_RegisterPersonaRequest.Size(m)
}
func (m *RegisterPersonaRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RegisterPersonaRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RegisterPersonaRequest proto.InternalMessageInfo

func (m *RegisterPersonaRequest) GetFirstName() string {
	if m != nil {
		return m.FirstName
	}
	return ""
}

func (m *RegisterPersonaRequest) GetLastName() string {
	if m != nil {
		return m.LastName
	}
	return ""
}

func (m *RegisterPersonaRequest) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *RegisterPersonaRequest) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

func (m *RegisterPersonaRequest) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

type RegisterPersonaResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RegisterPersonaResponse) Reset()         { *m = RegisterPersonaResponse{} }
func (m *RegisterPersonaResponse) String() string { return proto.CompactTextString(m) }
func (*RegisterPersonaResponse) ProtoMessage()    {}
func (*RegisterPersonaResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_a6083bb8799e537c, []int{8}
}

func (m *RegisterPersonaResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RegisterPersonaResponse.Unmarshal(m, b)
}
func (m *RegisterPersonaResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RegisterPersonaResponse.Marshal(b, m, deterministic)
}
func (m *RegisterPersonaResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RegisterPersonaResponse.Merge(m, src)
}
func (m *RegisterPersonaResponse) XXX_Size() int {
	return xxx_messageInfo_RegisterPersonaResponse.Size(m)
}
func (m *RegisterPersonaResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_RegisterPersonaResponse.DiscardUnknown(m)
}

var xxx_messageInfo_RegisterPersonaResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*Ping)(nil), "personaappapi.Ping")
	proto.RegisterType((*SetPingRequest)(nil), "personaappapi.SetPingRequest")
	proto.RegisterType((*SetPingResponse)(nil), "personaappapi.SetPingResponse")
	proto.RegisterType((*GetPingRequest)(nil), "personaappapi.GetPingRequest")
	proto.RegisterType((*GetPingResponse)(nil), "personaappapi.GetPingResponse")
	proto.RegisterType((*RegisterCompanyRequest)(nil), "personaappapi.RegisterCompanyRequest")
	proto.RegisterType((*RegisterCompanyResponse)(nil), "personaappapi.RegisterCompanyResponse")
	proto.RegisterType((*RegisterPersonaRequest)(nil), "personaappapi.RegisterPersonaRequest")
	proto.RegisterType((*RegisterPersonaResponse)(nil), "personaappapi.RegisterPersonaResponse")
}

func init() { proto.RegisterFile("personaapp.proto", fileDescriptor_a6083bb8799e537c) }

var fileDescriptor_a6083bb8799e537c = []byte{
	// 432 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x93, 0xcd, 0x8e, 0xd3, 0x30,
	0x10, 0xc7, 0x95, 0x6d, 0x96, 0x25, 0x53, 0xd8, 0x5d, 0x19, 0x04, 0x25, 0xe2, 0xa3, 0x44, 0x02,
	0x7a, 0x4a, 0xa5, 0x72, 0x01, 0x6e, 0x15, 0x87, 0xf4, 0x54, 0x55, 0x81, 0x7b, 0x71, 0xdb, 0x69,
	0x88, 0x48, 0x62, 0x13, 0x3b, 0xa0, 0x9e, 0x79, 0x0b, 0xce, 0x3c, 0x1a, 0x0f, 0x82, 0x62, 0x3b,
	0x09, 0x4d, 0x1b, 0x8a, 0xf6, 0x96, 0x19, 0xcf, 0xcc, 0xff, 0x37, 0x1f, 0x81, 0x6b, 0x8e, 0xb9,
	0x60, 0x19, 0xa5, 0x9c, 0xfb, 0x3c, 0x67, 0x92, 0x91, 0xbb, 0x8d, 0x87, 0xf2, 0xd8, 0x7d, 0x16,
	0x31, 0x16, 0x25, 0x38, 0x56, 0x8f, 0xab, 0x62, 0x3b, 0x96, 0x71, 0x8a, 0x42, 0xd2, 0xd4, 0xc4,
	0x7b, 0xbf, 0x2c, 0xb0, 0x17, 0x71, 0x16, 0x91, 0x6b, 0xe8, 0x7d, 0xc1, 0xdd, 0xc0, 0x1a, 0x5a,
	0x23, 0x27, 0x2c, 0x3f, 0xc9, 0x7d, 0x38, 0xff, 0x46, 0x93, 0x02, 0x07, 0x67, 0xca, 0xa7, 0x0d,
	0xf2, 0x16, 0x60, 0x9d, 0x23, 0x95, 0xb8, 0x59, 0x52, 0x39, 0xe8, 0x0d, 0xad, 0x51, 0x7f, 0xe2,
	0xfa, 0x5a, 0xc6, 0xaf, 0x64, 0xfc, 0x8f, 0x95, 0x4c, 0xe8, 0x98, 0xe8, 0xa9, 0x2c, 0x53, 0x0b,
	0xbe, 0xa9, 0x52, 0xed, 0xd3, 0xa9, 0x26, 0x7a, 0x2a, 0xbd, 0x37, 0x70, 0xf9, 0x01, 0x65, 0x09,
	0x1a, 0xe2, 0xd7, 0x02, 0x85, 0xfc, 0x5f, 0x5e, 0xef, 0x1d, 0x5c, 0xd5, 0x99, 0x82, 0xb3, 0x4c,
	0x20, 0x79, 0x05, 0x36, 0x8f, 0xb3, 0x48, 0xe5, 0xf6, 0x27, 0xf7, 0xfc, 0xbd, 0x91, 0xf9, 0x2a,
	0x54, 0x05, 0x78, 0x1e, 0x5c, 0x06, 0x27, 0x54, 0xcb, 0xfa, 0xc1, 0x4d, 0xeb, 0xff, 0xb0, 0xe0,
	0x41, 0x88, 0x51, 0x2c, 0x24, 0xe6, 0xef, 0x59, 0xca, 0x69, 0xb6, 0xab, 0x84, 0x9e, 0xc3, 0x9d,
	0xb5, 0xf6, 0x2c, 0x33, 0x9a, 0xa2, 0x51, 0xec, 0x1b, 0xdf, 0x9c, 0xa6, 0x58, 0xf6, 0x8b, 0x29,
	0x8d, 0x93, 0xaa, 0x5f, 0x65, 0x94, 0x5e, 0xfe, 0x99, 0x65, 0xa8, 0x56, 0xe3, 0x84, 0xda, 0x20,
	0x2e, 0xdc, 0xe6, 0x54, 0x88, 0xef, 0x2c, 0xdf, 0xa8, 0xc1, 0x3b, 0x61, 0x6d, 0x7b, 0x8f, 0xe0,
	0xe1, 0x01, 0x84, 0xee, 0xc4, 0xfb, 0xf9, 0x17, 0xe0, 0x42, 0x77, 0x51, 0x01, 0x3e, 0x06, 0x67,
	0x1b, 0xe7, 0x42, 0xce, 0x1b, 0xba, 0xc6, 0x51, 0xea, 0x25, 0xd4, 0x3c, 0x6a, 0xbc, 0xda, 0x6e,
	0xb8, 0x7b, 0x47, 0xb9, 0xed, 0x2e, 0xee, 0xf3, 0x6e, 0xee, 0x9a, 0x4d, 0x73, 0x4f, 0x7e, 0x9f,
	0x01, 0x18, 0xdf, 0x94, 0x73, 0x32, 0x83, 0x0b, 0x73, 0x03, 0xe4, 0x49, 0x6b, 0x1b, 0xfb, 0x57,
	0xe5, 0x3e, 0xed, 0x7a, 0x36, 0xab, 0x9d, 0xc1, 0x45, 0xd0, 0x51, 0x29, 0xf8, 0x77, 0xa5, 0xf6,
	0x91, 0x7c, 0x82, 0xab, 0xd6, 0xd4, 0xc9, 0x8b, 0x56, 0xca, 0xf1, 0xd3, 0x70, 0x5f, 0x9e, 0x0a,
	0x3b, 0x54, 0x30, 0xb3, 0xe8, 0x54, 0xd8, 0xdf, 0x6d, 0xa7, 0x42, 0x6b, 0xcc, 0xab, 0x5b, 0xea,
	0xa7, 0x7d, 0xfd, 0x27, 0x00, 0x00, 0xff, 0xff, 0x45, 0x2f, 0x3c, 0x18, 0x87, 0x04, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// PersonaAppClient is the client API for PersonaApp service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type PersonaAppClient interface {
	SetPing(ctx context.Context, in *SetPingRequest, opts ...grpc.CallOption) (*SetPingResponse, error)
	GetPing(ctx context.Context, in *GetPingRequest, opts ...grpc.CallOption) (*GetPingResponse, error)
	RegisterCompany(ctx context.Context, in *RegisterCompanyRequest, opts ...grpc.CallOption) (*RegisterCompanyResponse, error)
	RegisterPersona(ctx context.Context, in *RegisterPersonaRequest, opts ...grpc.CallOption) (*RegisterPersonaResponse, error)
}

type personaAppClient struct {
	cc *grpc.ClientConn
}

func NewPersonaAppClient(cc *grpc.ClientConn) PersonaAppClient {
	return &personaAppClient{cc}
}

func (c *personaAppClient) SetPing(ctx context.Context, in *SetPingRequest, opts ...grpc.CallOption) (*SetPingResponse, error) {
	out := new(SetPingResponse)
	err := c.cc.Invoke(ctx, "/personaappapi.PersonaApp/SetPing", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *personaAppClient) GetPing(ctx context.Context, in *GetPingRequest, opts ...grpc.CallOption) (*GetPingResponse, error) {
	out := new(GetPingResponse)
	err := c.cc.Invoke(ctx, "/personaappapi.PersonaApp/GetPing", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *personaAppClient) RegisterCompany(ctx context.Context, in *RegisterCompanyRequest, opts ...grpc.CallOption) (*RegisterCompanyResponse, error) {
	out := new(RegisterCompanyResponse)
	err := c.cc.Invoke(ctx, "/personaappapi.PersonaApp/RegisterCompany", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *personaAppClient) RegisterPersona(ctx context.Context, in *RegisterPersonaRequest, opts ...grpc.CallOption) (*RegisterPersonaResponse, error) {
	out := new(RegisterPersonaResponse)
	err := c.cc.Invoke(ctx, "/personaappapi.PersonaApp/RegisterPersona", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PersonaAppServer is the server API for PersonaApp service.
type PersonaAppServer interface {
	SetPing(context.Context, *SetPingRequest) (*SetPingResponse, error)
	GetPing(context.Context, *GetPingRequest) (*GetPingResponse, error)
	RegisterCompany(context.Context, *RegisterCompanyRequest) (*RegisterCompanyResponse, error)
	RegisterPersona(context.Context, *RegisterPersonaRequest) (*RegisterPersonaResponse, error)
}

// UnimplementedPersonaAppServer can be embedded to have forward compatible implementations.
type UnimplementedPersonaAppServer struct {
}

func (*UnimplementedPersonaAppServer) SetPing(ctx context.Context, req *SetPingRequest) (*SetPingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetPing not implemented")
}
func (*UnimplementedPersonaAppServer) GetPing(ctx context.Context, req *GetPingRequest) (*GetPingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPing not implemented")
}
func (*UnimplementedPersonaAppServer) RegisterCompany(ctx context.Context, req *RegisterCompanyRequest) (*RegisterCompanyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterCompany not implemented")
}
func (*UnimplementedPersonaAppServer) RegisterPersona(ctx context.Context, req *RegisterPersonaRequest) (*RegisterPersonaResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterPersona not implemented")
}

func RegisterPersonaAppServer(s *grpc.Server, srv PersonaAppServer) {
	s.RegisterService(&_PersonaApp_serviceDesc, srv)
}

func _PersonaApp_SetPing_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetPingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PersonaAppServer).SetPing(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/personaappapi.PersonaApp/SetPing",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PersonaAppServer).SetPing(ctx, req.(*SetPingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PersonaApp_GetPing_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PersonaAppServer).GetPing(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/personaappapi.PersonaApp/GetPing",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PersonaAppServer).GetPing(ctx, req.(*GetPingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PersonaApp_RegisterCompany_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterCompanyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PersonaAppServer).RegisterCompany(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/personaappapi.PersonaApp/RegisterCompany",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PersonaAppServer).RegisterCompany(ctx, req.(*RegisterCompanyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PersonaApp_RegisterPersona_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterPersonaRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PersonaAppServer).RegisterPersona(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/personaappapi.PersonaApp/RegisterPersona",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PersonaAppServer).RegisterPersona(ctx, req.(*RegisterPersonaRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _PersonaApp_serviceDesc = grpc.ServiceDesc{
	ServiceName: "personaappapi.PersonaApp",
	HandlerType: (*PersonaAppServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SetPing",
			Handler:    _PersonaApp_SetPing_Handler,
		},
		{
			MethodName: "GetPing",
			Handler:    _PersonaApp_GetPing_Handler,
		},
		{
			MethodName: "RegisterCompany",
			Handler:    _PersonaApp_RegisterCompany_Handler,
		},
		{
			MethodName: "RegisterPersona",
			Handler:    _PersonaApp_RegisterPersona_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "personaapp.proto",
}