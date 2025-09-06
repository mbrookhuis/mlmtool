// Package ecpsuma - this is a collection of tools use for SUSE Manager Operations
package mlmtool

import (
	"fmt"

	"github.com/spf13/cobra"

	_createActivationKey "mlmtool/pkg/usecases/createActivationKey"
	_sumanUseCase "mlmtool/pkg/usecases/susemanager"
	"mlmtool/pkg/util/checksumaserver"
	returncodes "mlmtool/pkg/util/returnCodes"
	_suman "mlmtool/pkg/util/suman"
)

var createActivationKeysCmd = &cobra.Command{
	Use:   "createActivationKeys",
	Short: "createActivationKeys for all POD related software channels",
	Long:  `createActivationKeys for all POD related software channels`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return executeCreateActivationKeys()
	},
}

// init start
func init() {
	rootCmd.AddCommand(createActivationKeysCmd)
}

// executeCreateActivationKeys
func executeCreateActivationKeys() error {
	// check if server is suse manager server
	if !checksumaserver.SumaServer() {
		return fmt.Errorf(returncodes.ErrNotRunningOnSuseManagerServer)
	}
	// Get suse manager credentials
	sumancfg, err := _suman.GetCredentials(cfg.FileSpacecmd)
	if err != nil {
		return fmt.Errorf("%s - suse manager credentials: %s\n%s", returncodes.ErrFetchingRequestData, cfg.FileSpacecmd, err)
	}
	// execute
	suseAPI := _sumanUseCase.NewSuseManagerAPI("rhn/manager/api", true, logger, cfg.RetryCount)
	sumanProxyUseCase := _sumanUseCase.NewProxy(&sumancfg, suseAPI, logger, cfg.RetryCount)
	suseUseCase := _sumanUseCase.NewSuseManager(sumanProxyUseCase, &sumancfg, logger)
	createActivationKey := _createActivationKey.NewCreateActivationKey(sumanProxyUseCase, suseUseCase, 120, logger)

	return createActivationKey.CreateActivationKey()
}
