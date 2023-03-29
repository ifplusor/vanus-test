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

package event

import (
	// standard libraries.
	stdrand "math/rand"
	"time"

	// third-party libraries.
	ce "github.com/cloudevents/sdk-go/v2"
	"github.com/rs/xid"

	// this project.
	"github.com/vanus-labs/vanus-test/internal/constant"
	"github.com/vanus-labs/vanus-test/internal/rand"
)

type Generator interface {
	Generate() ce.Event
}

type simpleGenerator struct {
	ceSource    string
	ceType      string
	maxDataSize int
}

// Make sure that simpleGenerator implements Generator.
var _ Generator = (*simpleGenerator)(nil)

func NewSimpleGenerator(ceSource, ceType string, maxPayloadSize int) Generator {
	return &simpleGenerator{
		ceSource:    ceSource,
		ceType:      ceType,
		maxDataSize: maxPayloadSize,
	}
}

func (sg *simpleGenerator) Generate() ce.Event {
	e := ce.NewEvent()
	e.SetID(xid.New().String())
	e.SetSource(sg.ceSource)
	e.SetType(sg.ceType)
	_ = e.SetData(ce.TextPlain, rand.RandString(8+stdrand.Intn(sg.maxDataSize)))
	return e
}

type scheduledGenerator struct {
	simpleGenerator
	maxDelayInSecond int64
}

// Make sure that scheduledGenerator implements Generator.
var _ Generator = (*scheduledGenerator)(nil)

func NewScheduledGenerator(ceSource, ceType string, maxPayloadSize int, maxDelay time.Duration) Generator {
	return &scheduledGenerator{
		simpleGenerator: simpleGenerator{
			ceSource:    ceSource,
			ceType:      ceType,
			maxDataSize: maxPayloadSize,
		},
		maxDelayInSecond: int64(maxDelay.Seconds()),
	}
}

func (sg *scheduledGenerator) Generate() ce.Event {
	e := sg.simpleGenerator.Generate()
	deliveryTime := time.Now().Add(time.Duration(stdrand.Int63n(sg.maxDelayInSecond)) * time.Second)
	e.SetExtension(constant.XVanusDeliveryTime, deliveryTime.Format(time.RFC3339Nano))
	return e
}
