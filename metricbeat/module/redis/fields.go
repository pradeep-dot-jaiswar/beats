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

// Code generated by beats/dev-tools/cmd/asset/asset.go - DO NOT EDIT.

package redis

import (
	"github.com/snappyflow/beats/v7/libbeat/asset"
)

func init() {
	if err := asset.SetFields("metricbeat", "redis", asset.ModuleFieldsPri, AssetRedis); err != nil {
		panic(err)
	}
}

// AssetRedis returns asset data.
// This is the base64 encoded gzipped contents of module/redis.
func AssetRedis() string {
	return "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"
}
