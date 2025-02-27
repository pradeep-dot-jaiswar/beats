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

package consumergroup

import (
	"io"
	"testing"
	"time"

	saramacluster "github.com/bsm/sarama-cluster"
	"github.com/pkg/errors"

	"github.com/snappyflow/beats/v7/libbeat/tests/compose"
	"github.com/snappyflow/beats/v7/metricbeat/mb"
	mbtest "github.com/snappyflow/beats/v7/metricbeat/mb/testing"
)

const (
	kafkaSASLConsumerUsername = "consumer"
	kafkaSASLConsumerPassword = "consumer-secret"
	kafkaSASLUsername         = "stats"
	kafkaSASLPassword         = "test-secret"
)

func TestData(t *testing.T) {
	service := compose.EnsureUp(t, "kafka",
		compose.UpWithTimeout(600*time.Second),
		compose.UpWithAdvertisedHostEnvFileForPort(9092),
	)

	c, err := startConsumer(t, service.HostForPort(9092), "metricbeat-test")
	if err != nil {
		t.Fatal(errors.Wrap(err, "starting kafka consumer"))
	}
	defer c.Close()

	ms := mbtest.NewReportingMetricSetV2Error(t, getConfig(service.HostForPort(9092)))
	for retries := 0; retries < 3; retries++ {
		err = mbtest.WriteEventsReporterV2Error(ms, t, "")
		if err == nil {
			return
		}
		time.Sleep(500 * time.Millisecond)
	}
	t.Fatal("write", err)
}

func TestFetch(t *testing.T) {
	service := compose.EnsureUp(t, "kafka",
		compose.UpWithTimeout(600*time.Second),
		compose.UpWithAdvertisedHostEnvFileForPort(9092),
	)

	c, err := startConsumer(t, service.HostForPort(9092), "metricbeat-test")
	if err != nil {
		t.Fatal(errors.Wrap(err, "starting kafka consumer"))
	}
	defer c.Close()

	f := mbtest.NewReportingMetricSetV2Error(t, getConfig(service.HostForPort(9092)))

	var data []mb.Event
	var errors []error
	for retries := 0; retries < 3; retries++ {
		data, errors = mbtest.ReportingFetchV2Error(f)
		if len(data) > 0 {
			continue
		}
		time.Sleep(500 * time.Millisecond)
	}
	if len(errors) > 0 {
		t.Fatalf("fetch %v", errors)
	}
	if len(data) == 0 {
		t.Fatalf("No consumer groups fetched")
	}
}

func startConsumer(t *testing.T, host string, topic string) (io.Closer, error) {
	brokers := []string{host}
	topics := []string{topic}
	config := saramacluster.NewConfig()
	config.Net.SASL.Enable = true
	config.Net.SASL.User = kafkaSASLConsumerUsername
	config.Net.SASL.Password = kafkaSASLConsumerPassword
	// The test panics unless CommitInterval is set due to the following bug in sarama:
	// https://github.com/Shopify/sarama/issues/1638
	// To work around the issue we need to set CommitInterval, but now sarama emits
	// a deprecation warning.
	config.Consumer.Offsets.CommitInterval = 1 * time.Second
	return saramacluster.NewConsumer(brokers, "test-group", topics, config)
}

func getConfig(host string) map[string]interface{} {
	return map[string]interface{}{
		"module":     "kafka",
		"metricsets": []string{"consumergroup"},
		"hosts":      []string{host},
		"username":   kafkaSASLUsername,
		"password":   kafkaSASLPassword,
	}
}
