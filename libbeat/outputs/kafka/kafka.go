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

package kafka

import (
	"errors"
	"time"

	"github.com/Shopify/sarama"
	"github.com/eapache/go-resiliency/breaker"
	"github.com/snappyflow/beats/v7/libbeat/beat"
	"github.com/snappyflow/beats/v7/libbeat/common"
	"github.com/snappyflow/beats/v7/libbeat/logp"
	"github.com/snappyflow/beats/v7/libbeat/outputs"
	"github.com/snappyflow/beats/v7/libbeat/outputs/codec"
	"github.com/snappyflow/beats/v7/libbeat/outputs/outil"
)

const (
	defaultWaitRetry = 1 * time.Second

	// NOTE: maxWaitRetry has no effect on mode, as logstash client currently does
	// not return ErrTempBulkFailure
	defaultMaxWaitRetry = 60 * time.Second

	logSelector = "kafka"
)

var (
	errNoTopicSet = errors.New("No topic configured")
	errNoHosts    = errors.New("No hosts configured")
)

func isRetriable(err error) bool {
	// retriableErrors are transient errors that may go away if we try again.
	// Taken from https://kafka.apache.org/25/javadoc/org/apache/kafka/common/errors/RetriableException.html.
	var retriableErrors = []error{
		sarama.ErrConsumerCoordinatorNotAvailable,
		sarama.ErrFetchSessionIDNotFound,
		sarama.ErrInvalidFetchSessionEpoch,
		sarama.ErrNotController,
		sarama.ErrNotCoordinatorForConsumer,
		sarama.ErrNotEnoughReplicasAfterAppend,
		sarama.ErrNotEnoughReplicas,
		sarama.ErrOffsetNotAvailable,
		sarama.ErrUnknownLeaderEpoch,
		breaker.ErrBreakerOpen,
	}

	for _, v := range retriableErrors {
		if errors.Is(err, v) {
			return true
		}
	}
	return false
}

func init() {
	sarama.Logger = kafkaLogger{log: logp.NewLogger(logSelector)}

	outputs.RegisterType("kafka", makeKafka)
}

func makeKafka(
	_ outputs.IndexManager,
	beat beat.Info,
	observer outputs.Observer,
	cfg *common.Config,
) (outputs.Group, error) {
	log := logp.NewLogger(logSelector)
	log.Debug("initialize kafka output")

	config, err := readConfig(cfg)
	if err != nil {
		return outputs.Fail(err)
	}

	topic, err := buildTopicSelector(cfg)
	if err != nil {
		return outputs.Fail(err)
	}

	libCfg, err := newSaramaConfig(log, config)
	if err != nil {
		return outputs.Fail(err)
	}

	hosts, err := outputs.ReadHostList(cfg)
	if err != nil {
		return outputs.Fail(err)
	}

	codec, err := codec.CreateEncoder(beat, config.Codec)
	if err != nil {
		return outputs.Fail(err)
	}

	client, err := newKafkaClient(observer, hosts, beat.IndexPrefix, config.Key, topic, codec, libCfg)
	if err != nil {
		return outputs.Fail(err)
	}

	retry := 0
	if config.MaxRetries < 0 {
		retry = -1
	}
	return outputs.Success(config.BulkMaxSize, retry, client)
}

func buildTopicSelector(cfg *common.Config) (outil.Selector, error) {
	return outil.BuildSelectorFromConfig(cfg, outil.Settings{
		Key:              "topic",
		MultiKey:         "topics",
		EnableSingleOnly: true,
		FailEmpty:        true,
		Case:             outil.SelectorKeepCase,
	})
}
