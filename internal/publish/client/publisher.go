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

package client

import (
	// standard libraries.
	"context"

	// third-party libraries.
	ce "github.com/cloudevents/sdk-go/v2"

	// first-party libraries.
	"github.com/vanus-labs/vanus/client"
	"github.com/vanus-labs/vanus/client/pkg/api"

	// this project.
	"github.com/vanus-labs/vanus-test/internal/publish"
)

type publisher struct {
	writer api.BusWriter
}

// Make sure publisher implements publish.Publisher and publish.BatchPublisher.
var (
	_ publish.Publisher      = (*publisher)(nil)
	_ publish.BatchPublisher = (*publisher)(nil)
)

func New(addrs []string, eventbus uint64) (publish.Publisher, error) {
	cli := client.Connect(addrs)

	writer := cli.Eventbus(context.Background(), api.WithID(eventbus)).Writer()

	cp := &publisher{
		writer: writer,
	}

	return cp, nil
}

func (p *publisher) Publish(ctx context.Context, e ce.Event) error {
	_, err := api.AppendOne(ctx, p.writer, &e)
	return err
}

func (p *publisher) PublishBatch(ctx context.Context, events []ce.Event) error {
	batch := make([]*ce.Event, len(events))
	for i, e := range events {
		batch[i] = &e
	}
	_, err := api.Append(ctx, p.writer, batch)
	return err
}
