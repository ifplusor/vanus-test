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

package event

import (
	// standard libraries.
	"context"
	"sync/atomic"
	"time"

	// third-party libraries.
	ce "github.com/cloudevents/sdk-go/v2"

	// first-party libraries.
	vanus "github.com/vanus-labs/sdk/golang"
	"github.com/vanus-labs/vanus/observability/log"

	// this project.
	sdkrecv "github.com/vanus-labs/vanus-test/internal/receive/sdk"
)

type ReceiveStats struct {
	Num  int64
	Cost int64
}

func Receive(
	ctx context.Context,
	subscription uint64,
	client vanus.Client,
	stats *ReceiveStats,
	opts ...sdkrecv.ReceiverOption,
) error {
	opts = append([]sdkrecv.ReceiverOption{
		sdkrecv.WithClient(client),
		sdkrecv.WithActive(),
	}, opts...)
	receiver, err := sdkrecv.New(subscription, opts...)
	if err != nil {
		return err
	}

	err = receiver.Receive(func(ctx context.Context, event ce.Event) error {
		atomic.AddInt64(&stats.Num, 1)
		atomic.AddInt64(&stats.Cost, time.Since(event.Time()).Microseconds())
		return nil
	})
	if err != nil {
		log.Error(ctx).Err(err).Msg("Failed to start events listening.")
	}
	return nil
}
