// Package ecpsuma - this is a collection of tools use for SUSE Manager Operations
package mlmtool

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	_createOsRelease "mlmtool/pkg/usecases/createOsRelease"
	_sumanUseCase "mlmtool/pkg/usecases/susemanager"
	"mlmtool/pkg/util/checksumaserver"
	cmdExecutor "mlmtool/pkg/util/cmdexecutor"
	returncodes "mlmtool/pkg/util/returnCodes"
	_suman "mlmtool/pkg/util/suman"
)

var createOsReleaseCmd = &cobra.Command{
	Use:   "createOsRelease",
	Short: "createOsRelease for given software channel",
	Long:  `createOsRelease for all given software channel`,
	RunE: func(cmd *cobra.Command, args []string) error {
		sumas, _ := cmd.Flags().GetString("suma-master")
		user, _ := cmd.Flags().GetString("user")
		password, _ := cmd.Flags().GetString("password")
		osRelease, _ := cmd.Flags().GetString("osrelease")
		return executeCreateOsRelease(osRelease, sumas, user, password)
	},
}

// init
func init() {
	rootCmd.AddCommand(createOsReleaseCmd)
	var sumam, user, password, osRelease string
	createOsReleaseCmd.Flags().StringVarP(&sumam, "suma-master", "m", "", "SUSE Manager Primary FQDN")
	createOsReleaseCmd.Flags().StringVarP(&user, "user", "u", "", "User to access SUSE Manager Primary")
	createOsReleaseCmd.Flags().StringVarP(&password, "password", "p", "", "User to access SUSE Manager Primary")
	createOsReleaseCmd.Flags().StringVarP(&osRelease, "osrelease", "o", "", "osRelease to be created.")
	createOsReleaseCmd.MarkFlagsRequiredTogether("user", "password")
	createOsReleaseCmd.MarkFlagsRequiredTogether("user", "suma-master")
	_ = createOsReleaseCmd.MarkFlagRequired("osrelease")
}

// executeCreateOsRelease
//
// param: osRelease
// param: sumam
// param: user
// param: password
func executeCreateOsRelease(osRelease string, sumam string, user string, password string) (err error) {

	var sumanHost, sumanPassword, sumanUser string
	logger.Debug("params", zap.Any("sumam", sumam), zap.Any("user", user), zap.Any("osRelease", osRelease))

	sumancfg := _sumanUseCase.SumanConfig{
		Host:     sumanHost,
		Login:    sumanUser,
		Password: sumanPassword,
		Insecure: true,
	}

	if len(strings.TrimSpace(sumam)) == 0 {
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
	} else {
		sumancfg.Login = user
		sumancfg.Password = password
		sumancfg.Host = sumam
	}

	suseAPI := _sumanUseCase.NewSuseManagerAPI("rhn/manager/api", true, logger, cfg.RetryCount)
	sumanProxyUseCase := _sumanUseCase.NewProxy(&sumancfg, suseAPI, logger, cfg.RetryCount)
	suseUseCase := _sumanUseCase.NewSuseManager(sumanProxyUseCase, &sumancfg, logger)
	cmdExec := cmdExecutor.NewCMDExecutor(logger)
	createOsRelease := _createOsRelease.NewCreateOsRelease(sumanProxyUseCase, suseUseCase, cmdExec, 120, logger, osRelease)

	return createOsRelease.CreateOsRelease()
}
