# Design and Architecture

## Job State

The job state in `go-job` is managed through a combination of job instances and their associated states. The state of a job instance is crucial for understanding its lifecycle and for debugging purposes.

<figure>
<img src="img/job-state.png" alt="job state" />
</figure>

Each job instance can transition through various states, such as `Pending`, `Running`, `Completed`, and `Failed`. These states are tracked in the job manager, allowing you to monitor the progress and outcome of each job instance.
