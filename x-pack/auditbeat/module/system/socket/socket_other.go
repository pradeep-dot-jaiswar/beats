// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

// +build !linux

package socket

import (
	"fmt"

	"github.com/snappyflow/beats/v7/metricbeat/mb"
)

const (
	moduleName    = "system"
	metricsetName = "socket"
)

func init() {
	mb.Registry.MustAddMetricSet(moduleName, metricsetName, New,
		mb.DefaultMetricSet(),
	)
}

// New returns an error.
func New(base mb.BaseMetricSet) (mb.MetricSet, error) {
	return nil, fmt.Errorf("the %v/%v dataset is only supported on Linux", moduleName, metricsetName)
}
