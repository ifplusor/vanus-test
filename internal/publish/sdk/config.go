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

type publisherConfig struct {
	endpoint string
	token    string
	client   vanus.Client

	tracer tracing.PublishTracer
}

type publisherOption func(*publisherConfig)

func WithEndpoint(endpoint string) publisherOption {
	return func(cfg *publisherConfig) {
		cfg.endpoint = endpoint
	}
}

func WithToken(token string) publisherOption {
	return func(cfg *publisherConfig) {
		cfg.token = token
	}
}

func WithClient(client vanus.Client) publisherOption {
	return func(cfg *publisherConfig) {
		cfg.client = client
	}
}

func WithTracer(tracer tracing.PublishTracer) publisherOption {
	return func(cfg *publisherConfig) {
		cfg.tracer = tracer
	}
}
