// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package talos

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"github.com/talos-systems/talos/internal/pkg/cluster"
	"github.com/talos-systems/talos/pkg/client"
)

// crashdumpCmd represents the crashdump command
var crashdumpCmd = &cobra.Command{
	Use:   "crashdump",
	Short: "Dump debug information about the cluster",
	Long:  ``,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return WithClient(func(ctx context.Context, c *client.Client) error {
			clientProvider := &cluster.ConfigClientProvider{
				DefaultClient: c,
			}
			defer clientProvider.Close() //nolint: errcheck

			worker := cluster.APICrashDumper{
				ClientProvider: clientProvider,
				Info:           &clusterState,
			}

			worker.CrashDump(ctx, os.Stdout)

			return nil
		})
	},
}

func init() {
	addCommand(crashdumpCmd)
	crashdumpCmd.Flags().StringVar(&clusterState.InitNode, "init-node", "", "specify IPs of init node")
	crashdumpCmd.Flags().StringSliceVar(&clusterState.ControlPlaneNodes, "control-plane-nodes", nil, "specify IPs of control plane nodes")
	crashdumpCmd.Flags().StringSliceVar(&clusterState.WorkerNodes, "worker-nodes", nil, "specify IPs of worker nodes")
}
