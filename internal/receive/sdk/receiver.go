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

	// this project.
	"github.com/vanus-labs/vanus-test/internal/receive"
	"github.com/vanus-labs/vanus-test/internal/tracing"
)

type receiver struct {
	subscriber vanus.Subscriber
	tracer     tracing.ReceiveTracer
}

// Make sure receiver implements receive.Receiver and receive.BatchReceiver.
var (
	_ receive.Receiver      = (*receiver)(nil)
	_ receive.BatchReceiver = (*receiver)(nil)
)

func New(subscription uint64, opts ...ReceiverOption) (receive.Receiver, error) {
	var cfg receiverConfig
	for _, opt := range opts {
		opt(&cfg)
	}

	client := cfg.client
	if client == nil {
		c, err := vanus.Connect(&vanus.ClientOptions{Endpoint: cfg.endpoint})
		if err != nil {
			return nil, err
		}
		client = c
	}

	subscriber := client.Subscriber(cfg.options(subscription)...)

	r := &receiver{
		subscriber: subscriber,
		tracer:     cfg.tracer,
	}

	return r, nil
}

func (r *receiver) Receive(listener receive.Listener) error {
	if r.tracer == nil {
		return r.subscriber.Listen(func(ctx context.Context, msgs ...vanus.Message) error {
			for _, msg := range msgs {
				event := msg.GetEvent()
				if err := listener(ctx, *event); err != nil {
					msg.Failed(err)
				} else {
					msg.Success()
				}
			}
			return nil
		})
	}

	return r.subscriber.Listen(func(ctx context.Context, msgs ...vanus.Message) error {
		for _, msg := range msgs {
			event := msg.GetEvent()

			ctx = r.tracer.BeforeReceive(ctx, event)
			err := listener(ctx, *event)
			r.tracer.AfterReceive(ctx, event, err)

			if err != nil {
				msg.Failed(err)
			} else {
				msg.Success()
			}
		}
		return nil
	})
}

func (r *receiver) ReceiveBatch(listener receive.BatchListener) error {
	// FIXME(james.yin)
	return r.subscriber.Listen(func(ctx context.Context, msgs ...vanus.Message) error {
		events := make([]ce.Event, 0, len(msgs))
		for _, msg := range msgs {
			events = append(events, *msg.GetEvent())
		}
		return listener(ctx, events)
	})
}
