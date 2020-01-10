// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package provision

import (
	"fmt"
	"net"

	"github.com/talos-systems/talos/pkg/config/types/v1alpha1/generate"
)

// ClusterRequest is the root object describing cluster to be provisioned.
type ClusterRequest struct {
	Name string

	Network NetworkRequest
	Nodes   NodeRequests

	Image             string
	KubernetesVersion string
}

// NetworkRequest describe cluster network.
type NetworkRequest struct {
	Name string
	CIDR net.IPNet
	MTU  int
}

// NodeRequests is a list of NodeRequest.
type NodeRequests []NodeRequest

// FindInitNode looks up init node, it returns an error if no init node is present or if it's duplicate.
func (reqs NodeRequests) FindInitNode() (req NodeRequest, err error) {
	found := false

	for i := range reqs {
		if reqs[i].Type == generate.TypeInit {
			if found {
				err = fmt.Errorf("duplicate init node in requests")
				return
			}

			req = reqs[i]
			found = true
		}
	}

	if !found {
		err = fmt.Errorf("no init node found in requests")
	}

	return
}

// MasterNodes returns subset of nodes which are Init/ControlPlane type.
func (reqs NodeRequests) MasterNodes() (nodes []NodeRequest) {
	for i := range reqs {
		if reqs[i].Type == generate.TypeInit || reqs[i].Type == generate.TypeControlPlane {
			nodes = append(nodes, reqs[i])
		}
	}

	return
}

// WorkerNodes returns subset of nodes which are Init/ControlPlane type.
func (reqs NodeRequests) WorkerNodes() (nodes []NodeRequest) {
	for i := range reqs {
		if reqs[i].Type == generate.TypeJoin {
			nodes = append(nodes, reqs[i])
		}
	}

	return
}

// NodeRequest describes a request for a node.
type NodeRequest struct {
	Type       generate.Type
	Name       string
	IP         net.IP
	ConfigData string

	// Share of CPUs, in 1e-9 fractions
	NanoCPUs int64
	// Memory limit in bytes
	Memory int64
}