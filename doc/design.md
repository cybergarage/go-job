# Design and Architecture

This document provides a detailed overview of \`go-job’s features and architecture.

## Terminology

This section defines the key terms and concepts used throughout the `go-job` system. Understanding these terms is essential for working effectively with the job scheduling and execution framework.

<table>
<colgroup>
<col style="width: 25%" />
<col style="width: 75%" />
</colgroup>
<thead>
<tr>
<th style="text-align: left;">Term</th>
<th style="text-align: left;">Definition</th>
</tr>
</thead>
<tbody>
<tr>
<td style="text-align: left;"><p>Job</p></td>
<td style="text-align: left;"><p>A reusable definition that specifies work to be performed, including the executor function, scheduling rules, and retry policies.</p></td>
</tr>
<tr>
<td style="text-align: left;"><p>Job Instance</p></td>
<td style="text-align: left;"><p>A specific execution of a job with concrete arguments, unique identifier, and state tracking throughout its lifecycle.</p></td>
</tr>
<tr>
<td style="text-align: left;"><p>Executor</p></td>
<td style="text-align: left;"><p>A Go function that implements the actual business logic for a specific job type.</p></td>
</tr>
<tr>
<td style="text-align: left;"><p>Processor</p></td>
<td style="text-align: left;"><p>A function that processes job instances, including completion and termination logic.</p></td>
</tr>
</tbody>
</table>

## Key Components

`go-job` is designed to handle job scheduling and execution efficiently. The architecture consists of several key components that work together to provide a robust job processing system.

The main components of `go-job` are:

<figure>
<img src="img/framework.png" alt="framework" />
</figure>

<table>
<colgroup>
<col style="width: 25%" />
<col style="width: 75%" />
</colgroup>
<thead>
<tr>
<th style="text-align: left;">Component</th>
<th style="text-align: left;">Description</th>
</tr>
</thead>
<tbody>
<tr>
<td style="text-align: left;"><p>Manager</p></td>
<td style="text-align: left;"><p>Coordinates job scheduling and execution across go-job components.</p></td>
</tr>
<tr>
<td style="text-align: left;"><p>Registry</p></td>
<td style="text-align: left;"><p>Holds job definitions and their associated executors.</p></td>
</tr>
<tr>
<td style="text-align: left;"><p>Worker</p></td>
<td style="text-align: left;"><p>Processes job instances by executing the registered functions.</p></td>
</tr>
<tr>
<td style="text-align: left;"><p>Queue</p></td>
<td style="text-align: left;"><p>Manages job instances, ensuring they are processed in the correct order.</p></td>
</tr>
<tr>
<td style="text-align: left;"><p>History</p></td>
<td style="text-align: left;"><p>Tracks state transitions of job instances, providing an execution history.</p></td>
</tr>
<tr>
<td style="text-align: left;"><p>Log</p></td>
<td style="text-align: left;"><p>Captures logs for each job instance, providing detailed execution information.</p></td>
</tr>
<tr>
<td style="text-align: left;"><p>Store</p></td>
<td style="text-align: left;"><p>Provides abstracted persistence for job metadata and execution state, enabling distributed operation and fault tolerance.</p></td>
</tr>
</tbody>
</table>

The queue, history, and log components can be shared between go-job servers using the store interface, which allows for a distributed architecture where multiple go-job servers can operate together, sharing job queue and instance information.

## Sequence Diagram

The `go-job` server is designed to be modular and extensible. Each component, including the registry, manager, and worker, can be independently developed and maintained.

The following sequence diagram illustrates the flow of job registration, scheduling, and processing.

<figure>
<img src="img/job-seqdgm.png" alt="job seqdgm" />
</figure>

The queue, history, and log components can be shared between go-job servers using distributed store plugins. This enables a distributed architecture where multiple go-job servers can operate together, sharing job instances and state information.

However, the registry that holds job definitions cannot be shared between go-job servers. Since Go has no built-in RPC mechanism to share job executors (which are function pointers), each go-job server must maintain its own local registry of job definitions.

## Job State Lifecycle

The job state in `go-job` is managed through a combination of job instances and their associated states. The state of a job instance is crucial for understanding its lifecycle and for debugging purposes.

<figure>
<img src="img/job-state.png" alt="job state" />
</figure>

<table>
<colgroup>
<col style="width: 25%" />
<col style="width: 75%" />
</colgroup>
<thead>
<tr>
<th style="text-align: left;">State</th>
<th style="text-align: left;">Description</th>
</tr>
</thead>
<tbody>
<tr>
<td style="text-align: left;"><p>Created</p></td>
<td style="text-align: left;"><p>The job instance has been created and is awaiting scheduling.</p></td>
</tr>
<tr>
<td style="text-align: left;"><p>Scheduled</p></td>
<td style="text-align: left;"><p>The job instance has been queued and is waiting to be processed by a worker.</p></td>
</tr>
<tr>
<td style="text-align: left;"><p>Processing</p></td>
<td style="text-align: left;"><p>The job instance is currently being executed by a worker.</p></td>
</tr>
<tr>
<td style="text-align: left;"><p>Terminated</p></td>
<td style="text-align: left;"><p>The job instance encountered an error or was forcibly stopped before completion.</p></td>
</tr>
<tr>
<td style="text-align: left;"><p>Completed</p></td>
<td style="text-align: left;"><p>The job instance finished successfully.</p></td>
</tr>
</tbody>
</table>

Each job instance can transition through various states, such as `Scheduled`, `Processing`, `Completed`, and `Terminated`. These states are tracked in the job manager, allowing you to monitor the progress and outcome of each job instance.
