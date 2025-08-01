// proto/job/v1/job_service.proto

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.29.3
// source: service.proto

package v1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
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

type JobState int32

const (
	JobState_JOB_STATE_UNSET      JobState = 0
	JobState_JOB_STATE_CREATED    JobState = 1
	JobState_JOB_STATE_SCHEDULED  JobState = 2
	JobState_JOB_STATE_PROCESSING JobState = 4
	JobState_JOB_STATE_CANCELLED  JobState = 8
	JobState_JOB_STATE_TIMED_OUT  JobState = 16
	JobState_JOB_STATE_COMPLETED  JobState = 32
	JobState_JOB_STATE_TERMINATED JobState = 64
)

// Enum value maps for JobState.
var (
	JobState_name = map[int32]string{
		0:  "JOB_STATE_UNSET",
		1:  "JOB_STATE_CREATED",
		2:  "JOB_STATE_SCHEDULED",
		4:  "JOB_STATE_PROCESSING",
		8:  "JOB_STATE_CANCELLED",
		16: "JOB_STATE_TIMED_OUT",
		32: "JOB_STATE_COMPLETED",
		64: "JOB_STATE_TERMINATED",
	}
	JobState_value = map[string]int32{
		"JOB_STATE_UNSET":      0,
		"JOB_STATE_CREATED":    1,
		"JOB_STATE_SCHEDULED":  2,
		"JOB_STATE_PROCESSING": 4,
		"JOB_STATE_CANCELLED":  8,
		"JOB_STATE_TIMED_OUT":  16,
		"JOB_STATE_COMPLETED":  32,
		"JOB_STATE_TERMINATED": 64,
	}
)

func (x JobState) Enum() *JobState {
	p := new(JobState)
	*p = x
	return p
}

func (x JobState) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (JobState) Descriptor() protoreflect.EnumDescriptor {
	return file_service_proto_enumTypes[0].Descriptor()
}

func (JobState) Type() protoreflect.EnumType {
	return &file_service_proto_enumTypes[0]
}

func (x JobState) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use JobState.Descriptor instead.
func (JobState) EnumDescriptor() ([]byte, []int) {
	return file_service_proto_rawDescGZIP(), []int{0}
}

type VersionRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *VersionRequest) Reset() {
	*x = VersionRequest{}
	mi := &file_service_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *VersionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VersionRequest) ProtoMessage() {}

func (x *VersionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_service_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VersionRequest.ProtoReflect.Descriptor instead.
func (*VersionRequest) Descriptor() ([]byte, []int) {
	return file_service_proto_rawDescGZIP(), []int{0}
}

type VersionResponse struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// ////////////////////////////
	// Basic information: 1-10
	// ////////////////////////////
	Version    string `protobuf:"bytes,1,opt,name=version,proto3" json:"version,omitempty"`
	ApiVersion string `protobuf:"bytes,2,opt,name=api_version,json=apiVersion,proto3" json:"api_version,omitempty"`
	// ////////////////////////////
	// Execution information: 11-20
	// ////////////////////////////
	Revision      *string `protobuf:"bytes,11,opt,name=revision,proto3,oneof" json:"revision,omitempty"`
	BuildDate     *string `protobuf:"bytes,12,opt,name=build_date,json=buildDate,proto3,oneof" json:"build_date,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *VersionResponse) Reset() {
	*x = VersionResponse{}
	mi := &file_service_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *VersionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VersionResponse) ProtoMessage() {}

func (x *VersionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_service_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VersionResponse.ProtoReflect.Descriptor instead.
func (*VersionResponse) Descriptor() ([]byte, []int) {
	return file_service_proto_rawDescGZIP(), []int{1}
}

func (x *VersionResponse) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

func (x *VersionResponse) GetApiVersion() string {
	if x != nil {
		return x.ApiVersion
	}
	return ""
}

func (x *VersionResponse) GetRevision() string {
	if x != nil && x.Revision != nil {
		return *x.Revision
	}
	return ""
}

func (x *VersionResponse) GetBuildDate() string {
	if x != nil && x.BuildDate != nil {
		return *x.BuildDate
	}
	return ""
}

type Job struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Kind of the job (e.g., "email", "data_processing")
	Kind string `protobuf:"bytes,1,opt,name=kind,proto3" json:"kind,omitempty"`
	// Description of the job
	Description string `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	// Registered at timestamp
	RegisteredAt *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=registered_at,json=registeredAt,proto3" json:"registered_at,omitempty"`
	// Schedule using cron expression
	CronSpec *string `protobuf:"bytes,11,opt,name=cron_spec,json=cronSpec,proto3,oneof" json:"cron_spec,omitempty"`
	// Schedule at a specific time
	ScheduleAt    *timestamppb.Timestamp `protobuf:"bytes,12,opt,name=schedule_at,json=scheduleAt,proto3,oneof" json:"schedule_at,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Job) Reset() {
	*x = Job{}
	mi := &file_service_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Job) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Job) ProtoMessage() {}

func (x *Job) ProtoReflect() protoreflect.Message {
	mi := &file_service_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Job.ProtoReflect.Descriptor instead.
func (*Job) Descriptor() ([]byte, []int) {
	return file_service_proto_rawDescGZIP(), []int{2}
}

func (x *Job) GetKind() string {
	if x != nil {
		return x.Kind
	}
	return ""
}

func (x *Job) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Job) GetRegisteredAt() *timestamppb.Timestamp {
	if x != nil {
		return x.RegisteredAt
	}
	return nil
}

func (x *Job) GetCronSpec() string {
	if x != nil && x.CronSpec != nil {
		return *x.CronSpec
	}
	return ""
}

func (x *Job) GetScheduleAt() *timestamppb.Timestamp {
	if x != nil {
		return x.ScheduleAt
	}
	return nil
}

type JobInstance struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Kind
	Kind string `protobuf:"bytes,1,opt,name=kind,proto3" json:"kind,omitempty"`
	// Unique instance identifier
	Uuid string `protobuf:"bytes,2,opt,name=uuid,proto3" json:"uuid,omitempty"`
	// Current state
	State JobState `protobuf:"varint,11,opt,name=state,proto3,enum=job.v1.JobState" json:"state,omitempty"`
	// Job arguments
	Arguments []string `protobuf:"bytes,12,rep,name=arguments,proto3" json:"arguments,omitempty"`
	// Execution results (if completed)
	Results []string `protobuf:"bytes,13,rep,name=results,proto3" json:"results,omitempty"`
	// Error information (if failed)
	Error         *string                `protobuf:"bytes,14,opt,name=error,proto3,oneof" json:"error,omitempty"`
	CreatedAt     *timestamppb.Timestamp `protobuf:"bytes,21,opt,name=created_at,json=createdAt,proto3,oneof" json:"created_at,omitempty"`
	ScheduledAt   *timestamppb.Timestamp `protobuf:"bytes,22,opt,name=scheduled_at,json=scheduledAt,proto3,oneof" json:"scheduled_at,omitempty"`
	ProcessedAt   *timestamppb.Timestamp `protobuf:"bytes,23,opt,name=processed_at,json=processedAt,proto3,oneof" json:"processed_at,omitempty"`
	CompletedAt   *timestamppb.Timestamp `protobuf:"bytes,24,opt,name=completed_at,json=completedAt,proto3,oneof" json:"completed_at,omitempty"`
	TerminatedAt  *timestamppb.Timestamp `protobuf:"bytes,25,opt,name=terminated_at,json=terminatedAt,proto3,oneof" json:"terminated_at,omitempty"`
	CancelledAt   *timestamppb.Timestamp `protobuf:"bytes,26,opt,name=cancelled_at,json=cancelledAt,proto3,oneof" json:"cancelled_at,omitempty"`
	TimedOutAt    *timestamppb.Timestamp `protobuf:"bytes,27,opt,name=timed_out_at,json=timedOutAt,proto3,oneof" json:"timed_out_at,omitempty"`
	AttemptCount  *int32                 `protobuf:"varint,31,opt,name=attempt_count,json=attemptCount,proto3,oneof" json:"attempt_count,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *JobInstance) Reset() {
	*x = JobInstance{}
	mi := &file_service_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *JobInstance) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JobInstance) ProtoMessage() {}

func (x *JobInstance) ProtoReflect() protoreflect.Message {
	mi := &file_service_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use JobInstance.ProtoReflect.Descriptor instead.
func (*JobInstance) Descriptor() ([]byte, []int) {
	return file_service_proto_rawDescGZIP(), []int{3}
}

func (x *JobInstance) GetKind() string {
	if x != nil {
		return x.Kind
	}
	return ""
}

func (x *JobInstance) GetUuid() string {
	if x != nil {
		return x.Uuid
	}
	return ""
}

func (x *JobInstance) GetState() JobState {
	if x != nil {
		return x.State
	}
	return JobState_JOB_STATE_UNSET
}

func (x *JobInstance) GetArguments() []string {
	if x != nil {
		return x.Arguments
	}
	return nil
}

func (x *JobInstance) GetResults() []string {
	if x != nil {
		return x.Results
	}
	return nil
}

func (x *JobInstance) GetError() string {
	if x != nil && x.Error != nil {
		return *x.Error
	}
	return ""
}

func (x *JobInstance) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *JobInstance) GetScheduledAt() *timestamppb.Timestamp {
	if x != nil {
		return x.ScheduledAt
	}
	return nil
}

func (x *JobInstance) GetProcessedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.ProcessedAt
	}
	return nil
}

func (x *JobInstance) GetCompletedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CompletedAt
	}
	return nil
}

func (x *JobInstance) GetTerminatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.TerminatedAt
	}
	return nil
}

func (x *JobInstance) GetCancelledAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CancelledAt
	}
	return nil
}

func (x *JobInstance) GetTimedOutAt() *timestamppb.Timestamp {
	if x != nil {
		return x.TimedOutAt
	}
	return nil
}

func (x *JobInstance) GetAttemptCount() int32 {
	if x != nil && x.AttemptCount != nil {
		return *x.AttemptCount
	}
	return 0
}

type ScheduleJobRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Kind to schedule (must be pre-registered)
	Kind string `protobuf:"bytes,1,opt,name=kind,proto3" json:"kind,omitempty"`
	// Arguments to pass to the job executor
	Arguments []string `protobuf:"bytes,11,rep,name=arguments,proto3" json:"arguments,omitempty"`
	// Priority (lower values = higher priority; -1 means unset)
	Priority      *int32 `protobuf:"varint,12,opt,name=priority,proto3,oneof" json:"priority,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ScheduleJobRequest) Reset() {
	*x = ScheduleJobRequest{}
	mi := &file_service_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ScheduleJobRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ScheduleJobRequest) ProtoMessage() {}

func (x *ScheduleJobRequest) ProtoReflect() protoreflect.Message {
	mi := &file_service_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ScheduleJobRequest.ProtoReflect.Descriptor instead.
func (*ScheduleJobRequest) Descriptor() ([]byte, []int) {
	return file_service_proto_rawDescGZIP(), []int{4}
}

func (x *ScheduleJobRequest) GetKind() string {
	if x != nil {
		return x.Kind
	}
	return ""
}

func (x *ScheduleJobRequest) GetArguments() []string {
	if x != nil {
		return x.Arguments
	}
	return nil
}

func (x *ScheduleJobRequest) GetPriority() int32 {
	if x != nil && x.Priority != nil {
		return *x.Priority
	}
	return 0
}

type ScheduleJobResponse struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Scheduled job instance
	Instance      *JobInstance `protobuf:"bytes,1,opt,name=instance,proto3" json:"instance,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ScheduleJobResponse) Reset() {
	*x = ScheduleJobResponse{}
	mi := &file_service_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ScheduleJobResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ScheduleJobResponse) ProtoMessage() {}

func (x *ScheduleJobResponse) ProtoReflect() protoreflect.Message {
	mi := &file_service_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ScheduleJobResponse.ProtoReflect.Descriptor instead.
func (*ScheduleJobResponse) Descriptor() ([]byte, []int) {
	return file_service_proto_rawDescGZIP(), []int{5}
}

func (x *ScheduleJobResponse) GetInstance() *JobInstance {
	if x != nil {
		return x.Instance
	}
	return nil
}

type ListRegisteredJobsRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListRegisteredJobsRequest) Reset() {
	*x = ListRegisteredJobsRequest{}
	mi := &file_service_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListRegisteredJobsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListRegisteredJobsRequest) ProtoMessage() {}

func (x *ListRegisteredJobsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_service_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListRegisteredJobsRequest.ProtoReflect.Descriptor instead.
func (*ListRegisteredJobsRequest) Descriptor() ([]byte, []int) {
	return file_service_proto_rawDescGZIP(), []int{6}
}

type ListRegisteredJobsResponse struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// List of registered jobs
	Jobs          []*Job `protobuf:"bytes,1,rep,name=jobs,proto3" json:"jobs,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListRegisteredJobsResponse) Reset() {
	*x = ListRegisteredJobsResponse{}
	mi := &file_service_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListRegisteredJobsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListRegisteredJobsResponse) ProtoMessage() {}

func (x *ListRegisteredJobsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_service_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListRegisteredJobsResponse.ProtoReflect.Descriptor instead.
func (*ListRegisteredJobsResponse) Descriptor() ([]byte, []int) {
	return file_service_proto_rawDescGZIP(), []int{7}
}

func (x *ListRegisteredJobsResponse) GetJobs() []*Job {
	if x != nil {
		return x.Jobs
	}
	return nil
}

type Query struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Filter by job kind
	Kind *string `protobuf:"bytes,1,opt,name=kind,proto3,oneof" json:"kind,omitempty"`
	// Filter by job instance UUID
	Uuid *string `protobuf:"bytes,2,opt,name=uuid,proto3,oneof" json:"uuid,omitempty"`
	// Filter by job state
	State         *JobState `protobuf:"varint,3,opt,name=state,proto3,enum=job.v1.JobState,oneof" json:"state,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Query) Reset() {
	*x = Query{}
	mi := &file_service_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Query) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Query) ProtoMessage() {}

func (x *Query) ProtoReflect() protoreflect.Message {
	mi := &file_service_proto_msgTypes[8]
	if x != nil {
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
	return file_service_proto_rawDescGZIP(), []int{8}
}

func (x *Query) GetKind() string {
	if x != nil && x.Kind != nil {
		return *x.Kind
	}
	return ""
}

func (x *Query) GetUuid() string {
	if x != nil && x.Uuid != nil {
		return *x.Uuid
	}
	return ""
}

func (x *Query) GetState() JobState {
	if x != nil && x.State != nil {
		return *x.State
	}
	return JobState_JOB_STATE_UNSET
}

type LookupInstancesRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Lookup query
	Query         *Query `protobuf:"bytes,1,opt,name=query,proto3" json:"query,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *LookupInstancesRequest) Reset() {
	*x = LookupInstancesRequest{}
	mi := &file_service_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LookupInstancesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LookupInstancesRequest) ProtoMessage() {}

func (x *LookupInstancesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_service_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LookupInstancesRequest.ProtoReflect.Descriptor instead.
func (*LookupInstancesRequest) Descriptor() ([]byte, []int) {
	return file_service_proto_rawDescGZIP(), []int{9}
}

func (x *LookupInstancesRequest) GetQuery() *Query {
	if x != nil {
		return x.Query
	}
	return nil
}

type LookupInstancesResponse struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// List of job instances
	Instances     []*JobInstance `protobuf:"bytes,1,rep,name=instances,proto3" json:"instances,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *LookupInstancesResponse) Reset() {
	*x = LookupInstancesResponse{}
	mi := &file_service_proto_msgTypes[10]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LookupInstancesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LookupInstancesResponse) ProtoMessage() {}

func (x *LookupInstancesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_service_proto_msgTypes[10]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LookupInstancesResponse.ProtoReflect.Descriptor instead.
func (*LookupInstancesResponse) Descriptor() ([]byte, []int) {
	return file_service_proto_rawDescGZIP(), []int{10}
}

func (x *LookupInstancesResponse) GetInstances() []*JobInstance {
	if x != nil {
		return x.Instances
	}
	return nil
}

var File_service_proto protoreflect.FileDescriptor

const file_service_proto_rawDesc = "" +
	"\n" +
	"\rservice.proto\x12\x06job.v1\x1a\x1fgoogle/protobuf/timestamp.proto\"\x10\n" +
	"\x0eVersionRequest\"\xad\x01\n" +
	"\x0fVersionResponse\x12\x18\n" +
	"\aversion\x18\x01 \x01(\tR\aversion\x12\x1f\n" +
	"\vapi_version\x18\x02 \x01(\tR\n" +
	"apiVersion\x12\x1f\n" +
	"\brevision\x18\v \x01(\tH\x00R\brevision\x88\x01\x01\x12\"\n" +
	"\n" +
	"build_date\x18\f \x01(\tH\x01R\tbuildDate\x88\x01\x01B\v\n" +
	"\t_revisionB\r\n" +
	"\v_build_date\"\xfe\x01\n" +
	"\x03Job\x12\x12\n" +
	"\x04kind\x18\x01 \x01(\tR\x04kind\x12 \n" +
	"\vdescription\x18\x02 \x01(\tR\vdescription\x12?\n" +
	"\rregistered_at\x18\x03 \x01(\v2\x1a.google.protobuf.TimestampR\fregisteredAt\x12 \n" +
	"\tcron_spec\x18\v \x01(\tH\x00R\bcronSpec\x88\x01\x01\x12@\n" +
	"\vschedule_at\x18\f \x01(\v2\x1a.google.protobuf.TimestampH\x01R\n" +
	"scheduleAt\x88\x01\x01B\f\n" +
	"\n" +
	"_cron_specB\x0e\n" +
	"\f_schedule_at\"\xc5\x06\n" +
	"\vJobInstance\x12\x12\n" +
	"\x04kind\x18\x01 \x01(\tR\x04kind\x12\x12\n" +
	"\x04uuid\x18\x02 \x01(\tR\x04uuid\x12&\n" +
	"\x05state\x18\v \x01(\x0e2\x10.job.v1.JobStateR\x05state\x12\x1c\n" +
	"\targuments\x18\f \x03(\tR\targuments\x12\x18\n" +
	"\aresults\x18\r \x03(\tR\aresults\x12\x19\n" +
	"\x05error\x18\x0e \x01(\tH\x00R\x05error\x88\x01\x01\x12>\n" +
	"\n" +
	"created_at\x18\x15 \x01(\v2\x1a.google.protobuf.TimestampH\x01R\tcreatedAt\x88\x01\x01\x12B\n" +
	"\fscheduled_at\x18\x16 \x01(\v2\x1a.google.protobuf.TimestampH\x02R\vscheduledAt\x88\x01\x01\x12B\n" +
	"\fprocessed_at\x18\x17 \x01(\v2\x1a.google.protobuf.TimestampH\x03R\vprocessedAt\x88\x01\x01\x12B\n" +
	"\fcompleted_at\x18\x18 \x01(\v2\x1a.google.protobuf.TimestampH\x04R\vcompletedAt\x88\x01\x01\x12D\n" +
	"\rterminated_at\x18\x19 \x01(\v2\x1a.google.protobuf.TimestampH\x05R\fterminatedAt\x88\x01\x01\x12B\n" +
	"\fcancelled_at\x18\x1a \x01(\v2\x1a.google.protobuf.TimestampH\x06R\vcancelledAt\x88\x01\x01\x12A\n" +
	"\ftimed_out_at\x18\x1b \x01(\v2\x1a.google.protobuf.TimestampH\aR\n" +
	"timedOutAt\x88\x01\x01\x12(\n" +
	"\rattempt_count\x18\x1f \x01(\x05H\bR\fattemptCount\x88\x01\x01B\b\n" +
	"\x06_errorB\r\n" +
	"\v_created_atB\x0f\n" +
	"\r_scheduled_atB\x0f\n" +
	"\r_processed_atB\x0f\n" +
	"\r_completed_atB\x10\n" +
	"\x0e_terminated_atB\x0f\n" +
	"\r_cancelled_atB\x0f\n" +
	"\r_timed_out_atB\x10\n" +
	"\x0e_attempt_count\"t\n" +
	"\x12ScheduleJobRequest\x12\x12\n" +
	"\x04kind\x18\x01 \x01(\tR\x04kind\x12\x1c\n" +
	"\targuments\x18\v \x03(\tR\targuments\x12\x1f\n" +
	"\bpriority\x18\f \x01(\x05H\x00R\bpriority\x88\x01\x01B\v\n" +
	"\t_priority\"F\n" +
	"\x13ScheduleJobResponse\x12/\n" +
	"\binstance\x18\x01 \x01(\v2\x13.job.v1.JobInstanceR\binstance\"\x1b\n" +
	"\x19ListRegisteredJobsRequest\"=\n" +
	"\x1aListRegisteredJobsResponse\x12\x1f\n" +
	"\x04jobs\x18\x01 \x03(\v2\v.job.v1.JobR\x04jobs\"\x82\x01\n" +
	"\x05Query\x12\x17\n" +
	"\x04kind\x18\x01 \x01(\tH\x00R\x04kind\x88\x01\x01\x12\x17\n" +
	"\x04uuid\x18\x02 \x01(\tH\x01R\x04uuid\x88\x01\x01\x12+\n" +
	"\x05state\x18\x03 \x01(\x0e2\x10.job.v1.JobStateH\x02R\x05state\x88\x01\x01B\a\n" +
	"\x05_kindB\a\n" +
	"\x05_uuidB\b\n" +
	"\x06_state\"=\n" +
	"\x16LookupInstancesRequest\x12#\n" +
	"\x05query\x18\x01 \x01(\v2\r.job.v1.QueryR\x05query\"L\n" +
	"\x17LookupInstancesResponse\x121\n" +
	"\tinstances\x18\x01 \x03(\v2\x13.job.v1.JobInstanceR\tinstances*\xce\x01\n" +
	"\bJobState\x12\x13\n" +
	"\x0fJOB_STATE_UNSET\x10\x00\x12\x15\n" +
	"\x11JOB_STATE_CREATED\x10\x01\x12\x17\n" +
	"\x13JOB_STATE_SCHEDULED\x10\x02\x12\x18\n" +
	"\x14JOB_STATE_PROCESSING\x10\x04\x12\x17\n" +
	"\x13JOB_STATE_CANCELLED\x10\b\x12\x17\n" +
	"\x13JOB_STATE_TIMED_OUT\x10\x10\x12\x17\n" +
	"\x13JOB_STATE_COMPLETED\x10 \x12\x18\n" +
	"\x14JOB_STATE_TERMINATED\x10@2\xc4\x02\n" +
	"\n" +
	"JobService\x12=\n" +
	"\n" +
	"GetVersion\x12\x16.job.v1.VersionRequest\x1a\x17.job.v1.VersionResponse\x12F\n" +
	"\vScheduleJob\x12\x1a.job.v1.ScheduleJobRequest\x1a\x1b.job.v1.ScheduleJobResponse\x12[\n" +
	"\x12ListRegisteredJobs\x12!.job.v1.ListRegisteredJobsRequest\x1a\".job.v1.ListRegisteredJobsResponse\x12R\n" +
	"\x0fLookupInstances\x12\x1e.job.v1.LookupInstancesRequest\x1a\x1f.job.v1.LookupInstancesResponseB*Z(github.com/cybergarage/go-job/api/job/v1b\x06proto3"

var (
	file_service_proto_rawDescOnce sync.Once
	file_service_proto_rawDescData []byte
)

func file_service_proto_rawDescGZIP() []byte {
	file_service_proto_rawDescOnce.Do(func() {
		file_service_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_service_proto_rawDesc), len(file_service_proto_rawDesc)))
	})
	return file_service_proto_rawDescData
}

var file_service_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_service_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_service_proto_goTypes = []any{
	(JobState)(0),                      // 0: job.v1.JobState
	(*VersionRequest)(nil),             // 1: job.v1.VersionRequest
	(*VersionResponse)(nil),            // 2: job.v1.VersionResponse
	(*Job)(nil),                        // 3: job.v1.Job
	(*JobInstance)(nil),                // 4: job.v1.JobInstance
	(*ScheduleJobRequest)(nil),         // 5: job.v1.ScheduleJobRequest
	(*ScheduleJobResponse)(nil),        // 6: job.v1.ScheduleJobResponse
	(*ListRegisteredJobsRequest)(nil),  // 7: job.v1.ListRegisteredJobsRequest
	(*ListRegisteredJobsResponse)(nil), // 8: job.v1.ListRegisteredJobsResponse
	(*Query)(nil),                      // 9: job.v1.Query
	(*LookupInstancesRequest)(nil),     // 10: job.v1.LookupInstancesRequest
	(*LookupInstancesResponse)(nil),    // 11: job.v1.LookupInstancesResponse
	(*timestamppb.Timestamp)(nil),      // 12: google.protobuf.Timestamp
}
var file_service_proto_depIdxs = []int32{
	12, // 0: job.v1.Job.registered_at:type_name -> google.protobuf.Timestamp
	12, // 1: job.v1.Job.schedule_at:type_name -> google.protobuf.Timestamp
	0,  // 2: job.v1.JobInstance.state:type_name -> job.v1.JobState
	12, // 3: job.v1.JobInstance.created_at:type_name -> google.protobuf.Timestamp
	12, // 4: job.v1.JobInstance.scheduled_at:type_name -> google.protobuf.Timestamp
	12, // 5: job.v1.JobInstance.processed_at:type_name -> google.protobuf.Timestamp
	12, // 6: job.v1.JobInstance.completed_at:type_name -> google.protobuf.Timestamp
	12, // 7: job.v1.JobInstance.terminated_at:type_name -> google.protobuf.Timestamp
	12, // 8: job.v1.JobInstance.cancelled_at:type_name -> google.protobuf.Timestamp
	12, // 9: job.v1.JobInstance.timed_out_at:type_name -> google.protobuf.Timestamp
	4,  // 10: job.v1.ScheduleJobResponse.instance:type_name -> job.v1.JobInstance
	3,  // 11: job.v1.ListRegisteredJobsResponse.jobs:type_name -> job.v1.Job
	0,  // 12: job.v1.Query.state:type_name -> job.v1.JobState
	9,  // 13: job.v1.LookupInstancesRequest.query:type_name -> job.v1.Query
	4,  // 14: job.v1.LookupInstancesResponse.instances:type_name -> job.v1.JobInstance
	1,  // 15: job.v1.JobService.GetVersion:input_type -> job.v1.VersionRequest
	5,  // 16: job.v1.JobService.ScheduleJob:input_type -> job.v1.ScheduleJobRequest
	7,  // 17: job.v1.JobService.ListRegisteredJobs:input_type -> job.v1.ListRegisteredJobsRequest
	10, // 18: job.v1.JobService.LookupInstances:input_type -> job.v1.LookupInstancesRequest
	2,  // 19: job.v1.JobService.GetVersion:output_type -> job.v1.VersionResponse
	6,  // 20: job.v1.JobService.ScheduleJob:output_type -> job.v1.ScheduleJobResponse
	8,  // 21: job.v1.JobService.ListRegisteredJobs:output_type -> job.v1.ListRegisteredJobsResponse
	11, // 22: job.v1.JobService.LookupInstances:output_type -> job.v1.LookupInstancesResponse
	19, // [19:23] is the sub-list for method output_type
	15, // [15:19] is the sub-list for method input_type
	15, // [15:15] is the sub-list for extension type_name
	15, // [15:15] is the sub-list for extension extendee
	0,  // [0:15] is the sub-list for field type_name
}

func init() { file_service_proto_init() }
func file_service_proto_init() {
	if File_service_proto != nil {
		return
	}
	file_service_proto_msgTypes[1].OneofWrappers = []any{}
	file_service_proto_msgTypes[2].OneofWrappers = []any{}
	file_service_proto_msgTypes[3].OneofWrappers = []any{}
	file_service_proto_msgTypes[4].OneofWrappers = []any{}
	file_service_proto_msgTypes[8].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_service_proto_rawDesc), len(file_service_proto_rawDesc)),
			NumEnums:      1,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_service_proto_goTypes,
		DependencyIndexes: file_service_proto_depIdxs,
		EnumInfos:         file_service_proto_enumTypes,
		MessageInfos:      file_service_proto_msgTypes,
	}.Build()
	File_service_proto = out.File
	file_service_proto_goTypes = nil
	file_service_proto_depIdxs = nil
}
