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

package cmd

import (
	// standard libraries.
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	// third-party libraries.
	ce "github.com/cloudevents/sdk-go/v2"
	"github.com/spf13/cobra"
	"github.com/thediveo/enumflag/v2"

	// this project.
	"github.com/vanus-labs/vanus-test/internal/publish"
	cepub "github.com/vanus-labs/vanus-test/internal/publish/ce"
	"github.com/vanus-labs/vanus-test/internal/publish/client"
	"github.com/vanus-labs/vanus-test/internal/publish/sdk"
)

type publishMode enumflag.Flag

const (
	ceMode publishMode = iota
	grpcMode
	sdkMode
	clientMode
)

var publishModeIds = map[publishMode][]string{
	ceMode:     {"ce"},
	sdkMode:    {"sdk"},
	clientMode: {"raw"},
}

func PublishCommand() *cobra.Command {
	var mode publishMode
	var addrs []string
	var token string
	var eventbus string
	var num int
	var size int
	var parallelism int

	cmd := &cobra.Command{
		Use:   "publish --addr address... --eventbus id",
		Short: "publish events.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseUint(eventbus, 16, 64)
			if err != nil {
				return err
			}

			var p publish.Publisher
			switch mode {
			case ceMode:
				p, err = cepub.New(addrs[0], id)
			case sdkMode:
				p, err = sdk.New(id, sdk.WithEndpoint(addrs[0]), sdk.WithToken(token))
			case clientMode:
				p, err = client.New(addrs, id)
			}
			if err != nil {
				return err
			}

			return runPublishTask(p, num, size, parallelism)
		},
	}

	cmd.Flags().Var(enumflag.New(&mode, "publishMode", publishModeIds, enumflag.EnumCaseInsensitive),
		"mode", "publish mode")
	cmd.Flags().StringArrayVar(&addrs, "addr", []string{}, "address of vanus server")
	cmd.Flags().StringVar(&token, "token", "admin", "token")
	cmd.Flags().StringVar(&eventbus, "eventbus", "", "id of eventbus to publish")
	cmd.Flags().IntVarP(&num, "num", "n", 10000, "number of events to publish")
	cmd.Flags().IntVar(&size, "size", 1024, "length of event payload")
	cmd.Flags().IntVar(&parallelism, "parallelism", 1, "parallelism when publish")

	cobra.MarkFlagRequired(cmd.Flags(), "addr")
	cobra.MarkFlagRequired(cmd.Flags(), "eventbus")

	return cmd
}

func runPublishTask(p publish.Publisher, num int, size int, parallelism int) error {
	data := func() string {
		str := ""
		for i := 0; i < size; i++ {
			str += "a"
		}
		return str
	}()

	var count, last int64
	var totalCost, lastCost int64
	go func() {
		t := time.NewTicker(time.Second)
		defer t.Stop()

		for range t.C {
			cu := atomic.LoadInt64(&count)
			ct := atomic.LoadInt64(&totalCost)
			if n := cu - last; n != 0 {
				log.Printf("TPS: %d\tlatency: %.3fms\n", n, float64(ct-lastCost)/float64(n)/1000)
			} else {
				log.Printf("TPS: %d\tlatency: NaN\n", n)
			}
			last = cu
			lastCost = ct
		}
	}()

	start := time.Now()

	var wg sync.WaitGroup
	wg.Add(parallelism)

	npt := (num + parallelism - 1) / parallelism
	for i := 0; i < parallelism; i++ {
		n := npt
		if i == parallelism-1 {
			n = num - (parallelism-1)*npt
		}

		go func(t, n int) {
			for j := 0; j < n; j++ {
				event := ce.NewEvent()
				event.SetType("ai.vanus.test")
				event.SetSource("https://vanus.ai")

				id := t*npt + j
				_ = event.SetData(ce.ApplicationJSON, map[string]interface{}{
					"id":   id,
					"data": data,
				})

				st := time.Now()
				err := p.Publish(context.Background(), event)
				cost := time.Since(st)
				if err != nil {
					log.Printf("Sent %d failed: %s", id, err.Error())
					time.Sleep(time.Second)
				} else {
					atomic.AddInt64(&count, 1)
					atomic.AddInt64(&totalCost, cost.Microseconds())
				}
			}
			wg.Done()
		}(i, n)
	}

	wg.Wait()

	cost := time.Since(start)
	fmt.Printf("cost: %d ms\n", cost.Milliseconds())
	fmt.Printf("failed: %d\n", int64(num)-count)
	fmt.Printf("tps: %f\n", float64(count)/cost.Seconds())
	return nil
}
