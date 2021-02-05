//指定proto3版本，默认是proto2

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.23.0
// 	protoc        v3.13.0
// source: string.proto

//指定包名

package pb

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

//定义类型
type StringRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	A string `protobuf:"bytes,1,opt,name=A,proto3" json:"A,omitempty"`
	B string `protobuf:"bytes,2,opt,name=B,proto3" json:"B,omitempty"`
}

func (x *StringRequest) Reset() {
	*x = StringRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_string_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StringRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StringRequest) ProtoMessage() {}

func (x *StringRequest) ProtoReflect() protoreflect.Message {
	mi := &file_string_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StringRequest.ProtoReflect.Descriptor instead.
func (*StringRequest) Descriptor() ([]byte, []int) {
	return file_string_proto_rawDescGZIP(), []int{0}
}

func (x *StringRequest) GetA() string {
	if x != nil {
		return x.A
	}
	return ""
}

func (x *StringRequest) GetB() string {
	if x != nil {
		return x.B
	}
	return ""
}

type StringResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Result string `protobuf:"bytes,1,opt,name=Result,proto3" json:"Result,omitempty"`
}

func (x *StringResponse) Reset() {
	*x = StringResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_string_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StringResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StringResponse) ProtoMessage() {}

func (x *StringResponse) ProtoReflect() protoreflect.Message {
	mi := &file_string_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StringResponse.ProtoReflect.Descriptor instead.
func (*StringResponse) Descriptor() ([]byte, []int) {
	return file_string_proto_rawDescGZIP(), []int{1}
}

func (x *StringResponse) GetResult() string {
	if x != nil {
		return x.Result
	}
	return ""
}

var File_string_proto protoreflect.FileDescriptor

var file_string_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02,
	0x70, 0x62, 0x22, 0x2b, 0x0a, 0x0d, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x0c, 0x0a, 0x01, 0x41, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x01,
	0x41, 0x12, 0x0c, 0x0a, 0x01, 0x42, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x01, 0x42, 0x22,
	0x28, 0x0a, 0x0e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x16, 0x0a, 0x06, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x32, 0x6c, 0x0a, 0x06, 0x53, 0x74, 0x72,
	0x69, 0x6e, 0x67, 0x12, 0x31, 0x0a, 0x06, 0x43, 0x6f, 0x6e, 0x63, 0x61, 0x74, 0x12, 0x11, 0x2e,
	0x70, 0x62, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x12, 0x2e, 0x70, 0x62, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x2f, 0x0a, 0x04, 0x44, 0x69, 0x66, 0x66, 0x12, 0x11,
	0x2e, 0x70, 0x62, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x12, 0x2e, 0x70, 0x62, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_string_proto_rawDescOnce sync.Once
	file_string_proto_rawDescData = file_string_proto_rawDesc
)

func file_string_proto_rawDescGZIP() []byte {
	file_string_proto_rawDescOnce.Do(func() {
		file_string_proto_rawDescData = protoimpl.X.CompressGZIP(file_string_proto_rawDescData)
	})
	return file_string_proto_rawDescData
}

var file_string_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_string_proto_goTypes = []interface{}{
	(*StringRequest)(nil),  // 0: pb.StringRequest
	(*StringResponse)(nil), // 1: pb.StringResponse
}
var file_string_proto_depIdxs = []int32{
	0, // 0: pb.String.Concat:input_type -> pb.StringRequest
	0, // 1: pb.String.Diff:input_type -> pb.StringRequest
	1, // 2: pb.String.Concat:output_type -> pb.StringResponse
	1, // 3: pb.String.Diff:output_type -> pb.StringResponse
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_string_proto_init() }
func file_string_proto_init() {
	if File_string_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_string_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StringRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_string_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StringResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_string_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_string_proto_goTypes,
		DependencyIndexes: file_string_proto_depIdxs,
		MessageInfos:      file_string_proto_msgTypes,
	}.Build()
	File_string_proto = out.File
	file_string_proto_rawDesc = nil
	file_string_proto_goTypes = nil
	file_string_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// StringClient is the client API for String service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type StringClient interface {
	Concat(ctx context.Context, in *StringRequest, opts ...grpc.CallOption) (*StringResponse, error)
	Diff(ctx context.Context, in *StringRequest, opts ...grpc.CallOption) (*StringResponse, error)
}

type stringClient struct {
	cc grpc.ClientConnInterface
}

func NewStringClient(cc grpc.ClientConnInterface) StringClient {
	return &stringClient{cc}
}

func (c *stringClient) Concat(ctx context.Context, in *StringRequest, opts ...grpc.CallOption) (*StringResponse, error) {
	out := new(StringResponse)
	err := c.cc.Invoke(ctx, "/pb.String/Concat", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *stringClient) Diff(ctx context.Context, in *StringRequest, opts ...grpc.CallOption) (*StringResponse, error) {
	out := new(StringResponse)
	err := c.cc.Invoke(ctx, "/pb.String/Diff", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// StringServer is the server API for String service.
type StringServer interface {
	Concat(context.Context, *StringRequest) (*StringResponse, error)
	Diff(context.Context, *StringRequest) (*StringResponse, error)
}

// UnimplementedStringServer can be embedded to have forward compatible implementations.
type UnimplementedStringServer struct {
}

func (*UnimplementedStringServer) Concat(context.Context, *StringRequest) (*StringResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Concat not implemented")
}
func (*UnimplementedStringServer) Diff(context.Context, *StringRequest) (*StringResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Diff not implemented")
}

func RegisterStringServer(s *grpc.Server, srv StringServer) {
	s.RegisterService(&_String_serviceDesc, srv)
}

func _String_Concat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StringRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StringServer).Concat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.String/Concat",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StringServer).Concat(ctx, req.(*StringRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _String_Diff_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StringRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StringServer).Diff(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.String/Diff",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StringServer).Diff(ctx, req.(*StringRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _String_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pb.String",
	HandlerType: (*StringServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Concat",
			Handler:    _String_Concat_Handler,
		},
		{
			MethodName: "Diff",
			Handler:    _String_Diff_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "string.proto",
}