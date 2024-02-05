package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/abemedia/appcast/integrations/apk"
	"github.com/abemedia/appcast/integrations/appinstaller"
	"github.com/abemedia/appcast/integrations/apt"
	"github.com/abemedia/appcast/integrations/sparkle"
	"github.com/abemedia/appcast/integrations/yum"
	"github.com/abemedia/appcast/pkg/config"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

func buildCmd() *cobra.Command {
	var configPath string

	cmd := &cobra.Command{
		Use:     "build",
		Short:   "Build appcast feed",
		Aliases: []string{"b"},
		RunE: func(cmd *cobra.Command, args []string) error {
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
