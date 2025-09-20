package mlmtool

import (
	"fmt"
	"os"

	model "mlmtool/pkg/models/inputfile"

	log "mlmtool/pkg/util/logger"
	_ "mlmtool/pkg/util/readconfig"
	ri "mlmtool/pkg/util/readconfig"
	_ "mlmtool/pkg/util/returnCodes"

	"github.com/spf13/cobra"
)

var cfgFile string
var AppConfig model.Config

var rootCmd = &cobra.Command{
	Use:   "mlmtool",
	Short: `mlmtool is a CLI client for core MLM tasks.`,
	Long: `mlmtool is a CLI client for core MLM tasks. This application uses a configuration file which can be specified with
the --config flag.`,
	SilenceUsage: true,

	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Initialize config before any command runs
		err := initConfig()
		if err != nil {
			return err
		}
		// Initialize logger after config is loaded
		err = log.InitLogger(AppConfig)
		if err != nil {
			return err
		}
		log.Debug("Configuration loaded, mlmtool started")
		return nil
	},
}

// Execute runs the root command of the application and handles any errors by exiting with a non-zero status.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// init sets up the application's configuration and flags, ensuring initialization and finalization callbacks are registered.
func init() {
	cobra.OnInitialize(func() {})
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "config.yaml", "config file (default is config.yaml)")
	cobra.OnFinalize(finalizeRun)
}

// initConfig initializes the configuration by reading the configuration file into the AppConfig variable.
// Returns an error if the configuration file cannot be read successfully.
func initConfig() error {
	err := ri.ReadConfig(cfgFile, &AppConfig)
	if err != nil {
		return err
	}
	//err = ri.ReadConfig(constants.GeneralConfigFile, &AppConfig)
	//if err != nil {
	//	return err
	//}
	return nil
}

// PostRun functions seem not to run reliably, at least when I tested
// See: https://github.com/spf13/cobra/issues/914
func finalizeRun() {
	fmt.Println("mlmtool finished")
}
