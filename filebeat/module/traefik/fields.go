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

package traefik

import (
	"github.com/snappyflow/beats/v7/libbeat/asset"
)

func init() {
	if err := asset.SetFields("filebeat", "traefik", asset.ModuleFieldsPri, AssetTraefik); err != nil {
		panic(err)
	}
}

// AssetTraefik returns asset data.
// This is the base64 encoded gzipped contents of module/traefik.
func AssetTraefik() string {
	return "eJyslb9upDAQxvt9ilH6IEV31RbXRIp0xTWn9MjBA2ut8XDjcSLe/mR2SQiyCZB1t3jn+33zx/Y9nLE/grDC2pwPAGLE4hHuni9f7g4AGn3FphND7gi/DgAAf0gHi1ATQ6fYG9eAnBCuQWCpgdpY9MUBoDZotT8OcffgVItTXlzSd3iEhil01y8JZFxPgxTUTG2eF9eUOeWqqkLv3z+n0Av4uB7JiTLOXxFDCaZWLoTo6N1MytDUVPDIpdHoxNQG+dN/Rodn7N+I9WxvwWdcv/1g7e/TIzz8fPgBF4b0QPWwUVmDTpKeGP8F9FJWFGb/GB1Zcs02O88nBBfaF+Ro4ErwSXzN5ASdLuPP2xVkcKBaHAswYmILdNLIi6rO0Udgu9NG0kRgO3q4EuDthIxjVcAMk/WmWE+MpR2S7kuPToqXXtAnXSpr1HynU3I6wkmkKxh9R85jEbWSMq1pWF2qKhwwMzItCZam22jBU+AKC6U1fz6ba8HD+ckOSh4c44pE3Bpmi3KiedtXFnvocJFUWJVuZhQXEmVbEJvGODUPXQOMtstXZG/I7ck4HbpupC6DWVakt3b382R7URJ8SmedjxqZM1fzyn5nNNbgVTO/ptcNdzkEbm19/ozlfcwfUMg8eVNJja+mmndjObVkehed1EH+YCU3N5MWEZSS2AggX9TB2tRlNgWl9/fQlkHpudlDmislZ6xByjwee8arIifGoZPvVOv6NjVIxRd6H9jghPvSeErdNrvAXyiOaEvVcJ6/j8wqfVyIjSF3o8ouib2X1Uh/q0ZmpWbZ3a6FOcHD/wAAAP//R4Rz0Q=="
}
