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

package logstash

import (
	"github.com/snappyflow/beats/v7/libbeat/asset"
)

func init() {
	if err := asset.SetFields("metricbeat", "logstash", asset.ModuleFieldsPri, AssetLogstash); err != nil {
		panic(err)
	}
}

// AssetLogstash returns asset data.
// This is the base64 encoded gzipped contents of module/logstash.
func AssetLogstash() string {
	return "eJyslM2OmzAQgO88xYhzwwNw6KmtmqpV95TLarWyYALeGA/yDKzy9isCZMEx+d055OCJP3/jGbOCHe5TMFSwKC4jANFiMIX477AURwA5cuZ0LZpsCt8jAIAxDRXljcEIwKFBxZhCoSIARhFtC07hOWY28TeIS5E6fokAthpNzumBswKrKpwZdCH7uiM5auphJeAwJ01plnI8LoZoi8Q+vP3z0sbwD58KlMQyS4wSymjFXqZWUvZbku6nI3j/qHThVC8qrvGzZwrp4jexwAl0NG3RsSZ7oyyja3WGSXj3Q7rHydoE2KP1W1sFjf0eX3Hen80/WNsteYlQdy/f26fJDvfv5PJA/oJPF6HSp4fXOgRebhocG1c7ypA5CRPON+5K+af+CFj/CL7LhEUJP/o6Xw8UqFCczjh56LFii1b8K7t7nn4eaOBXuSQxFdHL42TIFve1Y20zqrQthjIho8YKumTRghr/0/UVGv8bKegWja02gg6XB/1+l18D+sTlIwAA//9sYbU7"
}
