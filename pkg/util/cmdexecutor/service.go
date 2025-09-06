package cmdexecutor

// ICMDExecutor function definition for all command line scripts
type ICMDExecutor interface {
	ExecuteCommand(binaryName string, args []string) ([]string, error)
	CreateDirectory(path string) error
}
