// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package udp

import (
	"sync"
	"time"

	"github.com/snappyflow/beats/v7/filebeat/channel"
	"github.com/snappyflow/beats/v7/filebeat/harvester"
	"github.com/snappyflow/beats/v7/filebeat/input"
	"github.com/snappyflow/beats/v7/filebeat/inputsource"
	"github.com/snappyflow/beats/v7/filebeat/inputsource/udp"
	"github.com/snappyflow/beats/v7/libbeat/beat"
	"github.com/snappyflow/beats/v7/libbeat/common"
	"github.com/snappyflow/beats/v7/libbeat/logp"
)

func init() {
	err := input.Register("udp", NewInput)
	if err != nil {
		panic(err)
	}
}

// Input defines a udp input to receive event on a specific host:port.
type Input struct {
	sync.Mutex
	udp     *udp.Server
	started bool
	outlet  channel.Outleter
}

// NewInput creates a new udp input
func NewInput(
	cfg *common.Config,
	outlet channel.Connector,
	context input.Context,
) (input.Input, error) {

	out, err := outlet.Connect(cfg)
	if err != nil {
		return nil, err
	}

	config := defaultConfig
	if err = cfg.Unpack(&config); err != nil {
		return nil, err
	}

	forwarder := harvester.NewForwarder(out)
	callback := func(data []byte, metadata inputsource.NetworkMetadata) {
		forwarder.Send(beat.Event{
			Timestamp: time.Now(),
			Meta: common.MapStr{
				"truncated": metadata.Truncated,
			},
			Fields: common.MapStr{
				"message": string(data),
				"log": common.MapStr{
					"source": common.MapStr{
						"address": metadata.RemoteAddr.String(),
					},
				},
			},
		})
	}

	udp := udp.New(&config.Config, callback)

	return &Input{
		outlet:  out,
		udp:     udp,
		started: false,
	}, nil
}

// Run starts and start the UDP server and read events from the socket
func (p *Input) Run() {
	p.Lock()
	defer p.Unlock()

	if !p.started {
		logp.Info("Starting UDP input")
		err := p.udp.Start()
		if err != nil {
			logp.Err("Error running harvester: %v", err)
		}
		p.started = true
	}
}

// Stop stops the UDP input
func (p *Input) Stop() {
	defer p.outlet.Close()
	p.Lock()
	defer p.Unlock()

	logp.Info("Stopping UDP input")
	p.udp.Stop()
	p.started = false
}

// Wait suspends the UDP input
func (p *Input) Wait() {
	p.Stop()
}
