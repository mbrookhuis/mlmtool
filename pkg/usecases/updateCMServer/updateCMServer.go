// Package updatecmserver update the given server to the given osrelease
package updatecmserver

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	_sumanModels "mlmtool/pkg/models/susemanager"
	_sumanUseCase "mlmtool/pkg/usecases/susemanager"
	util "mlmtool/pkg/util/contains"
	returnCodes "mlmtool/pkg/util/returnCodes"
)

// UpdateCMServer data description
type UpdateCMServer struct {
	sumanProxy           _sumanUseCase.IProxy
	suse                 _sumanUseCase.ISuseManager
	suseOperationTimeout int
	logger               *zap.Logger
	osRelease            string
	updateServer         string
	highstate            bool
}

// ServerHostInfo structure of data needed
type ServerHostInfo struct {
	serverID           int
	currentBaseChannel string
	newBaseChannel     string
	serverType         string
	currentOsRelease   string
	spMig              bool
	testingMode        bool
}

// NewUpdateCMServer update given server
//
// param: sumanProxy
// param: suse
// param: suseOperationTimeout
// param: logger
// param: versionList
// param: runSync
// return: UpdateCMServer
func NewUpdateCMServer(sumanProxy _sumanUseCase.IProxy, suse _sumanUseCase.ISuseManager, suseOperationTimeout int, logger *zap.Logger, osRelease string, updateServer string, highstate bool) IUpdateCMServer {
	return &UpdateCMServer{
		sumanProxy:           sumanProxy,
		suse:                 suse,
		suseOperationTimeout: suseOperationTimeout,
		logger:               logger,
		osRelease:            osRelease,
		updateServer:         updateServer,
		highstate:            highstate,
	}
}

// UpdateCMServer update the server
//
// return: error
func (h *UpdateCMServer) UpdateCMServer() error {
	zf := []zapcore.Field{}
	// Login to the SUSE Manager Server and get the sessionkey to be used for api calls
	sessionKey, err := h.sumanProxy.SumanLogin()
	if err != nil {
		h.logger.Error(fmt.Sprintf("%v - error %v", returnCodes.ErrLoginSuseManager, err), zf...)
		return fmt.Errorf("%v - error %v", returnCodes.ErrLoginSuseManager, err)
	}
	// Fetch auth for further use.
	auth, err := h.suse.GetAuth(sessionKey)
	if err != nil {
		return err
	}
	defer func(sumanProxy _sumanUseCase.IProxy, auth _sumanUseCase.AuthParams) {
		err := sumanProxy.SumanLogout(auth)
		if err != nil {
			h.logger.Error(fmt.Sprintf("%v - error %v", returnCodes.ErrLogoutSuseManager, err), zf...)
		}
	}(h.sumanProxy, *auth)
	// Validate if the received data
	serverInfo, err := h.validateData(auth, zf...)
	if err != nil {
		return err
	}
	if strings.Contains(serverInfo.currentBaseChannel, h.osRelease) {
		zf = append(zf, zap.Any("Status: ", "No update needed. Correct osRelease already assigned."))
		h.logger.Info("CreateOsRelease status", zf...)
		return nil
	}
	h.logger.Debug(fmt.Sprintf("Current base channel: %v", serverInfo.currentBaseChannel), zf...)
	h.logger.Debug(fmt.Sprintf("New base channel: %v", serverInfo.newBaseChannel), zf...)
	h.logger.Debug(fmt.Sprintf("spmig %v", serverInfo.spMig), zf...)
	if serverInfo.spMig {
		err = h.doSPMig(serverInfo, auth, zf...)
		if err != nil {
			return err
		}
	} else {
		err = h.doUpdate(serverInfo, auth, zf...)
		if err != nil {
			return err
		}
	}
	return nil
}

// validateData
//
// param: auth// return: serverHostInfo
// return: error
func (h *UpdateCMServer) validateData(auth *_sumanUseCase.AuthParams, zf ...zapcore.Field) (ServerHostInfo, error) {
	/*
		check
		- server present
		- server type k3s* --> 3 active servers in rancher
		- osRelease present
		- can osRelease be used
	*/
	serverInfo, err := h.checkServer(auth, zf...)

	if err != nil {
		return serverInfo, err
	}
	if h.checkOsReleaseExists(auth, zf...) != nil {
		return serverInfo, err
	}
	serverInfo, err = h.updateValid(serverInfo, auth, zf...)
	if err != nil {
		return serverInfo, err
	}
	return serverInfo, nil
}

// checkServer
//
// param: auth// return: ServerHostInfo
// return: error
func (h *UpdateCMServer) checkServer(auth *_sumanUseCase.AuthParams, zf ...zapcore.Field) (ServerHostInfo, error) {
	h.logger.Debug("checkServer started", zf...)
	var serverInfo ServerHostInfo
	activeSystems, err := h.sumanProxy.SystemListActiveSystems(*auth)
	if err != nil {
		return serverInfo, err
	}
	// serverNotFound := true
	for _, activeSystem := range activeSystems {
		if util.PartOff(activeSystem.Name, h.updateServer) {
			// serverNotFound = false
			serverInfo.serverID = activeSystem.ID
			// serverDetail, err := h.sumanProxy.SystemID(*auth)
			baseChannel, err := h.sumanProxy.SystemGetSubscribedBaseChannel(*auth, serverInfo.serverID)
			if err != nil {
				return serverInfo, err
			}
			serverInfo.currentBaseChannel = baseChannel.Label
			formDataRaw, err := h.sumanProxy.GetSystemFormulaData(*auth, serverInfo.serverID, "mgts-srv")
			if err != nil {
				return serverInfo, err
			}
			var formData _sumanModels.MgtsSrvFormular
			byteArray, err := json.Marshal(formDataRaw)
			if err != nil {
				h.logger.Error(returnCodes.ErrFailedMarshalling, zf...)
				return serverInfo, err
			}
			err = json.Unmarshal(byteArray, &formData)
			if err != nil {
				h.logger.Error(returnCodes.ErrFailedUnMarshalling, zf...)
				return serverInfo, err
			}
			serverInfo.serverType = formData.Mgntenv.MgtsType
			return serverInfo, nil
		}
	}
	return serverInfo, errors.New(returnCodes.ErrSystemNotFound)
}

// checkOsReleaseExists
//
// param: auth
// param: zf
// return: error
func (h *UpdateCMServer) checkOsReleaseExists(auth *_sumanUseCase.AuthParams, zf ...zapcore.Field) error {
	// Check if the osRelease is already present in the softwareChannels. If this is the case an error will be reported.
	h.logger.Debug("Function CheckOsReleaseChannelExists started", zf...)
	baseChannels, err := h.getBaseChannels(auth)
	if err != nil {
		h.logger.Error(fmt.Sprintf("Error while getting list of baseChannels: %v", err.Error()), zf...)
		return err
	}
	if !util.SubInString(baseChannels, h.osRelease) {
		h.logger.Error("Given osRelease doesn't exist", zf...)
		return errors.New(returnCodes.ErrDataMissing)
	}
	h.logger.Info("Given osRelease exists.", zf...)
	return nil
}

// getBaseChannels
//
// param: auth
// return: []string of basechannel labels
// return: error
func (h *UpdateCMServer) getBaseChannels(auth *_sumanUseCase.AuthParams) ([]string, error) {
	// This will generate a slice containing all the software base channels.rm
	var resultString []string
	result, err := h.sumanProxy.ChannelListSoftwareChannels(*auth)
	if err != nil {
		return resultString, err
	}
	for _, getResult := range result {
		if getResult.ParentLabel == "" {
			resultString = append(resultString, getResult.Label)
		}
	}
	return resultString, nil
}

// updateValid
//
// param: currentOsRelease// return:
//
//nolint:funlen
func (h *UpdateCMServer) updateValid(serverInfo ServerHostInfo, auth *_sumanUseCase.AuthParams, zf ...zapcore.Field) (ServerHostInfo, error) {
	h.logger.Debug("updateValid started", zf...)
	if util.PartOff(h.osRelease, "special") {
		if !util.PartOff(serverInfo.currentBaseChannel, "special") {
			return serverInfo, fmt.Errorf("given osRelease is of type special and current is not. given %s, current %s", h.osRelease, serverInfo.currentBaseChannel)
		}
	}
	if util.PartOff(serverInfo.currentBaseChannel, "special") {
		if !util.PartOff(h.osRelease, "special") {
			releaseDate, err := time.Parse("060102 15:04:05", h.osRelease[5:11]+" 00:00:00")
			if err != nil {
				return serverInfo, err
			}
			checkDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.Local)
			h.logger.Debug(releaseDate.String())
			h.logger.Debug(checkDate.String())

			if checkDate.After(releaseDate) {
				return serverInfo, fmt.Errorf("given osRelease is not of type special and current is. given %s, current %s", h.osRelease, serverInfo.currentBaseChannel)
			}
		}
	}
	if util.PartOff(serverInfo.currentBaseChannel, "special") {
		serverInfo.currentOsRelease = serverInfo.currentBaseChannel[0:24]
	} else {
		serverInfo.currentOsRelease = serverInfo.currentBaseChannel[0:16]
	}
	if h.osRelease == serverInfo.currentOsRelease {
		// return errors.New(fmt.Sprintf("given osRelease is equal to existing. current %s, given %s", currentOsRelease, h.osRelease))
		h.logger.Warn("given osRelease is equal to existing.", zf...)
		return serverInfo, nil
	}
	if h.osRelease[0:2] != serverInfo.currentOsRelease[0:2] {
		return serverInfo, fmt.Errorf("OS type are not equal: current %s, given %s", serverInfo.currentOsRelease[0:4], h.osRelease[0:4])
	}
	osVersion, err := strconv.ParseInt(h.osRelease[2:4], 10, 64)
	if err != nil {
		return serverInfo, err
	}
	osVersionCurrent, err := strconv.ParseInt(serverInfo.currentOsRelease[2:4], 10, 64)
	if err != nil {
		return serverInfo, err
	}
	if osVersion < osVersionCurrent {
		return serverInfo, fmt.Errorf("given osRelease is older then existing. current %s, given %s", serverInfo.currentOsRelease, h.osRelease)
	}
	if util.PartOff(serverInfo.currentOsRelease, "special") && util.PartOff(h.osRelease, "special") {
		newDate, err := time.Parse("060102 15:04:05", h.osRelease[13:19]+" 00:00:00")
		if err != nil {
			return serverInfo, err
		}
		currentDate, err := time.Parse("060102 15:04:05", serverInfo.currentOsRelease[13:19]+" 00:00:00")
		if err != nil {
			return serverInfo, err
		}
		if currentDate.After(newDate) {
			return serverInfo, fmt.Errorf("given osRelease is older then existing. current %s, given %s", serverInfo.currentOsRelease, h.osRelease)
		}
	} else if util.PartOff(serverInfo.currentOsRelease, "special") && !util.PartOff(h.osRelease, "special") {
		newDate, err := time.Parse("060102 15:04:05", h.osRelease[5:11]+" 00:00:00")
		if err != nil {
			return serverInfo, err
		}
		currentDate, err := time.Parse("060102 15:04:05", serverInfo.currentOsRelease[13:19]+" 00:00:00")
		if err != nil {
			return serverInfo, err
		}
		if currentDate.After(newDate) {
			return serverInfo, fmt.Errorf("given osRelease is older then existing. current %s, given %s", serverInfo.currentOsRelease, h.osRelease)
		}
	} else {
		newDate, err := time.Parse("060102 15:04:05", h.osRelease[5:11]+" 00:00:00")
		if err != nil {
			return serverInfo, err
		}
		currentDate, err := time.Parse("060102 15:04:05", serverInfo.currentOsRelease[5:11]+" 00:00:00")
		if err != nil {
			return serverInfo, err
		}
		if currentDate.After(newDate) {
			return serverInfo, fmt.Errorf("given osRelease is older then existing. current %s, given %s", serverInfo.currentOsRelease, h.osRelease)
		}
	}
	serverInfo.newBaseChannel, serverInfo.spMig, err = h.getBaseChannelAndSPUpgrade(serverInfo, auth)
	if err != nil {
		return serverInfo, err
	}
	return serverInfo, nil
}

// doUpdate
//
// param: serverInfo
// param: auth
// param: zf
// return: error
func (h *UpdateCMServer) doUpdate(serverInfo ServerHostInfo, auth *_sumanUseCase.AuthParams, zf ...zapcore.Field) error {
	h.logger.Debug("doUpdate started", zf...)
	// set new channels
	err := h.setNewChannels(serverInfo, auth, zf...)
	if err != nil {
		return err
	}
	// run apply updates non sle-micro
	if h.osRelease[0:2] == "mi" {
		// run apply updates sle-micro
		err := h.sumanProxy.ScheduleScriptRun(*auth, serverInfo.serverID, 3000, "transactional-update up || echo")
		if err != nil {
			return err
		}
	} else {
		// run apply updates venv salt minion SLES
		err := h.sumanProxy.ScheduleScriptRun(*auth, serverInfo.serverID, 240, "zypper -n up venv-salt-minion")
		if err != nil {
			return err
		}
		if !serverInfo.testingMode {
			time.Sleep(90 * time.Second)
		}
		// run apply updates on SLES
		err = h.sumanProxy.ScheduleScriptRun(*auth, serverInfo.serverID, 3000, "zypper -n up --allow-vendor-change")
		if err != nil {
			return err
		}

		if !serverInfo.testingMode {
			time.Sleep(30 * time.Second)
		}
		// reboot server non sle-micro
		if h.highstate {
			err = h.sumanProxy.SystemScheduleApplyHighstate(*auth, serverInfo.serverID, 3000)
			if err != nil {
				return err
			}
		}
		err = h.sumanProxy.SystemScheduleReboot(*auth, serverInfo.serverID, 3000)
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *UpdateCMServer) doSPMig(serverInfo ServerHostInfo, auth *_sumanUseCase.AuthParams, zf ...zapcore.Field) error {
	h.logger.Debug("doSPMig started", zf...)
	err := h.setNewChannels(serverInfo, auth, zf...)
	if err != nil {
		return err
	}
	if h.osRelease[0:2] == "mi" {
		// run apply upgrade sle-micro
		err = h.sumanProxy.ScheduleScriptRun(*auth, serverInfo.serverID, 3000, "transactional-update dup || echo")
		if err != nil {
			return err
		}
	} else {
		// run apply updates venv salt minion SLES
		err := h.sumanProxy.ScheduleScriptRun(*auth, serverInfo.serverID, 240, "zypper -n up venv-salt-minion")
		if err != nil {
			return err
		}
		if !serverInfo.testingMode {
			time.Sleep(90 * time.Second)
		}
		// run apply distribution update on SLES
		err = h.sumanProxy.ScheduleScriptRun(*auth, serverInfo.serverID, 3000, "zypper -n dup --allow-vendor-change")
		if err != nil {
			return err
		}

		if !serverInfo.testingMode {
			time.Sleep(30 * time.Second)
		}
		// reboot server non sle-micro
		if h.highstate {
			err = h.sumanProxy.SystemScheduleApplyHighstate(*auth, serverInfo.serverID, 3000)
			if err != nil {
				return err
			}
		}
		err = h.sumanProxy.SystemScheduleReboot(*auth, serverInfo.serverID, 3000)
		if err != nil {
			return err
		}
	}
	return nil
}

// setNewChannels
//
// param: childChannels
// param: serverInfo
// param: auth
// return: error
func (h *UpdateCMServer) setNewChannels(serverInfo ServerHostInfo, auth *_sumanUseCase.AuthParams, zf ...zapcore.Field) error {
	h.logger.Debug("setNewChannels started", zf...)
	zf = append(zf, zap.Any("New Basechannel: ", serverInfo.newBaseChannel))
	h.logger.Debug("setNewChannels info", zf...)
	return h.suse.ChangeChannels(*auth, serverInfo.serverID, serverInfo.newBaseChannel)
}

// getBaseChannelAndSPUpgrade
//
// param: auth
// param: zf
// return:
func (h *UpdateCMServer) getBaseChannelAndSPUpgrade(serverInfo ServerHostInfo, auth *_sumanUseCase.AuthParams, zf ...zapcore.Field) (string, bool, error) {
	h.logger.Debug("getBaseChannelAndSPUpgrade", zf...)

	// check if SPMig is needed
	osVersion, err := strconv.ParseInt(h.osRelease[2:4], 10, 64)
	if err != nil {
		return "", false, err
	}
	osVersionCurrent, err := strconv.ParseInt(serverInfo.currentBaseChannel[2:4], 10, 64)
	if err != nil {
		return "", false, err
	}
	spMig := false
	if osVersion > osVersionCurrent {
		spMig = true
	} else {
		spMig = false
		// return serverInfo.currentBaseChannel, false, nil
	}
	baseChannels, err := h.getBaseChannels(auth)
	if err != nil {
		h.logger.Error(fmt.Sprintf("Error while getting list of baseChannels: %v", err.Error()), zf...)
		return "", false, err
	}
	for _, baseChannel := range baseChannels {
		if util.PartOff(baseChannel, h.osRelease) {
			return baseChannel, spMig, nil
		}
	}
	return "", spMig, fmt.Errorf("unable to get basechannel for osRelease %v", h.osRelease)
}
