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

package spool

import (
	"sync"

	"github.com/snappyflow/beats/v7/libbeat/beat"
	"github.com/snappyflow/beats/v7/libbeat/publisher"
	"github.com/snappyflow/beats/v7/libbeat/publisher/queue"
)

// forgetfulProducer forwards event to the inBroker. The forgetfulProducer
// provides no event ACK handling and no callbacks.
type forgetfulProducer struct {
	openState openState
}

// ackProducer forwards events to the inBroker. The ackBroker provides
// functionality for ACK/Drop callbacks.
type ackProducer struct {
	dropOnCancel bool
	seq          uint32
	state        produceState
	openState    openState
	pubCancel    chan producerCancelRequest
}

// openState tracks the producer->inBroker connection state.
type openState struct {
	ctx    *spoolCtx
	done   chan struct{}
	events chan pushRequest
}

// produceState holds the ackProducer internal callback and event ACK state
// shared between ackProducer instances and inBroker instances.
// The state is used to compute the number of per producer ACKed events and
// executing locally configured callbacks.
type produceState struct {
	ackCB     ackHandler
	dropCB    func(beat.Event)
	cancelled bool
	lastACK   uint32
}

type ackHandler func(count int)

type clientStates struct {
	mux     sync.Mutex
	clients []clientState
}

type clientState struct {
	seq   uint32        // event sequence number
	state *produceState // the producer it's state used to compute and signal the ACK count
}

func newProducer(
	ctx *spoolCtx,
	pubCancel chan producerCancelRequest,
	events chan pushRequest,
	ackCB ackHandler,
	dropCB func(beat.Event),
	dropOnCancel bool,
) queue.Producer {
	openState := openState{
		ctx:    ctx,
		done:   make(chan struct{}),
		events: events,
	}

	if ackCB == nil {
		return &forgetfulProducer{openState: openState}
	}

	p := &ackProducer{
		seq:          1,
		dropOnCancel: dropOnCancel,
		openState:    openState,
		pubCancel:    pubCancel,
	}
	p.state.ackCB = ackCB
	p.state.dropCB = dropCB
	return p
}

func (p *forgetfulProducer) Publish(event publisher.Event) bool {
	return p.openState.publish(p.makeRequest(event))
}

func (p *forgetfulProducer) TryPublish(event publisher.Event) bool {
	return p.openState.tryPublish(p.makeRequest(event))
}

func (p *forgetfulProducer) makeRequest(event publisher.Event) pushRequest {
	return pushRequest{event: event}
}

func (p *forgetfulProducer) Cancel() int {
	p.openState.Close()
	return 0
}

func (p *ackProducer) Publish(event publisher.Event) bool {
	return p.updSeq(p.openState.publish(p.makeRequest(event)))
}

func (p *ackProducer) TryPublish(event publisher.Event) bool {
	return p.updSeq(p.openState.tryPublish(p.makeRequest(event)))
}

func (p *ackProducer) Cancel() int {
	p.openState.Close()

	if p.dropOnCancel {
		ch := make(chan producerCancelResponse)
		p.pubCancel <- producerCancelRequest{
			state: &p.state,
			resp:  ch,
		}

		// wait for cancel to being processed
		resp := <-ch
		return resp.removed
	}
	return 0
}

func (p *ackProducer) updSeq(ok bool) bool {
	if ok {
		p.seq++
	}
	return ok
}

func (p *ackProducer) makeRequest(event publisher.Event) pushRequest {
	return pushRequest{event: event, seq: p.seq, state: &p.state}
}

func (st *openState) Close() {
	close(st.done)
}

func (st *openState) publish(req pushRequest) bool {
	select {
	case st.events <- req:
		return true
	case <-st.done:
		st.events = nil
		return false
	}
}

func (st *openState) tryPublish(req pushRequest) bool {
	select {
	case st.events <- req:
		return true
	case <-st.done:
		st.events = nil
		return false
	default:
		log := st.ctx.logger
		log.Debugf("Dropping event, queue is blocked (seq=%v) ", req.seq)
		return false
	}
}

func (s *clientStates) Add(st clientState) int {
	s.mux.Lock()
	s.clients = append(s.clients, st)
	l := len(s.clients)
	s.mux.Unlock()
	return l
}

func (s *clientStates) RemoveLast() {
	s.mux.Lock()
	s.clients = s.clients[:len(s.clients)-1]
	s.mux.Unlock()
}

func (s *clientStates) Pop(n int) (states []clientState) {
	s.mux.Lock()
	states, s.clients = s.clients[:n], s.clients[n:]
	s.mux.Unlock()
	return states
}
