package alias

import (
	"errors"
	"slices"

	"github.com/dangermike/roku_toy/aliasing"
	"github.com/dangermike/roku_toy/logging"
	"github.com/spf13/cobra"
)

func Alias() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "alias",
		Short: "add an alias for a Roku device",
		RunE:  aliasE,
	}

	cmd.Flags().BoolP("verbose", "v", false, "verbose logging")

	return cmd
}

func Unalias() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unalias",
		Short: "remove an alias for a Roku device",
		RunE:  unaliasE,
	}

	cmd.Flags().BoolP("verbose", "v", false, "verbose logging")

	return cmd
}

func aliasE(cmd *cobra.Command, args []string) error {
	debug, err := cmd.Flags().GetBool("verbose")
	if err != nil {
		return err
	}

	if len(args) != 2 {
		return errors.New("USN and name required")
	}

	ctx := logging.NewContext(cmd.Context(), logging.Configure(debug))

	aliases, err := aliasing.Load(ctx)
	if err != nil {
		return err
	}
	aliases = append(aliases, aliasing.Alias{USN: args[0], Name: args[1]})
	return aliasing.Save(ctx, aliases)
}

func unaliasE(cmd *cobra.Command, args []string) error {
	debug, err := cmd.Flags().GetBool("verbose")
	if err != nil {
		return err
	}

	if len(args) != 1 {
		return errors.New("USN or name required")
	}

	ctx := logging.NewContext(cmd.Context(), logging.Configure(debug))

	aliases, err := aliasing.Load(ctx)
	if err != nil {
		return err
	}

	aliases = slices.DeleteFunc(aliases, func(a aliasing.Alias) bool {
		return a.Name == args[0] || a.USN == args[0]
	})

	return aliasing.Save(ctx, aliases)
}
