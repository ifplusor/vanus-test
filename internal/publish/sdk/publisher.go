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
	// standard libraries.
	"context"

	// third-party libraries.
	ce "github.com/cloudevents/sdk-go/v2"

	// first-party libraries.
	vanus "github.com/vanus-labs/sdk/golang"

	// this package.
	"github.com/vanus-labs/vanus-test/internal/publish"
	"github.com/vanus-labs/vanus-test/internal/tracing"
)

type publisher struct {
	publisher vanus.Publisher
	tracer    tracing.PublishTracer
}

// Make sure publisher implements publish.Publisher and publish.BatchPublisher.
var (
	_ publish.Publisher      = (*publisher)(nil)
	_ publish.BatchPublisher = (*publisher)(nil)
)

func New(eventbus uint64, opts ...publisherOption) (publish.Publisher, error) {
	var cfg publisherConfig
	for _, opt := range opts {
		opt(&cfg)
	}

	client := cfg.client
	if client == nil {
		c, err := vanus.Connect(&vanus.ClientOptions{Endpoint: cfg.endpoint, Token: cfg.token})
		if err != nil {
			return nil, err
		}
		client = c
	}

	p := client.Publisher(vanus.WithEventbusID(eventbus))

	pp := &publisher{
		publisher: p,
		tracer:    cfg.tracer,
	}

	return pp, nil
}

func (p *publisher) Publish(ctx context.Context, e ce.Event) error {
	if p.tracer != nil {
		return p.publishWithTrace(ctx, &e)
	}
	return p.publisher.Publish(ctx, &e)
}

func (p *publisher) publishWithTrace(ctx context.Context, e *ce.Event) error {
	ctx = p.tracer.BeforePublish(ctx, e)
	err := p.publisher.Publish(ctx, e)
	p.tracer.AfterPublish(ctx, e, err)
	return err
}

func (p *publisher) PublishBatch(ctx context.Context, events []ce.Event) error {
	batch := make([]*ce.Event, len(events))
	for i, e := range events {
		batch[i] = &e
	}
	return p.publisher.Publish(ctx, batch...)
}
