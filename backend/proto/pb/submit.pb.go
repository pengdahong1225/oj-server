// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.26.0
// source: submit.proto

package pb

import (
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

// “提交”状态枚举，如果没有查询到状态，就意味着最近没有提交题目or题目提交过期了
type SubmitState int32

const (
	SubmitState_UPStateNormal    SubmitState = 0 // 初始状态
	SubmitState_UPStateCompiling SubmitState = 1 // 正在编译
	SubmitState_UPStateJudging   SubmitState = 2 // 已编译成功，正在判题
	SubmitState_UPStateExited    SubmitState = 3 // 已退出 -> 如何查询到这个状态，就意味着可以查询结果了
)

// Enum value maps for SubmitState.
var (
	SubmitState_name = map[int32]string{
		0: "UPStateNormal",
		1: "UPStateCompiling",
		2: "UPStateJudging",
		3: "UPStateExited",
	}
	SubmitState_value = map[string]int32{
		"UPStateNormal":    0,
		"UPStateCompiling": 1,
		"UPStateJudging":   2,
		"UPStateExited":    3,
	}
)

func (x SubmitState) Enum() *SubmitState {
	p := new(SubmitState)
	*p = x
	return p
}

func (x SubmitState) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (SubmitState) Descriptor() protoreflect.EnumDescriptor {
	return file_submit_proto_enumTypes[0].Descriptor()
}

func (SubmitState) Type() protoreflect.EnumType {
	return &file_submit_proto_enumTypes[0]
}

func (x SubmitState) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use SubmitState.Descriptor instead.
func (SubmitState) EnumDescriptor() ([]byte, []int) {
	return file_submit_proto_rawDescGZIP(), []int{0}
}

type SubmitForm struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProblemId int64  `protobuf:"varint,1,opt,name=problem_id,json=problemId,proto3" json:"problem_id,omitempty"`
	Title     string `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	Lang      string `protobuf:"bytes,3,opt,name=lang,proto3" json:"lang,omitempty"`
	Code      string `protobuf:"bytes,4,opt,name=code,proto3" json:"code,omitempty"`
	Uid       int64  `protobuf:"varint,5,opt,name=uid,proto3" json:"uid,omitempty"`
}

func (x *SubmitForm) Reset() {
	*x = SubmitForm{}
	if protoimpl.UnsafeEnabled {
		mi := &file_submit_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SubmitForm) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubmitForm) ProtoMessage() {}

func (x *SubmitForm) ProtoReflect() protoreflect.Message {
	mi := &file_submit_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubmitForm.ProtoReflect.Descriptor instead.
func (*SubmitForm) Descriptor() ([]byte, []int) {
	return file_submit_proto_rawDescGZIP(), []int{0}
}

func (x *SubmitForm) GetProblemId() int64 {
	if x != nil {
		return x.ProblemId
	}
	return 0
}

func (x *SubmitForm) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *SubmitForm) GetLang() string {
	if x != nil {
		return x.Lang
	}
	return ""
}

func (x *SubmitForm) GetCode() string {
	if x != nil {
		return x.Code
	}
	return ""
}

func (x *SubmitForm) GetUid() int64 {
	if x != nil {
		return x.Uid
	}
	return 0
}

type JudgeResult struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status     string            `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
	Error      string            `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`                                                                                             // 详细错误信息
	ExitStatus int64             `protobuf:"varint,3,opt,name=exitStatus,proto3" json:"exitStatus,omitempty"`                                                                                  // 程序返回值
	Time       int64             `protobuf:"varint,4,opt,name=time,proto3" json:"time,omitempty"`                                                                                              // 程序运行 CPU 时间，单位纳秒
	Memory     int64             `protobuf:"varint,5,opt,name=memory,proto3" json:"memory,omitempty"`                                                                                          // 程序运行内存，单位 byte
	RunTime    int64             `protobuf:"varint,6,opt,name=runTime,proto3" json:"runTime,omitempty"`                                                                                        // 程序运行现实时间，单位纳秒
	Files      map[string]string `protobuf:"bytes,7,rep,name=files,proto3" json:"files,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`     // copyOut 和 pipeCollector 指定的文件内容
	FileIds    map[string]string `protobuf:"bytes,8,rep,name=fileIds,proto3" json:"fileIds,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"` // copyFileCached 指定的文件 id
	FileError  []string          `protobuf:"bytes,9,rep,name=fileError,proto3" json:"fileError,omitempty"`                                                                                     // 文件错误详细信息
	// 以下字段不属于判题服务的返回字段
	Content  string    `protobuf:"bytes,10,opt,name=content,proto3" json:"content,omitempty"`
	TestCase *TestCase `protobuf:"bytes,11,opt,name=testCase,proto3" json:"testCase,omitempty"`
}

func (x *JudgeResult) Reset() {
	*x = JudgeResult{}
	if protoimpl.UnsafeEnabled {
		mi := &file_submit_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *JudgeResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JudgeResult) ProtoMessage() {}

func (x *JudgeResult) ProtoReflect() protoreflect.Message {
	mi := &file_submit_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use JudgeResult.ProtoReflect.Descriptor instead.
func (*JudgeResult) Descriptor() ([]byte, []int) {
	return file_submit_proto_rawDescGZIP(), []int{1}
}

func (x *JudgeResult) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *JudgeResult) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

func (x *JudgeResult) GetExitStatus() int64 {
	if x != nil {
		return x.ExitStatus
	}
	return 0
}

func (x *JudgeResult) GetTime() int64 {
	if x != nil {
		return x.Time
	}
	return 0
}

func (x *JudgeResult) GetMemory() int64 {
	if x != nil {
		return x.Memory
	}
	return 0
}

func (x *JudgeResult) GetRunTime() int64 {
	if x != nil {
		return x.RunTime
	}
	return 0
}

func (x *JudgeResult) GetFiles() map[string]string {
	if x != nil {
		return x.Files
	}
	return nil
}

func (x *JudgeResult) GetFileIds() map[string]string {
	if x != nil {
		return x.FileIds
	}
	return nil
}

func (x *JudgeResult) GetFileError() []string {
	if x != nil {
		return x.FileError
	}
	return nil
}

func (x *JudgeResult) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

func (x *JudgeResult) GetTestCase() *TestCase {
	if x != nil {
		return x.TestCase
	}
	return nil
}

var File_submit_proto protoreflect.FileDescriptor

var file_submit_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x73, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0d,
	0x70, 0x72, 0x6f, 0x62, 0x6c, 0x65, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x7b, 0x0a,
	0x0a, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x46, 0x6f, 0x72, 0x6d, 0x12, 0x1d, 0x0a, 0x0a, 0x70,
	0x72, 0x6f, 0x62, 0x6c, 0x65, 0x6d, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x09, 0x70, 0x72, 0x6f, 0x62, 0x6c, 0x65, 0x6d, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x69,
	0x74, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65,
	0x12, 0x12, 0x0a, 0x04, 0x6c, 0x61, 0x6e, 0x67, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x6c, 0x61, 0x6e, 0x67, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x69, 0x64, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x75, 0x69, 0x64, 0x22, 0xda, 0x03, 0x0a, 0x0b, 0x4a,
	0x75, 0x64, 0x67, 0x65, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x12, 0x1e, 0x0a, 0x0a, 0x65, 0x78, 0x69, 0x74,
	0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x65, 0x78,
	0x69, 0x74, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x69, 0x6d, 0x65,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06,
	0x6d, 0x65, 0x6d, 0x6f, 0x72, 0x79, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x6d, 0x65,
	0x6d, 0x6f, 0x72, 0x79, 0x12, 0x18, 0x0a, 0x07, 0x72, 0x75, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x72, 0x75, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x2d,
	0x0a, 0x05, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x18, 0x07, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x17, 0x2e,
	0x4a, 0x75, 0x64, 0x67, 0x65, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x2e, 0x46, 0x69, 0x6c, 0x65,
	0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x05, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x12, 0x33, 0x0a,
	0x07, 0x66, 0x69, 0x6c, 0x65, 0x49, 0x64, 0x73, 0x18, 0x08, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x19,
	0x2e, 0x4a, 0x75, 0x64, 0x67, 0x65, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x2e, 0x46, 0x69, 0x6c,
	0x65, 0x49, 0x64, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x07, 0x66, 0x69, 0x6c, 0x65, 0x49,
	0x64, 0x73, 0x12, 0x1c, 0x0a, 0x09, 0x66, 0x69, 0x6c, 0x65, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x18,
	0x09, 0x20, 0x03, 0x28, 0x09, 0x52, 0x09, 0x66, 0x69, 0x6c, 0x65, 0x45, 0x72, 0x72, 0x6f, 0x72,
	0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x0a, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x12, 0x25, 0x0a, 0x08, 0x74, 0x65,
	0x73, 0x74, 0x43, 0x61, 0x73, 0x65, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x09, 0x2e, 0x54,
	0x65, 0x73, 0x74, 0x43, 0x61, 0x73, 0x65, 0x52, 0x08, 0x74, 0x65, 0x73, 0x74, 0x43, 0x61, 0x73,
	0x65, 0x1a, 0x38, 0x0a, 0x0a, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12,
	0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65,
	0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a, 0x3a, 0x0a, 0x0c, 0x46,
	0x69, 0x6c, 0x65, 0x49, 0x64, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b,
	0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x2a, 0x5d, 0x0a, 0x0b, 0x53, 0x75, 0x62, 0x6d, 0x69,
	0x74, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x11, 0x0a, 0x0d, 0x55, 0x50, 0x53, 0x74, 0x61, 0x74,
	0x65, 0x4e, 0x6f, 0x72, 0x6d, 0x61, 0x6c, 0x10, 0x00, 0x12, 0x14, 0x0a, 0x10, 0x55, 0x50, 0x53,
	0x74, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x6d, 0x70, 0x69, 0x6c, 0x69, 0x6e, 0x67, 0x10, 0x01, 0x12,
	0x12, 0x0a, 0x0e, 0x55, 0x50, 0x53, 0x74, 0x61, 0x74, 0x65, 0x4a, 0x75, 0x64, 0x67, 0x69, 0x6e,
	0x67, 0x10, 0x02, 0x12, 0x11, 0x0a, 0x0d, 0x55, 0x50, 0x53, 0x74, 0x61, 0x74, 0x65, 0x45, 0x78,
	0x69, 0x74, 0x65, 0x64, 0x10, 0x03, 0x42, 0x07, 0x5a, 0x05, 0x2e, 0x2f, 0x3b, 0x70, 0x62, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_submit_proto_rawDescOnce sync.Once
	file_submit_proto_rawDescData = file_submit_proto_rawDesc
)

func file_submit_proto_rawDescGZIP() []byte {
	file_submit_proto_rawDescOnce.Do(func() {
		file_submit_proto_rawDescData = protoimpl.X.CompressGZIP(file_submit_proto_rawDescData)
	})
	return file_submit_proto_rawDescData
}

var file_submit_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_submit_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_submit_proto_goTypes = []any{
	(SubmitState)(0),    // 0: SubmitState
	(*SubmitForm)(nil),  // 1: SubmitForm
	(*JudgeResult)(nil), // 2: JudgeResult
	nil,                 // 3: JudgeResult.FilesEntry
	nil,                 // 4: JudgeResult.FileIdsEntry
	(*TestCase)(nil),    // 5: TestCase
}
var file_submit_proto_depIdxs = []int32{
	3, // 0: JudgeResult.files:type_name -> JudgeResult.FilesEntry
	4, // 1: JudgeResult.fileIds:type_name -> JudgeResult.FileIdsEntry
	5, // 2: JudgeResult.testCase:type_name -> TestCase
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_submit_proto_init() }
func file_submit_proto_init() {
	if File_submit_proto != nil {
		return
	}
	file_problem_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_submit_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*SubmitForm); i {
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
		file_submit_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*JudgeResult); i {
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
			RawDescriptor: file_submit_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_submit_proto_goTypes,
		DependencyIndexes: file_submit_proto_depIdxs,
		EnumInfos:         file_submit_proto_enumTypes,
		MessageInfos:      file_submit_proto_msgTypes,
	}.Build()
	File_submit_proto = out.File
	file_submit_proto_rawDesc = nil
	file_submit_proto_goTypes = nil
	file_submit_proto_depIdxs = nil
}
