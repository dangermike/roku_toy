package list

import (
	"fmt"

	"github.com/dangermike/roku_toy/aliasing"
	"github.com/dangermike/roku_toy/logging"
	"github.com/dangermike/roku_toy/roku"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "show all local Roku devices",
		RunE:  listE,
	}

	cmd.Flags().BoolP("verbose", "v", false, "verbose logging")

	return cmd
}

func listE(cmd *cobra.Command, args []string) error {
	debug, err := cmd.Flags().GetBool("verbose")
	if err != nil {
		return err
	}

	ctx := logging.NewContext(cmd.Context(), logging.Configure(debug))
	aliases := map[string]string{}
	al, err := aliasing.Load(ctx)
	if err != nil {
		return fmt.Errorf("failed to load aliases: %w", err)
	}
	for _, a := range al {
		aliases[a.USN] = a.Name
	}

	if err := roku.SSDP(ctx, func(dev *roku.Device) error {
		if alias, ok := aliases[dev.USN]; ok {
			fmt.Println(dev.USN, dev.Location, alias)
		} else {
			fmt.Println(dev.USN, dev.Location)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("failed to discover rokus: %w", err)
	}
	return nil
}
