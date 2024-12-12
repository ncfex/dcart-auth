// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        v3.12.4
// source: events.proto

package proto

import (
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
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

type BaseEvent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AggregateId   string               `protobuf:"bytes,1,opt,name=aggregate_id,json=aggregateId,proto3" json:"aggregate_id,omitempty"`
	AggregateType string               `protobuf:"bytes,2,opt,name=aggregate_type,json=aggregateType,proto3" json:"aggregate_type,omitempty"`
	EventType     string               `protobuf:"bytes,3,opt,name=event_type,json=eventType,proto3" json:"event_type,omitempty"`
	Version       int32                `protobuf:"varint,4,opt,name=version,proto3" json:"version,omitempty"`
	Timestamp     *timestamp.Timestamp `protobuf:"bytes,5,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
}

func (x *BaseEvent) Reset() {
	*x = BaseEvent{}
	mi := &file_events_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *BaseEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BaseEvent) ProtoMessage() {}

func (x *BaseEvent) ProtoReflect() protoreflect.Message {
	mi := &file_events_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BaseEvent.ProtoReflect.Descriptor instead.
func (*BaseEvent) Descriptor() ([]byte, []int) {
	return file_events_proto_rawDescGZIP(), []int{0}
}

func (x *BaseEvent) GetAggregateId() string {
	if x != nil {
		return x.AggregateId
	}
	return ""
}

func (x *BaseEvent) GetAggregateType() string {
	if x != nil {
		return x.AggregateType
	}
	return ""
}

func (x *BaseEvent) GetEventType() string {
	if x != nil {
		return x.EventType
	}
	return ""
}

func (x *BaseEvent) GetVersion() int32 {
	if x != nil {
		return x.Version
	}
	return 0
}

func (x *BaseEvent) GetTimestamp() *timestamp.Timestamp {
	if x != nil {
		return x.Timestamp
	}
	return nil
}

type EventMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AggregateId   string               `protobuf:"bytes,1,opt,name=aggregate_id,json=aggregateId,proto3" json:"aggregate_id,omitempty"`
	AggregateType string               `protobuf:"bytes,2,opt,name=aggregate_type,json=aggregateType,proto3" json:"aggregate_type,omitempty"`
	EventType     string               `protobuf:"bytes,3,opt,name=event_type,json=eventType,proto3" json:"event_type,omitempty"`
	Version       int32                `protobuf:"varint,4,opt,name=version,proto3" json:"version,omitempty"`
	Timestamp     *timestamp.Timestamp `protobuf:"bytes,5,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Payload       []byte               `protobuf:"bytes,6,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (x *EventMessage) Reset() {
	*x = EventMessage{}
	mi := &file_events_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EventMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EventMessage) ProtoMessage() {}

func (x *EventMessage) ProtoReflect() protoreflect.Message {
	mi := &file_events_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EventMessage.ProtoReflect.Descriptor instead.
func (*EventMessage) Descriptor() ([]byte, []int) {
	return file_events_proto_rawDescGZIP(), []int{1}
}

func (x *EventMessage) GetAggregateId() string {
	if x != nil {
		return x.AggregateId
	}
	return ""
}

func (x *EventMessage) GetAggregateType() string {
	if x != nil {
		return x.AggregateType
	}
	return ""
}

func (x *EventMessage) GetEventType() string {
	if x != nil {
		return x.EventType
	}
	return ""
}

func (x *EventMessage) GetVersion() int32 {
	if x != nil {
		return x.Version
	}
	return 0
}

func (x *EventMessage) GetTimestamp() *timestamp.Timestamp {
	if x != nil {
		return x.Timestamp
	}
	return nil
}

func (x *EventMessage) GetPayload() []byte {
	if x != nil {
		return x.Payload
	}
	return nil
}

type UserRegisteredEvent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Base         *BaseEvent `protobuf:"bytes,1,opt,name=base,proto3" json:"base,omitempty"`
	Username     string     `protobuf:"bytes,2,opt,name=username,proto3" json:"username,omitempty"`
	PasswordHash string     `protobuf:"bytes,3,opt,name=password_hash,json=passwordHash,proto3" json:"password_hash,omitempty"`
}

func (x *UserRegisteredEvent) Reset() {
	*x = UserRegisteredEvent{}
	mi := &file_events_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UserRegisteredEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserRegisteredEvent) ProtoMessage() {}

func (x *UserRegisteredEvent) ProtoReflect() protoreflect.Message {
	mi := &file_events_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserRegisteredEvent.ProtoReflect.Descriptor instead.
func (*UserRegisteredEvent) Descriptor() ([]byte, []int) {
	return file_events_proto_rawDescGZIP(), []int{2}
}

func (x *UserRegisteredEvent) GetBase() *BaseEvent {
	if x != nil {
		return x.Base
	}
	return nil
}

func (x *UserRegisteredEvent) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *UserRegisteredEvent) GetPasswordHash() string {
	if x != nil {
		return x.PasswordHash
	}
	return ""
}

type UserPasswordChangedEvent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Base            *BaseEvent `protobuf:"bytes,1,opt,name=base,proto3" json:"base,omitempty"`
	NewPasswordHash string     `protobuf:"bytes,2,opt,name=new_password_hash,json=newPasswordHash,proto3" json:"new_password_hash,omitempty"`
}

func (x *UserPasswordChangedEvent) Reset() {
	*x = UserPasswordChangedEvent{}
	mi := &file_events_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UserPasswordChangedEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserPasswordChangedEvent) ProtoMessage() {}

func (x *UserPasswordChangedEvent) ProtoReflect() protoreflect.Message {
	mi := &file_events_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserPasswordChangedEvent.ProtoReflect.Descriptor instead.
func (*UserPasswordChangedEvent) Descriptor() ([]byte, []int) {
	return file_events_proto_rawDescGZIP(), []int{3}
}

func (x *UserPasswordChangedEvent) GetBase() *BaseEvent {
	if x != nil {
		return x.Base
	}
	return nil
}

func (x *UserPasswordChangedEvent) GetNewPasswordHash() string {
	if x != nil {
		return x.NewPasswordHash
	}
	return ""
}

var File_events_proto protoreflect.FileDescriptor

var file_events_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05,
	0x65, 0x76, 0x65, 0x6e, 0x74, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xc8, 0x01, 0x0a, 0x09, 0x42, 0x61, 0x73, 0x65, 0x45,
	0x76, 0x65, 0x6e, 0x74, 0x12, 0x21, 0x0a, 0x0c, 0x61, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74,
	0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x61, 0x67, 0x67, 0x72,
	0x65, 0x67, 0x61, 0x74, 0x65, 0x49, 0x64, 0x12, 0x25, 0x0a, 0x0e, 0x61, 0x67, 0x67, 0x72, 0x65,
	0x67, 0x61, 0x74, 0x65, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0d, 0x61, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1d,
	0x0a, 0x0a, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x09, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x18, 0x0a,
	0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07,
	0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x38, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x22, 0xe5, 0x01, 0x0a, 0x0c, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x61, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x65, 0x5f,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x61, 0x67, 0x67, 0x72, 0x65, 0x67,
	0x61, 0x74, 0x65, 0x49, 0x64, 0x12, 0x25, 0x0a, 0x0e, 0x61, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61,
	0x74, 0x65, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x61,
	0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1d, 0x0a, 0x0a,
	0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x09, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x76,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x76, 0x65,
	0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x38, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12,
	0x18, 0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x22, 0x7c, 0x0a, 0x13, 0x55, 0x73, 0x65,
	0x72, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x65, 0x64, 0x45, 0x76, 0x65, 0x6e, 0x74,
	0x12, 0x24, 0x0a, 0x04, 0x62, 0x61, 0x73, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10,
	0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x2e, 0x42, 0x61, 0x73, 0x65, 0x45, 0x76, 0x65, 0x6e, 0x74,
	0x52, 0x04, 0x62, 0x61, 0x73, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61,
	0x6d, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x5f, 0x68,
	0x61, 0x73, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x70, 0x61, 0x73, 0x73, 0x77,
	0x6f, 0x72, 0x64, 0x48, 0x61, 0x73, 0x68, 0x22, 0x6c, 0x0a, 0x18, 0x55, 0x73, 0x65, 0x72, 0x50,
	0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x64, 0x45, 0x76,
	0x65, 0x6e, 0x74, 0x12, 0x24, 0x0a, 0x04, 0x62, 0x61, 0x73, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x10, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x2e, 0x42, 0x61, 0x73, 0x65, 0x45, 0x76,
	0x65, 0x6e, 0x74, 0x52, 0x04, 0x62, 0x61, 0x73, 0x65, 0x12, 0x2a, 0x0a, 0x11, 0x6e, 0x65, 0x77,
	0x5f, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x6e, 0x65, 0x77, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72,
	0x64, 0x48, 0x61, 0x73, 0x68, 0x42, 0x49, 0x5a, 0x47, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x6e, 0x63, 0x66, 0x65, 0x78, 0x2f, 0x64, 0x63, 0x61, 0x72, 0x74, 0x2d,
	0x61, 0x75, 0x74, 0x68, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x61, 0x64,
	0x61, 0x70, 0x74, 0x65, 0x72, 0x73, 0x2f, 0x73, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x61, 0x72, 0x79,
	0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x69, 0x6e, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_events_proto_rawDescOnce sync.Once
	file_events_proto_rawDescData = file_events_proto_rawDesc
)

func file_events_proto_rawDescGZIP() []byte {
	file_events_proto_rawDescOnce.Do(func() {
		file_events_proto_rawDescData = protoimpl.X.CompressGZIP(file_events_proto_rawDescData)
	})
	return file_events_proto_rawDescData
}

var file_events_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_events_proto_goTypes = []any{
	(*BaseEvent)(nil),                // 0: event.BaseEvent
	(*EventMessage)(nil),             // 1: event.EventMessage
	(*UserRegisteredEvent)(nil),      // 2: event.UserRegisteredEvent
	(*UserPasswordChangedEvent)(nil), // 3: event.UserPasswordChangedEvent
	(*timestamp.Timestamp)(nil),      // 4: google.protobuf.Timestamp
}
var file_events_proto_depIdxs = []int32{
	4, // 0: event.BaseEvent.timestamp:type_name -> google.protobuf.Timestamp
	4, // 1: event.EventMessage.timestamp:type_name -> google.protobuf.Timestamp
	0, // 2: event.UserRegisteredEvent.base:type_name -> event.BaseEvent
	0, // 3: event.UserPasswordChangedEvent.base:type_name -> event.BaseEvent
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_events_proto_init() }
func file_events_proto_init() {
	if File_events_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_events_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_events_proto_goTypes,
		DependencyIndexes: file_events_proto_depIdxs,
		MessageInfos:      file_events_proto_msgTypes,
	}.Build()
	File_events_proto = out.File
	file_events_proto_rawDesc = nil
	file_events_proto_goTypes = nil
	file_events_proto_depIdxs = nil
}