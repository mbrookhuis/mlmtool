// Package mlmtool - this is a collection of tools use for SUSE Manager Operations
package mlmtool

import (
	_model "mlmtool/pkg/models/createSoftwareProject"
	_createSoftwareProject "mlmtool/pkg/usecases/createSoftwareProject"

	_sumanUseCase "mlmtool/pkg/usecases/susemanager"
	"mlmtool/pkg/util/logger"

	"github.com/spf13/cobra"
)

var createSoftwareProjectCmd = &cobra.Command{
	Use:   "createSoftwareProject",
	Short: "createSoftwareProject for given software channel",
	Long:  `createSoftwareProject for all given software channel`,
	RunE: func(cmd *cobra.Command, args []string) error {
		project, _ := cmd.Flags().GetString("project")
		environment, _ := cmd.Flags().GetString("environment")
		baseChannel, _ := cmd.Flags().GetString("basechannel")
		addChannel, _ := cmd.Flags().GetString("addchannel")
		deleteChannel, _ := cmd.Flags().GetString("deletechannel")
		description, _ := cmd.Flags().GetString("description")
		return executeCreateSoftwareProject(project, environment, baseChannel, addChannel, deleteChannel, description)
	},
}

// init initializes the createSoftwareProjectCmd by adding it to the rootCmd and defining its flags.
func init() {
	rootCmd.AddCommand(createSoftwareProjectCmd)
	var project, environment, baseChannel, addChannel, deleteChannel, description string
	createSoftwareProjectCmd.Flags().StringVarP(&project, "project", "p", "",
		"name of the project to be created. Required")
	createSoftwareProjectCmd.Flags().StringVarP(&environment, "environment", "e", "",
		"Comma delimited list without spaces of the environments to be created. Required")
	createSoftwareProjectCmd.Flags().StringVarP(&baseChannel, "basechannel", "b", "",
		"The base channel on which the project should be based.")
	createSoftwareProjectCmd.Flags().StringVarP(&addChannel, "addchannel", "a", "",
		"Comma delimited list without spaces of the channels to be added. Can be used together with --basechannel")
	createSoftwareProjectCmd.Flags().StringVarP(&deleteChannel, "deletechannel", "d", "",
		"Comma delimited list without spaces of the channels to be removed from the project.")
	createSoftwareProjectCmd.Flags().StringVarP(&description, "description", "m", "",
		"Description of the project to be created.")
	_ = createSoftwareProjectCmd.MarkFlagRequired("SoftwareProject")
}

// executeCreateSoftwareProject initializes and executes the process to create or update a software project.
// It configures SUSE Manager API, processes the provided parameters, and invokes the necessary workflows.
// Returns an error if any step, including parameter validation, SUSE Manager login, or project creation fails.
func executeCreateSoftwareProject(project string, environment string, baseChannel string, addChannel string, deleteChannel string, description string) (err error) {
	logger.Debug("CreateSoftwareProject started")
	logger.Debug("params: ")
	logger.Debug("   project: ", project)
	logger.Debug("   environment: ", environment)
	logger.Debug("   baseChannel: ", baseChannel)
	logger.Debug("   addChannel: ", addChannel)
	logger.Debug("   deleteChannel: ", deleteChannel)
	logger.Debug("   description: ", description)

	var sumancfg _sumanUseCase.SumanConfig
	sumancfg.Login = AppConfig.Suman.User
	sumancfg.Password = AppConfig.Suman.Password
	sumancfg.Host = AppConfig.Suman.Server
	sumancfg.Insecure = AppConfig.Suman.SslCertificateCheck

	var inputData _model.InputData
	inputData.Project = project
	inputData.Environment = environment
	inputData.BaseChannel = baseChannel
	inputData.AddChannel = addChannel
	inputData.DeleteChannel = deleteChannel
	inputData.Description = description

	suseAPI := _sumanUseCase.NewSuseManagerAPI("rhn/manager/api", true, AppConfig.Suman.RetryCount)
	sumanProxyUseCase := _sumanUseCase.NewProxy(&sumancfg, suseAPI, AppConfig.Suman.RetryCount)
	suseUseCase := _sumanUseCase.NewSuseManager(sumanProxyUseCase, &sumancfg)
	createSoftwareProject := _createSoftwareProject.NewCreateSoftwareProject(sumanProxyUseCase, suseUseCase, AppConfig.Suman.Timeout, AppConfig, inputData)

	return createSoftwareProject.CreateSoftwareProject()
}
