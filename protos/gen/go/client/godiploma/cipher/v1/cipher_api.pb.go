// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: client/godiploma/cipher/v1/cipher_api.proto

package cipherv1

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type CreateStegoImageRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	RequestId     string                 `protobuf:"bytes,1,opt,name=request_id,json=requestId,proto3" json:"request_id,omitempty"`
	Method        Method                 `protobuf:"varint,2,opt,name=method,proto3,enum=client.godiploma.cipher.v1.Method" json:"method,omitempty"`
	Plaintext     string                 `protobuf:"bytes,3,opt,name=plaintext,proto3" json:"plaintext,omitempty"`
	Files         []*File                `protobuf:"bytes,4,rep,name=files,proto3" json:"files,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateStegoImageRequest) Reset() {
	*x = CreateStegoImageRequest{}
	mi := &file_client_godiploma_cipher_v1_cipher_api_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateStegoImageRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateStegoImageRequest) ProtoMessage() {}

func (x *CreateStegoImageRequest) ProtoReflect() protoreflect.Message {
	mi := &file_client_godiploma_cipher_v1_cipher_api_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateStegoImageRequest.ProtoReflect.Descriptor instead.
func (*CreateStegoImageRequest) Descriptor() ([]byte, []int) {
	return file_client_godiploma_cipher_v1_cipher_api_proto_rawDescGZIP(), []int{0}
}

func (x *CreateStegoImageRequest) GetRequestId() string {
	if x != nil {
		return x.RequestId
	}
	return ""
}

func (x *CreateStegoImageRequest) GetMethod() Method {
	if x != nil {
		return x.Method
	}
	return Method_METHOD_UNSPECIFIED
}

func (x *CreateStegoImageRequest) GetPlaintext() string {
	if x != nil {
		return x.Plaintext
	}
	return ""
}

func (x *CreateStegoImageRequest) GetFiles() []*File {
	if x != nil {
		return x.Files
	}
	return nil
}

type CreateStegoImageResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Files         []*File                `protobuf:"bytes,1,rep,name=files,proto3" json:"files,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateStegoImageResponse) Reset() {
	*x = CreateStegoImageResponse{}
	mi := &file_client_godiploma_cipher_v1_cipher_api_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateStegoImageResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateStegoImageResponse) ProtoMessage() {}

func (x *CreateStegoImageResponse) ProtoReflect() protoreflect.Message {
	mi := &file_client_godiploma_cipher_v1_cipher_api_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateStegoImageResponse.ProtoReflect.Descriptor instead.
func (*CreateStegoImageResponse) Descriptor() ([]byte, []int) {
	return file_client_godiploma_cipher_v1_cipher_api_proto_rawDescGZIP(), []int{1}
}

func (x *CreateStegoImageResponse) GetFiles() []*File {
	if x != nil {
		return x.Files
	}
	return nil
}

type ExtractRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	RequestId     string                 `protobuf:"bytes,1,opt,name=request_id,json=requestId,proto3" json:"request_id,omitempty"`
	Method        Method                 `protobuf:"varint,2,opt,name=method,proto3,enum=client.godiploma.cipher.v1.Method" json:"method,omitempty"`
	Files         []*File                `protobuf:"bytes,4,rep,name=files,proto3" json:"files,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ExtractRequest) Reset() {
	*x = ExtractRequest{}
	mi := &file_client_godiploma_cipher_v1_cipher_api_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ExtractRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExtractRequest) ProtoMessage() {}

func (x *ExtractRequest) ProtoReflect() protoreflect.Message {
	mi := &file_client_godiploma_cipher_v1_cipher_api_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExtractRequest.ProtoReflect.Descriptor instead.
func (*ExtractRequest) Descriptor() ([]byte, []int) {
	return file_client_godiploma_cipher_v1_cipher_api_proto_rawDescGZIP(), []int{2}
}

func (x *ExtractRequest) GetRequestId() string {
	if x != nil {
		return x.RequestId
	}
	return ""
}

func (x *ExtractRequest) GetMethod() Method {
	if x != nil {
		return x.Method
	}
	return Method_METHOD_UNSPECIFIED
}

func (x *ExtractRequest) GetFiles() []*File {
	if x != nil {
		return x.Files
	}
	return nil
}

type ExtractResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Plaintext     []string               `protobuf:"bytes,1,rep,name=plaintext,proto3" json:"plaintext,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ExtractResponse) Reset() {
	*x = ExtractResponse{}
	mi := &file_client_godiploma_cipher_v1_cipher_api_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ExtractResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExtractResponse) ProtoMessage() {}

func (x *ExtractResponse) ProtoReflect() protoreflect.Message {
	mi := &file_client_godiploma_cipher_v1_cipher_api_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExtractResponse.ProtoReflect.Descriptor instead.
func (*ExtractResponse) Descriptor() ([]byte, []int) {
	return file_client_godiploma_cipher_v1_cipher_api_proto_rawDescGZIP(), []int{3}
}

func (x *ExtractResponse) GetPlaintext() []string {
	if x != nil {
		return x.Plaintext
	}
	return nil
}

var File_client_godiploma_cipher_v1_cipher_api_proto protoreflect.FileDescriptor

const file_client_godiploma_cipher_v1_cipher_api_proto_rawDesc = "" +
	"\n" +
	"+client/godiploma/cipher/v1/cipher_api.proto\x12\x1aclient.godiploma.cipher.v1\x1a'client/godiploma/cipher/v1/cipher.proto\x1a%client/godiploma/cipher/v1/file.proto\x1a\x17validate/validate.proto\"\xf3\x01\n" +
	"\x17CreateStegoImageRequest\x12'\n" +
	"\n" +
	"request_id\x18\x01 \x01(\tB\b\xfaB\x05r\x03\xb0\x01\x01R\trequestId\x12F\n" +
	"\x06method\x18\x02 \x01(\x0e2\".client.godiploma.cipher.v1.MethodB\n" +
	"\xfaB\a\x82\x01\x04\x10\x01 \x00R\x06method\x12%\n" +
	"\tplaintext\x18\x03 \x01(\tB\a\xfaB\x04r\x02\x10\x01R\tplaintext\x12@\n" +
	"\x05files\x18\x04 \x03(\v2 .client.godiploma.cipher.v1.FileB\b\xfaB\x05\x92\x01\x02\b\x01R\x05files\"R\n" +
	"\x18CreateStegoImageResponse\x126\n" +
	"\x05files\x18\x01 \x03(\v2 .client.godiploma.cipher.v1.FileR\x05files\"\xc3\x01\n" +
	"\x0eExtractRequest\x12'\n" +
	"\n" +
	"request_id\x18\x01 \x01(\tB\b\xfaB\x05r\x03\xb0\x01\x01R\trequestId\x12F\n" +
	"\x06method\x18\x02 \x01(\x0e2\".client.godiploma.cipher.v1.MethodB\n" +
	"\xfaB\a\x82\x01\x04\x10\x01 \x00R\x06method\x12@\n" +
	"\x05files\x18\x04 \x03(\v2 .client.godiploma.cipher.v1.FileB\b\xfaB\x05\x92\x01\x02\b\x01R\x05files\"/\n" +
	"\x0fExtractResponse\x12\x1c\n" +
	"\tplaintext\x18\x01 \x03(\tR\tplaintext2\xf2\x01\n" +
	"\rCipherService\x12}\n" +
	"\x10CreateStegoImage\x123.client.godiploma.cipher.v1.CreateStegoImageRequest\x1a4.client.godiploma.cipher.v1.CreateStegoImageResponse\x12b\n" +
	"\aExtract\x12*.client.godiploma.cipher.v1.ExtractRequest\x1a+.client.godiploma.cipher.v1.ExtractResponseB\x91\x02\n" +
	"\x1ecom.client.godiploma.cipher.v1B\x0eCipherApiProtoP\x01ZTgithub.com/bogdanpashtet/godiploma/protos/gen/go/client/godiploma/cipher/v1;cipherv1\xa2\x02\x03CGC\xaa\x02\x1aClient.Godiploma.Cipher.V1\xca\x02\x1aClient\\Godiploma\\Cipher\\V1\xe2\x02&Client\\Godiploma\\Cipher\\V1\\GPBMetadata\xea\x02\x1dClient::Godiploma::Cipher::V1b\x06proto3"

var (
	file_client_godiploma_cipher_v1_cipher_api_proto_rawDescOnce sync.Once
	file_client_godiploma_cipher_v1_cipher_api_proto_rawDescData []byte
)

func file_client_godiploma_cipher_v1_cipher_api_proto_rawDescGZIP() []byte {
	file_client_godiploma_cipher_v1_cipher_api_proto_rawDescOnce.Do(func() {
		file_client_godiploma_cipher_v1_cipher_api_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_client_godiploma_cipher_v1_cipher_api_proto_rawDesc), len(file_client_godiploma_cipher_v1_cipher_api_proto_rawDesc)))
	})
	return file_client_godiploma_cipher_v1_cipher_api_proto_rawDescData
}

var file_client_godiploma_cipher_v1_cipher_api_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_client_godiploma_cipher_v1_cipher_api_proto_goTypes = []any{
	(*CreateStegoImageRequest)(nil),  // 0: client.godiploma.cipher.v1.CreateStegoImageRequest
	(*CreateStegoImageResponse)(nil), // 1: client.godiploma.cipher.v1.CreateStegoImageResponse
	(*ExtractRequest)(nil),           // 2: client.godiploma.cipher.v1.ExtractRequest
	(*ExtractResponse)(nil),          // 3: client.godiploma.cipher.v1.ExtractResponse
	(Method)(0),                      // 4: client.godiploma.cipher.v1.Method
	(*File)(nil),                     // 5: client.godiploma.cipher.v1.File
}
var file_client_godiploma_cipher_v1_cipher_api_proto_depIdxs = []int32{
	4, // 0: client.godiploma.cipher.v1.CreateStegoImageRequest.method:type_name -> client.godiploma.cipher.v1.Method
	5, // 1: client.godiploma.cipher.v1.CreateStegoImageRequest.files:type_name -> client.godiploma.cipher.v1.File
	5, // 2: client.godiploma.cipher.v1.CreateStegoImageResponse.files:type_name -> client.godiploma.cipher.v1.File
	4, // 3: client.godiploma.cipher.v1.ExtractRequest.method:type_name -> client.godiploma.cipher.v1.Method
	5, // 4: client.godiploma.cipher.v1.ExtractRequest.files:type_name -> client.godiploma.cipher.v1.File
	0, // 5: client.godiploma.cipher.v1.CipherService.CreateStegoImage:input_type -> client.godiploma.cipher.v1.CreateStegoImageRequest
	2, // 6: client.godiploma.cipher.v1.CipherService.Extract:input_type -> client.godiploma.cipher.v1.ExtractRequest
	1, // 7: client.godiploma.cipher.v1.CipherService.CreateStegoImage:output_type -> client.godiploma.cipher.v1.CreateStegoImageResponse
	3, // 8: client.godiploma.cipher.v1.CipherService.Extract:output_type -> client.godiploma.cipher.v1.ExtractResponse
	7, // [7:9] is the sub-list for method output_type
	5, // [5:7] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_client_godiploma_cipher_v1_cipher_api_proto_init() }
func file_client_godiploma_cipher_v1_cipher_api_proto_init() {
	if File_client_godiploma_cipher_v1_cipher_api_proto != nil {
		return
	}
	file_client_godiploma_cipher_v1_cipher_proto_init()
	file_client_godiploma_cipher_v1_file_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_client_godiploma_cipher_v1_cipher_api_proto_rawDesc), len(file_client_godiploma_cipher_v1_cipher_api_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_client_godiploma_cipher_v1_cipher_api_proto_goTypes,
		DependencyIndexes: file_client_godiploma_cipher_v1_cipher_api_proto_depIdxs,
		MessageInfos:      file_client_godiploma_cipher_v1_cipher_api_proto_msgTypes,
	}.Build()
	File_client_godiploma_cipher_v1_cipher_api_proto = out.File
	file_client_godiploma_cipher_v1_cipher_api_proto_goTypes = nil
	file_client_godiploma_cipher_v1_cipher_api_proto_depIdxs = nil
}
