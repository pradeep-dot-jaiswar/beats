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

package cmd

import (
	"flag"

	"github.com/spf13/pflag"

	cmd "github.com/snappyflow/beats/v7/libbeat/cmd"
	"github.com/snappyflow/beats/v7/libbeat/cmd/instance"
	"github.com/snappyflow/beats/v7/libbeat/common"
	"github.com/snappyflow/beats/v7/libbeat/publisher/processing"
	"github.com/snappyflow/beats/v7/packetbeat/beater"

	// Register fields and protocol modules.
	_ "github.com/snappyflow/beats/v7/packetbeat/include"
)

const (
	// Name of this beat.
	Name = "packetbeat"

	// ecsVersion specifies the version of ECS that Packetbeat is implementing.
	ecsVersion = "1.5.0"
)

// withECSVersion is a modifier that adds ecs.version to events.
var withECSVersion = processing.WithFields(common.MapStr{
	"ecs": common.MapStr{
		"version": ecsVersion,
	},
})

// RootCmd to handle beats cli
var RootCmd *cmd.BeatsRootCmd

func init() {
	var runFlags = pflag.NewFlagSet(Name, pflag.ExitOnError)
	runFlags.AddGoFlag(flag.CommandLine.Lookup("I"))
	runFlags.AddGoFlag(flag.CommandLine.Lookup("t"))
	runFlags.AddGoFlag(flag.CommandLine.Lookup("O"))
	runFlags.AddGoFlag(flag.CommandLine.Lookup("l"))
	runFlags.AddGoFlag(flag.CommandLine.Lookup("dump"))

	settings := instance.Settings{
		RunFlags:      runFlags,
		Name:          Name,
		HasDashboards: true,
		Processing:    processing.MakeDefaultSupport(true, withECSVersion, processing.WithHost, processing.WithAgentMeta()),
	}
	RootCmd = cmd.GenRootCmdWithSettings(beater.New, settings)
	RootCmd.AddCommand(genDevicesCommand())
}
