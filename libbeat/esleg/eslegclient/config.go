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

package eslegclient

import (
	"fmt"
	"time"

	"github.com/snappyflow/beats/v7/libbeat/common"
	"github.com/snappyflow/beats/v7/libbeat/common/transport/kerberos"
	"github.com/snappyflow/beats/v7/libbeat/common/transport/tlscommon"
)

type config struct {
	Hosts    []string          `config:"hosts" validate:"required"`
	Protocol string            `config:"protocol"`
	Path     string            `config:"path"`
	Params   map[string]string `config:"parameters"`
	Headers  map[string]string `config:"headers"`

	TLS      *tlscommon.Config `config:"ssl"`
	Kerberos *kerberos.Config  `config:"kerberos"`

	ProxyURL     string `config:"proxy_url"`
	ProxyDisable bool   `config:"proxy_disable"`

	Username string `config:"username"`
	Password string `config:"password"`
	APIKey   string `config:"api_key"`

	CompressionLevel int           `config:"compression_level" validate:"min=0, max=9"`
	EscapeHTML       bool          `config:"escape_html"`
	Timeout          time.Duration `config:"timeout"`
}

func defaultConfig() config {
	return config{
		Protocol:         "",
		Path:             "",
		ProxyURL:         "",
		ProxyDisable:     false,
		Params:           nil,
		Username:         "",
		Password:         "",
		APIKey:           "",
		Timeout:          90 * time.Second,
		CompressionLevel: 0,
		EscapeHTML:       false,
		TLS:              nil,
	}
}

func (c *config) Validate() error {
	if c.ProxyURL != "" && !c.ProxyDisable {
		if _, err := common.ParseURL(c.ProxyURL); err != nil {
			return err
		}
	}

	if c.APIKey != "" && (c.Username != "" || c.Password != "") {
		return fmt.Errorf("cannot set both api_key and username/password")
	}

	return nil
}
