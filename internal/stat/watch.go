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

package stat

import (
	// standard libraries.
	"context"
	"sync/atomic"
	"time"

	// first-party libraries.
	"github.com/vanus-labs/vanus/observability/log"
)

func WatchIncrease(ctx context.Context, values map[string]*int64) {
	go func() {
		t := time.NewTicker(time.Second)
		defer t.Stop()

		prev := make(map[string]int64, len(values))
		for {
			select {
			case <-ctx.Done():
				log.Info(ctx).Interface("values", values).Msg("TPS printer is exit.")
				return
			case <-t.C:
				val := make(map[string]interface{}, len(values))
				for k, v := range values {
					cur := atomic.LoadInt64(v)
					val[k] = cur - prev[k]
					prev[k] = cur
				}
				log.Info(ctx).Interface("TPS", val).Send()
			}
		}
	}()
}

func WatchValues(ctx context.Context, values map[string]*int64) {
	go func() {
		t := time.NewTicker(time.Second)
		defer t.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Info(ctx).Msg("Total printer is exit.")
				return
			case <-t.C:
				cur := make(map[string]interface{}, len(values))
				for k, v := range values {
					cur[k] = atomic.LoadInt64(v)
				}
				log.Info(ctx).Interface("Total", values).Send()
			}
		}
	}()
}
