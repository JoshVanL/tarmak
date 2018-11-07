// Copyright Jetstack Ltd. See LICENSE for details.
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jetstack/tarmak/pkg/tarmak"
	"github.com/jetstack/tarmak/pkg/tarmak/snapshot/etcd"
)

var clusterSnapshotEtcdSaveCmd = &cobra.Command{
	Use:   "save [target path prefix]",
	Short: "save etcd snapshot to target path prefix, i.e 'backup-'",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("expecting single target path, got=%d", len(args))
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		t := tarmak.New(globalFlags)
		s := etcd.New(t, args[0])
		t.CancellationContext().WaitOrCancel(t.NewCmdSnapshot(cmd.Flags(), args, s).Save)
	},
}

func init() {
	clusterSnapshotEtcdCmd.AddCommand(clusterSnapshotEtcdSaveCmd)
}
