// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

package localremote

import (
	"github.com/snappyflow/beats/v7/x-pack/elastic-agent/pkg/artifact"
	"github.com/snappyflow/beats/v7/x-pack/elastic-agent/pkg/artifact/download"
	"github.com/snappyflow/beats/v7/x-pack/elastic-agent/pkg/artifact/download/composed"
	"github.com/snappyflow/beats/v7/x-pack/elastic-agent/pkg/artifact/download/fs"
	"github.com/snappyflow/beats/v7/x-pack/elastic-agent/pkg/artifact/download/http"
	"github.com/snappyflow/beats/v7/x-pack/elastic-agent/pkg/artifact/download/snapshot"
	"github.com/snappyflow/beats/v7/x-pack/elastic-agent/pkg/core/logger"
	"github.com/snappyflow/beats/v7/x-pack/elastic-agent/pkg/release"
)

// NewVerifier creates a downloader which first checks local directory
// and then fallbacks to remote if configured.
func NewVerifier(log *logger.Logger, config *artifact.Config) (download.Verifier, error) {
	verifiers := make([]download.Verifier, 0, 3)

	fsVer, err := fs.NewVerifier(config)
	if err != nil {
		return nil, err
	}
	verifiers = append(verifiers, fsVer)

	// try snapshot repo before official
	if release.Snapshot() {
		snapshotVerifier, err := snapshot.NewVerifier(config)
		if err != nil {
			log.Error(err)
		} else {
			verifiers = append(verifiers, snapshotVerifier)
		}
	}

	remoteVer, err := http.NewVerifier(config)
	if err != nil {
		return nil, err
	}
	verifiers = append(verifiers, remoteVer)

	return composed.NewVerifier(verifiers...), nil
}
