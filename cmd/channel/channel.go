package channel

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/dangermike/roku_toy/aliasing"
	"github.com/dangermike/roku_toy/logging"
	"github.com/dangermike/roku_toy/roku"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "channel",
		Short: "get or set channel",
	}

	cmd.AddCommand(cmdSet(), cmdGet(), cmdList())

	return cmd
}

func cmdSet() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set",
		Short: "Set the current channel by name or number",
		RunE:  setE,
	}
	AddFlags(cmd.Flags())
	return cmd
}

func cmdGet() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get the current channel",
		RunE:  getE,
	}
	AddFlags(cmd.Flags())
	return cmd
}

func cmdList() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List installed channels",
		RunE:  listE,
	}
	AddFlags(cmd.Flags())
	return cmd
}

func setE(cmd *cobra.Command, args []string) error {
	cfg, err := ParseFlags(cmd.Flags())
	if err != nil {
		return nil
	}
	if len(args) != 1 {
		return errors.New("channel name or ID required")
	}
	ctx := logging.NewContext(cmd.Context(), logging.Configure(cfg.Debug))
	device, err := GetDevice(ctx, cfg.Device, cfg.FirstDevice)
	if err != nil {
		return err
	}

	if ID, err := strconv.Atoi(args[0]); err == nil {
		return device.Launch(ctx, ID)
	}
	return device.LaunchByName(ctx, args[0])
}

func getE(cmd *cobra.Command, args []string) error {
	cfg, err := ParseFlags(cmd.Flags())
	if err != nil {
		return nil
	}
	ctx := logging.NewContext(cmd.Context(), logging.Configure(cfg.Debug))
	device, err := GetDevice(ctx, cfg.Device, cfg.FirstDevice)
	if err != nil {
		return err
	}

	app, err := device.ActiveApp(ctx)
	if err != nil {
		return nil
	}
	fmt.Printf("%s (%d)\n", app.Name, app.ID)
	return nil
}

func listE(cmd *cobra.Command, args []string) error {
	cfg, err := ParseFlags(cmd.Flags())
	if err != nil {
		return nil
	}
	ctx := logging.NewContext(cmd.Context(), logging.Configure(cfg.Debug))
	device, err := GetDevice(ctx, cfg.Device, cfg.FirstDevice)
	if err != nil {
		return err
	}

	apps, err := device.QueryApps(ctx)
	if err != nil {
		return err
	}

	fmt.Printf("%s (%d)\n", "home", 0)
	for _, app := range apps {
		fmt.Printf("%s (%d)\n", app.Name, app.ID)
	}
	return nil
}

func AddFlags(flags *pflag.FlagSet) {
	flags.BoolP("first", "1", false, "select device first device found on the network")
	flags.StringP("device", "d", "", "select device by name or USN (required if more than one device on the network)")
	flags.BoolP("verbose", "v", false, "verbose logging")
}

type Cfg struct {
	Debug       bool
	Device      string
	FirstDevice bool
}

func ParseFlags(flags *pflag.FlagSet) (Cfg, error) {
	var cfg Cfg

	return cfg, errors.Join(
		GetFlagT(&cfg.Debug, flags, "verbose", (*pflag.FlagSet).GetBool),
		GetFlagT(&cfg.FirstDevice, flags, "first", (*pflag.FlagSet).GetBool),
		GetFlagT(&cfg.Device, flags, "device", (*pflag.FlagSet).GetString),
	)
}

func GetFlagT[T any](target *T, flags *pflag.FlagSet, field string, extractor func(*pflag.FlagSet, string) (T, error)) error {
	var err error
	*target, err = extractor(flags, field)
	return err
}

func GetDevice(ctx context.Context, device string, first bool) (*roku.Device, error) {
	aliases := map[string]string{}
	if len(device) > 0 {
		al, err := aliasing.Load(ctx)
		if err != nil {
			return nil, err
		}
		for _, al := range al {
			aliases[al.USN] = al.Name
		}
	}

	log := logging.FromContext(ctx)

	var target *roku.Device

	if err := roku.SSDP(ctx, func(dev *roku.Device) error {
		if len(device) > 0 && dev.USN != device && aliases[dev.USN] != device {
			log.Debug("skipping", zap.String("USN", dev.USN))
			return nil
		}
		if target != nil {
			return errors.New("more than one roku device found. device name is required")
		}
		target = dev
		if first || len(device) > 0 {
			return ErrDeviceFound
		}
		return nil
	}); errors.Is(err, ErrDeviceFound) {
		return target, nil
	} else if err != nil {
		return nil, err
	}

	if target == nil {
		return nil, errors.New("no matching roku found")
	}
	return target, nil
}

var ErrDeviceFound = errors.New("found device")
