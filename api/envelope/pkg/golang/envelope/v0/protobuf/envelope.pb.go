// Code generated by protoc-gen-go. DO NOT EDIT.
// source: envelope.proto

package protobuf

import (
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
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

type Envelope struct {
	Version              uint64              `protobuf:"varint,1,opt,name=Version,proto3" json:"Version,omitempty"`
	Type                 []byte              `protobuf:"bytes,2,opt,name=Type,proto3" json:"Type,omitempty"`
	Identifier           []byte              `protobuf:"bytes,3,opt,name=Identifier,proto3" json:"Identifier,omitempty"`
	MetaNet              *MetaNet            `protobuf:"bytes,4,opt,name=MetaNet,proto3" json:"MetaNet,omitempty"`
	EncryptedPayloads    []*EncryptedPayload `protobuf:"bytes,5,rep,name=EncryptedPayloads,proto3" json:"EncryptedPayloads,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *Envelope) Reset()         { *m = Envelope{} }
func (m *Envelope) String() string { return proto.CompactTextString(m) }
func (*Envelope) ProtoMessage()    {}
func (*Envelope) Descriptor() ([]byte, []int) {
	return fileDescriptor_ee266e8c558e9dc5, []int{0}
}

func (m *Envelope) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Envelope.Unmarshal(m, b)
}
func (m *Envelope) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Envelope.Marshal(b, m, deterministic)
}
func (m *Envelope) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Envelope.Merge(m, src)
}
func (m *Envelope) XXX_Size() int {
	return xxx_messageInfo_Envelope.Size(m)
}
func (m *Envelope) XXX_DiscardUnknown() {
	xxx_messageInfo_Envelope.DiscardUnknown(m)
}

var xxx_messageInfo_Envelope proto.InternalMessageInfo

func (m *Envelope) GetVersion() uint64 {
	if m != nil {
		return m.Version
	}
	return 0
}

func (m *Envelope) GetType() []byte {
	if m != nil {
		return m.Type
	}
	return nil
}

func (m *Envelope) GetIdentifier() []byte {
	if m != nil {
		return m.Identifier
	}
	return nil
}

func (m *Envelope) GetMetaNet() *MetaNet {
	if m != nil {
		return m.MetaNet
	}
	return nil
}

func (m *Envelope) GetEncryptedPayloads() []*EncryptedPayload {
	if m != nil {
		return m.EncryptedPayloads
	}
	return nil
}

type MetaNet struct {
	Index                uint32   `protobuf:"varint,1,opt,name=Index,proto3" json:"Index,omitempty"`
	Parent               []byte   `protobuf:"bytes,2,opt,name=Parent,proto3" json:"Parent,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MetaNet) Reset()         { *m = MetaNet{} }
func (m *MetaNet) String() string { return proto.CompactTextString(m) }
func (*MetaNet) ProtoMessage()    {}
func (*MetaNet) Descriptor() ([]byte, []int) {
	return fileDescriptor_ee266e8c558e9dc5, []int{1}
}

func (m *MetaNet) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MetaNet.Unmarshal(m, b)
}
func (m *MetaNet) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MetaNet.Marshal(b, m, deterministic)
}
func (m *MetaNet) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MetaNet.Merge(m, src)
}
func (m *MetaNet) XXX_Size() int {
	return xxx_messageInfo_MetaNet.Size(m)
}
func (m *MetaNet) XXX_DiscardUnknown() {
	xxx_messageInfo_MetaNet.DiscardUnknown(m)
}

var xxx_messageInfo_MetaNet proto.InternalMessageInfo

func (m *MetaNet) GetIndex() uint32 {
	if m != nil {
		return m.Index
	}
	return 0
}

func (m *MetaNet) GetParent() []byte {
	if m != nil {
		return m.Parent
	}
	return nil
}

type EncryptedPayload struct {
	Sender               uint32      `protobuf:"varint,1,opt,name=Sender,proto3" json:"Sender,omitempty"`
	Receivers            []*Receiver `protobuf:"bytes,2,rep,name=Receivers,proto3" json:"Receivers,omitempty"`
	Payload              []byte      `protobuf:"bytes,3,opt,name=Payload,proto3" json:"Payload,omitempty"`
	EncryptionType       uint32      `protobuf:"varint,4,opt,name=EncryptionType,proto3" json:"EncryptionType,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *EncryptedPayload) Reset()         { *m = EncryptedPayload{} }
func (m *EncryptedPayload) String() string { return proto.CompactTextString(m) }
func (*EncryptedPayload) ProtoMessage()    {}
func (*EncryptedPayload) Descriptor() ([]byte, []int) {
	return fileDescriptor_ee266e8c558e9dc5, []int{2}
}

func (m *EncryptedPayload) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EncryptedPayload.Unmarshal(m, b)
}
func (m *EncryptedPayload) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EncryptedPayload.Marshal(b, m, deterministic)
}
func (m *EncryptedPayload) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EncryptedPayload.Merge(m, src)
}
func (m *EncryptedPayload) XXX_Size() int {
	return xxx_messageInfo_EncryptedPayload.Size(m)
}
func (m *EncryptedPayload) XXX_DiscardUnknown() {
	xxx_messageInfo_EncryptedPayload.DiscardUnknown(m)
}

var xxx_messageInfo_EncryptedPayload proto.InternalMessageInfo

func (m *EncryptedPayload) GetSender() uint32 {
	if m != nil {
		return m.Sender
	}
	return 0
}

func (m *EncryptedPayload) GetReceivers() []*Receiver {
	if m != nil {
		return m.Receivers
	}
	return nil
}

func (m *EncryptedPayload) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (m *EncryptedPayload) GetEncryptionType() uint32 {
	if m != nil {
		return m.EncryptionType
	}
	return 0
}

type Receiver struct {
	Index                uint32   `protobuf:"varint,1,opt,name=Index,proto3" json:"Index,omitempty"`
	EncryptedKey         []byte   `protobuf:"bytes,2,opt,name=EncryptedKey,proto3" json:"EncryptedKey,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Receiver) Reset()         { *m = Receiver{} }
func (m *Receiver) String() string { return proto.CompactTextString(m) }
func (*Receiver) ProtoMessage()    {}
func (*Receiver) Descriptor() ([]byte, []int) {
	return fileDescriptor_ee266e8c558e9dc5, []int{3}
}

func (m *Receiver) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Receiver.Unmarshal(m, b)
}
func (m *Receiver) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Receiver.Marshal(b, m, deterministic)
}
func (m *Receiver) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Receiver.Merge(m, src)
}
func (m *Receiver) XXX_Size() int {
	return xxx_messageInfo_Receiver.Size(m)
}
func (m *Receiver) XXX_DiscardUnknown() {
	xxx_messageInfo_Receiver.DiscardUnknown(m)
}

var xxx_messageInfo_Receiver proto.InternalMessageInfo

func (m *Receiver) GetIndex() uint32 {
	if m != nil {
		return m.Index
	}
	return 0
}

func (m *Receiver) GetEncryptedKey() []byte {
	if m != nil {
		return m.EncryptedKey
	}
	return nil
}

func init() {
	proto.RegisterType((*Envelope)(nil), "protobuf.Envelope")
	proto.RegisterType((*MetaNet)(nil), "protobuf.MetaNet")
	proto.RegisterType((*EncryptedPayload)(nil), "protobuf.EncryptedPayload")
	proto.RegisterType((*Receiver)(nil), "protobuf.Receiver")
}

func init() { proto.RegisterFile("envelope.proto", fileDescriptor_ee266e8c558e9dc5) }

var fileDescriptor_ee266e8c558e9dc5 = []byte{
	// 339 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x52, 0x4d, 0x4b, 0xeb, 0x40,
	0x14, 0x25, 0x6d, 0xfa, 0xf1, 0x6e, 0xfb, 0xca, 0xeb, 0xf0, 0x78, 0x0c, 0x6f, 0x21, 0x21, 0x0b,
	0x09, 0x08, 0x49, 0xa9, 0x0b, 0x57, 0x22, 0x88, 0x05, 0x8b, 0x28, 0x65, 0x14, 0x17, 0xee, 0xd2,
	0xe6, 0x36, 0x0e, 0xad, 0x33, 0x61, 0x3a, 0x2d, 0xc6, 0x1f, 0xe3, 0x5f, 0xf2, 0x2f, 0x49, 0xc7,
	0x99, 0x46, 0x2b, 0xae, 0x32, 0xe7, 0x9e, 0x33, 0xe7, 0xcc, 0xb9, 0x04, 0x7a, 0x28, 0x36, 0xb8,
	0x94, 0x05, 0xc6, 0x85, 0x92, 0x5a, 0x92, 0xb6, 0xf9, 0x4c, 0xd7, 0xf3, 0xf0, 0xcd, 0x83, 0xf6,
	0xc8, 0x92, 0x84, 0x42, 0xeb, 0x1e, 0xd5, 0x8a, 0x4b, 0x41, 0xbd, 0xc0, 0x8b, 0x7c, 0xe6, 0x20,
	0x21, 0xe0, 0xdf, 0x95, 0x05, 0xd2, 0x5a, 0xe0, 0x45, 0x5d, 0x66, 0xce, 0xe4, 0x00, 0x60, 0x9c,
	0xa1, 0xd0, 0x7c, 0xce, 0x51, 0xd1, 0xba, 0x61, 0x3e, 0x4d, 0xc8, 0x11, 0xb4, 0xae, 0x51, 0xa7,
	0x37, 0xa8, 0xa9, 0x1f, 0x78, 0x51, 0x67, 0xd8, 0x8f, 0x5d, 0x6c, 0x6c, 0x09, 0xe6, 0x14, 0xe4,
	0x12, 0xfa, 0x23, 0x31, 0x53, 0x65, 0xa1, 0x31, 0x9b, 0xa4, 0xe5, 0x52, 0xa6, 0xd9, 0x8a, 0x36,
	0x82, 0x7a, 0xd4, 0x19, 0xfe, 0xaf, 0xae, 0xed, 0x4b, 0xd8, 0xf7, 0x4b, 0xe1, 0xc9, 0x2e, 0x96,
	0xfc, 0x85, 0xc6, 0x58, 0x64, 0xf8, 0x6c, 0xda, 0xfc, 0x66, 0x1f, 0x80, 0xfc, 0x83, 0xe6, 0x24,
	0x55, 0x28, 0xb4, 0x6d, 0x63, 0x51, 0xf8, 0xea, 0xc1, 0x9f, 0x7d, 0xbb, 0xad, 0xf8, 0x16, 0x45,
	0x86, 0xca, 0x7a, 0x58, 0x44, 0x06, 0xf0, 0x8b, 0xe1, 0x0c, 0xf9, 0x06, 0xd5, 0x8a, 0xd6, 0xcc,
	0x3b, 0x49, 0xf5, 0x4e, 0x47, 0xb1, 0x4a, 0xb4, 0x5d, 0xae, 0x35, 0xb5, 0xbb, 0x72, 0x90, 0x1c,
	0x42, 0xcf, 0xe6, 0x72, 0x29, 0xcc, 0x9a, 0x7d, 0x93, 0xb5, 0x37, 0x0d, 0x2f, 0xa0, 0xed, 0xec,
	0x7e, 0xa8, 0x16, 0x42, 0x77, 0xd7, 0xe0, 0x0a, 0x4b, 0x5b, 0xf0, 0xcb, 0xec, 0xfc, 0xec, 0xe1,
	0x34, 0xe7, 0xfa, 0x71, 0x3d, 0x8d, 0x67, 0xf2, 0x29, 0xd1, 0x72, 0x81, 0x82, 0xbf, 0x60, 0x96,
	0xb8, 0x5f, 0x24, 0x29, 0x16, 0x79, 0x92, 0xcb, 0x65, 0x2a, 0xf2, 0x6a, 0xb6, 0x19, 0x24, 0xae,
	0xdc, 0xb4, 0x69, 0x4e, 0xc7, 0xef, 0x01, 0x00, 0x00, 0xff, 0xff, 0x2a, 0x7e, 0xf4, 0x1b, 0x55,
	0x02, 0x00, 0x00,
}
