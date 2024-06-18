// Copyright 2016 Prometheus Team
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        (unknown)
// source: remote.proto

package _go

import (
	_ "buf.build/gen/go/gogo/protobuf/protocolbuffers/go/gogoproto"
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

type ReadRequest_ResponseType int32

const (
	// Server will return a single ReadResponse message with matched series that includes list of raw samples.
	// It's recommended to use streamed response types instead.
	//
	// Response headers:
	// Content-Type: "application/x-protobuf"
	// Content-Encoding: "snappy"
	ReadRequest_SAMPLES ReadRequest_ResponseType = 0
	// Server will stream a delimited ChunkedReadResponse message that
	// contains XOR or HISTOGRAM(!) encoded chunks for a single series.
	// Each message is following varint size and fixed size bigendian
	// uint32 for CRC32 Castagnoli checksum.
	//
	// Response headers:
	// Content-Type: "application/x-streamed-protobuf; proto=prometheus.ChunkedReadResponse"
	// Content-Encoding: ""
	ReadRequest_STREAMED_XOR_CHUNKS ReadRequest_ResponseType = 1
)

// Enum value maps for ReadRequest_ResponseType.
var (
	ReadRequest_ResponseType_name = map[int32]string{
		0: "SAMPLES",
		1: "STREAMED_XOR_CHUNKS",
	}
	ReadRequest_ResponseType_value = map[string]int32{
		"SAMPLES":             0,
		"STREAMED_XOR_CHUNKS": 1,
	}
)

func (x ReadRequest_ResponseType) Enum() *ReadRequest_ResponseType {
	p := new(ReadRequest_ResponseType)
	*p = x
	return p
}

func (x ReadRequest_ResponseType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ReadRequest_ResponseType) Descriptor() protoreflect.EnumDescriptor {
	return file_remote_proto_enumTypes[0].Descriptor()
}

func (ReadRequest_ResponseType) Type() protoreflect.EnumType {
	return &file_remote_proto_enumTypes[0]
}

func (x ReadRequest_ResponseType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ReadRequest_ResponseType.Descriptor instead.
func (ReadRequest_ResponseType) EnumDescriptor() ([]byte, []int) {
	return file_remote_proto_rawDescGZIP(), []int{1, 0}
}

type WriteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Timeseries []*TimeSeries     `protobuf:"bytes,1,rep,name=timeseries,proto3" json:"timeseries,omitempty"`
	Metadata   []*MetricMetadata `protobuf:"bytes,3,rep,name=metadata,proto3" json:"metadata,omitempty"`
}

func (x *WriteRequest) Reset() {
	*x = WriteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_remote_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WriteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WriteRequest) ProtoMessage() {}

func (x *WriteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_remote_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WriteRequest.ProtoReflect.Descriptor instead.
func (*WriteRequest) Descriptor() ([]byte, []int) {
	return file_remote_proto_rawDescGZIP(), []int{0}
}

func (x *WriteRequest) GetTimeseries() []*TimeSeries {
	if x != nil {
		return x.Timeseries
	}
	return nil
}

func (x *WriteRequest) GetMetadata() []*MetricMetadata {
	if x != nil {
		return x.Metadata
	}
	return nil
}

// ReadRequest represents a remote read request.
type ReadRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Queries []*Query `protobuf:"bytes,1,rep,name=queries,proto3" json:"queries,omitempty"`
	// accepted_response_types allows negotiating the content type of the response.
	//
	// Response types are taken from the list in the FIFO order. If no response type in `accepted_response_types` is
	// implemented by server, error is returned.
	// For request that do not contain `accepted_response_types` field the SAMPLES response type will be used.
	AcceptedResponseTypes []ReadRequest_ResponseType `protobuf:"varint,2,rep,packed,name=accepted_response_types,json=acceptedResponseTypes,proto3,enum=prometheus.ReadRequest_ResponseType" json:"accepted_response_types,omitempty"`
}

func (x *ReadRequest) Reset() {
	*x = ReadRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_remote_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReadRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReadRequest) ProtoMessage() {}

func (x *ReadRequest) ProtoReflect() protoreflect.Message {
	mi := &file_remote_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReadRequest.ProtoReflect.Descriptor instead.
func (*ReadRequest) Descriptor() ([]byte, []int) {
	return file_remote_proto_rawDescGZIP(), []int{1}
}

func (x *ReadRequest) GetQueries() []*Query {
	if x != nil {
		return x.Queries
	}
	return nil
}

func (x *ReadRequest) GetAcceptedResponseTypes() []ReadRequest_ResponseType {
	if x != nil {
		return x.AcceptedResponseTypes
	}
	return nil
}

// ReadResponse is a response when response_type equals SAMPLES.
type ReadResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// In same order as the request's queries.
	Results []*QueryResult `protobuf:"bytes,1,rep,name=results,proto3" json:"results,omitempty"`
}

func (x *ReadResponse) Reset() {
	*x = ReadResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_remote_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReadResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReadResponse) ProtoMessage() {}

func (x *ReadResponse) ProtoReflect() protoreflect.Message {
	mi := &file_remote_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReadResponse.ProtoReflect.Descriptor instead.
func (*ReadResponse) Descriptor() ([]byte, []int) {
	return file_remote_proto_rawDescGZIP(), []int{2}
}

func (x *ReadResponse) GetResults() []*QueryResult {
	if x != nil {
		return x.Results
	}
	return nil
}

type Query struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StartTimestampMs int64           `protobuf:"varint,1,opt,name=start_timestamp_ms,json=startTimestampMs,proto3" json:"start_timestamp_ms,omitempty"`
	EndTimestampMs   int64           `protobuf:"varint,2,opt,name=end_timestamp_ms,json=endTimestampMs,proto3" json:"end_timestamp_ms,omitempty"`
	Matchers         []*LabelMatcher `protobuf:"bytes,3,rep,name=matchers,proto3" json:"matchers,omitempty"`
	Hints            *ReadHints      `protobuf:"bytes,4,opt,name=hints,proto3" json:"hints,omitempty"`
}

func (x *Query) Reset() {
	*x = Query{}
	if protoimpl.UnsafeEnabled {
		mi := &file_remote_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Query) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Query) ProtoMessage() {}

func (x *Query) ProtoReflect() protoreflect.Message {
	mi := &file_remote_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Query.ProtoReflect.Descriptor instead.
func (*Query) Descriptor() ([]byte, []int) {
	return file_remote_proto_rawDescGZIP(), []int{3}
}

func (x *Query) GetStartTimestampMs() int64 {
	if x != nil {
		return x.StartTimestampMs
	}
	return 0
}

func (x *Query) GetEndTimestampMs() int64 {
	if x != nil {
		return x.EndTimestampMs
	}
	return 0
}

func (x *Query) GetMatchers() []*LabelMatcher {
	if x != nil {
		return x.Matchers
	}
	return nil
}

func (x *Query) GetHints() *ReadHints {
	if x != nil {
		return x.Hints
	}
	return nil
}

type QueryResult struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Samples within a time series must be ordered by time.
	Timeseries []*TimeSeries `protobuf:"bytes,1,rep,name=timeseries,proto3" json:"timeseries,omitempty"`
}

func (x *QueryResult) Reset() {
	*x = QueryResult{}
	if protoimpl.UnsafeEnabled {
		mi := &file_remote_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QueryResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueryResult) ProtoMessage() {}

func (x *QueryResult) ProtoReflect() protoreflect.Message {
	mi := &file_remote_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueryResult.ProtoReflect.Descriptor instead.
func (*QueryResult) Descriptor() ([]byte, []int) {
	return file_remote_proto_rawDescGZIP(), []int{4}
}

func (x *QueryResult) GetTimeseries() []*TimeSeries {
	if x != nil {
		return x.Timeseries
	}
	return nil
}

// ChunkedReadResponse is a response when response_type equals STREAMED_XOR_CHUNKS.
// We strictly stream full series after series, optionally split by time. This means that a single frame can contain
// partition of the single series, but once a new series is started to be streamed it means that no more chunks will
// be sent for previous one. Series are returned sorted in the same way TSDB block are internally.
type ChunkedReadResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ChunkedSeries []*ChunkedSeries `protobuf:"bytes,1,rep,name=chunked_series,json=chunkedSeries,proto3" json:"chunked_series,omitempty"`
	// query_index represents an index of the query from ReadRequest.queries these chunks relates to.
	QueryIndex int64 `protobuf:"varint,2,opt,name=query_index,json=queryIndex,proto3" json:"query_index,omitempty"`
}

func (x *ChunkedReadResponse) Reset() {
	*x = ChunkedReadResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_remote_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChunkedReadResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChunkedReadResponse) ProtoMessage() {}

func (x *ChunkedReadResponse) ProtoReflect() protoreflect.Message {
	mi := &file_remote_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChunkedReadResponse.ProtoReflect.Descriptor instead.
func (*ChunkedReadResponse) Descriptor() ([]byte, []int) {
	return file_remote_proto_rawDescGZIP(), []int{5}
}

func (x *ChunkedReadResponse) GetChunkedSeries() []*ChunkedSeries {
	if x != nil {
		return x.ChunkedSeries
	}
	return nil
}

func (x *ChunkedReadResponse) GetQueryIndex() int64 {
	if x != nil {
		return x.QueryIndex
	}
	return 0
}

var File_remote_proto protoreflect.FileDescriptor

var file_remote_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a,
	0x70, 0x72, 0x6f, 0x6d, 0x65, 0x74, 0x68, 0x65, 0x75, 0x73, 0x1a, 0x0b, 0x74, 0x79, 0x70, 0x65,
	0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x14, 0x67, 0x6f, 0x67, 0x6f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2f, 0x67, 0x6f, 0x67, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x90, 0x01,
	0x0a, 0x0c, 0x57, 0x72, 0x69, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x3c,
	0x0a, 0x0a, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x65, 0x72, 0x69, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x16, 0x2e, 0x70, 0x72, 0x6f, 0x6d, 0x65, 0x74, 0x68, 0x65, 0x75, 0x73, 0x2e,
	0x54, 0x69, 0x6d, 0x65, 0x53, 0x65, 0x72, 0x69, 0x65, 0x73, 0x42, 0x04, 0xc8, 0xde, 0x1f, 0x00,
	0x52, 0x0a, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x65, 0x72, 0x69, 0x65, 0x73, 0x12, 0x3c, 0x0a, 0x08,
	0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1a,
	0x2e, 0x70, 0x72, 0x6f, 0x6d, 0x65, 0x74, 0x68, 0x65, 0x75, 0x73, 0x2e, 0x4d, 0x65, 0x74, 0x72,
	0x69, 0x63, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x42, 0x04, 0xc8, 0xde, 0x1f, 0x00,
	0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x4a, 0x04, 0x08, 0x02, 0x10, 0x03,
	0x22, 0xce, 0x01, 0x0a, 0x0b, 0x52, 0x65, 0x61, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x2b, 0x0a, 0x07, 0x71, 0x75, 0x65, 0x72, 0x69, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x11, 0x2e, 0x70, 0x72, 0x6f, 0x6d, 0x65, 0x74, 0x68, 0x65, 0x75, 0x73, 0x2e, 0x51,
	0x75, 0x65, 0x72, 0x79, 0x52, 0x07, 0x71, 0x75, 0x65, 0x72, 0x69, 0x65, 0x73, 0x12, 0x5c, 0x0a,
	0x17, 0x61, 0x63, 0x63, 0x65, 0x70, 0x74, 0x65, 0x64, 0x5f, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0e, 0x32, 0x24,
	0x2e, 0x70, 0x72, 0x6f, 0x6d, 0x65, 0x74, 0x68, 0x65, 0x75, 0x73, 0x2e, 0x52, 0x65, 0x61, 0x64,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x54, 0x79, 0x70, 0x65, 0x52, 0x15, 0x61, 0x63, 0x63, 0x65, 0x70, 0x74, 0x65, 0x64, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x54, 0x79, 0x70, 0x65, 0x73, 0x22, 0x34, 0x0a, 0x0c, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0b, 0x0a, 0x07, 0x53,
	0x41, 0x4d, 0x50, 0x4c, 0x45, 0x53, 0x10, 0x00, 0x12, 0x17, 0x0a, 0x13, 0x53, 0x54, 0x52, 0x45,
	0x41, 0x4d, 0x45, 0x44, 0x5f, 0x58, 0x4f, 0x52, 0x5f, 0x43, 0x48, 0x55, 0x4e, 0x4b, 0x53, 0x10,
	0x01, 0x22, 0x41, 0x0a, 0x0c, 0x52, 0x65, 0x61, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x31, 0x0a, 0x07, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x17, 0x2e, 0x70, 0x72, 0x6f, 0x6d, 0x65, 0x74, 0x68, 0x65, 0x75, 0x73, 0x2e,
	0x51, 0x75, 0x65, 0x72, 0x79, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x52, 0x07, 0x72, 0x65, 0x73,
	0x75, 0x6c, 0x74, 0x73, 0x22, 0xc2, 0x01, 0x0a, 0x05, 0x51, 0x75, 0x65, 0x72, 0x79, 0x12, 0x2c,
	0x0a, 0x12, 0x73, 0x74, 0x61, 0x72, 0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x5f, 0x6d, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x10, 0x73, 0x74, 0x61, 0x72,
	0x74, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x4d, 0x73, 0x12, 0x28, 0x0a, 0x10,
	0x65, 0x6e, 0x64, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x5f, 0x6d, 0x73,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0e, 0x65, 0x6e, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x4d, 0x73, 0x12, 0x34, 0x0a, 0x08, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x65,
	0x72, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x70, 0x72, 0x6f, 0x6d, 0x65,
	0x74, 0x68, 0x65, 0x75, 0x73, 0x2e, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x4d, 0x61, 0x74, 0x63, 0x68,
	0x65, 0x72, 0x52, 0x08, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x72, 0x73, 0x12, 0x2b, 0x0a, 0x05,
	0x68, 0x69, 0x6e, 0x74, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x70, 0x72,
	0x6f, 0x6d, 0x65, 0x74, 0x68, 0x65, 0x75, 0x73, 0x2e, 0x52, 0x65, 0x61, 0x64, 0x48, 0x69, 0x6e,
	0x74, 0x73, 0x52, 0x05, 0x68, 0x69, 0x6e, 0x74, 0x73, 0x22, 0x45, 0x0a, 0x0b, 0x51, 0x75, 0x65,
	0x72, 0x79, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x36, 0x0a, 0x0a, 0x74, 0x69, 0x6d, 0x65,
	0x73, 0x65, 0x72, 0x69, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x70,
	0x72, 0x6f, 0x6d, 0x65, 0x74, 0x68, 0x65, 0x75, 0x73, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x53, 0x65,
	0x72, 0x69, 0x65, 0x73, 0x52, 0x0a, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x65, 0x72, 0x69, 0x65, 0x73,
	0x22, 0x78, 0x0a, 0x13, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x65, 0x64, 0x52, 0x65, 0x61, 0x64, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x40, 0x0a, 0x0e, 0x63, 0x68, 0x75, 0x6e, 0x6b,
	0x65, 0x64, 0x5f, 0x73, 0x65, 0x72, 0x69, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x19, 0x2e, 0x70, 0x72, 0x6f, 0x6d, 0x65, 0x74, 0x68, 0x65, 0x75, 0x73, 0x2e, 0x43, 0x68, 0x75,
	0x6e, 0x6b, 0x65, 0x64, 0x53, 0x65, 0x72, 0x69, 0x65, 0x73, 0x52, 0x0d, 0x63, 0x68, 0x75, 0x6e,
	0x6b, 0x65, 0x64, 0x53, 0x65, 0x72, 0x69, 0x65, 0x73, 0x12, 0x1f, 0x0a, 0x0b, 0x71, 0x75, 0x65,
	0x72, 0x79, 0x5f, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a,
	0x71, 0x75, 0x65, 0x72, 0x79, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x42, 0x3b, 0x5a, 0x39, 0x62, 0x75,
	0x66, 0x2e, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x67, 0x6f, 0x2f, 0x70,
	0x72, 0x6f, 0x6d, 0x65, 0x74, 0x68, 0x65, 0x75, 0x73, 0x2f, 0x70, 0x72, 0x6f, 0x6d, 0x65, 0x74,
	0x68, 0x65, 0x75, 0x73, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x62, 0x75, 0x66,
	0x66, 0x65, 0x72, 0x73, 0x2f, 0x67, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_remote_proto_rawDescOnce sync.Once
	file_remote_proto_rawDescData = file_remote_proto_rawDesc
)

func file_remote_proto_rawDescGZIP() []byte {
	file_remote_proto_rawDescOnce.Do(func() {
		file_remote_proto_rawDescData = protoimpl.X.CompressGZIP(file_remote_proto_rawDescData)
	})
	return file_remote_proto_rawDescData
}

var file_remote_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_remote_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_remote_proto_goTypes = []any{
	(ReadRequest_ResponseType)(0), // 0: prometheus.ReadRequest.ResponseType
	(*WriteRequest)(nil),          // 1: prometheus.WriteRequest
	(*ReadRequest)(nil),           // 2: prometheus.ReadRequest
	(*ReadResponse)(nil),          // 3: prometheus.ReadResponse
	(*Query)(nil),                 // 4: prometheus.Query
	(*QueryResult)(nil),           // 5: prometheus.QueryResult
	(*ChunkedReadResponse)(nil),   // 6: prometheus.ChunkedReadResponse
	(*TimeSeries)(nil),            // 7: prometheus.TimeSeries
	(*MetricMetadata)(nil),        // 8: prometheus.MetricMetadata
	(*LabelMatcher)(nil),          // 9: prometheus.LabelMatcher
	(*ReadHints)(nil),             // 10: prometheus.ReadHints
	(*ChunkedSeries)(nil),         // 11: prometheus.ChunkedSeries
}
var file_remote_proto_depIdxs = []int32{
	7,  // 0: prometheus.WriteRequest.timeseries:type_name -> prometheus.TimeSeries
	8,  // 1: prometheus.WriteRequest.metadata:type_name -> prometheus.MetricMetadata
	4,  // 2: prometheus.ReadRequest.queries:type_name -> prometheus.Query
	0,  // 3: prometheus.ReadRequest.accepted_response_types:type_name -> prometheus.ReadRequest.ResponseType
	5,  // 4: prometheus.ReadResponse.results:type_name -> prometheus.QueryResult
	9,  // 5: prometheus.Query.matchers:type_name -> prometheus.LabelMatcher
	10, // 6: prometheus.Query.hints:type_name -> prometheus.ReadHints
	7,  // 7: prometheus.QueryResult.timeseries:type_name -> prometheus.TimeSeries
	11, // 8: prometheus.ChunkedReadResponse.chunked_series:type_name -> prometheus.ChunkedSeries
	9,  // [9:9] is the sub-list for method output_type
	9,  // [9:9] is the sub-list for method input_type
	9,  // [9:9] is the sub-list for extension type_name
	9,  // [9:9] is the sub-list for extension extendee
	0,  // [0:9] is the sub-list for field type_name
}

func init() { file_remote_proto_init() }
func file_remote_proto_init() {
	if File_remote_proto != nil {
		return
	}
	file_types_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_remote_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*WriteRequest); i {
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
		file_remote_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*ReadRequest); i {
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
		file_remote_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*ReadResponse); i {
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
		file_remote_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*Query); i {
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
		file_remote_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*QueryResult); i {
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
		file_remote_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*ChunkedReadResponse); i {
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
			RawDescriptor: file_remote_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_remote_proto_goTypes,
		DependencyIndexes: file_remote_proto_depIdxs,
		EnumInfos:         file_remote_proto_enumTypes,
		MessageInfos:      file_remote_proto_msgTypes,
	}.Build()
	File_remote_proto = out.File
	file_remote_proto_rawDesc = nil
	file_remote_proto_goTypes = nil
	file_remote_proto_depIdxs = nil
}
