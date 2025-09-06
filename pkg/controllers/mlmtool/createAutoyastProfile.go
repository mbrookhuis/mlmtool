// Package ecpsuma - this is a collection of tools use for SUSE Manager Operations
package mlmtool

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	_createAutoyastProfile "mlmtool/pkg/usecases/createAutoyastProfile"
	_sumanUseCase "mlmtool/pkg/usecases/susemanager"
	"mlmtool/pkg/util/checksumaserver"
	_consts "mlmtool/pkg/util/consts"
	returncodes "mlmtool/pkg/util/returnCodes"
	_suman "mlmtool/pkg/util/suman"
)

var createAutoyastProfile = &cobra.Command{
	Use:   "createAutoyastProfile",
	Short: "createAutoyastProfile for given profile",
	Long:  `createAutoyastProfile for given profile`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var sumancfg _sumanUseCase.SumanConfig
		var locationXML, profileName string
		var replaceExisting bool
		sumancfg.Host, _ = cmd.Flags().GetString("suma-master")
		sumancfg.Login, _ = cmd.Flags().GetString("user")
		sumancfg.Password, _ = cmd.Flags().GetString("password")
		locationXML, _ = cmd.Flags().GetString("locationxml")
		profileName, _ = cmd.Flags().GetString("profilename")
		replaceExisting, _ = cmd.Flags().GetBool("replace")
		return executeCreateAutoyastProfile(sumancfg, locationXML, profileName, replaceExisting)
	},
}

// init function
func init() {
	rootCmd.AddCommand(createAutoyastProfile)
	var sumam, user, password, locationXML, profileName string
	var replaceExisting bool
	createAutoyastProfile.Flags().StringVarP(&sumam, "suma-master", "m", "", "SUSE Manager Primary FQDN")
	createAutoyastProfile.Flags().StringVarP(&user, "user", "u", "", "User to access SUSE Manager Primary")
	createAutoyastProfile.Flags().StringVarP(&password, "password", "p", "", "User to access SUSE Manager Primary")
	createAutoyastProfile.Flags().StringVarP(&locationXML, "locationxml", "l", _consts.DefaultAutoyastDir, fmt.Sprintf("Base directory storing xml. Default: %v", _consts.DefaultAutoyastDir))
	createAutoyastProfile.Flags().StringVarP(&profileName, "profilename", "n", "", fmt.Sprintf("profilename name. Required. Select from the following: %s", _consts.AutoyastTypes))
	createAutoyastProfile.Flags().BoolVarP(&replaceExisting, "replace", "r", false, "Replace autoyast profile if it already exists. Default: false")
	createAutoyastProfile.MarkFlagsRequiredTogether("user", "suma-master")
	_ = createAutoyastProfile.MarkFlagRequired("profilename")
}

// execute updating the routes of all servers
func executeCreateAutoyastProfile(sumancfg _sumanUseCase.SumanConfig, locationXML string, profileName string, replaceExisting bool) (err error) {

	logger.Debug("params", zap.Any("SumaPrimHost", sumancfg.Host), zap.Any("SumaPrimUser", sumancfg.Login), zap.Any("locationXML", locationXML), zap.Any("profileName", profileName))

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

	// execute
	suseAPI := _sumanUseCase.NewSuseManagerAPI("rhn/manager/api", true, logger, cfg.RetryCount)
	sumanProxyUseCase := _sumanUseCase.NewProxy(&sumancfg, suseAPI, logger, cfg.RetryCount)
	suseUseCase := _sumanUseCase.NewSuseManager(sumanProxyUseCase, &sumancfg, logger)
	createAutoyastProfile := _createAutoyastProfile.NewCreateAutoyastProfile(sumanProxyUseCase, suseUseCase, 120, logger, locationXML, profileName, replaceExisting)

	return createAutoyastProfile.CreateAutoyastProfile()
}
