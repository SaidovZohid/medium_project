// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: like_service.proto

package post_service

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_like_service_proto protoreflect.FileDescriptor

var file_like_service_proto_rawDesc = []byte{
	0x0a, 0x12, 0x6c, 0x69, 0x6b, 0x65, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x67, 0x65, 0x6e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0a,
	0x6c, 0x69, 0x6b, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x32, 0xd0, 0x01, 0x0a, 0x0b, 0x4c,
	0x69, 0x6b, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x32, 0x0a, 0x0e, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x4f, 0x72, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x12, 0x0e, 0x2e, 0x67,
	0x65, 0x6e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4c, 0x69, 0x6b, 0x65, 0x1a, 0x0e, 0x2e, 0x67,
	0x65, 0x6e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4c, 0x69, 0x6b, 0x65, 0x22, 0x00, 0x12, 0x31,
	0x0a, 0x03, 0x47, 0x65, 0x74, 0x12, 0x18, 0x2e, 0x67, 0x65, 0x6e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x47, 0x65, 0x74, 0x4c, 0x69, 0x6b, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x0e, 0x2e, 0x67, 0x65, 0x6e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4c, 0x69, 0x6b, 0x65, 0x22,
	0x00, 0x12, 0x5a, 0x0a, 0x15, 0x47, 0x65, 0x74, 0x4c, 0x69, 0x6b, 0x65, 0x73, 0x44, 0x69, 0x73,
	0x6c, 0x69, 0x6b, 0x65, 0x73, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x19, 0x2e, 0x67, 0x65, 0x6e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x4c, 0x69, 0x6b, 0x65, 0x73, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x24, 0x2e, 0x67, 0x65, 0x6e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x4c, 0x69, 0x6b, 0x65, 0x73, 0x44, 0x69, 0x73, 0x6c, 0x69, 0x6b, 0x65, 0x73, 0x43, 0x6f,
	0x75, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x17, 0x5a,
	0x15, 0x67, 0x65, 0x6e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x70, 0x6f, 0x73, 0x74, 0x5f, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_like_service_proto_goTypes = []interface{}{
	(*Like)(nil),                       // 0: genproto.Like
	(*GetLikeRequest)(nil),             // 1: genproto.GetLikeRequest
	(*GetLikesRequest)(nil),            // 2: genproto.GetLikesRequest
	(*LikesDislikesCountResponse)(nil), // 3: genproto.LikesDislikesCountResponse
}
var file_like_service_proto_depIdxs = []int32{
	0, // 0: genproto.LikeService.CreateOrUpdate:input_type -> genproto.Like
	1, // 1: genproto.LikeService.Get:input_type -> genproto.GetLikeRequest
	2, // 2: genproto.LikeService.GetLikesDislikesCount:input_type -> genproto.GetLikesRequest
	0, // 3: genproto.LikeService.CreateOrUpdate:output_type -> genproto.Like
	0, // 4: genproto.LikeService.Get:output_type -> genproto.Like
	3, // 5: genproto.LikeService.GetLikesDislikesCount:output_type -> genproto.LikesDislikesCountResponse
	3, // [3:6] is the sub-list for method output_type
	0, // [0:3] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_like_service_proto_init() }
func file_like_service_proto_init() {
	if File_like_service_proto != nil {
		return
	}
	file_like_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_like_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_like_service_proto_goTypes,
		DependencyIndexes: file_like_service_proto_depIdxs,
	}.Build()
	File_like_service_proto = out.File
	file_like_service_proto_rawDesc = nil
	file_like_service_proto_goTypes = nil
	file_like_service_proto_depIdxs = nil
}
