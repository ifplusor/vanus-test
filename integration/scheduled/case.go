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

package scheduled

import (
	// standard libraries.
	"context"
	"errors"
	"fmt"
	"time"

	// first-party libraries.
	vanus "github.com/vanus-labs/sdk/golang"
	"github.com/vanus-labs/vanus/observability/log"
	ctrlpb "github.com/vanus-labs/vanus/proto/pkg/controller"
	metapb "github.com/vanus-labs/vanus/proto/pkg/meta"

	// this project.
	"github.com/vanus-labs/vanus-test/internal/constant"
	"github.com/vanus-labs/vanus-test/internal/event"
	"github.com/vanus-labs/vanus-test/internal/rand"
	sdkrecv "github.com/vanus-labs/vanus-test/internal/receive/sdk"
	"github.com/vanus-labs/vanus-test/internal/stat"
	"github.com/vanus-labs/vanus-test/internal/tracing/pairing"
)

const (
	caseName = "ai.vanus.integration.scheduled"
)

var (
	parallelism          = 4
	maximumPayloadSize   = 4 * 1024
	maximumDelayInSecond = int32(120)
)

func Test(ctx context.Context, endpoint string, num int64, timeout time.Duration) (err error) {
	client, err := vanus.Connect(&vanus.ClientOptions{Endpoint: endpoint, Token: "admin"})
	if err != nil {
		return err
	}

	// Create an eventbus.
	ebName := fmt.Sprintf("regression-scheduled-%s", rand.RandString(5))
	eb, err := client.Controller().Eventbus().Create(ctx, "default", ebName)
	if err != nil {
		return err
	}

	// Delete the eventbus when exit.
	defer func(ctx context.Context) {
		if err2 := client.Controller().Eventbus().Delete(ctx, vanus.WithEventbusID(eb.Id)); err2 != nil {
			if err != nil {
				err = errors.Join(err, err2)
			} else {
				err = err2
			}
		}
	}(ctx)

	// Create a subscription.
	req := &ctrlpb.SubscriptionRequest{
		NamespaceId: eb.NamespaceId,
		Name:        ebName,
		Description: "regression test: scheduled",
		EventbusId:  eb.Id,
		Sink:        "localhost:8080", // FIXME
		Protocol:    metapb.Protocol_HTTP,
		Config: &metapb.SubscriptionConfig{
			OffsetType: metapb.SubscriptionConfig_EARLIEST,
		},
		Disable: true,
	}
	sub, err := client.Controller().Subscription().Create(ctx, req)
	if err != nil {
		return err
	}

	// Delete the subscription when exit.
	defer func(ctx context.Context) {
		if err2 := client.Controller().Subscription().Delete(ctx,
			vanus.WithSubscriptionID(vanus.NewID(sub.Id))); err2 != nil {
			if err != nil {
				err = errors.Join(err, err2)
			} else {
				err = err2
			}
		}
	}(ctx)

	var publishStats event.PublishStats
	var receiveStats event.ReceiveStats

	tracer := pairing.New()

	ctx, cancel := context.WithCancel(ctx)

	go func() {
		opts := []sdkrecv.ReceiverOption{
			sdkrecv.WithTracer(tracer),
			sdkrecv.WithBatchSize(32),
			sdkrecv.WithParallelism(8),
		}
		if err2 := event.Receive(ctx, sub.Id, client, &receiveStats, opts...); err2 != nil {
			panic(err2)
		}
	}()

	stat.WatchValues(ctx, map[string]*int64{
		"Sending":   &publishStats.Success,
		"Receiving": &receiveStats.Num,
	})

	maxDelay := time.Duration(maximumDelayInSecond) * time.Second
	generator := event.NewScheduledGenerator(constant.TestRepo, caseName, maximumPayloadSize, maxDelay)
	if err = event.Publish(ctx, eb.Id, num, parallelism, generator, client, tracer, &publishStats); err != nil {
		cancel()
		return err
	}

	success := tracer.Wait(timeout + maxDelay)

	cancel()

	if !success {
		log.Error(ctx).Int64("success", publishStats.Success).Int64("received", receiveStats.Num).
			Msgf("failed to run %s case because of timeout after sending finished", caseName)

		tracer.Report(context.Background())

		return fmt.Errorf("run %s case failed", caseName)
	}

	log.Info(ctx).Int64("success", publishStats.Success).Int64("sent_failed", publishStats.Failed).
		Int64("received", receiveStats.Num).Msgf("success to run %s case", caseName)

	return
}
