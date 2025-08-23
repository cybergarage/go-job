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
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	labelKind = "kind"
)

var (
	// Current number of registered jobs.
	mRegisteredJobs = prometheus.NewGauge(
		prometheus.GaugeOpts{ // nolint: exhaustruct
			Name: "go_job_registered",
			Help: "Current number of registered jobs",
		})

	// Current number of queued jobs by kind.
	mQueuedJobs = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{ // nolint: exhaustruct
			Name: "go_job_queued",
			Help: "Current number of queued jobs by kind",
		},
		[]string{labelKind},
	)

	// Total number of executed jobs by kind.
	mExecutedJobs = prometheus.NewCounterVec(
		prometheus.CounterOpts{ // nolint: exhaustruct
			Name: "go_job_executed_total",
			Help: "Total number of executed jobs by kind",
		},
		[]string{labelKind},
	)

	// Total number of successfully completed jobs by kind.
	mCompletedJobs = prometheus.NewCounterVec(
		prometheus.CounterOpts{ // nolint: exhaustruct
			Name: "go_job_completed_total",
			Help: "Total number of successfully completed jobs by kind",
		},
		[]string{labelKind},
	)

	// Total number of terminated jobs by kind.
	mTerminatedJobs = prometheus.NewCounterVec(
		prometheus.CounterOpts{ // nolint: exhaustruct
			Name: "go_job_terminated_total",
			Help: "Total number of terminated jobs by kind",
		},
		[]string{labelKind},
	)

	// Total number of canceled jobs by kind.
	mCanceledJobs = prometheus.NewCounterVec(
		prometheus.CounterOpts{ // nolint: exhaustruct
			Name: "go_job_canceled_total",
			Help: "Total number of canceled jobs by kind",
		},
		[]string{labelKind},
	)

	// Total number of timed out jobs by kind.
	mTimedOutJobs = prometheus.NewCounterVec(
		prometheus.CounterOpts{ // nolint: exhaustruct
			Name: "go_job_timedout_total",
			Help: "Total number of timed out jobs by kind",
		},
		[]string{labelKind},
	)

	// Histogram of job execution durations in seconds, labeled by job type.
	mJobDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{ // nolint: exhaustruct
			Name:    "go_job_duration_seconds",
			Help:    "Histogram of job execution durations in seconds by job type",
			Buckets: prometheus.DefBuckets,
		},
		[]string{labelKind},
	)

	// Current number of workers.
	mWorkers = prometheus.NewGauge(prometheus.GaugeOpts{ // nolint: exhaustruct
		Name: "go_job_workers",
		Help: "Current number of workers",
	})
)

func init() { // Register all metrics with Prometheus
	prometheus.MustRegister(
		mRegisteredJobs,
		mQueuedJobs,
		mExecutedJobs,
		mCompletedJobs,
		mTerminatedJobs,
		mCanceledJobs,
		mTimedOutJobs,
		mJobDuration,
		mWorkers,
	)
}

type metricsServer struct {
	// Embed the Prometheus metrics registry.
	registry   *prometheus.Registry
	httpServer *http.Server
	Addr       string
}

func newMetricsServer() *metricsServer {
	return &metricsServer{
		registry:   prometheus.NewRegistry(),
		httpServer: nil,
		Addr:       "",
	}
}

func (ms *metricsServer) Start(port int) error {
	err := ms.Stop()
	if err != nil {
		return err
	}

	addr := net.JoinHostPort(ms.Addr, strconv.Itoa(port))
	ms.httpServer = &http.Server{ // nolint:exhaustruct
		Addr:              addr,
		ReadHeaderTimeout: 5 * time.Second,
		Handler:           promhttp.Handler(),
	}

	c := make(chan error)
	go func() {
		c <- ms.httpServer.ListenAndServe()
	}()

	select {
	case err = <-c:
	case <-time.After(time.Millisecond * 500):
		err = nil
	}

	return err
}

func (ms *metricsServer) Stop() error {
	if ms.httpServer == nil {
		return nil
	}

	err := ms.httpServer.Close()
	if err != nil {
		return err
	}

	return nil
}
