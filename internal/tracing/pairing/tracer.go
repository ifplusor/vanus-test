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

package pairing

import (
	// standard libraries.
	"context"
	"os"
	"sync"
	"time"

	// third-party libraries.
	ce "github.com/cloudevents/sdk-go/v2"

	// first-party libraries.
	"github.com/vanus-labs/vanus/observability/log"

	// this project.
	"github.com/vanus-labs/vanus-test/internal/tracing"
)

type Tracer struct {
	mutex   sync.Mutex
	cache   map[string]*ce.Event
	pending map[string]*ce.Event
}

// Make sure matchingTracer implements tracing.PublishTracer and tracing.ReceiveTracer.
var (
	_ tracing.PublishTracer = (*Tracer)(nil)
	_ tracing.ReceiveTracer = (*Tracer)(nil)
)

func New() *Tracer {
	return &Tracer{
		cache:   make(map[string]*ce.Event),
		pending: make(map[string]*ce.Event),
	}
}

func (t *Tracer) BeforePublish(ctx context.Context, event *ce.Event) context.Context {
	return ctx
}

func (t *Tracer) AfterPublish(ctx context.Context, event *ce.Event, err error) {
	if err != nil {
		return
	}

	t.mutex.Lock()
	t.cache[event.ID()] = event
	t.mutex.Unlock()
}

func (t *Tracer) BeforeReceive(ctx context.Context, event *ce.Event) context.Context {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	e, exist := t.cache[event.ID()]
	if !exist {
		t.pending[event.ID()] = event
		log.Warn(ctx).Stringer("event", event).Msg("Received a event, but which isn't found in cache.")
	} else if string(e.Data()) != string(event.Data()) {
		log.Error(ctx).Msg("Received a event, but data was corrupted.")
		os.Exit(1)
	}

	return ctx
}

func (t *Tracer) AfterReceive(ctx context.Context, event *ce.Event, err error) {
	if err != nil {
		return
	}

	t.mutex.Lock()
	delete(t.cache, event.ID())
	t.mutex.Unlock()
}

func (t *Tracer) Wait(duration time.Duration) bool {
	now := time.Now()
	for time.Since(now) < duration {
		t.mutex.Lock()
		if len(t.cache) == 0 {
			t.mutex.Unlock()
			return true
		}
		for k := range t.pending {
			delete(t.cache, k)
		}
		t.mutex.Unlock()
		time.Sleep(time.Second)
	}
	return false
}

func (t *Tracer) Report(ctx context.Context) {
	log.Info(ctx).Msg("lost events are below:")
	for _, v := range t.cache {
		log.Info(ctx).Stringer("event", v).Msg("lost event")
	}
}
