package mlmtool

import (
	"fmt"

	"github.com/spf13/cobra"

	_consts "mlmtool/pkg/util/constants"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "display version of mlmtool",
	Long:  `display version of mlmtool`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("mlmtool Version: %v\n", _consts.MlmToolVersion)
	},
}

// init function
func init() {
	rootCmd.AddCommand(versionCmd)
}
