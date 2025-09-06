// Package createautoyastprofile  - update the routes
package createautoyastprofile

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"

	_sumanUseCase "mlmtool/pkg/usecases/susemanager"
	_consts "mlmtool/pkg/util/consts"
	utilContains "mlmtool/pkg/util/contains"
	returnCodes "mlmtool/pkg/util/returnCodes"
)

// CreateAutoyastProfile - call updateRoutes
type CreateAutoyastProfile struct {
	sumanProxy           _sumanUseCase.IProxy
	suse                 _sumanUseCase.ISuseManager
	suseOperationTimeout int
	logger               *zap.Logger
	locationXML          string
	profileName          string
	replaceExisting      bool
}

// NewCreateAutoyastProfile
//
// param: sumanProxy
// param: suse
// param: suseOperationTimeout
// param: logger
// param: locationXML
// param: profileName
// param: replaceExisting
// return:
func NewCreateAutoyastProfile(sumanProxy _sumanUseCase.IProxy, suse _sumanUseCase.ISuseManager, suseOperationTimeout int, logger *zap.Logger, locationXML string, profileName string, replaceExisting bool) ICreateAutoyastProfile {
	return &CreateAutoyastProfile{
		sumanProxy:           sumanProxy,
		suse:                 suse,
		suseOperationTimeout: suseOperationTimeout,
		logger:               logger,
		locationXML:          locationXML,
		profileName:          profileName,
		replaceExisting:      replaceExisting,
	}
}

// CreateAutoyastProfile
//
// return: error
func (h *CreateAutoyastProfile) CreateAutoyastProfile() error {
	h.logger.Info("createAutoyastProfile started")

	err := h.validateData()
	if err != nil {
		return err
	}
	autoyastXML, err := h.getAutoyastXML()
	if err != nil {
		h.logger.Error(fmt.Sprintf("Error getting data from file %v. Error: %v", h.locationXML, err))
		return err
	}
	sessionKey, err := h.sumanProxy.SumanLogin()
	if err != nil {
		h.logger.Error(fmt.Sprintf("Error while login to susemanager Primary. Err: %v", err))
		return err
	}
	auth, err := h.suse.GetAuth(sessionKey)
	if err != nil {
		return err
	}
	defer func(sumanProxy _sumanUseCase.IProxy, auth _sumanUseCase.AuthParams) {
		err := sumanProxy.SumanLogout(auth)
		if err != nil {
			h.logger.Warn("Unable to logout from SUSE Manager", zap.Any("error", err))
		}
	}(h.sumanProxy, *auth)
	profileExist, err := h.checkProfileIsPresent(auth)
	if err != nil {
		h.logger.Warn("Error checking profile", zap.Error(err))
	}
	if profileExist && h.replaceExisting {
		err := h.deleteProfile(auth)
		if err != nil {
			h.logger.Warn(fmt.Sprintf("Unable to delete profile %v", h.profileName), zap.Any("error", err))
			return err
		}
	} else if profileExist && !h.replaceExisting {
		message := fmt.Errorf("the profile for %s already exist and option replace is not set", h.profileName)
		h.logger.Error(message.Error())
		return message
	}
	err = h.createProfile(autoyastXML, auth)
	if err != nil {
		h.logger.Error(fmt.Sprintf("Error getting data from file %v. Error: %v", h.locationXML, err))
		return err
	}
	err = h.addProfileVar(auth)
	if err != nil {
		h.logger.Error(fmt.Sprintf("Error adding variables to profile %v. Error: %v", h.profileName, err))
		return err
	}
	return nil
}

// validateData
//
// return: error when direcoty doesn't exist or value for name is wrong, otherwise nil
func (h *CreateAutoyastProfile) validateData() error {
	h.logger.Debug("validateData started")
	// check if the profileName exists
	if utilContains.PartOff(_consts.AutoyastTypes, h.profileName) {
		h.logger.Info(fmt.Sprintf("profilename %v is correct", h.profileName))
	} else {
		h.logger.Error(fmt.Sprintf("profilename %v is not correct. Should be one of the following: %v", h.profileName, _consts.AutoyastTypes))
		return errors.New(returnCodes.ErrDataWrongFormat)
	}
	// check if directory exists
	dirXML := fmt.Sprintf("%v/%v", h.locationXML, h.profileName)
	fileXML := fmt.Sprintf("%v/autoyast.xml", dirXML)
	err, errMes := utilContains.Exists(fileXML)
	if err {
		h.logger.Info(fmt.Sprintf("autoyast present in %v", dirXML))
	} else {
		h.logger.Error(fmt.Sprintf("No autoyast.xml present in %v. Error: %v", dirXML, errMes))
		return errors.New(returnCodes.ErrFileNotPresent)
	}
	return nil
}

// getAutoyastXML
//
// return: content of file
// return: error
func (h *CreateAutoyastProfile) getAutoyastXML() (string, error) {
	h.logger.Debug("getAutoyastXML started")
	var uconfig string
	ufile, err := os.ReadFile(filepath.Clean(fmt.Sprintf("%v/%v/autoyast.xml", h.locationXML, h.profileName)))
	if err != nil {
		h.logger.Error(returnCodes.ErrOpeningFile)
		h.logger.Error(fmt.Sprintf("file: %v, error: %v", h.locationXML, err))
		return uconfig, err
	}
	uconfig = string(ufile)
	return uconfig, nil
}

// checkProfileIsPresent
//
// param: auth
// return:
func (h *CreateAutoyastProfile) checkProfileIsPresent(auth *_sumanUseCase.AuthParams) (bool, error) {
	h.logger.Debug("checkProfileIsPresent started")
	profileList, err := h.sumanProxy.KickstartListKickstarts(*auth)
	if err != nil {
		h.logger.Error(returnCodes.ErrDataWrongFormat)
		return false, err
	}
	for _, profileName := range profileList {
		if profileName.Name == h.profileName {
			return true, nil
		}
	}
	return false, nil
}

// deleteProfile
//
// param: auth
func (h *CreateAutoyastProfile) deleteProfile(auth *_sumanUseCase.AuthParams) error {
	h.logger.Debug("deleteProfile")
	_, err := h.sumanProxy.KickstartDeleteProfile(*auth, h.profileName)
	if err != nil {
		h.logger.Error(returnCodes.ErrDataWrongFormat)
		return err
	}
	return nil
}

// createProfile
//
// param: autoyastXML
// param: auth
func (h *CreateAutoyastProfile) createProfile(autoyastXML string, auth *_sumanUseCase.AuthParams) error {
	h.logger.Debug("createProfile started")
	_, err := h.sumanProxy.KickstartImportRawFile(*auth, h.profileName, _consts.VirtType, _consts.AutoyastDistribution, autoyastXML)
	if err != nil {
		h.logger.Error(returnCodes.ErrFailedCreateInventory)
		h.logger.Error(fmt.Sprintf("unable to create autoyast profile for: %v, error: %v", h.profileName, err))
		return err
	}
	return nil
}

// addProfileVar
//
// param: auth
func (h *CreateAutoyastProfile) addProfileVar(auth *_sumanUseCase.AuthParams) error {
	h.logger.Debug("addProfileVar started")
	for _, pN := range _consts.ProfileVariables {
		if pN.ProfileName == h.profileName {
			_, err := h.sumanProxy.KickstartProfileSetVariables(*auth, h.profileName, pN.ProfileVars)
			if err != nil {
				h.logger.Error(returnCodes.ErrFailedCreateInventory)
				h.logger.Error(fmt.Sprintf("unable to add variables to autoyast profile for: %v, error: %v", h.profileName, err))
				return err
			}
			return nil
		}
	}
	return nil
}
