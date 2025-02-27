// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

// +build !integration

package cloudfoundry

import (
	"testing"
	"time"

	"github.com/cloudfoundry-community/go-cfclient"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/snappyflow/beats/v7/libbeat/logp"
)

func TestClientCacheWrap(t *testing.T) {
	ttl := 500 * time.Millisecond
	guid := mustCreateFakeGuid()
	app := cfclient.App{
		Guid:   guid,
		Memory: 1, // use this field to track if from cache or from client
	}
	fakeClient := &fakeCFClient{app, 0}
	cache := newClientCacheWrap(fakeClient, ttl, ttl, logp.NewLogger("cloudfoundry"))

	missingAppGuid := mustCreateFakeGuid()

	// should err; different app client doesn't have
	one, err := cache.GetAppByGuid(missingAppGuid)
	assert.Nil(t, one)
	assert.True(t, cfclient.IsAppNotFoundError(err))
	assert.Equal(t, 1, fakeClient.callCount)

	// calling again; the miss should be cached
	one, err = cache.GetAppByGuid(missingAppGuid)
	assert.Nil(t, one)
	assert.True(t, cfclient.IsAppNotFoundError(err))
	assert.Equal(t, 1, fakeClient.callCount)

	// fetched from client for the first time
	one, err = cache.GetAppByGuid(guid)
	assert.NoError(t, err)
	assert.Equal(t, app, *one)
	assert.Equal(t, 2, fakeClient.callCount)

	// updated app in fake client, new fetch should not have updated app
	updatedApp := cfclient.App{
		Guid:   guid,
		Memory: 2,
	}
	fakeClient.app = updatedApp
	two, err := cache.GetAppByGuid(guid)
	assert.NoError(t, err)
	assert.Equal(t, app, *two)
	assert.Equal(t, 2, fakeClient.callCount)

	// wait the ttl, then it should have updated app
	time.Sleep(ttl)
	three, err := cache.GetAppByGuid(guid)
	assert.NoError(t, err)
	assert.Equal(t, updatedApp, *three)
	assert.Equal(t, 3, fakeClient.callCount)
}

type fakeCFClient struct {
	app       cfclient.App
	callCount int
}

func (f *fakeCFClient) GetAppByGuid(guid string) (cfclient.App, error) {
	f.callCount++
	if f.app.Guid != guid {
		return cfclient.App{}, notFoundError()
	}
	return f.app, nil
}

func mustCreateFakeGuid() string {
	uuid, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	return uuid.String()
}

// notFoundError returns a cloud foundry error that satisfies cfclient.IsAppNotFoundError(err)
func notFoundError() error {
	return cfclient.CloudFoundryError{Code: 100004}
}
