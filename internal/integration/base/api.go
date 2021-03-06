// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// +build integration_api

package base

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/talos-systems/talos/internal/app/machined/pkg/runtime"
	"github.com/talos-systems/talos/internal/pkg/cluster/check"
	"github.com/talos-systems/talos/internal/pkg/provision"
	"github.com/talos-systems/talos/internal/pkg/provision/access"
	"github.com/talos-systems/talos/pkg/client"
	"github.com/talos-systems/talos/pkg/client/config"
)

// APISuite is a base suite for API tests
type APISuite struct {
	suite.Suite
	TalosSuite

	Client *client.Client
}

// SetupSuite initializes Talos API client
func (apiSuite *APISuite) SetupSuite() {
	cfg, err := config.Open(apiSuite.TalosConfig)
	apiSuite.Require().NoError(err)

	opts := []client.OptionFunc{
		client.WithConfig(cfg),
	}

	if apiSuite.Endpoint != "" {
		opts = append(opts, client.WithEndpoints(apiSuite.Endpoint))
	}

	apiSuite.Client, err = client.New(context.TODO(), opts...)
	apiSuite.Require().NoError(err)
}

// DiscoverNodes provides list of Talos nodes in the cluster.
//
// As there's no way to provide this functionality via Talos API, it works the following way:
// 1. If there's a provided cluster info, it's used.
// 2. If integration test was compiled with k8s support, k8s is used.
func (apiSuite *APISuite) DiscoverNodes() []string {
	discoveredNodes := apiSuite.TalosSuite.DiscoverNodes()
	if discoveredNodes != nil {
		return discoveredNodes
	}

	var err error

	apiSuite.discoveredNodes, err = discoverNodesK8s(apiSuite.Client, &apiSuite.TalosSuite)
	apiSuite.Require().NoError(err, "k8s discovery failed")

	if apiSuite.discoveredNodes == nil {
		// still no nodes, skip the test
		apiSuite.T().Skip("no nodes were discovered")
	}

	return apiSuite.discoveredNodes
}

// Capabilities describes current cluster allowed actions.
type Capabilities struct {
	RunsTalosKernel bool
	SupportsReboot  bool
	SupportsRecover bool
}

// Capabilities returns a set of capabilities to skip tests for different environments.
func (apiSuite *APISuite) Capabilities() Capabilities {
	v, err := apiSuite.Client.Version(context.Background())
	apiSuite.Require().NoError(err)

	caps := Capabilities{}

	if v.Messages[0].Platform != nil {
		switch v.Messages[0].Platform.Mode {
		case runtime.ModeContainer.String():
		default:
			caps.RunsTalosKernel = true
			caps.SupportsReboot = true
			caps.SupportsRecover = true
		}
	}

	return caps
}

// AssertClusterHealthy verifies that cluster is healthy using provisioning checks.
func (apiSuite *APISuite) AssertClusterHealthy(ctx context.Context) {
	if apiSuite.Cluster == nil {
		// can't assert if cluster state was provided
		apiSuite.T().Skip("cluster health can't be verified when cluster state is not provided")
	}

	clusterAccess := access.NewAdapter(apiSuite.Cluster, provision.WithTalosClient(apiSuite.Client))
	defer clusterAccess.Close() //nolint: errcheck

	apiSuite.Require().NoError(check.Wait(ctx, clusterAccess, check.DefaultClusterChecks(), check.StderrReporter()))
}

// ReadUptime reads node uptime.
//
// Context provided might have specific node attached for API call.
func (apiSuite *APISuite) ReadUptime(ctx context.Context) (float64, error) {
	// set up a short timeout around uptime read calls to work around
	// cases when rebooted node doesn't answer for a long time on requests
	reqCtx, reqCtxCancel := context.WithTimeout(ctx, 10*time.Second)
	defer reqCtxCancel()

	reader, errCh, err := apiSuite.Client.Read(reqCtx, "/proc/uptime")
	if err != nil {
		return 0, err
	}

	defer reader.Close() //nolint: errcheck

	var uptime float64

	n, err := fmt.Fscanf(reader, "%f", &uptime)
	if err != nil {
		return 0, err
	}

	if n != 1 {
		return 0, fmt.Errorf("not all fields scanned: %d", n)
	}

	_, err = io.Copy(ioutil.Discard, reader)
	if err != nil {
		return 0, err
	}

	for err = range errCh {
		if err != nil {
			return 0, err
		}
	}

	return uptime, reader.Close()
}

// TearDownSuite closes Talos API client
func (apiSuite *APISuite) TearDownSuite() {
	if apiSuite.Client != nil {
		apiSuite.Assert().NoError(apiSuite.Client.Close())
	}
}
