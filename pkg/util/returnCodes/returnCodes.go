// Package returncodes Package errorCode - containing a list of constants to be used
package returncodes

const (
	/*
		Errors
	*/

	/*
		general
	*/

	// Error generic consts
	Error = "error"

	// ErrNotOk - Not OK Error Text
	ErrNotOk = "not ok status received"

	/*
		Data processing
	*/

	// ErrFailedMarshalling failed to marshal data
	ErrFailedMarshalling = "failed to marshal data"

	// ErrFailedUnMarshalling failed to marshal data
	ErrFailedUnMarshalling = "failed to unmarshal data"

	// ErrProcessingData error processing requested data
	ErrProcessingData = "error processing requested data"

	// ErrFetchingRequestData error fetching requested data
	ErrFetchingRequestData = "error fetching requested data"

	// ErrRunningService error running service
	ErrRunningService = "error running service"

	// ErrOpeningFile error opening file
	ErrOpeningFile = "error opening file"

	// ErrFileUnableToCreate unable to create directory or file
	ErrFileUnableToCreate = "unable to create directory or file"

	// ErrFileNotPresent the given file is not present
	ErrFileNotPresent = "the given file is not present"

	// ErrFailedCreateInventory failed to insert data to db
	ErrFailedCreateInventory = "failed to insert inventory data to pbs collection"

	// ErrInventoryAlreadyExist data already exist
	ErrInventoryAlreadyExist = "inventory already exist for serial number"

	// ErrFetchingInventoryData error fetching inventory data
	ErrFetchingInventoryData = "error fetching inventory data"

	// ErrNoIdentifierProvided no SID provided
	ErrNoIdentifierProvided = "no system identifier provided"

	// ErrDataWrongFormat the given data has the wrong format
	ErrDataWrongFormat = "the given data has the wrong format"

	// ErrDataMissing no date is provided
	ErrDataMissing = "no data is provided"

	// NoValidIP the given IP is not valid
	NoValidIP = "the given ip is not valid"

	// NoValidDnsName the given DNS name is not valid
	NoValidDnsName = "the given DNS name is not valid"
	/*
		SUSE Manager Related
	*/

	// ErrLoginSuseManager error while logging in to suse manager
	ErrLoginSuseManager = "error while logging in to suse manager"

	// ErrLogoutSuseManager error while logging out off to suse manager
	ErrLogoutSuseManager = "error while logging out off suse manager"

	// ErrHandlingSuseManagerResponse error while handling suse manager response
	ErrHandlingSuseManagerResponse = "error while handling suse manager response"

	// ErrHTTPSuseManagerResponse error http suse manager response
	ErrHTTPSuseManagerResponse = "error http suse manager response"

	// ErrNotRunningOnSuseManagerServer not running on suse manager server
	ErrNotRunningOnSuseManagerServer = "not running on suse manager server"

	// ErrSystemNotFound error system not found
	ErrSystemNotFound = "error system not found"

	// ErrSystemGroupNotFound error systemgroup not found
	ErrSystemGroupNotFound = "error systemgroup not found"

	// ErrConversionTime represents an error message indicating failure in time conversion operations.
	ErrConversionTime = "error converting time"

	/*
	   Running scripts
	*/

	// FailedCMD failed execute os command
	FailedCMD = "failed to execute command "

	// ErrFailedCreatingDirectory failed to create the directory
	ErrFailedCreatingDirectory = "failed to create directory requested"

	// ExecuteScript generic identifier for os methods
	ExecuteScript = "Execute Script"

	/*
		Information
	*/

	// InfFinishedRunningService FinishedRunningService finished running service
	InfFinishedRunningService = "finished running service"

	/*
		suma server check
	*/
	// NotRunningOnSumaPrim mlmtool is not execute on SUSE Manager Server Primary
	NotRunningOnSumaPrim = "not running on suse manager server primary"
	// NotRunningOnSumaSec mlmtool is not execute on SUSE Manager Server Secondary
	NotRunningOnSumaSec = "not running on suse manager server secondary"
	// NotRunningOnSuma mlmtool is not execute on SUSE Manager Server
	NotRunningOnSuma = "not running on suse manager server"
)
