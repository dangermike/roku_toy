package device

import (
	"github.com/spf13/cobra"

	"github.com/dangermike/roku_toy/cmd/device/alias"
	"github.com/dangermike/roku_toy/cmd/device/list"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "device",
		Short: "discover and manage Roku devices",
	}

	cmd.AddCommand(list.Cmd(), alias.Alias(), alias.Unalias())

	return cmd
}
