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

// +build integration

package info

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/snappyflow/beats/v7/libbeat/tests/compose"
	mbtest "github.com/snappyflow/beats/v7/metricbeat/mb/testing"
)

func TestFetch(t *testing.T) {
	service := compose.EnsureUp(t, "haproxy")

	f := mbtest.NewReportingMetricSetV2Error(t, getConfig(service.HostForPort(14567)))
	events, errs := mbtest.ReportingFetchV2Error(f)

	assert.Empty(t, errs)
	if !assert.NotEmpty(t, events) {
		t.FailNow()
	}

	t.Logf("%s/%s event: %+v", f.Module().Name(), f.Name(),
		events[0].BeatEvent("haproxy", "info").Fields.StringToPrint())

}

func TestData(t *testing.T) {
	service := compose.EnsureUp(t, "haproxy")

	config := getConfig(service.HostForPort(14567))
	f := mbtest.NewReportingMetricSetV2Error(t, config)
	err := mbtest.WriteEventsReporterV2Error(f, t, ".")
	if err != nil {
		t.Fatal("write", err)
	}

}

func getConfig(host string) map[string]interface{} {
	return map[string]interface{}{
		"module":     "haproxy",
		"metricsets": []string{"info"},
		"hosts":      []string{"tcp://" + host},
	}
}
