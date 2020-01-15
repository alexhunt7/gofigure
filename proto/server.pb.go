// Code generated by protoc-gen-go. DO NOT EDIT.
// source: server.proto

package proto

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
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

type Empty struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Empty) Reset()         { *m = Empty{} }
func (m *Empty) String() string { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()    {}
func (*Empty) Descriptor() ([]byte, []int) {
	return fileDescriptor_ad098daeda4239f7, []int{0}
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
	proto.RegisterType((*Empty)(nil), "gofigure.Empty")
}

func init() { proto.RegisterFile("server.proto", fileDescriptor_ad098daeda4239f7) }

var fileDescriptor_ad098daeda4239f7 = []byte{
	// 248 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x29, 0x4e, 0x2d, 0x2a,
	0x4b, 0x2d, 0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x48, 0xcf, 0x4f, 0xcb, 0x4c, 0x2f,
	0x2d, 0x4a, 0x95, 0xe2, 0x4a, 0xcb, 0xcc, 0x49, 0x85, 0x88, 0x4a, 0x71, 0xa5, 0x56, 0xa4, 0x26,
	0x43, 0xd8, 0x4a, 0xec, 0x5c, 0xac, 0xae, 0xb9, 0x05, 0x25, 0x95, 0x46, 0xd7, 0x98, 0xb8, 0x38,
	0xdc, 0xa1, 0xaa, 0x85, 0xec, 0xb9, 0x38, 0x5d, 0x32, 0x8b, 0x52, 0x93, 0x4b, 0xf2, 0x8b, 0x2a,
	0x85, 0x44, 0xf5, 0x60, 0xa6, 0xe8, 0xb9, 0x65, 0xe6, 0xa4, 0x06, 0xa5, 0x16, 0x96, 0xa6, 0x16,
	0x97, 0x48, 0x49, 0x22, 0x84, 0xe1, 0x6a, 0x83, 0x52, 0x8b, 0x4b, 0x73, 0x4a, 0x94, 0x18, 0x84,
	0x4c, 0xb9, 0x58, 0x40, 0x6a, 0x71, 0xe9, 0x15, 0x41, 0x17, 0x86, 0x6a, 0xb3, 0xe6, 0x62, 0x73,
	0x49, 0xcd, 0x49, 0x2d, 0x49, 0x15, 0x12, 0x47, 0x32, 0x1d, 0x2c, 0x02, 0xd3, 0x2a, 0x86, 0x29,
	0x01, 0xd5, 0x6c, 0xc4, 0xc5, 0x12, 0x5c, 0x92, 0x58, 0x22, 0x24, 0x84, 0x6a, 0x78, 0x40, 0x62,
	0x49, 0x06, 0xb2, 0x85, 0x20, 0x35, 0xc8, 0xee, 0x74, 0xad, 0x48, 0x4d, 0x46, 0x76, 0x27, 0x88,
	0x8f, 0xc5, 0x9d, 0x10, 0x61, 0xa8, 0x36, 0x2d, 0x90, 0xb6, 0xcc, 0x12, 0x21, 0x7e, 0x24, 0x79,
	0x50, 0x28, 0x4a, 0xa1, 0x0b, 0x28, 0x31, 0x38, 0xa9, 0x46, 0x29, 0xa7, 0x67, 0x96, 0x64, 0x94,
	0x26, 0xe9, 0x25, 0xe7, 0xe7, 0xea, 0x27, 0xe6, 0xa4, 0x56, 0x64, 0x94, 0xe6, 0x95, 0x98, 0xeb,
	0xc3, 0x14, 0xea, 0x83, 0x23, 0x22, 0x89, 0x0d, 0x4c, 0x19, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff,
	0x49, 0x91, 0x9c, 0x5b, 0xc1, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// GofigureClient is the client API for Gofigure service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type GofigureClient interface {
	Directory(ctx context.Context, in *FileRequest, opts ...grpc.CallOption) (*DirectoryResult, error)
	File(ctx context.Context, in *FileRequest, opts ...grpc.CallOption) (*FileResult, error)
	Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*DeleteResult, error)
	Stat(ctx context.Context, in *FilePath, opts ...grpc.CallOption) (*StatResult, error)
	Exec(ctx context.Context, in *ExecRequest, opts ...grpc.CallOption) (*ExecResult, error)
	Exit(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
}

type gofigureClient struct {
	cc *grpc.ClientConn
}

func NewGofigureClient(cc *grpc.ClientConn) GofigureClient {
	return &gofigureClient{cc}
}

func (c *gofigureClient) Directory(ctx context.Context, in *FileRequest, opts ...grpc.CallOption) (*DirectoryResult, error) {
	out := new(DirectoryResult)
	err := c.cc.Invoke(ctx, "/gofigure.Gofigure/Directory", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gofigureClient) File(ctx context.Context, in *FileRequest, opts ...grpc.CallOption) (*FileResult, error) {
	out := new(FileResult)
	err := c.cc.Invoke(ctx, "/gofigure.Gofigure/File", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gofigureClient) Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*DeleteResult, error) {
	out := new(DeleteResult)
	err := c.cc.Invoke(ctx, "/gofigure.Gofigure/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gofigureClient) Stat(ctx context.Context, in *FilePath, opts ...grpc.CallOption) (*StatResult, error) {
	out := new(StatResult)
	err := c.cc.Invoke(ctx, "/gofigure.Gofigure/Stat", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gofigureClient) Exec(ctx context.Context, in *ExecRequest, opts ...grpc.CallOption) (*ExecResult, error) {
	out := new(ExecResult)
	err := c.cc.Invoke(ctx, "/gofigure.Gofigure/Exec", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gofigureClient) Exit(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/gofigure.Gofigure/Exit", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GofigureServer is the server API for Gofigure service.
type GofigureServer interface {
	Directory(context.Context, *FileRequest) (*DirectoryResult, error)
	File(context.Context, *FileRequest) (*FileResult, error)
	Delete(context.Context, *DeleteRequest) (*DeleteResult, error)
	Stat(context.Context, *FilePath) (*StatResult, error)
	Exec(context.Context, *ExecRequest) (*ExecResult, error)
	Exit(context.Context, *Empty) (*Empty, error)
}

func RegisterGofigureServer(s *grpc.Server, srv GofigureServer) {
	s.RegisterService(&_Gofigure_serviceDesc, srv)
}

func _Gofigure_Directory_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GofigureServer).Directory(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gofigure.Gofigure/Directory",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GofigureServer).Directory(ctx, req.(*FileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Gofigure_File_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GofigureServer).File(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gofigure.Gofigure/File",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GofigureServer).File(ctx, req.(*FileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Gofigure_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GofigureServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gofigure.Gofigure/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GofigureServer).Delete(ctx, req.(*DeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Gofigure_Stat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FilePath)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GofigureServer).Stat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gofigure.Gofigure/Stat",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GofigureServer).Stat(ctx, req.(*FilePath))
	}
	return interceptor(ctx, in, info, handler)
}

func _Gofigure_Exec_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExecRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GofigureServer).Exec(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gofigure.Gofigure/Exec",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GofigureServer).Exec(ctx, req.(*ExecRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Gofigure_Exit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GofigureServer).Exit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gofigure.Gofigure/Exit",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GofigureServer).Exit(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _Gofigure_serviceDesc = grpc.ServiceDesc{
	ServiceName: "gofigure.Gofigure",
	HandlerType: (*GofigureServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Directory",
			Handler:    _Gofigure_Directory_Handler,
		},
		{
			MethodName: "File",
			Handler:    _Gofigure_File_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _Gofigure_Delete_Handler,
		},
		{
			MethodName: "Stat",
			Handler:    _Gofigure_Stat_Handler,
		},
		{
			MethodName: "Exec",
			Handler:    _Gofigure_Exec_Handler,
		},
		{
			MethodName: "Exit",
			Handler:    _Gofigure_Exit_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "server.proto",
}
