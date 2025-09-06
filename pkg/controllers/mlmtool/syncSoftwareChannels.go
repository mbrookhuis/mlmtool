// Package ecpsuma - this is a collection of tools use for SUSE Manager Operations
package mlmtool

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	_updateGeneralModels "mlmtool/pkg/models/susemanager"
	_sumanUseCase "mlmtool/pkg/usecases/susemanager"
	_syncSoftwareChannels "mlmtool/pkg/usecases/syncSoftwareChannels"
	"mlmtool/pkg/util/checksumaserver"
	returncodes "mlmtool/pkg/util/returnCodes"
	_suman "mlmtool/pkg/util/suman"
)

var syncSoftwareChannels = &cobra.Command{
	Use:   "syncSoftwareChannels",
	Short: "syncSoftwareChannels: sync software channels between SUSE Manager Primary and Secondary",
	Long:  `syncSoftwareChannels: sync software channels between SUSE Manager Primary and Secondary`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("doing syncSoftwareChannels")
		return executeSyncSoftwareChannels()
	},
}

// init
func init() {
	rootCmd.AddCommand(syncSoftwareChannels)
}

// executeSyncSoftwareChannels add and update the defined software channels
func executeSyncSoftwareChannels() error {
	// check if server is suse manager server
	if !checksumaserver.Secondary() {
		return fmt.Errorf(returncodes.NotRunningOnSumaSec)
	}
	// get suse manager credentials
	ufile, err := os.ReadFile(cfg.FileUyuni)
	if err != nil {
		return fmt.Errorf("%s: %s\n%s", returncodes.ErrOpeningFile, cfg.FileUyuni, err)
	}
	uconfig := _updateGeneralModels.UyunihubYaml{}
	err = yaml.Unmarshal(ufile, &uconfig)
	if err != nil {
		return fmt.Errorf("%s: %s\n%s", returncodes.ErrFailedUnMarshalling, cfg.FileUyuni, err)
	}
	sumancfgPrim := _sumanUseCase.SumanConfig{
		Host:     uconfig.ServerMgmt.MasterURL,
		Login:    uconfig.ServerMgmt.MasterUser,
		Password: uconfig.ServerMgmt.MasterPw,
		Insecure: true,
	}
	sumancfgSec, err := _suman.GetCredentials(cfg.FileSpacecmd)
	if err != nil {
		return fmt.Errorf("%s - suse manager credentials: %s\n%s", returncodes.ErrFetchingRequestData, cfg.FileSpacecmd, err)
	}
	suseAPI := _sumanUseCase.NewSuseManagerAPI("rhn/manager/api", true, logger, cfg.RetryCount)
	sumanProxyUseCasePrim := _sumanUseCase.NewProxy(&sumancfgPrim, suseAPI, logger, cfg.RetryCount)
	sumanProxyUseCaseSec := _sumanUseCase.NewProxy(&sumancfgSec, suseAPI, logger, cfg.RetryCount)
	suseUseCasePrim := _sumanUseCase.NewSuseManager(sumanProxyUseCasePrim, &sumancfgPrim, logger)
	suseUseCaseSec := _sumanUseCase.NewSuseManager(sumanProxyUseCaseSec, &sumancfgSec, logger)
	syncSoftwareChannels := _syncSoftwareChannels.NewSyncSoftwareChannels(sumanProxyUseCasePrim, sumanProxyUseCaseSec, suseUseCasePrim, suseUseCaseSec, 120, logger)

	return syncSoftwareChannels.SyncSoftwareChannels()
}
