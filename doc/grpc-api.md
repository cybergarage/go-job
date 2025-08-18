# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [service.proto](#service-proto)
    - [Job](#job-v1-Job)
    - [JobInstance](#job-v1-JobInstance)
    - [ListRegisteredJobsRequest](#job-v1-ListRegisteredJobsRequest)
    - [ListRegisteredJobsResponse](#job-v1-ListRegisteredJobsResponse)
    - [LookupInstancesRequest](#job-v1-LookupInstancesRequest)
    - [LookupInstancesResponse](#job-v1-LookupInstancesResponse)
    - [Query](#job-v1-Query)
    - [ScheduleJobRequest](#job-v1-ScheduleJobRequest)
    - [ScheduleJobResponse](#job-v1-ScheduleJobResponse)
    - [VersionRequest](#job-v1-VersionRequest)
    - [VersionResponse](#job-v1-VersionResponse)
  
    - [JobState](#job-v1-JobState)
  
    - [JobService](#job-v1-JobService)
  
- [Scalar Value Types](#scalar-value-types)



<a name="service-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## service.proto
proto/job/v1/job_service.proto


<a name="job-v1-Job"></a>

### Job



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| kind | [string](#string) |  | Kind of the job (e.g., &#34;email&#34;, &#34;data_processing&#34;) |
| description | [string](#string) |  | Description of the job |
| registered_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  | Registered at timestamp |
| cron_spec | [string](#string) | optional | Schedule using cron expression |
| schedule_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) | optional | Schedule at a specific time |






<a name="job-v1-JobInstance"></a>

### JobInstance



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| kind | [string](#string) |  | Kind of the job (e.g., &#34;email&#34;, &#34;data_processing&#34;) |
| uuid | [string](#string) |  | Unique instance identifier |
| state | [JobState](#job-v1-JobState) |  | Current state |
| arguments | [string](#string) | repeated | Job arguments |
| results | [string](#string) | repeated | Execution results (if completed) |
| error | [string](#string) | optional | Error information (if terminated) |
| created_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) | optional |  |
| scheduled_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) | optional |  |
| processed_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) | optional |  |
| completed_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) | optional |  |
| terminated_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) | optional |  |
| cancelled_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) | optional |  |
| timed_out_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) | optional |  |
| attempts | [int32](#int32) | optional | Total attempt count (initial execution &#43; retries) |






<a name="job-v1-ListRegisteredJobsRequest"></a>

### ListRegisteredJobsRequest







<a name="job-v1-ListRegisteredJobsResponse"></a>

### ListRegisteredJobsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| jobs | [Job](#job-v1-Job) | repeated | List of registered jobs |






<a name="job-v1-LookupInstancesRequest"></a>

### LookupInstancesRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| query | [Query](#job-v1-Query) |  | Lookup query |






<a name="job-v1-LookupInstancesResponse"></a>

### LookupInstancesResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| instances | [JobInstance](#job-v1-JobInstance) | repeated | List of job instances |






<a name="job-v1-Query"></a>

### Query



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| kind | [string](#string) | optional | Filter by job kind |
| uuid | [string](#string) | optional | Filter by job instance UUID |
| state | [JobState](#job-v1-JobState) | optional | Filter by job state |






<a name="job-v1-ScheduleJobRequest"></a>

### ScheduleJobRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| kind | [string](#string) |  | Kind to schedule (must be pre-registered) |
| arguments | [string](#string) | repeated | Arguments to pass to the job executor |
| priority | [int32](#int32) | optional | Priority (lower values = higher priority; -1 means unset) |






<a name="job-v1-ScheduleJobResponse"></a>

### ScheduleJobResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| instance | [JobInstance](#job-v1-JobInstance) |  | Scheduled job instance |






<a name="job-v1-VersionRequest"></a>

### VersionRequest







<a name="job-v1-VersionResponse"></a>

### VersionResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| version | [string](#string) |  |  |
| api_version | [string](#string) |  |  |





 


<a name="job-v1-JobState"></a>

### JobState


| Name | Number | Description |
| ---- | ------ | ----------- |
| JOB_STATE_UNSET | 0 |  |
| JOB_STATE_CREATED | 1 |  |
| JOB_STATE_SCHEDULED | 2 |  |
| JOB_STATE_PROCESSING | 4 |  |
| JOB_STATE_CANCELLED | 8 |  |
| JOB_STATE_TIMED_OUT | 16 |  |
| JOB_STATE_COMPLETED | 32 |  |
| JOB_STATE_TERMINATED | 64 |  |


 

 


<a name="job-v1-JobService"></a>

### JobService
JobService provides job scheduling and management capabilities.
This service allows you to register jobs, schedule them for execution,
and monitor their progress through various states.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GetVersion | [VersionRequest](#job-v1-VersionRequest) | [VersionResponse](#job-v1-VersionResponse) | GetVersion returns the service version information. |
| ScheduleJob | [ScheduleJobRequest](#job-v1-ScheduleJobRequest) | [ScheduleJobResponse](#job-v1-ScheduleJobResponse) | ScheduleJob schedules a job for execution with the specified parameters. The job must be pre-registered in the system before it can be scheduled. |
| ListRegisteredJobs | [ListRegisteredJobsRequest](#job-v1-ListRegisteredJobsRequest) | [ListRegisteredJobsResponse](#job-v1-ListRegisteredJobsResponse) | ListRegisteredJobs returns all currently registered jobs in the system. |
| LookupInstances | [LookupInstancesRequest](#job-v1-LookupInstancesRequest) | [LookupInstancesResponse](#job-v1-LookupInstancesResponse) | LookupInstances searches for job instances based on the provided query criteria. |

 



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

