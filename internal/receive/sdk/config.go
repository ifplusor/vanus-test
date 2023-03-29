// Copyright 2023 Linkall Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sdk

import (
	// first-party libraries.
	vanus "github.com/vanus-labs/sdk/golang"

	// this project.
	"github.com/vanus-labs/vanus-test/internal/tracing"
)

type receiverConfig struct {
	endpoint string
	client   vanus.Client

	active bool
	port   int

	parallelism int
	batch       int

	tracer tracing.ReceiveTracer
}

func (cfg *receiverConfig) options(subscription uint64) []vanus.SubscriptionOption {
	opts := []vanus.SubscriptionOption{
		vanus.WithSubscriptionID(vanus.NewID(subscription)),
	}

	if cfg.active {
		opts = append(opts, vanus.WithActiveMode(true))
	} else {
		opts = append(opts, vanus.WithListenPort(cfg.port))
	}

	if cfg.parallelism > 0 {
		opts = append(opts, vanus.WithParallelism(cfg.parallelism))
	}
	if cfg.batch > 0 {
		opts = append(opts, vanus.WithMaxBatchSize(cfg.batch))
	}

	return opts
}

type ReceiverOption func(*receiverConfig)

func WithEndpoint(endpoint string) ReceiverOption {
	return func(cfg *receiverConfig) {
		cfg.endpoint = endpoint
	}
}

func WithClient(client vanus.Client) ReceiverOption {
	return func(cfg *receiverConfig) {
		cfg.client = client
	}
}

func WithActive() ReceiverOption {
	return func(cfg *receiverConfig) {
		cfg.active = true
	}
}

func WithPassive(port int) ReceiverOption {
	return func(cfg *receiverConfig) {
		cfg.active = false
		cfg.port = port
	}
}

func WithParallelism(parallelism int) ReceiverOption {
	return func(cfg *receiverConfig) {
		cfg.parallelism = parallelism
	}
}

func WithBatchSize(batch int) ReceiverOption {
	return func(cfg *receiverConfig) {
		cfg.batch = batch
	}
}

func WithTracer(tracer tracing.ReceiveTracer) ReceiverOption {
	return func(cfg *receiverConfig) {
		cfg.tracer = tracer
	}
}
