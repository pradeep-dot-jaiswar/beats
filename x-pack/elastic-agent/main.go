// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/snappyflow/beats/v7/x-pack/elastic-agent/pkg/agent/cmd"
)

// Setups and Runs agent.
func main() {
	rand.Seed(time.Now().UnixNano())
	command := cmd.NewCommand()
	if err := command.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
