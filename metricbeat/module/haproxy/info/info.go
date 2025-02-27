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

package info

import (
	"github.com/pkg/errors"

	"github.com/snappyflow/beats/v7/libbeat/logp"
	"github.com/snappyflow/beats/v7/metricbeat/helper"
	"github.com/snappyflow/beats/v7/metricbeat/mb"
	"github.com/snappyflow/beats/v7/metricbeat/module/haproxy"
)

const (
	statsMethod = "info"
)

var (
	debugf = logp.MakeDebug("haproxy-info")
)

// init registers the haproxy info MetricSet.
func init() {
	mb.Registry.MustAddMetricSet("haproxy", "info", New,
		mb.WithHostParser(haproxy.HostParser),
		mb.DefaultMetricSet(),
	)
}

// MetricSet for haproxy info.
type MetricSet struct {
	mb.BaseMetricSet
	*helper.HTTP
}

// New creates a haproxy info MetricSet.
func New(base mb.BaseMetricSet) (mb.MetricSet, error) {
	return &MetricSet{BaseMetricSet: base}, nil
}

// Fetch fetches info stats from the haproxy service.
func (m *MetricSet) Fetch(reporter mb.ReporterV2) error {
	hapc, err := haproxy.NewHaproxyClient(m.HostData().URI, m.BaseMetricSet)
	if err != nil {
		return errors.Wrap(err, "failed creating haproxy client")
	}

	res, err := hapc.GetInfo()
	if err != nil {
		return errors.Wrap(err, "failed fetching haproxy info")
	}

	event, err := eventMapping(res, reporter)
	if err != nil {
		return errors.Wrap(err, "error in mapping")
	}
	reporter.Event(event)
	return nil
}
