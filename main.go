package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/dangermike/roku_toy/roku"

	"github.com/dangermike/roku_toy/cmd"
)

func main() {
	if err := cmd.Cmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func mainEx() error {
	ctx := context.Background()
	rokus := map[string]*roku.Device{}
	if err := roku.SSDP(ctx, func(dev *roku.Device) error {
		fmt.Println(dev)
		rokus[dev.USN] = dev
		return nil
	}); err != nil {
		return fmt.Errorf("failed to discover rokus: %w", err)
	}
	for k, v := range rokus {
		fmt.Println(k, v)
		apps, err := v.QueryApps()
		if err != nil {
			return err
		}
		for _, app := range apps {
			fmt.Println(app)
		}

		app, err := v.ActiveApp()
		if err != nil {
			return err
		}

		if app.Name != "Home" {
			fmt.Println("nobody's home")
		}

		for _, x := range []int{0, 13, 12, 0} {
			fmt.Println("launching", x)
			if err := v.Launch(x); err != nil {
				return err
			}
			time.Sleep(2 * time.Second)
		}

	}
	return nil
}
