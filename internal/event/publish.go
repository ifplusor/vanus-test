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
	"sync"
	"sync/atomic"
	"time"

	// third-party libraries.
	"github.com/rs/zerolog"

	// first-party libraries.
	vanus "github.com/vanus-labs/sdk/golang"
	"github.com/vanus-labs/vanus/observability/log"

	// this project.
	sdkpub "github.com/vanus-labs/vanus-test/internal/publish/sdk"
	"github.com/vanus-labs/vanus-test/internal/tracing"
)

type PublishStats struct {
	Success int64
	Failed  int64
	Cost    int64
}

func Publish(
	ctx context.Context,
	eventbus uint64,
	num int64,
	parallelism int,
	generator Generator,
	client vanus.Client,
	tracer tracing.PublishTracer,
	stats *PublishStats,
) error {
	publisher, err := sdkpub.New(eventbus, sdkpub.WithClient(client), sdkpub.WithTracer(tracer))
	if err != nil {
		return err
	}

	failedLogger := log.With().Logger().Sample(&zerolog.BurstSampler{
		Burst:  10,
		Period: 300 * time.Millisecond,
	})

	start := time.Now()

	wg := sync.WaitGroup{}
	wg.Add(parallelism)
	for p := 0; p < parallelism; p++ {
		go func() {
			defer wg.Done()

			for atomic.AddInt64(&num, -1) >= 0 {
				event := generator.Generate()
				err := publisher.Publish(context.Background(), event)
				if err != nil {
					atomic.AddInt64(&num, 1)
					failedLogger.Warn().Err(err).Msg("Failed to send events.")
					atomic.AddInt64(&stats.Failed, 1)
				} else {
					atomic.AddInt64(&stats.Success, 1)
				}
			}
		}()
	}
	wg.Wait()

	log.Info(ctx).Int64("success", stats.Success).Int64("failed", stats.Failed).Dur("cost", time.Now().Sub(start)).
		Msg("Finished to sent all events.")

	return nil
}
