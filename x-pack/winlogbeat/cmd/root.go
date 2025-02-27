// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

package cmd

import (
	"github.com/snappyflow/beats/v7/winlogbeat/cmd"
	xpackcmd "github.com/snappyflow/beats/v7/x-pack/libbeat/cmd"

	// Register fields.
	_ "github.com/snappyflow/beats/v7/x-pack/winlogbeat/include"
)

// Name of this beat.
var Name = cmd.Name

// RootCmd to handle beats cli
var RootCmd = cmd.RootCmd

func init() {
	xpackcmd.AddXPack(RootCmd, cmd.Name)
}
