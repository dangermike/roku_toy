package cmd

import (
	"github.com/spf13/cobra"

	"github.com/dangermike/roku_toy/cmd/channel"
	"github.com/dangermike/roku_toy/cmd/device"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{}

	cmd.AddCommand(device.Cmd(), channel.Cmd())

	return cmd
}
