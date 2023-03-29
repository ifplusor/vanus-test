// Copyright 2023 Linkall Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	// standard libraries.
	"context"
	"os"
	"time"

	// third-party libraries.
	"github.com/spf13/cobra"
	"gopkg.in/natefinch/lumberjack.v2"

	// first-party libraries.
	"github.com/vanus-labs/vanus/observability/log"

	// this project.
	"github.com/vanus-labs/vanus-test/integration/pubsub"
	"github.com/vanus-labs/vanus-test/integration/scheduled"
)

var gatewayEndpoint = os.Getenv("VANUS_GATEWAY")

func main() {
	var output string
	var num int64
	var timeout string

	cmd := &cobra.Command{
		Use:   "vanus-integration [case]",
		Short: "the integration test tool of vanus",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if output != "" {
				configLogger(output)
			}

			d, err := time.ParseDuration(timeout)
			if err != nil {
				return err
			}

			switch args[0] {
			case "pubsub":
				return pubsub.Test(context.Background(), gatewayEndpoint, num, d)
			case "scheduled":
				return scheduled.Test(context.Background(), gatewayEndpoint, num, d)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&output, "output", "", "file path of output, default is console")

	cmd.Flags().Int64Var(&num, "num", 100, "number of events to publish")
	cmd.Flags().StringVar(&timeout, "wait-timeout", "30s", "timeout of receive waiting")

	if err := cmd.Execute(); err != nil {
		os.Exit(-1)
	}
}

func configLogger(path string) {
	logger := &lumberjack.Logger{
		Filename:   path,
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28,   // days
		Compress:   true, // disabled by default
	}

	log.SetOutput(logger)
}
