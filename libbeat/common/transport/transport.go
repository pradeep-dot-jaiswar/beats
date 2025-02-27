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

package transport

import (
	"errors"
	"net"

	"github.com/snappyflow/beats/v7/libbeat/logp"
)

type Dialer interface {
	Dial(network, address string) (net.Conn, error)
}

type DialerFunc func(network, address string) (net.Conn, error)

var (
	ErrNotConnected = errors.New("client is not connected")
)

func (d DialerFunc) Dial(network, address string) (net.Conn, error) {
	return d(network, address)
}

func Dial(c Config, network, address string) (net.Conn, error) {
	d, err := MakeDialer(c)
	if err != nil {
		return nil, err
	}
	return d.Dial(network, address)
}

func MakeDialer(c Config) (Dialer, error) {
	var err error
	dialer := NetDialer(c.Timeout)
	dialer, err = ProxyDialer(logp.NewLogger(logSelector), c.Proxy, dialer)
	if err != nil {
		return nil, err
	}
	if c.Stats != nil {
		dialer = StatsDialer(dialer, c.Stats)
	}

	if c.TLS != nil {
		return TLSDialer(dialer, c.TLS, c.Timeout)
	}
	return dialer, nil
}
