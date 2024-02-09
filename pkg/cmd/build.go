package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"

	"github.com/kubri/kubri/integrations/apk"
	"github.com/kubri/kubri/integrations/appinstaller"
	"github.com/kubri/kubri/integrations/apt"
	"github.com/kubri/kubri/integrations/sparkle"
	"github.com/kubri/kubri/integrations/yum"
	"github.com/kubri/kubri/pkg/config"
)

func buildCmd() *cobra.Command {
	var configPath string

	cmd := &cobra.Command{
		Use:     "build",
		Short:   "Publish packages for common package managers and software update frameworks",
		Aliases: []string{"b"},
		RunE: func(cmd *cobra.Command, _ []string) error {
			p, err := config.Load(configPath)
			if err != nil {
				return err
			}

			integrations := []struct {
				name string
				fn   func(ctx context.Context) error
			}{
				{"APK", fn(apk.Build, p.Apk)},
				{"App Installer", fn(appinstaller.Build, p.Appinstaller)},
				{"APT", fn(apt.Build, p.Apt)},
				{"YUM", fn(yum.Build, p.Yum)},
				{"Sparkle", fn(sparkle.Build, p.Sparkle)},
			}

			var n int
			start := time.Now()
			g, ctx := errgroup.WithContext(cmd.Context())

			for _, integration := range integrations {
				if integration.fn != nil {
					integration := integration
					n++
					g.Go(func() error {
						log.Print("Publishing " + integration.name + " packages...")
						if err := integration.fn(ctx); err != nil {
							return fmt.Errorf("failed to publish "+integration.name+" packages: %w", err)
						}
						log.Print("Completed publishing " + integration.name + " packages.")
						return nil
					})
				}
			}

			if n == 0 {
				return errors.New("no integrations configured")
			}

			if err := g.Wait(); err != nil {
				return err
			}

			log.Printf("Completed in %s", time.Since(start).Truncate(time.Millisecond))

			return nil
		},
	}

	cmd.Flags().StringVarP(&configPath, "config", "c", "", "load configuration from a file")

	return cmd
}

func fn[C any, F func(ctx context.Context, c *C) error](fn F, config *C) func(ctx context.Context) error {
	if config == nil {
		return nil
	}
	return func(ctx context.Context) error {
		return fn(ctx, config)
	}
}
