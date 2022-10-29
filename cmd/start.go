package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/Stoakes/authz-checker/internal/authzchecker"
	"github.com/Stoakes/authz-checker/internal/config"

	"github.com/Stoakes/go-pkg/log"
	"github.com/Stoakes/go-pkg/toolsz"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	flags := startCmd.Flags()

	flags.StringP("log_level", "l", "info", "Zap compliant log level parameter: debug, info, warn, error...")
	flags.Bool("debug", false, "Enable non structured logger, text output and log_level to debug")
	flags.String("config", "", "Path to authz-checker JSON configuration page")
}

var startCmd = &cobra.Command{
	Use:     "start",
	Short:   "Start the authz-checker server",
	Long:    ``,
	Example: ``,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := viper.BindPFlags(cmd.Flags()); err != nil {
			return err
		}
		if err := viper.BindPFlags(cmd.PersistentFlags()); err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithCancel(cmd.Context())
		defer cancel()

		// initialize logger
		zapLogLevel := log.Setup(ctx, log.Options{
			Debug:         viper.GetBool("debug"),
			LogLevel:      viper.GetString("log_level"),
			AppName:       "authz-checker",
			DisableFields: true,
		})

		// errorsChan collects errors from workloads
		// channel length = number of major goroutines running: authz-checker server, tools server,
		errorsChan := make(chan workloadResult, 2)

		/* Read config file */
		var configFilePath string
		var err error
		if len(viper.GetString("config")) != 0 {
			configFilePath, err = filepath.Abs(viper.GetString("config"))
			if err != nil {
				return err
			}
			log.Bg().Debugf("Reading config from %s", configFilePath)
		}
		appConfig, err := config.ReadAppConfigFromJSON(configFilePath)
		if err != nil {
			return fmt.Errorf("Cannot read app config from %s: %s", configFilePath, err.Error())
		}

		/* Create and start authz checker server */
		appServer, err := authzchecker.New(authzchecker.Parameters{
			Port:      appConfig.Port,
			AppConfig: appConfig,
		})
		if err != nil {
			return err
		}
		go startWorkload(ctx, appServer, "server", errorsChan)

		toolsServer := toolsz.New(appConfig.ToolsPort, zapLogLevel, nil)
		go startWorkload(ctx, toolsServer, "toolsz", errorsChan)

		/* Signal handler */
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		defer func() {
			signal.Stop(c)
		}()

		select {
		case <-c:
			log.For(ctx).Info("Termination signal received")
			cancel()
		case result := <-errorsChan:
			if result.err != nil {
				log.For(ctx).Errorf("%s stopped with %s. Shutting down", result.name, result.err)
			} else {
				log.For(ctx).Errorf("%s stopped without errors. Shutting down", result.name, result.err)
			}
			cancel()
			time.Sleep(time.Second)

			return result.err
		case <-ctx.Done():
		}

		return nil
	},
}

// workload wraps major components for simpler start
type workload interface {
	Start(ctx context.Context) error
}

// workloadResult helps collecting major components result
type workloadResult struct {
	name string
	err  error
}

// startWorkload starts a workload and takes care of error handling
func startWorkload(ctx context.Context, s workload, name string, errChan chan workloadResult) {
	err := s.Start(ctx)
	errChan <- workloadResult{name: name, err: err}
}
