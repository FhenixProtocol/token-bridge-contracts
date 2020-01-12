// Code generated by protoc-gen-go. DO NOT EDIT.
// source: valmessage.proto

package valmessage

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	common "github.com/offchainlabs/arbitrum/packages/arb-util/common"
	protocol "github.com/offchainlabs/arbitrum/packages/arb-util/protocol"
	valprotocol "github.com/offchainlabs/arbitrum/packages/arb-validator/valprotocol"
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

type TokenTypeBuf struct {
	Value                []byte   `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TokenTypeBuf) Reset()         { *m = TokenTypeBuf{} }
func (m *TokenTypeBuf) String() string { return proto.CompactTextString(m) }
func (*TokenTypeBuf) ProtoMessage()    {}
func (*TokenTypeBuf) Descriptor() ([]byte, []int) {
	return fileDescriptor_b34ccd35396e2606, []int{0}
}

func (m *TokenTypeBuf) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TokenTypeBuf.Unmarshal(m, b)
}
func (m *TokenTypeBuf) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TokenTypeBuf.Marshal(b, m, deterministic)
}
func (m *TokenTypeBuf) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TokenTypeBuf.Merge(m, src)
}
func (m *TokenTypeBuf) XXX_Size() int {
	return xxx_messageInfo_TokenTypeBuf.Size(m)
}
func (m *TokenTypeBuf) XXX_DiscardUnknown() {
	xxx_messageInfo_TokenTypeBuf.DiscardUnknown(m)
}

var xxx_messageInfo_TokenTypeBuf proto.InternalMessageInfo

func (m *TokenTypeBuf) GetValue() []byte {
	if m != nil {
		return m.Value
	}
	return nil
}

type VMConfiguration struct {
	GracePeriod           uint64                    `protobuf:"varint,1,opt,name=grace_period,json=gracePeriod,proto3" json:"grace_period,omitempty"`
	EscrowRequired        *common.BigIntegerBuf     `protobuf:"bytes,2,opt,name=escrow_required,json=escrowRequired,proto3" json:"escrow_required,omitempty"`
	EscrowCurrency        *valprotocol.AddressBuf   `protobuf:"bytes,3,opt,name=escrow_currency,json=escrowCurrency,proto3" json:"escrow_currency,omitempty"`
	AssertKeys            []*valprotocol.AddressBuf `protobuf:"bytes,4,rep,name=assert_keys,json=assertKeys,proto3" json:"assert_keys,omitempty"`
	MaxExecutionStepCount uint32                    `protobuf:"varint,5,opt,name=max_execution_step_count,json=maxExecutionStepCount,proto3" json:"max_execution_step_count,omitempty"`
	Owner                 *valprotocol.AddressBuf   `protobuf:"bytes,6,opt,name=owner,proto3" json:"owner,omitempty"`
	XXX_NoUnkeyedLiteral  struct{}                  `json:"-"`
	XXX_unrecognized      []byte                    `json:"-"`
	XXX_sizecache         int32                     `json:"-"`
}

func (m *VMConfiguration) Reset()         { *m = VMConfiguration{} }
func (m *VMConfiguration) String() string { return proto.CompactTextString(m) }
func (*VMConfiguration) ProtoMessage()    {}
func (*VMConfiguration) Descriptor() ([]byte, []int) {
	return fileDescriptor_b34ccd35396e2606, []int{1}
}

func (m *VMConfiguration) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_VMConfiguration.Unmarshal(m, b)
}
func (m *VMConfiguration) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_VMConfiguration.Marshal(b, m, deterministic)
}
func (m *VMConfiguration) XXX_Merge(src proto.Message) {
	xxx_messageInfo_VMConfiguration.Merge(m, src)
}
func (m *VMConfiguration) XXX_Size() int {
	return xxx_messageInfo_VMConfiguration.Size(m)
}
func (m *VMConfiguration) XXX_DiscardUnknown() {
	xxx_messageInfo_VMConfiguration.DiscardUnknown(m)
}

var xxx_messageInfo_VMConfiguration proto.InternalMessageInfo

func (m *VMConfiguration) GetGracePeriod() uint64 {
	if m != nil {
		return m.GracePeriod
	}
	return 0
}

func (m *VMConfiguration) GetEscrowRequired() *common.BigIntegerBuf {
	if m != nil {
		return m.EscrowRequired
	}
	return nil
}

func (m *VMConfiguration) GetEscrowCurrency() *valprotocol.AddressBuf {
	if m != nil {
		return m.EscrowCurrency
	}
	return nil
}

func (m *VMConfiguration) GetAssertKeys() []*valprotocol.AddressBuf {
	if m != nil {
		return m.AssertKeys
	}
	return nil
}

func (m *VMConfiguration) GetMaxExecutionStepCount() uint32 {
	if m != nil {
		return m.MaxExecutionStepCount
	}
	return 0
}

func (m *VMConfiguration) GetOwner() *valprotocol.AddressBuf {
	if m != nil {
		return m.Owner
	}
	return nil
}

type UnanimousAssertionValidatorNotification struct {
	Accepted             bool     `protobuf:"varint,1,opt,name=accepted,proto3" json:"accepted,omitempty"`
	Signatures           [][]byte `protobuf:"bytes,2,rep,name=signatures,proto3" json:"signatures,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UnanimousAssertionValidatorNotification) Reset() {
	*m = UnanimousAssertionValidatorNotification{}
}
func (m *UnanimousAssertionValidatorNotification) String() string { return proto.CompactTextString(m) }
func (*UnanimousAssertionValidatorNotification) ProtoMessage()    {}
func (*UnanimousAssertionValidatorNotification) Descriptor() ([]byte, []int) {
	return fileDescriptor_b34ccd35396e2606, []int{2}
}

func (m *UnanimousAssertionValidatorNotification) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UnanimousAssertionValidatorNotification.Unmarshal(m, b)
}
func (m *UnanimousAssertionValidatorNotification) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UnanimousAssertionValidatorNotification.Marshal(b, m, deterministic)
}
func (m *UnanimousAssertionValidatorNotification) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UnanimousAssertionValidatorNotification.Merge(m, src)
}
func (m *UnanimousAssertionValidatorNotification) XXX_Size() int {
	return xxx_messageInfo_UnanimousAssertionValidatorNotification.Size(m)
}
func (m *UnanimousAssertionValidatorNotification) XXX_DiscardUnknown() {
	xxx_messageInfo_UnanimousAssertionValidatorNotification.DiscardUnknown(m)
}

var xxx_messageInfo_UnanimousAssertionValidatorNotification proto.InternalMessageInfo

func (m *UnanimousAssertionValidatorNotification) GetAccepted() bool {
	if m != nil {
		return m.Accepted
	}
	return false
}

func (m *UnanimousAssertionValidatorNotification) GetSignatures() [][]byte {
	if m != nil {
		return m.Signatures
	}
	return nil
}

type SignedMessage struct {
	Message              *valprotocol.MessageBuf `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	Signature            []byte                  `protobuf:"bytes,2,opt,name=signature,proto3" json:"signature,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                `json:"-"`
	XXX_unrecognized     []byte                  `json:"-"`
	XXX_sizecache        int32                   `json:"-"`
}

func (m *SignedMessage) Reset()         { *m = SignedMessage{} }
func (m *SignedMessage) String() string { return proto.CompactTextString(m) }
func (*SignedMessage) ProtoMessage()    {}
func (*SignedMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_b34ccd35396e2606, []int{3}
}

func (m *SignedMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SignedMessage.Unmarshal(m, b)
}
func (m *SignedMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SignedMessage.Marshal(b, m, deterministic)
}
func (m *SignedMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SignedMessage.Merge(m, src)
}
func (m *SignedMessage) XXX_Size() int {
	return xxx_messageInfo_SignedMessage.Size(m)
}
func (m *SignedMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_SignedMessage.DiscardUnknown(m)
}

var xxx_messageInfo_SignedMessage proto.InternalMessageInfo

func (m *SignedMessage) GetMessage() *valprotocol.MessageBuf {
	if m != nil {
		return m.Message
	}
	return nil
}

func (m *SignedMessage) GetSignature() []byte {
	if m != nil {
		return m.Signature
	}
	return nil
}

type UnanimousAssertionValidatorRequest struct {
	BeforeHash           *common.HashBuf               `protobuf:"bytes,1,opt,name=beforeHash,proto3" json:"beforeHash,omitempty"`
	BeforeInbox          *common.HashBuf               `protobuf:"bytes,2,opt,name=beforeInbox,proto3" json:"beforeInbox,omitempty"`
	SequenceNum          uint64                        `protobuf:"varint,3,opt,name=sequenceNum,proto3" json:"sequenceNum,omitempty"`
	TimeBounds           *protocol.TimeBoundsBlocksBuf `protobuf:"bytes,4,opt,name=timeBounds,proto3" json:"timeBounds,omitempty"`
	SignedMessages       []*SignedMessage              `protobuf:"bytes,5,rep,name=signedMessages,proto3" json:"signedMessages,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                      `json:"-"`
	XXX_unrecognized     []byte                        `json:"-"`
	XXX_sizecache        int32                         `json:"-"`
}

func (m *UnanimousAssertionValidatorRequest) Reset()         { *m = UnanimousAssertionValidatorRequest{} }
func (m *UnanimousAssertionValidatorRequest) String() string { return proto.CompactTextString(m) }
func (*UnanimousAssertionValidatorRequest) ProtoMessage()    {}
func (*UnanimousAssertionValidatorRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_b34ccd35396e2606, []int{4}
}

func (m *UnanimousAssertionValidatorRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UnanimousAssertionValidatorRequest.Unmarshal(m, b)
}
func (m *UnanimousAssertionValidatorRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UnanimousAssertionValidatorRequest.Marshal(b, m, deterministic)
}
func (m *UnanimousAssertionValidatorRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UnanimousAssertionValidatorRequest.Merge(m, src)
}
func (m *UnanimousAssertionValidatorRequest) XXX_Size() int {
	return xxx_messageInfo_UnanimousAssertionValidatorRequest.Size(m)
}
func (m *UnanimousAssertionValidatorRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_UnanimousAssertionValidatorRequest.DiscardUnknown(m)
}

var xxx_messageInfo_UnanimousAssertionValidatorRequest proto.InternalMessageInfo

func (m *UnanimousAssertionValidatorRequest) GetBeforeHash() *common.HashBuf {
	if m != nil {
		return m.BeforeHash
	}
	return nil
}

func (m *UnanimousAssertionValidatorRequest) GetBeforeInbox() *common.HashBuf {
	if m != nil {
		return m.BeforeInbox
	}
	return nil
}

func (m *UnanimousAssertionValidatorRequest) GetSequenceNum() uint64 {
	if m != nil {
		return m.SequenceNum
	}
	return 0
}

func (m *UnanimousAssertionValidatorRequest) GetTimeBounds() *protocol.TimeBoundsBlocksBuf {
	if m != nil {
		return m.TimeBounds
	}
	return nil
}

func (m *UnanimousAssertionValidatorRequest) GetSignedMessages() []*SignedMessage {
	if m != nil {
		return m.SignedMessages
	}
	return nil
}

type ValidatorRequest struct {
	RequestId *common.HashBuf `protobuf:"bytes,1,opt,name=request_id,json=requestId,proto3" json:"request_id,omitempty"`
	// Types that are valid to be assigned to Request:
	//	*ValidatorRequest_Unanimous
	//	*ValidatorRequest_UnanimousNotification
	Request              isValidatorRequest_Request `protobuf_oneof:"request"`
	XXX_NoUnkeyedLiteral struct{}                   `json:"-"`
	XXX_unrecognized     []byte                     `json:"-"`
	XXX_sizecache        int32                      `json:"-"`
}

func (m *ValidatorRequest) Reset()         { *m = ValidatorRequest{} }
func (m *ValidatorRequest) String() string { return proto.CompactTextString(m) }
func (*ValidatorRequest) ProtoMessage()    {}
func (*ValidatorRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_b34ccd35396e2606, []int{5}
}

func (m *ValidatorRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ValidatorRequest.Unmarshal(m, b)
}
func (m *ValidatorRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ValidatorRequest.Marshal(b, m, deterministic)
}
func (m *ValidatorRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ValidatorRequest.Merge(m, src)
}
func (m *ValidatorRequest) XXX_Size() int {
	return xxx_messageInfo_ValidatorRequest.Size(m)
}
func (m *ValidatorRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ValidatorRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ValidatorRequest proto.InternalMessageInfo

func (m *ValidatorRequest) GetRequestId() *common.HashBuf {
	if m != nil {
		return m.RequestId
	}
	return nil
}

type isValidatorRequest_Request interface {
	isValidatorRequest_Request()
}

type ValidatorRequest_Unanimous struct {
	Unanimous *UnanimousAssertionValidatorRequest `protobuf:"bytes,2,opt,name=unanimous,proto3,oneof"`
}

type ValidatorRequest_UnanimousNotification struct {
	UnanimousNotification *UnanimousAssertionValidatorNotification `protobuf:"bytes,3,opt,name=unanimousNotification,proto3,oneof"`
}

func (*ValidatorRequest_Unanimous) isValidatorRequest_Request() {}

func (*ValidatorRequest_UnanimousNotification) isValidatorRequest_Request() {}

func (m *ValidatorRequest) GetRequest() isValidatorRequest_Request {
	if m != nil {
		return m.Request
	}
	return nil
}

func (m *ValidatorRequest) GetUnanimous() *UnanimousAssertionValidatorRequest {
	if x, ok := m.GetRequest().(*ValidatorRequest_Unanimous); ok {
		return x.Unanimous
	}
	return nil
}

func (m *ValidatorRequest) GetUnanimousNotification() *UnanimousAssertionValidatorNotification {
	if x, ok := m.GetRequest().(*ValidatorRequest_UnanimousNotification); ok {
		return x.UnanimousNotification
	}
	return nil
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*ValidatorRequest) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*ValidatorRequest_Unanimous)(nil),
		(*ValidatorRequest_UnanimousNotification)(nil),
	}
}

type UnanimousAssertionFollowerResponse struct {
	Accepted             bool            `protobuf:"varint,1,opt,name=accepted,proto3" json:"accepted,omitempty"`
	Signature            []byte          `protobuf:"bytes,2,opt,name=signature,proto3" json:"signature,omitempty"`
	AssertionHash        *common.HashBuf `protobuf:"bytes,3,opt,name=assertion_hash,json=assertionHash,proto3" json:"assertion_hash,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *UnanimousAssertionFollowerResponse) Reset()         { *m = UnanimousAssertionFollowerResponse{} }
func (m *UnanimousAssertionFollowerResponse) String() string { return proto.CompactTextString(m) }
func (*UnanimousAssertionFollowerResponse) ProtoMessage()    {}
func (*UnanimousAssertionFollowerResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_b34ccd35396e2606, []int{6}
}

func (m *UnanimousAssertionFollowerResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UnanimousAssertionFollowerResponse.Unmarshal(m, b)
}
func (m *UnanimousAssertionFollowerResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UnanimousAssertionFollowerResponse.Marshal(b, m, deterministic)
}
func (m *UnanimousAssertionFollowerResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UnanimousAssertionFollowerResponse.Merge(m, src)
}
func (m *UnanimousAssertionFollowerResponse) XXX_Size() int {
	return xxx_messageInfo_UnanimousAssertionFollowerResponse.Size(m)
}
func (m *UnanimousAssertionFollowerResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_UnanimousAssertionFollowerResponse.DiscardUnknown(m)
}

var xxx_messageInfo_UnanimousAssertionFollowerResponse proto.InternalMessageInfo

func (m *UnanimousAssertionFollowerResponse) GetAccepted() bool {
	if m != nil {
		return m.Accepted
	}
	return false
}

func (m *UnanimousAssertionFollowerResponse) GetSignature() []byte {
	if m != nil {
		return m.Signature
	}
	return nil
}

func (m *UnanimousAssertionFollowerResponse) GetAssertionHash() *common.HashBuf {
	if m != nil {
		return m.AssertionHash
	}
	return nil
}

type FollowerResponse struct {
	RequestId            *common.HashBuf                     `protobuf:"bytes,1,opt,name=request_id,json=requestId,proto3" json:"request_id,omitempty"`
	Unanimous            *UnanimousAssertionFollowerResponse `protobuf:"bytes,3,opt,name=unanimous,proto3" json:"unanimous,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                            `json:"-"`
	XXX_unrecognized     []byte                              `json:"-"`
	XXX_sizecache        int32                               `json:"-"`
}

func (m *FollowerResponse) Reset()         { *m = FollowerResponse{} }
func (m *FollowerResponse) String() string { return proto.CompactTextString(m) }
func (*FollowerResponse) ProtoMessage()    {}
func (*FollowerResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_b34ccd35396e2606, []int{7}
}

func (m *FollowerResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FollowerResponse.Unmarshal(m, b)
}
func (m *FollowerResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FollowerResponse.Marshal(b, m, deterministic)
}
func (m *FollowerResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FollowerResponse.Merge(m, src)
}
func (m *FollowerResponse) XXX_Size() int {
	return xxx_messageInfo_FollowerResponse.Size(m)
}
func (m *FollowerResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_FollowerResponse.DiscardUnknown(m)
}

var xxx_messageInfo_FollowerResponse proto.InternalMessageInfo

func (m *FollowerResponse) GetRequestId() *common.HashBuf {
	if m != nil {
		return m.RequestId
	}
	return nil
}

func (m *FollowerResponse) GetUnanimous() *UnanimousAssertionFollowerResponse {
	if m != nil {
		return m.Unanimous
	}
	return nil
}

func init() {
	proto.RegisterType((*TokenTypeBuf)(nil), "valmessage.TokenTypeBuf")
	proto.RegisterType((*VMConfiguration)(nil), "valmessage.VMConfiguration")
	proto.RegisterType((*UnanimousAssertionValidatorNotification)(nil), "valmessage.UnanimousAssertionValidatorNotification")
	proto.RegisterType((*SignedMessage)(nil), "valmessage.SignedMessage")
	proto.RegisterType((*UnanimousAssertionValidatorRequest)(nil), "valmessage.UnanimousAssertionValidatorRequest")
	proto.RegisterType((*ValidatorRequest)(nil), "valmessage.ValidatorRequest")
	proto.RegisterType((*UnanimousAssertionFollowerResponse)(nil), "valmessage.UnanimousAssertionFollowerResponse")
	proto.RegisterType((*FollowerResponse)(nil), "valmessage.FollowerResponse")
}

func init() { proto.RegisterFile("valmessage.proto", fileDescriptor_b34ccd35396e2606) }

var fileDescriptor_b34ccd35396e2606 = []byte{
	// 707 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x53, 0xdd, 0x4e, 0xdb, 0x48,
	0x18, 0x25, 0x09, 0x01, 0xf2, 0x25, 0xfc, 0x68, 0x76, 0x11, 0x5e, 0xb4, 0xac, 0xb2, 0xd6, 0x4a,
	0x9b, 0x1b, 0x12, 0xb1, 0x48, 0xdb, 0xde, 0xb4, 0x2a, 0xa6, 0xad, 0x40, 0x2d, 0xa8, 0x32, 0x94,
	0x8b, 0xde, 0xb8, 0x13, 0xfb, 0x8b, 0x33, 0x8a, 0x3d, 0xe3, 0xce, 0x78, 0x42, 0xf2, 0x0a, 0xbd,
	0xea, 0x55, 0x9f, 0xa2, 0x0f, 0xd4, 0xc7, 0xa9, 0xfc, 0x13, 0xc7, 0x0d, 0x34, 0xa5, 0x57, 0xf6,
	0x9c, 0x39, 0xdf, 0xef, 0x39, 0x03, 0x3b, 0x63, 0x1a, 0x84, 0xa8, 0x14, 0xf5, 0xb1, 0x1b, 0x49,
	0x11, 0x0b, 0x02, 0x73, 0x64, 0xff, 0x60, 0x4c, 0x83, 0x14, 0x75, 0x45, 0xd0, 0x2b, 0xfd, 0x67,
	0xd4, 0xfd, 0xbd, 0xe2, 0x6e, 0xe1, 0xe2, 0x37, 0x57, 0x84, 0xa1, 0xe0, 0xbd, 0xec, 0x93, 0x81,
	0xe6, 0x3f, 0xd0, 0xba, 0x16, 0x23, 0xe4, 0xd7, 0xd3, 0x08, 0x2d, 0x3d, 0x20, 0xbf, 0x43, 0x7d,
	0x4c, 0x03, 0x8d, 0x46, 0xa5, 0x5d, 0xe9, 0xb4, 0xec, 0xec, 0x60, 0x7e, 0xad, 0xc2, 0xf6, 0xcd,
	0xc5, 0xa9, 0xe0, 0x03, 0xe6, 0x6b, 0x49, 0x63, 0x26, 0x38, 0xf9, 0x1b, 0x5a, 0xbe, 0xa4, 0x2e,
	0x3a, 0x11, 0x4a, 0x26, 0xbc, 0x34, 0x60, 0xd5, 0x6e, 0xa6, 0xd8, 0x9b, 0x14, 0x22, 0x4f, 0x61,
	0x1b, 0x95, 0x2b, 0xc5, 0xad, 0x23, 0xf1, 0x83, 0x66, 0x12, 0x3d, 0xa3, 0xda, 0xae, 0x74, 0x9a,
	0xff, 0xed, 0x76, 0xf3, 0x26, 0x2c, 0xe6, 0x9f, 0xf3, 0x18, 0x7d, 0x94, 0x96, 0x1e, 0xd8, 0x5b,
	0x19, 0xdb, 0xce, 0xc9, 0xe4, 0x59, 0x11, 0xef, 0x6a, 0x29, 0x91, 0xbb, 0x53, 0xa3, 0x96, 0xc6,
	0xef, 0x75, 0xcb, 0x73, 0x9f, 0x78, 0x9e, 0x44, 0xa5, 0x4a, 0x19, 0x4e, 0x73, 0x3a, 0x79, 0x0c,
	0x4d, 0xaa, 0x14, 0xca, 0xd8, 0x19, 0xe1, 0x54, 0x19, 0xab, 0xed, 0xda, 0xb2, 0x68, 0xc8, 0xb8,
	0xaf, 0x70, 0xaa, 0xc8, 0x23, 0x30, 0x42, 0x3a, 0x71, 0x70, 0x82, 0xae, 0x4e, 0xe6, 0x75, 0x54,
	0x8c, 0x91, 0xe3, 0x0a, 0xcd, 0x63, 0xa3, 0xde, 0xae, 0x74, 0x36, 0xed, 0xdd, 0x90, 0x4e, 0x5e,
	0xcc, 0xae, 0xaf, 0x62, 0x8c, 0x4e, 0x93, 0x4b, 0x72, 0x08, 0x75, 0x71, 0xcb, 0x51, 0x1a, 0x6b,
	0xcb, 0x5b, 0xcd, 0x58, 0x26, 0xc2, 0xbf, 0x6f, 0x39, 0xe5, 0x2c, 0x14, 0x5a, 0x9d, 0xa4, 0xe5,
	0x99, 0xe0, 0x37, 0x34, 0x60, 0x1e, 0x8d, 0x85, 0xbc, 0x14, 0x31, 0x1b, 0x30, 0x37, 0xdb, 0xf8,
	0x3e, 0x6c, 0x50, 0xd7, 0xc5, 0x28, 0xc6, 0x6c, 0xdb, 0x1b, 0x76, 0x71, 0x26, 0x7f, 0x01, 0x28,
	0xe6, 0x73, 0x1a, 0x6b, 0x89, 0xca, 0xa8, 0xb6, 0x6b, 0x9d, 0x96, 0x5d, 0x42, 0xcc, 0xf7, 0xb0,
	0x79, 0xc5, 0x7c, 0x8e, 0xde, 0x45, 0xe6, 0x22, 0x72, 0x04, 0xeb, 0xb9, 0xa1, 0xd2, 0x5c, 0x8b,
	0x8d, 0xe6, 0xb4, 0xa4, 0xd1, 0x19, 0x8f, 0xfc, 0x09, 0x8d, 0x22, 0x63, 0x2a, 0x64, 0xcb, 0x9e,
	0x03, 0xe6, 0x97, 0x2a, 0x98, 0x4b, 0x26, 0x49, 0x44, 0x45, 0x15, 0x93, 0x1e, 0x40, 0x1f, 0x07,
	0x42, 0xe2, 0x19, 0x55, 0xc3, 0xbc, 0xf4, 0xf6, 0xcc, 0x0e, 0x09, 0x96, 0x0a, 0x31, 0xa7, 0x90,
	0x23, 0x68, 0x66, 0xa7, 0x73, 0xde, 0x17, 0x93, 0xdc, 0x40, 0x77, 0x22, 0xca, 0x1c, 0xd2, 0x86,
	0xa6, 0x4a, 0xca, 0x71, 0x17, 0x2f, 0x75, 0x98, 0x7a, 0x66, 0xd5, 0x2e, 0x43, 0xe4, 0x09, 0x40,
	0xcc, 0x42, 0xb4, 0x84, 0xe6, 0x5e, 0x62, 0x8b, 0x24, 0xe7, 0x41, 0xb7, 0x98, 0xfe, 0xba, 0xb8,
	0xb3, 0x02, 0xe1, 0x8e, 0x32, 0x73, 0xcc, 0x03, 0xc8, 0x09, 0x6c, 0xa9, 0xf2, 0x36, 0x95, 0x51,
	0x4f, 0x9d, 0xf5, 0x47, 0xb7, 0xf4, 0x72, 0xbf, 0xdb, 0xb7, 0xbd, 0x10, 0x60, 0x7e, 0xac, 0xc2,
	0xce, 0x9d, 0xe5, 0x74, 0x01, 0x64, 0xf6, 0xeb, 0x30, 0xef, 0x47, 0xcb, 0x69, 0xe4, 0x94, 0x73,
	0x8f, 0x5c, 0x42, 0x43, 0xcf, 0x56, 0x9e, 0x6f, 0xa6, 0x5b, 0x6e, 0xe1, 0xe7, 0x7a, 0x9c, 0xad,
	0xd8, 0xf3, 0x14, 0x64, 0x04, 0xbb, 0xc5, 0xa1, 0x6c, 0xbd, 0xfc, 0xd9, 0x1d, 0x3f, 0x30, 0x77,
	0x39, 0xf4, 0x6c, 0xc5, 0xbe, 0x3f, 0xa7, 0xd5, 0x80, 0xf5, 0x7c, 0x12, 0xf3, 0x73, 0xe5, 0x3e,
	0xef, 0xbc, 0x14, 0x41, 0x20, 0x6e, 0x51, 0xda, 0xa8, 0x22, 0xc1, 0x15, 0x2e, 0x7d, 0x00, 0x4b,
	0xcd, 0x49, 0xfe, 0x87, 0x2d, 0x3a, 0x4b, 0xeb, 0x0c, 0x13, 0xe7, 0xd5, 0xee, 0x5f, 0xee, 0x66,
	0x41, 0x4b, 0x10, 0xf3, 0x53, 0x05, 0x76, 0xee, 0xb4, 0xf1, 0xab, 0x2a, 0xbd, 0x2e, 0xab, 0x54,
	0x7b, 0x88, 0x4a, 0x8b, 0x25, 0x4b, 0x1a, 0x59, 0xcf, 0xdf, 0x59, 0x3e, 0x8b, 0x87, 0xba, 0x9f,
	0x54, 0xec, 0x89, 0xc1, 0xc0, 0x1d, 0x52, 0xc6, 0x03, 0xda, 0x57, 0x3d, 0x2a, 0xfb, 0x2c, 0x96,
	0x3a, 0xec, 0x45, 0xd4, 0x1d, 0x25, 0x2e, 0x4b, 0x90, 0xc3, 0xf1, 0x4c, 0x96, 0xde, 0xbc, 0x66,
	0x7f, 0x2d, 0xf5, 0xfa, 0xf1, 0xb7, 0x00, 0x00, 0x00, 0xff, 0xff, 0x6c, 0x1e, 0xad, 0xae, 0x6b,
	0x06, 0x00, 0x00,
}
