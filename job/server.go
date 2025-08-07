// Copyright (C) 2025 The go-job Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package job

import (
	"context"
	"errors"
	"net"
	"strconv"

	v1 "github.com/cybergarage/go-job/job/api/gen/go/v1"
	logger "github.com/cybergarage/go-logger/log"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Server is an interface that defines methods for managing the job server.
type Server interface {
	// Manager returns the job manager associated with the server.
	Manager() Manager
	// Start starts the job server.
	Start() error
	// Stop stops the job server.
	Stop() error
	// Restart restarts the job server.
	Restart() error
}

type server struct {
	v1.UnimplementedJobServiceServer

	grpcServer *grpc.Server
	manager    Manager
	addr       string
	port       int
}

// NewServer returns a new job server instance.
func NewServer(opts ...any) (Server, error) {
	mgr, err := NewManager(opts...)
	if err != nil {
		return nil, err
	}
	return &server{
		manager:                       mgr,
		addr:                          DefaultGrpcAddr,
		port:                          DefaultGrpcPort,
		grpcServer:                    nil,
		UnimplementedJobServiceServer: v1.UnimplementedJobServiceServer{},
	}, nil
}

// Manager returns the job manager associated with the server.
func (server *server) Manager() Manager {
	return server.manager
}

func (server *server) bindAddr() string {
	return net.JoinHostPort(server.addr, strconv.Itoa(server.port))
}

func (server *server) grpcStart() error {
	var err error
	listener, err := net.Listen("tcp", server.bindAddr())
	if err != nil {
		return err
	}

	loggingUnaryInterceptor := func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		resp, err := handler(ctx, req)
		if err == nil {
			logger.Infof("gRPC Request: %s", info.FullMethod)
		} else {
			logger.Errorf("gRPC Request: %s", info.FullMethod)
		}
		return resp, err
	}

	server.grpcServer = grpc.NewServer(grpc.UnaryInterceptor(loggingUnaryInterceptor))
	v1.RegisterJobServiceServer(server.grpcServer, server)
	go func() {
		if err := server.grpcServer.Serve(listener); err != nil {
			logger.Error(err)
		}
	}()

	return nil
}

// Stop stops the Grpc server.
func (server *server) grpcStop() error {
	if server.grpcServer != nil {
		server.grpcServer.GracefulStop()
		server.grpcServer = nil
	}

	return nil
}

// Start starts the job server.
func (server *server) Start() error {
	starters := []func() error{
		server.manager.Start,
		server.grpcStart,
	}
	var errs error
	for _, starter := range starters {
		if err := starter(); err != nil {
			errs = errors.Join(errs, err)
		}
	}

	if errs != nil {
		return errs
	}

	logger.Infof("%s/%s (%s) started", ProductName, Version, server.bindAddr())

	return nil
}

// Stop stops the job server.
func (server *server) Stop() error {
	stoppers := []func() error{
		server.manager.Stop,
		server.grpcStop,
	}
	var errs error
	for _, stopper := range stoppers {
		if err := stopper(); err != nil {
			errs = errors.Join(errs, err)
		}
	}

	if errs != nil {
		return errs
	}

	logger.Infof("%s/%s (%s) terminated", ProductName, Version, server.bindAddr())

	return nil
}

// Restart restarts the job server.
func (server *server) Restart() error {
	if err := server.Stop(); err != nil {
		return err
	}
	return server.Start()
}

// GetVersion returns the version of the job server.
func (server *server) GetVersion(ctx context.Context, req *v1.VersionRequest) (*v1.VersionResponse, error) {
	return &v1.VersionResponse{
		Version:    Version,
		ApiVersion: DefaultAPIVersion,
	}, nil
}

// ScheduleJob schedules a job with the specified kind, priority, and arguments.
func (server *server) ScheduleJob(ctx context.Context, req *v1.ScheduleJobRequest) (*v1.ScheduleJobResponse, error) {
	kind := req.GetKind()

	opts := []any{}
	priority := req.Priority
	if priority != nil {
		opts = append(opts, WithPriority(Priority(*priority)))
	}
	if req.GetArguments() != nil {
		args := []any{}
		for _, arg := range req.GetArguments() {
			args = append(args, arg)
		}
		if 0 < len(args) {
			opts = append(opts, WithArguments(args...))
		}
	}

	postJob, err := server.Manager().ScheduleRegisteredJob(kind, opts...)
	if err != nil {
		return nil, err
	}

	return &v1.ScheduleJobResponse{
		Instance: &v1.JobInstance{
			Kind:         kind,
			Uuid:         postJob.UUID().String(),
			State:        v1.JobState_JOB_STATE_SCHEDULED,
			Arguments:    nil,
			Results:      nil,
			Error:        nil,
			CreatedAt:    nil,
			ScheduledAt:  nil,
			ProcessedAt:  nil,
			CompletedAt:  nil,
			TerminatedAt: nil,
			CancelledAt:  nil,
			TimedOutAt:   nil,
			AttemptCount: nil,
		},
	}, nil
}

// ListRegisteredJobs returns a list of registered jobs.
func (server *server) ListRegisteredJobs(ctx context.Context, req *v1.ListRegisteredJobsRequest) (*v1.ListRegisteredJobsResponse, error) {
	allJobs, err := server.Manager().ListJobs()
	if err != nil {
		return nil, err
	}

	jobs := []*v1.Job{}
	for _, job := range allJobs {
		jobs = append(jobs, &v1.Job{
			Kind:         job.Kind(),
			Description:  job.Description(),
			RegisteredAt: timestamppb.New(job.RegisteredAt()),
			CronSpec:     nil,
			ScheduleAt:   nil,
		})
	}

	return &v1.ListRegisteredJobsResponse{
		Jobs: jobs,
	}, nil
}

// LookupInstances looks up all job instances which match the specified query.
func (server *server) LookupInstances(ctx context.Context, req *v1.LookupInstancesRequest) (*v1.LookupInstancesResponse, error) {
	queryOpts := []QueryOption{}
	queryKind := req.GetQuery().Kind
	if queryKind != nil && 0 < len(*queryKind) {
		queryOpts = append(queryOpts, WithQueryKind(*queryKind))
	}
	queryUUID := req.GetQuery().Uuid
	if queryUUID != nil && 0 < len(*queryUUID) {
		uuid, err := NewUUIDFrom(*queryUUID)
		if err != nil {
			return nil, err
		}
		queryOpts = append(queryOpts, WithQueryUUID(uuid))
	}
	queryState := req.GetQuery().State
	if queryState != nil {
		state, err := newStateFrom(*queryState)
		if err != nil {
			return nil, err
		}
		queryOpts = append(queryOpts, WithQueryState(state))
	}

	allInstances, err := server.Manager().LookupInstances(NewQuery(queryOpts...))
	if err != nil {
		return nil, err
	}

	instances := []*v1.JobInstance{}
	for _, instance := range allInstances {
		state, err := instance.State().ProtoState()
		if err != nil {
			return nil, err
		}
		instances = append(instances, &v1.JobInstance{
			Kind:         instance.Kind(),
			Uuid:         instance.UUID().String(),
			State:        state,
			Arguments:    nil,
			Results:      nil,
			Error:        nil,
			CreatedAt:    nil,
			ScheduledAt:  nil,
			ProcessedAt:  nil,
			CompletedAt:  nil,
			TerminatedAt: nil,
			CancelledAt:  nil,
			TimedOutAt:   nil,
			AttemptCount: nil})
	}

	return &v1.LookupInstancesResponse{
		Instances: instances,
	}, nil
}
