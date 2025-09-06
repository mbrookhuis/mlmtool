// Package ecpsuma - this is a collection of tools use for SUSE Manager Operations
package mlmtool

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	_sumanUseCase "mlmtool/pkg/usecases/susemanager"
	_updateCMServer "mlmtool/pkg/usecases/updateCMServer"
	"mlmtool/pkg/util/checksumaserver"
	returncodes "mlmtool/pkg/util/returnCodes"
	_suman "mlmtool/pkg/util/suman"
)

var updateCMServerCmd = &cobra.Command{
	Use:   "updateCMServer",
	Short: "updateCMServer to add or modify version used for SUSE Manager Primary",
	Long:  `updateCMServer to add or modify version used for SUSE Manager Primary`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var osRelease, updateServer string
		var highstate bool
		var sumancfg _sumanUseCase.SumanConfig
		sumancfg.Host, _ = cmd.Flags().GetString("suma-master")
		sumancfg.Login, _ = cmd.Flags().GetString("user")
		sumancfg.Password, _ = cmd.Flags().GetString("password")
		osRelease, _ = cmd.Flags().GetString("osrelease")
		updateServer, _ = cmd.Flags().GetString("server")
		highstate, _ = cmd.Flags().GetBool("highstate")
		return executeUpdateCMServer(sumancfg, osRelease, updateServer, highstate)
	},
}

// init
func init() {
	rootCmd.AddCommand(updateCMServerCmd)
	var sumam, user, password, osRelease, updateServer string
	var highstate bool
	updateCMServerCmd.Flags().StringVarP(&sumam, "suma-master", "m", "", "SUSE Manager Primary FQDN")
	updateCMServerCmd.Flags().StringVarP(&user, "user", "u", "", "User to access SUSE Manager Primary")
	updateCMServerCmd.Flags().StringVarP(&password, "password", "p", "", "User to access SUSE Manager Primary")
	updateCMServerCmd.Flags().StringVarP(&osRelease, "osrelease", "o", "", "The OS-Release the server should be updated to.")
	updateCMServerCmd.Flags().StringVarP(&updateServer, "server", "s", "", "The Server to be updated.")
	updateCMServerCmd.Flags().BoolVarP(&highstate, "highstate", "i", false, "run highstate after updates are done but before restart")
	updateCMServerCmd.MarkFlagsRequiredTogether("user", "password")
	updateCMServerCmd.MarkFlagsRequiredTogether("user", "suma-master")
	_ = updateCMServerCmd.MarkFlagRequired("osrelease")
	_ = updateCMServerCmd.MarkFlagRequired("server")
}

// executeUpdateCMServer
//
// param: sumancfg
// param: inputParams
func executeUpdateCMServer(sumancfg _sumanUseCase.SumanConfig, osRelease string, updateServer string, highstate bool) (err error) {

	logger.Debug("params", zap.Any("SumaPrimHost", sumancfg.Host), zap.Any("SumaPrimUser", sumancfg.Login), zap.Any("Server", updateServer), zap.Any("osRelease", osRelease))
	// Check if at least one of the version is defined. If not, return an error
	// all seems to be OK. Continue with the execution

	if len(strings.TrimSpace(sumancfg.Host)) == 0 {
		// it seems that there are no parameters given for sumam, user and password. So checking if server is a suma and read credentials
		// check if server is suse manager server
		if !checksumaserver.Primary() {
			return fmt.Errorf(returncodes.NotRunningOnSumaPrim)
		}

		// get suse manager credentials
		sumancfg, err = _suman.GetCredentials(cfg.FileSpacecmd)
		if err != nil {
			return fmt.Errorf("%s - suse manager credentials: %s\n%s", returncodes.ErrFetchingRequestData, cfg.FileSpacecmd, err)
		}
	}

	suseAPI := _sumanUseCase.NewSuseManagerAPI("rhn/manager/api", true, logger, cfg.RetryCount)
	sumanProxyUseCase := _sumanUseCase.NewProxy(&sumancfg, suseAPI, logger, cfg.RetryCount)
	suseUseCase := _sumanUseCase.NewSuseManager(sumanProxyUseCase, &sumancfg, logger)
	updateCMServer := _updateCMServer.NewUpdateCMServer(sumanProxyUseCase, suseUseCase, 120, logger, osRelease, updateServer, highstate)

	return updateCMServer.UpdateCMServer()
}
