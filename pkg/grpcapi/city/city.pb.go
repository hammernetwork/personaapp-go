// Code generated by protoc-gen-go. DO NOT EDIT.
// source: city/city.proto

package personaappapi_city

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	wrappers "github.com/golang/protobuf/ptypes/wrappers"
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

// Get cities
type GetCitiesRequest struct {
	Rating               *wrappers.Int32Value  `protobuf:"bytes,1,opt,name=rating,proto3" json:"rating,omitempty"`
	Filter               *wrappers.StringValue `protobuf:"bytes,2,opt,name=filter,proto3" json:"filter,omitempty"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *GetCitiesRequest) Reset()         { *m = GetCitiesRequest{} }
func (m *GetCitiesRequest) String() string { return proto.CompactTextString(m) }
func (*GetCitiesRequest) ProtoMessage()    {}
func (*GetCitiesRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_0fe24b41a967713c, []int{0}
}

func (m *GetCitiesRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetCitiesRequest.Unmarshal(m, b)
}
func (m *GetCitiesRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetCitiesRequest.Marshal(b, m, deterministic)
}
func (m *GetCitiesRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetCitiesRequest.Merge(m, src)
}
func (m *GetCitiesRequest) XXX_Size() int {
	return xxx_messageInfo_GetCitiesRequest.Size(m)
}
func (m *GetCitiesRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetCitiesRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetCitiesRequest proto.InternalMessageInfo

func (m *GetCitiesRequest) GetRating() *wrappers.Int32Value {
	if m != nil {
		return m.Rating
	}
	return nil
}

func (m *GetCitiesRequest) GetFilter() *wrappers.StringValue {
	if m != nil {
		return m.Filter
	}
	return nil
}

type GetCitiesResponse struct {
	Cities               []*City  `protobuf:"bytes,1,rep,name=cities,proto3" json:"cities,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetCitiesResponse) Reset()         { *m = GetCitiesResponse{} }
func (m *GetCitiesResponse) String() string { return proto.CompactTextString(m) }
func (*GetCitiesResponse) ProtoMessage()    {}
func (*GetCitiesResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_0fe24b41a967713c, []int{1}
}

func (m *GetCitiesResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetCitiesResponse.Unmarshal(m, b)
}
func (m *GetCitiesResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetCitiesResponse.Marshal(b, m, deterministic)
}
func (m *GetCitiesResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetCitiesResponse.Merge(m, src)
}
func (m *GetCitiesResponse) XXX_Size() int {
	return xxx_messageInfo_GetCitiesResponse.Size(m)
}
func (m *GetCitiesResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetCitiesResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetCitiesResponse proto.InternalMessageInfo

func (m *GetCitiesResponse) GetCities() []*City {
	if m != nil {
		return m.Cities
	}
	return nil
}

// Upsert city
type UpsertCityRequest struct {
	Id                   *wrappers.StringValue `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name                 string                `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	CountryCode          int32                 `protobuf:"varint,3,opt,name=country_code,json=countryCode,proto3" json:"country_code,omitempty"`
	Rating               int32                 `protobuf:"varint,4,opt,name=rating,proto3" json:"rating,omitempty"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *UpsertCityRequest) Reset()         { *m = UpsertCityRequest{} }
func (m *UpsertCityRequest) String() string { return proto.CompactTextString(m) }
func (*UpsertCityRequest) ProtoMessage()    {}
func (*UpsertCityRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_0fe24b41a967713c, []int{2}
}

func (m *UpsertCityRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UpsertCityRequest.Unmarshal(m, b)
}
func (m *UpsertCityRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UpsertCityRequest.Marshal(b, m, deterministic)
}
func (m *UpsertCityRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpsertCityRequest.Merge(m, src)
}
func (m *UpsertCityRequest) XXX_Size() int {
	return xxx_messageInfo_UpsertCityRequest.Size(m)
}
func (m *UpsertCityRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_UpsertCityRequest.DiscardUnknown(m)
}

var xxx_messageInfo_UpsertCityRequest proto.InternalMessageInfo

func (m *UpsertCityRequest) GetId() *wrappers.StringValue {
	if m != nil {
		return m.Id
	}
	return nil
}

func (m *UpsertCityRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *UpsertCityRequest) GetCountryCode() int32 {
	if m != nil {
		return m.CountryCode
	}
	return 0
}

func (m *UpsertCityRequest) GetRating() int32 {
	if m != nil {
		return m.Rating
	}
	return 0
}

type UpsertCityResponse struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UpsertCityResponse) Reset()         { *m = UpsertCityResponse{} }
func (m *UpsertCityResponse) String() string { return proto.CompactTextString(m) }
func (*UpsertCityResponse) ProtoMessage()    {}
func (*UpsertCityResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_0fe24b41a967713c, []int{3}
}

func (m *UpsertCityResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UpsertCityResponse.Unmarshal(m, b)
}
func (m *UpsertCityResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UpsertCityResponse.Marshal(b, m, deterministic)
}
func (m *UpsertCityResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpsertCityResponse.Merge(m, src)
}
func (m *UpsertCityResponse) XXX_Size() int {
	return xxx_messageInfo_UpsertCityResponse.Size(m)
}
func (m *UpsertCityResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_UpsertCityResponse.DiscardUnknown(m)
}

var xxx_messageInfo_UpsertCityResponse proto.InternalMessageInfo

func (m *UpsertCityResponse) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

// Delete city
type DeleteCityRequest struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DeleteCityRequest) Reset()         { *m = DeleteCityRequest{} }
func (m *DeleteCityRequest) String() string { return proto.CompactTextString(m) }
func (*DeleteCityRequest) ProtoMessage()    {}
func (*DeleteCityRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_0fe24b41a967713c, []int{4}
}

func (m *DeleteCityRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeleteCityRequest.Unmarshal(m, b)
}
func (m *DeleteCityRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeleteCityRequest.Marshal(b, m, deterministic)
}
func (m *DeleteCityRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeleteCityRequest.Merge(m, src)
}
func (m *DeleteCityRequest) XXX_Size() int {
	return xxx_messageInfo_DeleteCityRequest.Size(m)
}
func (m *DeleteCityRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_DeleteCityRequest.DiscardUnknown(m)
}

var xxx_messageInfo_DeleteCityRequest proto.InternalMessageInfo

func (m *DeleteCityRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type DeleteCityResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DeleteCityResponse) Reset()         { *m = DeleteCityResponse{} }
func (m *DeleteCityResponse) String() string { return proto.CompactTextString(m) }
func (*DeleteCityResponse) ProtoMessage()    {}
func (*DeleteCityResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_0fe24b41a967713c, []int{5}
}

func (m *DeleteCityResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeleteCityResponse.Unmarshal(m, b)
}
func (m *DeleteCityResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeleteCityResponse.Marshal(b, m, deterministic)
}
func (m *DeleteCityResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeleteCityResponse.Merge(m, src)
}
func (m *DeleteCityResponse) XXX_Size() int {
	return xxx_messageInfo_DeleteCityResponse.Size(m)
}
func (m *DeleteCityResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_DeleteCityResponse.DiscardUnknown(m)
}

var xxx_messageInfo_DeleteCityResponse proto.InternalMessageInfo

type City struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name                 string   `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	CountryCode          int32    `protobuf:"varint,3,opt,name=country_code,json=countryCode,proto3" json:"country_code,omitempty"`
	Rating               int32    `protobuf:"varint,4,opt,name=rating,proto3" json:"rating,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *City) Reset()         { *m = City{} }
func (m *City) String() string { return proto.CompactTextString(m) }
func (*City) ProtoMessage()    {}
func (*City) Descriptor() ([]byte, []int) {
	return fileDescriptor_0fe24b41a967713c, []int{6}
}

func (m *City) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_City.Unmarshal(m, b)
}
func (m *City) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_City.Marshal(b, m, deterministic)
}
func (m *City) XXX_Merge(src proto.Message) {
	xxx_messageInfo_City.Merge(m, src)
}
func (m *City) XXX_Size() int {
	return xxx_messageInfo_City.Size(m)
}
func (m *City) XXX_DiscardUnknown() {
	xxx_messageInfo_City.DiscardUnknown(m)
}

var xxx_messageInfo_City proto.InternalMessageInfo

func (m *City) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *City) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *City) GetCountryCode() int32 {
	if m != nil {
		return m.CountryCode
	}
	return 0
}

func (m *City) GetRating() int32 {
	if m != nil {
		return m.Rating
	}
	return 0
}

func init() {
	proto.RegisterType((*GetCitiesRequest)(nil), "personaappapi.city.GetCitiesRequest")
	proto.RegisterType((*GetCitiesResponse)(nil), "personaappapi.city.GetCitiesResponse")
	proto.RegisterType((*UpsertCityRequest)(nil), "personaappapi.city.UpsertCityRequest")
	proto.RegisterType((*UpsertCityResponse)(nil), "personaappapi.city.UpsertCityResponse")
	proto.RegisterType((*DeleteCityRequest)(nil), "personaappapi.city.DeleteCityRequest")
	proto.RegisterType((*DeleteCityResponse)(nil), "personaappapi.city.DeleteCityResponse")
	proto.RegisterType((*City)(nil), "personaappapi.city.City")
}

func init() { proto.RegisterFile("city/city.proto", fileDescriptor_0fe24b41a967713c) }

var fileDescriptor_0fe24b41a967713c = []byte{
	// 402 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x92, 0xdd, 0x6a, 0xdb, 0x30,
	0x14, 0xc7, 0xb1, 0x93, 0x99, 0xe5, 0x64, 0x64, 0xb3, 0x18, 0xc3, 0x64, 0x1f, 0x64, 0x5e, 0x36,
	0x72, 0x31, 0x9c, 0xe1, 0xec, 0x05, 0x96, 0x6c, 0x84, 0xdd, 0x0d, 0x8f, 0x8d, 0x42, 0x2f, 0x8a,
	0x63, 0x9f, 0x18, 0x81, 0x2b, 0xa9, 0xb2, 0x4c, 0xc9, 0x45, 0xdf, 0xa1, 0x0f, 0xd1, 0x07, 0x2d,
	0x96, 0xd5, 0xc4, 0xc4, 0xa6, 0xb9, 0xe9, 0x8d, 0xb1, 0x8e, 0x7e, 0xe7, 0x7f, 0x3e, 0xfe, 0x82,
	0x97, 0x09, 0x55, 0xbb, 0x79, 0xf5, 0x09, 0x84, 0xe4, 0x8a, 0x13, 0x22, 0x50, 0x16, 0x9c, 0xc5,
	0xb1, 0x10, 0xb1, 0xa0, 0x41, 0x75, 0x33, 0xfe, 0x90, 0x71, 0x9e, 0xe5, 0x38, 0xd7, 0xc4, 0xa6,
	0xdc, 0xce, 0xaf, 0x65, 0x2c, 0x2a, 0xae, 0xce, 0xf1, 0x6f, 0xe0, 0xd5, 0x1a, 0xd5, 0x8a, 0x2a,
	0x8a, 0x45, 0x84, 0x57, 0x25, 0x16, 0x8a, 0x2c, 0xc0, 0x91, 0xb1, 0xa2, 0x2c, 0xf3, 0xac, 0x89,
	0x35, 0x1b, 0x86, 0x6f, 0x83, 0x5a, 0x24, 0x78, 0x10, 0x09, 0x7e, 0x33, 0xb5, 0x08, 0xff, 0xc7,
	0x79, 0x89, 0x91, 0x41, 0xc9, 0x77, 0x70, 0xb6, 0x34, 0x57, 0x28, 0x3d, 0x5b, 0x27, 0xbd, 0x6b,
	0x25, 0xfd, 0x55, 0x92, 0xb2, 0xcc, 0x64, 0xd5, 0xac, 0xff, 0x0b, 0xdc, 0x46, 0xf9, 0x42, 0x70,
	0x56, 0x20, 0xf9, 0x06, 0x4e, 0xa2, 0x23, 0x9e, 0x35, 0xe9, 0xcd, 0x86, 0xa1, 0x17, 0xb4, 0x07,
	0x0b, 0x56, 0x54, 0xed, 0x22, 0xc3, 0xf9, 0xb7, 0x16, 0xb8, 0xff, 0x44, 0x81, 0x52, 0xe9, 0xb0,
	0x99, 0xe3, 0x2b, 0xd8, 0x34, 0x35, 0x33, 0x3c, 0xde, 0x8e, 0x4d, 0x53, 0x42, 0xa0, 0xcf, 0xe2,
	0x4b, 0xd4, 0xed, 0x0f, 0x22, 0xfd, 0x4f, 0x3e, 0xc2, 0x8b, 0x84, 0x97, 0x4c, 0xc9, 0xdd, 0x45,
	0xc2, 0x53, 0xf4, 0x7a, 0x13, 0x6b, 0xf6, 0x2c, 0x1a, 0x9a, 0xd8, 0x8a, 0xa7, 0x48, 0xde, 0xec,
	0x97, 0xd5, 0xd7, 0x97, 0xe6, 0xe4, 0x4f, 0x81, 0x34, 0x3b, 0x32, 0xa3, 0x8d, 0xf6, 0x2d, 0x0d,
	0xaa, 0xa2, 0xfe, 0x27, 0x70, 0x7f, 0x62, 0x8e, 0x0a, 0x9b, 0x7d, 0x1f, 0x43, 0xaf, 0x81, 0x34,
	0xa1, 0x5a, 0xca, 0x47, 0xe8, 0x57, 0xe7, 0x63, 0xfa, 0x89, 0xe7, 0x08, 0xef, 0x6c, 0x18, 0xfd,
	0xa9, 0xd7, 0xff, 0x43, 0x08, 0x5d, 0xf1, 0x0c, 0x06, 0x7b, 0xd3, 0xc8, 0xb4, 0xcb, 0x9c, 0xe3,
	0x27, 0x35, 0xfe, 0x7c, 0x82, 0x32, 0xeb, 0x39, 0x07, 0x38, 0x2c, 0x8d, 0x74, 0x26, 0xb5, 0x6c,
	0x1e, 0x7f, 0x39, 0x85, 0x1d, 0xc4, 0x0f, 0x6b, 0xec, 0x16, 0x6f, 0x79, 0xd1, 0x2d, 0xde, 0x76,
	0x63, 0xf9, 0x1e, 0x5c, 0xce, 0x72, 0xca, 0xb0, 0xc1, 0x2f, 0x9f, 0xaf, 0xa5, 0x48, 0x2a, 0x6c,
	0xe3, 0xe8, 0x67, 0xb7, 0xb8, 0x0f, 0x00, 0x00, 0xff, 0xff, 0xad, 0x33, 0x89, 0xdf, 0xb4, 0x03,
	0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// PersonaAppCityClient is the client API for PersonaAppCity service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type PersonaAppCityClient interface {
	// Cities
	GetCities(ctx context.Context, in *GetCitiesRequest, opts ...grpc.CallOption) (*GetCitiesResponse, error)
	UpsertCity(ctx context.Context, in *UpsertCityRequest, opts ...grpc.CallOption) (*UpsertCityResponse, error)
	DeleteCity(ctx context.Context, in *DeleteCityRequest, opts ...grpc.CallOption) (*DeleteCityResponse, error)
}

type personaAppCityClient struct {
	cc *grpc.ClientConn
}

func NewPersonaAppCityClient(cc *grpc.ClientConn) PersonaAppCityClient {
	return &personaAppCityClient{cc}
}

func (c *personaAppCityClient) GetCities(ctx context.Context, in *GetCitiesRequest, opts ...grpc.CallOption) (*GetCitiesResponse, error) {
	out := new(GetCitiesResponse)
	err := c.cc.Invoke(ctx, "/personaappapi.city.PersonaAppCity/GetCities", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *personaAppCityClient) UpsertCity(ctx context.Context, in *UpsertCityRequest, opts ...grpc.CallOption) (*UpsertCityResponse, error) {
	out := new(UpsertCityResponse)
	err := c.cc.Invoke(ctx, "/personaappapi.city.PersonaAppCity/UpsertCity", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *personaAppCityClient) DeleteCity(ctx context.Context, in *DeleteCityRequest, opts ...grpc.CallOption) (*DeleteCityResponse, error) {
	out := new(DeleteCityResponse)
	err := c.cc.Invoke(ctx, "/personaappapi.city.PersonaAppCity/DeleteCity", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PersonaAppCityServer is the server API for PersonaAppCity service.
type PersonaAppCityServer interface {
	// Cities
	GetCities(context.Context, *GetCitiesRequest) (*GetCitiesResponse, error)
	UpsertCity(context.Context, *UpsertCityRequest) (*UpsertCityResponse, error)
	DeleteCity(context.Context, *DeleteCityRequest) (*DeleteCityResponse, error)
}

// UnimplementedPersonaAppCityServer can be embedded to have forward compatible implementations.
type UnimplementedPersonaAppCityServer struct {
}

func (*UnimplementedPersonaAppCityServer) GetCities(ctx context.Context, req *GetCitiesRequest) (*GetCitiesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCities not implemented")
}
func (*UnimplementedPersonaAppCityServer) UpsertCity(ctx context.Context, req *UpsertCityRequest) (*UpsertCityResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpsertCity not implemented")
}
func (*UnimplementedPersonaAppCityServer) DeleteCity(ctx context.Context, req *DeleteCityRequest) (*DeleteCityResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteCity not implemented")
}

func RegisterPersonaAppCityServer(s *grpc.Server, srv PersonaAppCityServer) {
	s.RegisterService(&_PersonaAppCity_serviceDesc, srv)
}

func _PersonaAppCity_GetCities_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCitiesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PersonaAppCityServer).GetCities(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/personaappapi.city.PersonaAppCity/GetCities",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PersonaAppCityServer).GetCities(ctx, req.(*GetCitiesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PersonaAppCity_UpsertCity_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpsertCityRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PersonaAppCityServer).UpsertCity(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/personaappapi.city.PersonaAppCity/UpsertCity",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PersonaAppCityServer).UpsertCity(ctx, req.(*UpsertCityRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PersonaAppCity_DeleteCity_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteCityRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PersonaAppCityServer).DeleteCity(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/personaappapi.city.PersonaAppCity/DeleteCity",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PersonaAppCityServer).DeleteCity(ctx, req.(*DeleteCityRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _PersonaAppCity_serviceDesc = grpc.ServiceDesc{
	ServiceName: "personaappapi.city.PersonaAppCity",
	HandlerType: (*PersonaAppCityServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetCities",
			Handler:    _PersonaAppCity_GetCities_Handler,
		},
		{
			MethodName: "UpsertCity",
			Handler:    _PersonaAppCity_UpsertCity_Handler,
		},
		{
			MethodName: "DeleteCity",
			Handler:    _PersonaAppCity_DeleteCity_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "city/city.proto",
}
