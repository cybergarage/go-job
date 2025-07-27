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

package server

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/cybergarage/go-job/job"
	"github.com/cybergarage/go-logger/log"
	"github.com/spf13/cobra"
)

var cfgFile string

var rootCmd = &cobra.Command{ // nolint:exhaustruct
	Use:               "jobd",
	Version:           job.Version,
	Short:             "",
	Long:              "",
	DisableAutoGenTag: true,
}

var versionCmd = &cobra.Command{ // nolint:exhaustruct
	Use:               "version",
	Short:             "Print " + strings.ToLower(job.ProductName) + " version",
	DisableAutoGenTag: true,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(job.Version)
		os.Exit(0)
	},
}

func GetRootCommand() *cobra.Command {
	return rootCmd
}

func Execute() {
	log.SetSharedLogger(log.NewStdoutLogger(log.LevelInfo))

	server, err := job.NewServer()
	if err != nil {
		log.Errorf("%s couldn't be created (%s)", job.ProductName, err.Error())
		os.Exit(1)
	}

	if err := server.Start(); err != nil {
		log.Errorf("%s couldn't be started (%s)", job.ProductName, err.Error())
		os.Exit(1)
	}

	sigCh := make(chan os.Signal, 1)

	signal.Notify(sigCh,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM)

	exitCh := make(chan int)

	go func() {
		for {
			s := <-sigCh
			switch s {
			case syscall.SIGHUP:
				log.Infof("caught %s, restarting...", s.String())
				if err := server.Restart(); err != nil {
					log.Errorf("%s couldn't be restarted (%s)", job.ProductName, err.Error())
					os.Exit(1)
				}
			case syscall.SIGINT, syscall.SIGTERM:
				log.Infof("caught %s, terminating...", s.String())
				if err := server.Stop(); err != nil {
					log.Errorf("%s couldn't be terminated (%s)", job.ProductName, err.Error())
					os.Exit(1)
				}
				exitCh <- 0
			}
		}
	}()

	code := <-exitCh

	os.Exit(code)
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./job.yaml)")
}
