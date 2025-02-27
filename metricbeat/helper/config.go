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

package helper

import (
	"time"

	"github.com/snappyflow/beats/v7/libbeat/common/transport/tlscommon"
)

// Config for an HTTP helper
type Config struct {
	TLS             *tlscommon.Config `config:"ssl"`
	ConnectTimeout  time.Duration     `config:"connect_timeout"`
	Timeout         time.Duration     `config:"timeout"`
	Headers         map[string]string `config:"headers"`
	BearerTokenFile string            `config:"bearer_token_file"`
}

func defaultConfig() Config {
	return Config{
		ConnectTimeout: 2 * time.Second,
		Timeout:        10 * time.Second,
	}
}
