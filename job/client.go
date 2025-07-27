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

// Client represens a gRPC client.
type Client struct {
	host string
	port int
	conn *grpc.ClientConn
}

// NewClient returns a new gRPC client.
func NewClient() *Client {
	client := &Client{
		host: "",
		port: DefaultGrpcPort,
		conn: nil,
	}
	return client
}

// SetPort sets a port number.
func (client *Client) SetPort(port int) {
	client.port = port
}

// SetHost sets a host name.
func (client *Client) SetHost(host string) {
	client.host = host
}

// Open opens a gRPC connection.
func (client *Client) Open() error {
	addr := net.JoinHostPort(client.host, strconv.Itoa(client.port))
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	client.conn = conn
	return nil
}

// Close closes the gRPC connection.
func (client *Client) Close() error {
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
func (client *Client) GetVersion() (string, error) {
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
func (client *Client) ScheduleJob(kind string, args ...any) (Instance, error) {
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
	state, err := NewStateFrom(ji.GetState())
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
func (client *Client) ListRegisteredJobs() ([]Job, error) {
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
