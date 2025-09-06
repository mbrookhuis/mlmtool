// Package createosrelease create osRelease
package createosrelease

import (
	"fmt"
	"time"

	"go.uber.org/zap"

	_osReleaseModels "mlmtool/pkg/models/createOsRelease"
	_sumamodels "mlmtool/pkg/models/susemanager"
	_sumanUseCase "mlmtool/pkg/usecases/susemanager"
	cmdExecutor "mlmtool/pkg/util/cmdexecutor"
	_consts "mlmtool/pkg/util/consts"
	"mlmtool/pkg/util/contains"
	returnCodes "mlmtool/pkg/util/returnCodes"
)

// CreateOsRelease - general
type CreateOsRelease struct {
	sumanProxy           _sumanUseCase.IProxy
	suse                 _sumanUseCase.ISuseManager
	cmdExec              cmdExecutor.ICMDExecutor
	suseOperationTimeout int
	logger               *zap.Logger
	osRelease            string
}

// NewCreateOsRelease - call createOsRelease
//
// param: sumanProxy
// param: suse
// param: suseOperationTimeout
// param: logger
// param: osRelease
// param: special
// return:
func NewCreateOsRelease(sumanProxy _sumanUseCase.IProxy, suse _sumanUseCase.ISuseManager, cmdExec cmdExecutor.ICMDExecutor, suseOperationTimeout int, logger *zap.Logger, osRelease string) ICreateOsRelease {
	return &CreateOsRelease{
		sumanProxy:           sumanProxy,
		suse:                 suse,
		cmdExec:              cmdExec,
		suseOperationTimeout: suseOperationTimeout,
		logger:               logger,
		osRelease:            osRelease,
	}
}

// CreateOsRelease starting
//
// return: error
func (h *CreateOsRelease) CreateOsRelease() error {
	h.logger.Info("CreateOsRelease started")
	// Login to the SUSE Manager Server and get the sessionkey to be used for api calls
	sessionKey, err := h.sumanProxy.SumanLogin()
	if err != nil {
		h.logger.Error(returnCodes.ErrLoginSuseManager, zap.Any("error", err))
		return err
	}
	// Fetch auth for further use.
	auth, err := h.suse.GetAuth(sessionKey)
	if err != nil {
		return err
	}
	defer func(sumanProxy _sumanUseCase.IProxy, auth _sumanUseCase.AuthParams) {
		err := sumanProxy.SumanLogout(auth)
		if err != nil {
			h.logger.Warn(returnCodes.ErrLogoutSuseManager, zap.Any("error", err))
		}
	}(h.sumanProxy, *auth)
	// Validate if the requested osRelease can be created.
	err = h.validateGivenOsRelease(auth)
	if err != nil {
		return err
	}
	// create osRelease
	err = h.createOsRelease(auth)
	if err != nil {
		return err
	}
	return nil
}

// validateGivenOsRelease - check if given osRelease
//
// param: auth
func (h *CreateOsRelease) validateGivenOsRelease(auth *_sumanUseCase.AuthParams) error {
	// Check if there are already channels present for this osRelease. Abort if this is the case.
	err := h.CheckOsReleaseFormat()
	if err != nil {
		return err
	}
	// Check if the project already exist for this osRelease. Abort if this is the case.
	err = h.CheckOsReleaseCMProjectExists(auth)
	if err != nil {
		return err
	}
	// Check if there are already channels present for this osRelease. Abort if this is the case.
	err = h.CheckOsReleaseChannelExists(auth)
	if err != nil {
		return err
	}
	// Check if there is already a distribution present for this osRelease. Abort if this is the case.
	err = h.CheckOsReleaseDistroExists(auth)
	if err != nil {
		return err
	}
	// Check if the label (first 4 digits) is  correct and defined
	err = h.CheckOsReleaseLabel()
	if err != nil {
		return err
	}
	// Check if the given date is correct
	err = h.CheckOsReleaseDate()
	if err != nil {
		return err
	}
	// Check if the env option makes sense (so no 002 if 001 doesn't exist)
	err = h.CheckOsReleaseEnv()
	if err != nil {
		return err
	}
	err = h.CheckOsReleaseDefaultExtraChannels(auth)
	if err != nil {
		return err
	}

	return nil
}

// CheckOsReleaseFormat - check if format is correct
func (h *CreateOsRelease) CheckOsReleaseFormat() error {
	// Check if the osRelease given has the correct format. If this is not the case an error will be reported.
	h.logger.Debug("Function CheckOsReleaseFormat started", zap.Any("osRelease", h.osRelease))
	if len(h.osRelease) != 16 || h.osRelease[4:5] != "-" || h.osRelease[11:12] != "-" {
		h.logger.Error("Given osRelease has wrong format. Should be ppvv-dddddd-r001", zap.Any("osRelease", h.osRelease))
		return fmt.Errorf(returnCodes.ErrDataWrongFormat)
	}
	h.logger.Debug("Function CheckOsReleaseFormat finished")
	return nil
}

// CheckOsReleaseChannelExists - check if the channel already exists
//
// param: auth
func (h *CreateOsRelease) CheckOsReleaseChannelExists(auth *_sumanUseCase.AuthParams) error {
	// Check if the osRelease is already present in the softwareChannels. If this is the case an error will be reported.
	h.logger.Debug("Function CheckOsReleaseChannelExists started", zap.Any("osRelease", h.osRelease))
	baseChannels, err := h.GetBaseChannels(auth)
	if err != nil {
		h.logger.Error("Error while getting list of baseChannels", zap.Any("error", err))
		return err
	}
	if contains.SubInString(baseChannels, h.osRelease) {
		h.logger.Error("Given osRelease channels already exists.", zap.Any("osRelease", h.osRelease))
		return fmt.Errorf(returnCodes.ErrInventoryAlreadyExist)
	}
	h.logger.Debug("Function CheckOsReleaseChannelExists finished")
	return nil
}

// CheckOsReleaseCMProjectExists - check is project already exists
//
// param: auth
func (h *CreateOsRelease) CheckOsReleaseCMProjectExists(auth *_sumanUseCase.AuthParams) error {
	// Check if the osRelease is already present in the softwareChannels. If this is the case an error will be reported.
	h.logger.Debug("Function CheckOsReleaseCMProjectExists started", zap.Any("osRelease", h.osRelease))
	projects, err := h.sumanProxy.ContentManagementListProjects(*auth)
	if err != nil {
		h.logger.Error("Error while getting list of existing projects", zap.Any("error", err))
		return err
	}
	for _, project := range projects {
		if project.Label == h.osRelease[0:11] {
			h.logger.Error("Given osRelease project already exists.", zap.Any("osRelease", h.osRelease))
			return fmt.Errorf(returnCodes.ErrInventoryAlreadyExist)
		}
	}
	h.logger.Debug("Function CheckOsReleaseCMProjectExists finished")
	return nil
}

// CheckOsReleaseDistroExists - check if distro already exists
//
// param: auth
func (h *CreateOsRelease) CheckOsReleaseDistroExists(auth *_sumanUseCase.AuthParams) error {
	// Check if the distribution used by this osRelease is already present. If this is the case an error will be reported.
	h.logger.Debug("Function CheckOsReleaseDistroExists started", zap.Any("osRelease", h.osRelease))
	_, err := h.sumanProxy.KickstartTreeGetDetails(*auth, h.osRelease)
	if err == nil {
		h.logger.Error("Given osRelease distribution already exists.", zap.Any("osRelease", h.osRelease))
		return fmt.Errorf(returnCodes.ErrInventoryAlreadyExist)
	}
	h.logger.Debug("Function CheckOsReleaseDistroExists finished")
	return nil
}

// CheckOsReleaseLabel - check is label is correct
func (h *CreateOsRelease) CheckOsReleaseLabel() error {
	// Check if the osRelease is already present in the softwareChannels. If this is the case an error will be reported.
	h.logger.Debug("Function CheckOsReleaseLabel started", zap.Any("osRelease", h.osRelease))
	osReleaseLabel := h.osRelease[0:4]
	if !contains.Contains(_consts.CorrectLabels, osReleaseLabel) {
		h.logger.Error("Given osRelease is not correct. The product does not exist.", zap.Any("osRelease", h.osRelease), zap.Any("Product label", osReleaseLabel))
		return fmt.Errorf(returnCodes.ErrDataWrongFormat)
	}
	h.logger.Debug("Function CheckOsReleaseLabel finished")
	return nil
}

// CheckOsReleaseDate - check if date is correct
func (h *CreateOsRelease) CheckOsReleaseDate() error {
	// Check if the osRelease is already present in the softwareChannels. If this is the case an error will be reported.
	h.logger.Debug("Function CheckOsReleaseDate started", zap.Any("osRelease", h.osRelease))
	date := h.osRelease[5:11]
	_, err := time.Parse("060102", date)
	if err != nil {
		h.logger.Error("Given osRelease date is invalid.", zap.Any("osRelease", h.osRelease), zap.Any("error", err))
		return fmt.Errorf(returnCodes.ErrDataWrongFormat)
	}
	h.logger.Debug("Function CheckOsReleaseDate finished")
	return nil
}

// CheckOsReleaseEnv - check is env is correct
func (h *CreateOsRelease) CheckOsReleaseEnv() error {
	// Check if the osRelease is already present in the softwareChannels. If this is the case an error will be reported.
	h.logger.Debug("Function CheckOsReleaseenv started", zap.Any("osRelease", h.osRelease))
	osReleaseEnv := h.osRelease[12:16]
	if !contains.Contains(_consts.CorrectEnvironments, osReleaseEnv) {
		h.logger.Error("Given osRelease env is invalid.", zap.Any("osRelease", h.osRelease))
		return fmt.Errorf(returnCodes.ErrDataWrongFormat)
	}
	h.logger.Debug("Function CheckOsReleaseEnv finished")
	return nil
}

// CheckOsReleaseDefaultExtraChannels
//
// param: auth
func (h *CreateOsRelease) CheckOsReleaseDefaultExtraChannels(auth *_sumanUseCase.AuthParams) error {
	// Check if the default extra channels are present.
	h.logger.Debug("Function CheckOsReleaseDefaultExtraChannels started", zap.Any("osRelease", h.osRelease))

	channelExist, err := h.sumanProxy.ChannelSoftwareIsExisting(*auth, _consts.BaseChannelExtra)
	if err != nil {
		return err
	}
	if !channelExist {
		_, err = h.sumanProxy.ChannelSoftwareCreate(*auth, _consts.BaseChannelExtra, _consts.BaseChannelExtra, _consts.BaseChannelExtra, "channel-x86_64", "")
		if err != nil {
			h.logger.Error(returnCodes.ErrFailedCreateInventory, zap.Any("Failed creating channel", _consts.BaseChannelExtra), zap.Any("error", err))
			return fmt.Errorf(returnCodes.ErrFailedCreateInventory)
		}
	}
	for _, dataRel := range _consts.ListOsReleaseData {
		h.logger.Debug(h.osRelease[0:4])
		h.logger.Debug(dataRel.Label)
		if dataRel.Label == h.osRelease[0:4] {
			for _, childChannel := range dataRel.ChildChannelsExtra {
				h.logger.Debug(childChannel)
				channelExistChild, err := h.sumanProxy.ChannelSoftwareIsExisting(*auth, childChannel)
				if err != nil {
					return err
				}
				if !channelExistChild {
					_, err = h.sumanProxy.ChannelSoftwareCreate(*auth, childChannel, childChannel, childChannel, "channel-x86_64", _consts.BaseChannelExtra)
					if err != nil {
						h.logger.Error(returnCodes.ErrFailedCreateInventory, zap.Any("Failed creating channel", childChannel), zap.Any("error", err))
						return fmt.Errorf(returnCodes.ErrFailedCreateInventory)
					}
				}
			}
		}
	}

	h.logger.Debug("Function CheckOsReleaseEnv finished")
	return nil
}

// createOsRelease - create osRelease
//
// param: auth
func (h *CreateOsRelease) createOsRelease(auth *_sumanUseCase.AuthParams) error {
	/*
			Create CL project
			- create project
			- add channels
			- check if filter exists, if not create filter
		    - assign filter
			- create environment
		*	- sync channels
		*	Create distribution
	*/
	h.logger.Debug("Function createOsRelease started", zap.Any("osRelease", h.osRelease))
	osReleaseData, err := h.getDataProjectOsRelease()
	if err != nil {
		return err
	}
	// create extra repo
	err = h.createExtraRepo(auth)
	if err != nil {
		return err
	}
	// create project
	err = h.createProjectOsRelease(auth)
	if err != nil {
		return err
	}
	// add channels
	err = h.addChannelsProjectOsRelease(auth, osReleaseData)
	if err != nil {
		return err
	}
	// check and create filter
	err = h.createFilterOsRelease(auth)
	if err != nil {
		return err
	}
	// create environment
	err = h.createEnvironmentOsRelease(auth)
	if err != nil {
		return err
	}
	err = h.syncOsRelease(auth)
	if err != nil {
		return err
	}
	err = h.createDistributionOsRelease(auth, osReleaseData)
	if err != nil {
		return err
	}
	h.logger.Debug("Function createOsRelease finished", zap.Any("osRelease", h.osRelease))
	return nil
}

// getDataProjectOsRelease - get needed data
//
// return:
func (h *CreateOsRelease) getDataProjectOsRelease() (_osReleaseModels.OsReleaseRecord, error) {
	// Create os project
	h.logger.Debug("Function getDataOsRelease started", zap.Any("osRelease", h.osRelease))
	for _, dataRel := range _consts.ListOsReleaseData {
		if dataRel.Label == h.osRelease[0:4] {
			return dataRel, nil
		}
	}
	var dummy _osReleaseModels.OsReleaseRecord
	h.logger.Error("No data present for the given OS release", zap.Any("osRelease", h.osRelease))
	h.logger.Debug("Function getDataOsRelease finished", zap.Any("osRelease", h.osRelease))
	return dummy, fmt.Errorf(returnCodes.ErrFailedCreateInventory)
}

// createExtraRepo
//
// param: auth
func (h *CreateOsRelease) createExtraRepo(auth *_sumanUseCase.AuthParams) error {
	h.logger.Debug("Function createExtraRepo started", zap.Any("osRelease", h.osRelease))
	// cmdExec := cmdExecutor.NewCMDExecutor(h.logger)
	channelLabel := fmt.Sprintf("%v-extra", h.osRelease)
	path := _consts.ExtraRepoDir + h.osRelease
	err := h.createExtraRepoFS()
	if err != nil {
		h.logger.Debug(fmt.Sprintf("%v, error %v", returnCodes.FailedCMD, err))
		return err
	}
	// create repo
	_, err = h.sumanProxy.ChannelSoftwareCreateRepo(*auth, channelLabel, "yum", fmt.Sprintf("file://%v", path))
	if err != nil {
		h.logger.Warn("Repository already exists")
	}
	// create channel
	_, err = h.sumanProxy.ChannelSoftwareCreate(*auth, channelLabel, channelLabel, channelLabel, "channel-x86_64", _consts.BaseChannelExtra)
	if err != nil {
		h.logger.Error(returnCodes.ErrFailedCreateInventory, zap.Any("Failed creating channel", channelLabel), zap.Any("error", err))
		return fmt.Errorf(returnCodes.ErrFailedCreateInventory)
	}
	// add repository
	_, err = h.sumanProxy.ChannelSoftwareAssociateRepo(*auth, channelLabel, channelLabel)
	if err != nil {
		h.logger.Error(returnCodes.ErrFailedCreateInventory, zap.Any("Failed associating channel and repository", channelLabel), zap.Any("error", err))
		return fmt.Errorf(returnCodes.ErrFailedCreateInventory)
	}
	// sync channel/repo
	_, err = h.sumanProxy.ChannelSoftwareSyncRepo(*auth, channelLabel)
	if err != nil {
		h.logger.Error(returnCodes.ErrFailedCreateInventory, zap.Any("Failed associating channel and repository", channelLabel), zap.Any("error", err))
		return fmt.Errorf(returnCodes.ErrFailedCreateInventory)
	}
	h.logger.Debug("Function createExtraRepo finished", zap.Any("osRelease", h.osRelease))
	return nil
}

func (h *CreateOsRelease) createExtraRepoFS() error {
	h.logger.Debug("Function createExtraRepoFS started")
	// create dir
	path := fmt.Sprintf("%v%v-extra", _consts.ExtraRepoDir, h.osRelease)
	err := h.cmdExec.CreateDirectory(path)
	if err != nil {
		h.logger.Error(fmt.Sprintf("%v %v with %v", returnCodes.ErrFailedCreatingDirectory, path, err))
		return err
	}
	// run createrepo
	commandRepo := "/usr/bin/createrepo"
	args := []string{path}
	_, err = h.cmdExec.ExecuteCommand(commandRepo, args)
	if err != nil {
		h.logger.Error(fmt.Sprintf("%v %v with %v", returnCodes.FailedCMD, commandRepo, err))
		return err
	}
	h.logger.Debug("Function createExtraRepoFS finished")
	return nil
}

// createProjectOsRelease - create project
//
// param: auth
func (h *CreateOsRelease) createProjectOsRelease(auth *_sumanUseCase.AuthParams) error {
	// Create os project
	h.logger.Debug("Function createProjectOsRelease started", zap.Any("osRelease", h.osRelease))
	name := h.osRelease[0:11]
	_, err := h.sumanProxy.ContentManagementCreate(*auth, name, name, h.osRelease)
	if err != nil {
		return err
	}
	h.logger.Debug("Function createProjectOsRelease finished", zap.Any("osRelease", h.osRelease))
	return nil
}

// addChannelsProjectOsRelease - add channels to project
//
// param: auth
// param: osReleaseData
func (h *CreateOsRelease) addChannelsProjectOsRelease(auth *_sumanUseCase.AuthParams, osReleaseData _osReleaseModels.OsReleaseRecord) error {
	// add channels to project
	h.logger.Debug("Function addChannelsProjectOsRelease started", zap.Any("osRelease", h.osRelease))
	// get needed channels
	var childParentChannel, childExtraChannel []string
	extraChannel := fmt.Sprintf("%s-extra", h.osRelease)
	allChildren, err := h.sumanProxy.ChannelSoftwareListChildren(*auth, osReleaseData.ParentChannel)
	if err != nil {
		return err
	}
	for _, child := range allChildren {
		childParentChannel = append(childParentChannel, child.Label)
	}
	extraChildren, err := h.sumanProxy.ChannelSoftwareListChildren(*auth, _consts.BaseChannelExtra)
	if err != nil {
		return err
	}
	for _, child := range extraChildren {
		childExtraChannel = append(childExtraChannel, child.Label)
	}
	name := h.osRelease[0:11]
	_, err = h.sumanProxy.ContentManagementAttachSource(*auth, name, "software", osReleaseData.ParentChannel)
	if err != nil {
		return err
	}
	for _, dataRel := range _consts.ListOsReleaseData {
		if dataRel.Label == h.osRelease[0:4] {
			for _, childChannel := range dataRel.ChildChannelsDefault {
				if contains.Contains(childParentChannel, childChannel) {
					_, err = h.sumanProxy.ContentManagementAttachSource(*auth, name, "software", childChannel)
					if err != nil {
						return err
					}
				} else {
					h.logger.Warn("Channel doesn't exist", zap.Any("wanted channel", childChannel))
				}
			}
			for _, childChannel := range dataRel.ChildChannelsExtra {
				if contains.Contains(childExtraChannel, childChannel) {
					_, err = h.sumanProxy.ContentManagementAttachSource(*auth, name, "software", childChannel)
					if err != nil {
						return err
					}
				} else {
					h.logger.Warn("Channel doesn't exist", zap.Any("wanted channel", childChannel))
				}
			}
			_, err = h.sumanProxy.ContentManagementAttachSource(*auth, name, "software", extraChannel)
			if err != nil {
				return err
			}
		}
	}
	h.logger.Debug("Function addChannelsProjectOsRelease finished", zap.Any("osRelease", h.osRelease))
	return nil
}

// createFilterOsRelease - create filter
//
// param: auth
func (h *CreateOsRelease) createFilterOsRelease(auth *_sumanUseCase.AuthParams) error {
	// check if filter exists
	h.logger.Debug("Function createFilterOsRelease started", zap.Any("osRelease", h.osRelease))
	filter := "release-" + h.osRelease[0:11]
	date := h.osRelease[5:11]
	name := h.osRelease[0:11]
	// get needed channels
	allFilters, err := h.sumanProxy.ContentManagementListFilters(*auth)
	if err != nil {
		return err
	}
	var filterID int
	filterNotPresent := true
	for _, detailsFilter := range allFilters {
		if detailsFilter.Name == filter {
			h.logger.Debug("Filter is already present", zap.Any("osRelease", h.osRelease), zap.Any("filter", filter))
			filterNotPresent = false
			filterID = detailsFilter.ID
			break
		}
	}
	// if not exists, create filter
	var filterInfo _sumamodels.ContentManagementFilter
	if filterNotPresent {
		filterDate, err := time.Parse("060102 15:04:05", date+" 00:00:00")
		if err != nil {
			return err
		}
		layout := "2006-01-02T15:04:05Z"
		filterDateNew := filterDate.Format(layout)
		var filterCriteria _sumamodels.FilterCriteria
		filterCriteria.Matcher = "greatereq"
		filterCriteria.Value = filterDateNew
		filterCriteria.Field = "issue_date"
		filterInfo, err = h.sumanProxy.ContentManagementCreateFilter(*auth, filter, "deny", "erratum", filterCriteria)
		if err != nil {
			return err
		}
		filterID = filterInfo.ID
	}
	_, err = h.sumanProxy.ContentManagementAttachFilter(*auth, name, filterID)
	if err != nil {
		return err
	}
	h.logger.Debug("Function createFilterOsRelease finished", zap.Any("osRelease", h.osRelease))
	return nil
}

// createEnvironmentOsRelease - create environment
//
// param: auth
func (h *CreateOsRelease) createEnvironmentOsRelease(auth *_sumanUseCase.AuthParams) error {
	// Create os project
	h.logger.Debug("Function createEnvironmentOsRelease started", zap.Any("osRelease", h.osRelease))
	name := h.osRelease[0:11]
	release := h.osRelease[12:16]
	_, err := h.sumanProxy.ContentManagementCreateEnvironment(*auth, name, "", release, release, "Release"+h.osRelease) //  projectLabel string, predecessorLabel string, envlabel string, name string, description string) (sumamodels.ContentManagementEnvironment, error)
	if err != nil {
		return err
	}
	h.logger.Debug("Function createEnvironmentOsRelease finished", zap.Any("osRelease", h.osRelease))
	return nil
}

// syncOsRelease - sync the release
//
// param: auth
func (h *CreateOsRelease) syncOsRelease(auth *_sumanUseCase.AuthParams) error {
	// Create os project
	h.logger.Debug("Function syncOsRelease started", zap.Any("osRelease", h.osRelease))
	name := h.osRelease[0:11]
	_, err := h.sumanProxy.ContentManagementBuildProject(*auth, name) //  projectLabel string, predecessorLabel string, envlabel string, name string, description string) (sumamodels.ContentManagementEnvironment, error)
	if err != nil {
		return err
	}
	h.logger.Debug("Function syncOsRelease finished", zap.Any("osRelease", h.osRelease))
	return nil
}

// createDistributionOsRelease - create distribution
//
// param: auth
// param: osReleaseData
func (h *CreateOsRelease) createDistributionOsRelease(auth *_sumanUseCase.AuthParams, osReleaseData _osReleaseModels.OsReleaseRecord) error {
	// Create os project
	h.logger.Debug("Function createDistributionOsRelease started", zap.Any("osRelease", h.osRelease))
	kernelOptions := fmt.Sprintf("useonlinerepo=1 insecure=1 audit=1 rootdelay=5 install=http://%s/ks/dist/%s self_update=0", auth.Host, h.osRelease)
	_, err := h.sumanProxy.KickstartTreeCreateKernelOptions(*auth, h.osRelease, osReleaseData.TreePath, h.osRelease+"-"+osReleaseData.ParentChannel, "sles15generic", kernelOptions, "") //  projectLabel string, predecessorLabel string, envlabel string, name string, description string) (sumamodels.ContentManagementEnvironment, error)
	if err != nil {
		_, err1 := h.sumanProxy.KickstartTreeCreate(*auth, h.osRelease, osReleaseData.TreePath, h.osRelease+"-"+osReleaseData.ParentChannel, "sles15generic")
		if err1 != nil {
			return err1
		}
		fmt.Println("==================================================================================")
		fmt.Println("Please add the following to the kernel options in the distribution for this osRelease")
		fmt.Println("useonlinerepo=1 insecure=1 audit=1 rootdelay=5")
		fmt.Println("==================================================================================")
		h.logger.Info("==================================================================================", zap.Any("osRelease", h.osRelease))
		h.logger.Info("Please add the following to the kernel options in the distribution for this osRelease", zap.Any("osRelease", h.osRelease))
		h.logger.Info("useonlinerepo=1 insecure=1 audit=1 rootdelay=5", zap.Any("osRelease", h.osRelease))
		h.logger.Info("==================================================================================", zap.Any("osRelease", h.osRelease))
		h.logger.Debug("Function createProjectOsRelease finished", zap.Any("osRelease", h.osRelease))
	}
	return nil
}

// GetBaseChannels - get base channels
//
// param: auth
// return:
func (h *CreateOsRelease) GetBaseChannels(auth *_sumanUseCase.AuthParams) ([]string, error) {
	// This will generate a slice containing all the software base channels.rm
	var resultString []string
	result, err := h.sumanProxy.ChannelListSoftwareChannels(*auth)
	if err != nil {
		h.logger.Error("Unable to get list of software channels", zap.Any("error", err))
		return resultString, err
	}
	for _, getResult := range result {
		if getResult.ParentLabel == "" {
			resultString = append(resultString, getResult.Label)
		}
	}
	return resultString, nil
}
