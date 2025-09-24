// Package mlmtool - this is a collection of tools use for SUSE Manager Operations
package mlmtool

import (
	_model "mlmtool/pkg/models/syncStage"
	_syncStage "mlmtool/pkg/usecases/syncStage"

	_sumanUseCase "mlmtool/pkg/usecases/susemanager"
	"mlmtool/pkg/util/logger"

	"github.com/spf13/cobra"
)

/*
	Project       string
	Environment   string
	Backup   bool
	Wait    bool
	Description   string
*/

var syncStageCmd = &cobra.Command{
	Use:   "syncStage",
	Short: "syncStage for given software channel",
	Long:  `syncStage for all given software channel`,
	RunE: func(cmd *cobra.Command, args []string) error {
		project, _ := cmd.Flags().GetString("project")
		environment, _ := cmd.Flags().GetString("environment")
		backup, _ := cmd.Flags().GetBool("backup")
		wait, _ := cmd.Flags().GetBool("wait")
		description, _ := cmd.Flags().GetString("description")
		return executeSyncStage(project, environment, backup, wait, description)
	},
}

// init initializes the syncStageCmd by adding it to the rootCmd and defining its flags.
func init() {
	rootCmd.AddCommand(syncStageCmd)
	var project, environment, description string
	var wait bool
	syncStageCmd.Flags().StringVarP(&project, "project", "p", "",
		"name of the project to be created. Required")
	syncStageCmd.Flags().StringVarP(&environment, "environment", "e", "",
		"Comma delimited list without spaces of the environments to be created. Required")
	syncStageCmd.Flags().BoolVarP(&wait, "wait", "w", false,
		"Wait with finish, until sync is completed. Otherwise the sync runs in the background")
	syncStageCmd.Flags().StringVarP(&description, "description", "d", "",
		"Description of the project to be created.")
	_ = syncStageCmd.MarkFlagRequired("project")
	_ = syncStageCmd.MarkFlagRequired("environment")
}

func executeSyncStage(project string, environment string, backup bool, wait bool, description string) (err error) {
	logger.Debug("syncStage started")
	logger.Debug("params: ")
	logger.Debug("   project: ", project)
	logger.Debug("   environment: ", environment)
	logger.Debug("   wait: ", wait)
	logger.Debug("   description: ", description)

	var sumancfg _sumanUseCase.SumanConfig
	sumancfg.Login = AppConfig.Suman.User
	sumancfg.Password = AppConfig.Suman.Password
	sumancfg.Host = AppConfig.Suman.Server
	sumancfg.Insecure = AppConfig.Suman.SslCertificateCheck

	var inputData _model.InputData
	inputData.Project = project
	inputData.Environment = environment
	inputData.Wait = wait
	inputData.Description = description

	suseAPI := _sumanUseCase.NewSuseManagerAPI("rhn/manager/api", true, AppConfig.Suman.RetryCount)
	sumanProxyUseCase := _sumanUseCase.NewProxy(&sumancfg, suseAPI, AppConfig.Suman.RetryCount)
	suseUseCase := _sumanUseCase.NewSuseManager(sumanProxyUseCase, &sumancfg)
	syncStage := _syncStage.NewSyncStage(sumanProxyUseCase, suseUseCase, AppConfig.Suman.Timeout, AppConfig, inputData)

	return syncStage.SyncStage()
}
