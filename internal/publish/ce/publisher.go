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

package ce

import (
	// standard libraries.
	"context"
	"fmt"
	"strings"

	// third-party libraries.
	ce "github.com/cloudevents/sdk-go/v2"
	"github.com/vanus-labs/vanus-test/internal/publish"
)

const (
	httpPrefix  = "http://"
	httpsPrefix = "https://"
)

type publisher struct {
	target string
	client ce.Client
}

// Make sure publisher implements publish.Publisher and publish.BatchPublisher.
var _ publish.Publisher = (*publisher)(nil)

func New(endpoint string, eventbus uint64) (publish.Publisher, error) {
	client, err := ce.NewClientHTTP()
	if err != nil {
		return nil, err
	}

	var target string
	if strings.HasPrefix(endpoint, httpPrefix) || strings.HasPrefix(endpoint, httpsPrefix) {
		target = fmt.Sprintf("%s/%016X", endpoint, eventbus)
	} else {
		target = fmt.Sprintf("%s%s/%016X", httpPrefix, endpoint, eventbus)
	}

	cp := &publisher{
		target: target,
		client: client,
	}

	return cp, nil
}

func (cp *publisher) Publish(ctx context.Context, e ce.Event) error {
	ctx = ce.ContextWithTarget(ctx, cp.target)
	re := cp.client.Send(ctx, e)
	if !ce.IsACK(re) {
		return re
	}
	return nil
}
