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
	"fmt"
	"net"
	"strconv"

	v1 "github.com/cybergarage/go-job/job/api/gen/go/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// gRPC client implementation for client.
type grpcClient struct {
	host string
	port int
	conn *grpc.ClientConn
}

// NewClient returns a new gRPC client.
func NewGrpcClient() Client {
	client := &grpcClient{
		host: "",
		port: DefaultGrpcPort,
		conn: nil,
	}
	return client
}

// Name returns the name of the client.
func (client *grpcClient) Name() string {
	return "gRPC"
}

// SetPort sets a port number.
func (client *grpcClient) SetPort(port int) {
	client.port = port
}

// SetHost sets a host name.
func (client *grpcClient) SetHost(host string) {
	client.host = host
}

// Open opens a gRPC connection.
func (client *grpcClient) Open() error {
	addr := net.JoinHostPort(client.host, strconv.Itoa(client.port))
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	client.conn = conn
	return nil
}

// Close closes the gRPC connection.
func (client *grpcClient) Close() error {
	if client.conn == nil {
		return nil
	}
	err := client.conn.Close()
	if err != nil {
		return err
	}
	client.conn = nil
	return nil
}

// GetVersion retrieves the version of the job service.
func (client *grpcClient) GetVersion() (string, error) {
	c := v1.NewJobServiceClient(client.conn)
	req := &v1.VersionRequest{}
	res, err := c.GetVersion(context.Background(), req)
	if err != nil {
		return "", err
	}
	return res.GetVersion(), nil
}

// ScheduleJob schedules a job with the specified kind, priority, and arguments.
// The priority is lower for higher priority jobs, similar to Unix nice values.
func (client *grpcClient) ScheduleJob(kind string, args ...any) (Instance, error) {
	reqArgs := make([]string, len(args))
	for i, arg := range args {
		reqArgs[i] = fmt.Sprintf("%v", arg)
	}
	c := v1.NewJobServiceClient(client.conn)
	req := &v1.ScheduleJobRequest{
		Kind:      kind,
		Arguments: reqArgs,
		Priority:  nil,
	}
	res, err := c.ScheduleJob(context.Background(), req)
	if err != nil {
		return nil, err
	}
	ji := res.GetInstance()
	uuid, err := NewUUIDFrom(ji.GetUuid())
	if err != nil {
		return nil, err
	}
	state, err := newStateFrom(ji.GetState())
	if err != nil {
		return nil, err
	}
	return NewInstance(
		WithUUID(uuid),
		WithKind(kind),
		WithState(state),
	)
}

// ListRegisteredJobs lists all registered jobs.
func (client *grpcClient) ListRegisteredJobs() ([]Job, error) {
	c := v1.NewJobServiceClient(client.conn)
	req := &v1.ListRegisteredJobsRequest{}
	res, err := c.ListRegisteredJobs(context.Background(), req)
	if err != nil {
		return nil, err
	}
	pbJobs := make([]Job, len(res.GetJobs()))
	for i, pbJob := range res.GetJobs() {
		job, err := NewJob(
			WithKind(pbJob.GetKind()),
			WithDescription(pbJob.GetDescription()),
			WithRegisteredAt(pbJob.GetRegisteredAt().AsTime()),
		)
		if err != nil {
			return nil, err
		}
		pbJobs[i] = job
	}
	return pbJobs, nil
}

// LookupInstances looks up job instances based on the provided query.
func (client *grpcClient) LookupInstances(query Query) ([]Instance, error) {
	c := v1.NewJobServiceClient(client.conn)

	pbQuery := &v1.Query{
		Kind:  nil,
		Uuid:  nil,
		State: nil,
	}
	kind, ok := query.Kind()
	if ok {
		pbQuery.Kind = &kind
	}
	id, ok := query.UUID()
	if ok {
		idStr := id.String()
		pbQuery.Uuid = &idStr
	}
	state, ok := query.State()
	if ok {
		pbState, err := state.ProtoState()
		if err != nil {
			return nil, err
		}
		pbQuery.State = &pbState
	}

	req := &v1.LookupInstancesRequest{
		Query: pbQuery,
	}
	res, err := c.LookupInstances(context.Background(), req)
	if err != nil {
		return nil, err
	}

	pbInstances := make([]Instance, len(res.GetInstances()))
	for i, pbInstance := range res.GetInstances() {
		uuid, err := NewUUIDFrom(pbInstance.GetUuid())
		if err != nil {
			return nil, err
		}
		state, err := newStateFrom(pbInstance.GetState())
		if err != nil {
			return nil, err
		}
		pbInstances[i], err = NewInstance(
			WithUUID(uuid),
			WithKind(pbInstance.GetKind()),
			WithState(state),
		)
		if err != nil {
			return nil, err
		}
	}
	return pbInstances, nil
}
