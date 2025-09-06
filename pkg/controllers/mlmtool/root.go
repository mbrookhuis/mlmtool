package mlmtool

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"mlmtool/pkg/config"
	logging "mlmtool/pkg/util/logger"
	returncodes "mlmtool/pkg/util/returnCodes"
)

var (
	cfg     *config.Config
	debug   bool
	logger  *zap.Logger
	slogger *zap.SugaredLogger

	rootCmd = &cobra.Command{
		Use:   "mlmtool",
		Short: `mlmtool is a CLI client for core ECP tasks.`,
		Long:  `mlmtool manages CeML and POD related bootstrapping, updating, networking/dns related tasks`,
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		// cobra does handle printing the error message, if any
		// fmt.Println(err)
		if slogger != nil {
			slogger.Fatalf("%s: %s\n%s", returncodes.ErrRunningService, cfg.Name, err)
		}
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	cobra.OnFinalize(finalizeRun)
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Development mode. Debug logs. For tests debug logging, export DEBUG environment variable")
}

func initConfig() {
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, _ []string) {

		name := cmd.Name()

		// Do not reinitilize cfg or logger to local scope!
		cfg = config.New(name, debug)
		logger = logging.New(cfg)
		slogger = logger.Sugar()

		// ctx := context.WithValue(cmd.Context(), "logger", logger)
		// cmd.SetContext(ctx)

		slogger.Debugf("Custom logger initialized: %s", cfg.Name)
	}
	rootCmd.PersistentPostRun = func(cmd *cobra.Command, args []string) {

		slogger.Infof("%s successfully: %s", returncodes.InfFinishedRunningService, cfg.Name)
	}
}

// PostRun functions seem not to run reliably, at least when I tested
// See: https://github.com/spf13/cobra/issues/914
func finalizeRun() {
	slogger.Debugf("PostRun closing logger: %s", cfg.Name)
	defer func() {
		_ = logger.Sync()
	}()
}
