// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v3.12.4
// source: problem.proto

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

type Problem struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id          int64          `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	CreateAt    string         `protobuf:"bytes,2,opt,name=create_at,json=createAt,proto3" json:"create_at,omitempty"`
	DeleteAt    string         `protobuf:"bytes,3,opt,name=delete_at,json=deleteAt,proto3" json:"delete_at,omitempty"`
	Title       string         `protobuf:"bytes,4,opt,name=title,proto3" json:"title,omitempty"`
	Description string         `protobuf:"bytes,5,opt,name=description,proto3" json:"description,omitempty"`
	Level       int32          `protobuf:"varint,6,opt,name=level,proto3" json:"level,omitempty"`
	Tags        []string       `protobuf:"bytes,7,rep,name=tags,proto3" json:"tags,omitempty"`
	CreateBy    int64          `protobuf:"varint,8,opt,name=create_by,json=createBy,proto3" json:"create_by,omitempty"`
	Config      *ProblemConfig `protobuf:"bytes,9,opt,name=config,proto3" json:"config,omitempty"`
	Status      int32          `protobuf:"varint,10,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *Problem) Reset() {
	*x = Problem{}
	mi := &file_problem_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Problem) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Problem) ProtoMessage() {}

func (x *Problem) ProtoReflect() protoreflect.Message {
	mi := &file_problem_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Problem.ProtoReflect.Descriptor instead.
func (*Problem) Descriptor() ([]byte, []int) {
	return file_problem_proto_rawDescGZIP(), []int{0}
}

func (x *Problem) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Problem) GetCreateAt() string {
	if x != nil {
		return x.CreateAt
	}
	return ""
}

func (x *Problem) GetDeleteAt() string {
	if x != nil {
		return x.DeleteAt
	}
	return ""
}

func (x *Problem) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *Problem) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Problem) GetLevel() int32 {
	if x != nil {
		return x.Level
	}
	return 0
}

func (x *Problem) GetTags() []string {
	if x != nil {
		return x.Tags
	}
	return nil
}

func (x *Problem) GetCreateBy() int64 {
	if x != nil {
		return x.CreateBy
	}
	return 0
}

func (x *Problem) GetConfig() *ProblemConfig {
	if x != nil {
		return x.Config
	}
	return nil
}

func (x *Problem) GetStatus() int32 {
	if x != nil {
		return x.Status
	}
	return 0
}

type ProblemConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TestCases    []*TestCase `protobuf:"bytes,1,rep,name=test_cases,json=testCases,proto3" json:"test_cases,omitempty"`
	CompileLimit *Limit      `protobuf:"bytes,2,opt,name=compile_limit,json=compileLimit,proto3" json:"compile_limit,omitempty"`
	RunLimit     *Limit      `protobuf:"bytes,3,opt,name=run_limit,json=runLimit,proto3" json:"run_limit,omitempty"`
	Level        int32       `protobuf:"varint,4,opt,name=level,proto3" json:"level,omitempty"`
	Name         string      `protobuf:"bytes,5,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *ProblemConfig) Reset() {
	*x = ProblemConfig{}
	mi := &file_problem_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ProblemConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProblemConfig) ProtoMessage() {}

func (x *ProblemConfig) ProtoReflect() protoreflect.Message {
	mi := &file_problem_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProblemConfig.ProtoReflect.Descriptor instead.
func (*ProblemConfig) Descriptor() ([]byte, []int) {
	return file_problem_proto_rawDescGZIP(), []int{1}
}

func (x *ProblemConfig) GetTestCases() []*TestCase {
	if x != nil {
		return x.TestCases
	}
	return nil
}

func (x *ProblemConfig) GetCompileLimit() *Limit {
	if x != nil {
		return x.CompileLimit
	}
	return nil
}

func (x *ProblemConfig) GetRunLimit() *Limit {
	if x != nil {
		return x.RunLimit
	}
	return nil
}

func (x *ProblemConfig) GetLevel() int32 {
	if x != nil {
		return x.Level
	}
	return 0
}

func (x *ProblemConfig) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type TestCase struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Input  string `protobuf:"bytes,1,opt,name=input,proto3" json:"input,omitempty"`
	Output string `protobuf:"bytes,2,opt,name=output,proto3" json:"output,omitempty"`
}

func (x *TestCase) Reset() {
	*x = TestCase{}
	mi := &file_problem_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TestCase) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TestCase) ProtoMessage() {}

func (x *TestCase) ProtoReflect() protoreflect.Message {
	mi := &file_problem_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TestCase.ProtoReflect.Descriptor instead.
func (*TestCase) Descriptor() ([]byte, []int) {
	return file_problem_proto_rawDescGZIP(), []int{2}
}

func (x *TestCase) GetInput() string {
	if x != nil {
		return x.Input
	}
	return ""
}

func (x *TestCase) GetOutput() string {
	if x != nil {
		return x.Output
	}
	return ""
}

type Limit struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CpuLimit    int64 `protobuf:"varint,1,opt,name=cpu_limit,json=cpuLimit,proto3" json:"cpu_limit,omitempty"`
	ClockLimit  int64 `protobuf:"varint,2,opt,name=clock_limit,json=clockLimit,proto3" json:"clock_limit,omitempty"`
	MemoryLimit int64 `protobuf:"varint,3,opt,name=memory_limit,json=memoryLimit,proto3" json:"memory_limit,omitempty"`
	StackLimit  int64 `protobuf:"varint,4,opt,name=stack_limit,json=stackLimit,proto3" json:"stack_limit,omitempty"`
	ProcLimit   int64 `protobuf:"varint,5,opt,name=proc_limit,json=procLimit,proto3" json:"proc_limit,omitempty"`
}

func (x *Limit) Reset() {
	*x = Limit{}
	mi := &file_problem_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Limit) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Limit) ProtoMessage() {}

func (x *Limit) ProtoReflect() protoreflect.Message {
	mi := &file_problem_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Limit.ProtoReflect.Descriptor instead.
func (*Limit) Descriptor() ([]byte, []int) {
	return file_problem_proto_rawDescGZIP(), []int{3}
}

func (x *Limit) GetCpuLimit() int64 {
	if x != nil {
		return x.CpuLimit
	}
	return 0
}

func (x *Limit) GetClockLimit() int64 {
	if x != nil {
		return x.ClockLimit
	}
	return 0
}

func (x *Limit) GetMemoryLimit() int64 {
	if x != nil {
		return x.MemoryLimit
	}
	return 0
}

func (x *Limit) GetStackLimit() int64 {
	if x != nil {
		return x.StackLimit
	}
	return 0
}

func (x *Limit) GetProcLimit() int64 {
	if x != nil {
		return x.ProcLimit
	}
	return 0
}

var File_problem_proto protoreflect.FileDescriptor

var file_problem_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x70, 0x72, 0x6f, 0x62, 0x6c, 0x65, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x92, 0x02, 0x0a, 0x07, 0x50, 0x72, 0x6f, 0x62, 0x6c, 0x65, 0x6d, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x63,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x5f, 0x61, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x41, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x64, 0x65, 0x6c, 0x65,
	0x74, 0x65, 0x5f, 0x61, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x64, 0x65, 0x6c,
	0x65, 0x74, 0x65, 0x41, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64,
	0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x14, 0x0a,
	0x05, 0x6c, 0x65, 0x76, 0x65, 0x6c, 0x18, 0x06, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x6c, 0x65,
	0x76, 0x65, 0x6c, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x61, 0x67, 0x73, 0x18, 0x07, 0x20, 0x03, 0x28,
	0x09, 0x52, 0x04, 0x74, 0x61, 0x67, 0x73, 0x12, 0x1b, 0x0a, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x5f, 0x62, 0x79, 0x18, 0x08, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x63, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x42, 0x79, 0x12, 0x26, 0x0a, 0x06, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x18, 0x09,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x50, 0x72, 0x6f, 0x62, 0x6c, 0x65, 0x6d, 0x43, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x52, 0x06, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x16, 0x0a, 0x06,
	0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x22, 0xb5, 0x01, 0x0a, 0x0d, 0x50, 0x72, 0x6f, 0x62, 0x6c, 0x65, 0x6d,
	0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x28, 0x0a, 0x0a, 0x74, 0x65, 0x73, 0x74, 0x5f, 0x63,
	0x61, 0x73, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x09, 0x2e, 0x54, 0x65, 0x73,
	0x74, 0x43, 0x61, 0x73, 0x65, 0x52, 0x09, 0x74, 0x65, 0x73, 0x74, 0x43, 0x61, 0x73, 0x65, 0x73,
	0x12, 0x2b, 0x0a, 0x0d, 0x63, 0x6f, 0x6d, 0x70, 0x69, 0x6c, 0x65, 0x5f, 0x6c, 0x69, 0x6d, 0x69,
	0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x06, 0x2e, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x52,
	0x0c, 0x63, 0x6f, 0x6d, 0x70, 0x69, 0x6c, 0x65, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x12, 0x23, 0x0a,
	0x09, 0x72, 0x75, 0x6e, 0x5f, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x06, 0x2e, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x52, 0x08, 0x72, 0x75, 0x6e, 0x4c, 0x69, 0x6d,
	0x69, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x65, 0x76, 0x65, 0x6c, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x05, 0x6c, 0x65, 0x76, 0x65, 0x6c, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x38, 0x0a, 0x08,
	0x54, 0x65, 0x73, 0x74, 0x43, 0x61, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x6e, 0x70, 0x75,
	0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x12, 0x16,
	0x0a, 0x06, 0x6f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x6f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x22, 0xa8, 0x01, 0x0a, 0x05, 0x4c, 0x69, 0x6d, 0x69, 0x74,
	0x12, 0x1b, 0x0a, 0x09, 0x63, 0x70, 0x75, 0x5f, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x08, 0x63, 0x70, 0x75, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x12, 0x1f, 0x0a,
	0x0b, 0x63, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x0a, 0x63, 0x6c, 0x6f, 0x63, 0x6b, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x12, 0x21,
	0x0a, 0x0c, 0x6d, 0x65, 0x6d, 0x6f, 0x72, 0x79, 0x5f, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x0b, 0x6d, 0x65, 0x6d, 0x6f, 0x72, 0x79, 0x4c, 0x69, 0x6d, 0x69,
	0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x73, 0x74, 0x61, 0x63, 0x6b, 0x5f, 0x6c, 0x69, 0x6d, 0x69, 0x74,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x73, 0x74, 0x61, 0x63, 0x6b, 0x4c, 0x69, 0x6d,
	0x69, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x63, 0x5f, 0x6c, 0x69, 0x6d, 0x69, 0x74,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x63, 0x4c, 0x69, 0x6d, 0x69,
	0x74, 0x42, 0x07, 0x5a, 0x05, 0x2e, 0x2f, 0x3b, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_problem_proto_rawDescOnce sync.Once
	file_problem_proto_rawDescData = file_problem_proto_rawDesc
)

func file_problem_proto_rawDescGZIP() []byte {
	file_problem_proto_rawDescOnce.Do(func() {
		file_problem_proto_rawDescData = protoimpl.X.CompressGZIP(file_problem_proto_rawDescData)
	})
	return file_problem_proto_rawDescData
}

var file_problem_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_problem_proto_goTypes = []any{
	(*Problem)(nil),       // 0: Problem
	(*ProblemConfig)(nil), // 1: ProblemConfig
	(*TestCase)(nil),      // 2: TestCase
	(*Limit)(nil),         // 3: Limit
}
var file_problem_proto_depIdxs = []int32{
	1, // 0: Problem.config:type_name -> ProblemConfig
	2, // 1: ProblemConfig.test_cases:type_name -> TestCase
	3, // 2: ProblemConfig.compile_limit:type_name -> Limit
	3, // 3: ProblemConfig.run_limit:type_name -> Limit
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_problem_proto_init() }
func file_problem_proto_init() {
	if File_problem_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_problem_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_problem_proto_goTypes,
		DependencyIndexes: file_problem_proto_depIdxs,
		MessageInfos:      file_problem_proto_msgTypes,
	}.Build()
	File_problem_proto = out.File
	file_problem_proto_rawDesc = nil
	file_problem_proto_goTypes = nil
	file_problem_proto_depIdxs = nil
}
