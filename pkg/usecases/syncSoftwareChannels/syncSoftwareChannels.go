// Package syncsoftwarechannels - update the software channels on SUSE Manager
package syncsoftwarechannels

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"

	_updateGeneralModels "mlmtool/pkg/models/susemanager"
	_sumanUseCase "mlmtool/pkg/usecases/susemanager"
	utilContains "mlmtool/pkg/util/contains"
)

const (
	fileUyuni = "/opt/uyunihub/uyunihub.yaml"
)

// SyncSoftwareChannels - data call syncSoftwareChannels
type SyncSoftwareChannels struct {
	sumanProxyPrim       _sumanUseCase.IProxy
	sumanProxySec        _sumanUseCase.IProxy
	susePrim             _sumanUseCase.ISuseManager
	suseSec              _sumanUseCase.ISuseManager
	suseoperationtimeout int
	logger               *zap.Logger
}

// NewSyncSoftwareChannels - call syncSoftwareChannels
//
// param: sumanProxyPrim
// param: sumanProxySec
// param: susePrim
// param: suseSec
// param: suseoperationtimeout
// param: logger
// return: ISyncSoftwareChannels
func NewSyncSoftwareChannels(sumanProxyPrim _sumanUseCase.IProxy, sumanProxySec _sumanUseCase.IProxy, susePrim _sumanUseCase.ISuseManager, suseSec _sumanUseCase.ISuseManager, suseoperationtimeout int, logger *zap.Logger) ISyncSoftwareChannels {
	return &SyncSoftwareChannels{
		sumanProxyPrim:       sumanProxyPrim,
		sumanProxySec:        sumanProxySec,
		susePrim:             susePrim,
		suseSec:              suseSec,
		suseoperationtimeout: suseoperationtimeout,
		logger:               logger,
	}
}

// SyncSoftwareChannels - synchronize all software channels
//
// return: error
func (h *SyncSoftwareChannels) SyncSoftwareChannels() error {

	// Suse-Login Primary
	sessionKeyPrim, err := h.sumanProxyPrim.SumanLogin()
	if err != nil {
		h.logger.Error("Error while login to susemanager Primary", zap.Any("error", err))
		return err
	}

	// Fetch auth for further use.
	authPrim, err := h.susePrim.GetAuth(sessionKeyPrim)
	if err != nil {
		return err
	}
	defer func(sumanProxyPrim _sumanUseCase.IProxy, auth _sumanUseCase.AuthParams) {
		err := sumanProxyPrim.SumanLogout(auth)
		if err != nil {
			h.logger.Error("Error during logout from susemanager Server", zap.Any("error", err))
		}
	}(h.sumanProxyPrim, *authPrim)
	// Fetch the groupId of general

	// Suse-Login Secondary
	sessionKeySec, err := h.sumanProxySec.SumanLogin()
	if err != nil {
		h.logger.Error("Error while login to susemanager Secondary", zap.Any("error", err))
		return err
	}
	// Fetch auth for further use.
	authSec, err := h.suseSec.GetAuth(sessionKeySec)
	if err != nil {
		return err
	}
	// Suse-Logout
	defer func(sumanProxySec _sumanUseCase.IProxy, auth _sumanUseCase.AuthParams) {
		err := sumanProxySec.SumanLogout(auth)
		if err != nil {
			h.logger.Error("Error during logout from susemanager Server", zap.Any("error", err))
		}
	}(h.sumanProxySec, *authSec)
	hostName := authSec.Host
	h.getNeededBaseChannels(hostName, authPrim, authSec)
	h.logger.Info("Finished syncSoftwareChannels")
	return nil
}

// GetBaseChannels - get all basechannels labels
//

// param: auth
// return: list of base channel labels
func (h *SyncSoftwareChannels) GetBaseChannels(auth *_sumanUseCase.AuthParams) []string {
	var resultString []string
	result, err := h.sumanProxyPrim.ChannelListSoftwareChannels(*auth)
	if err != nil {
		h.logger.Fatal("Unable to get list of software channels", zap.Any("error", err))
	}
	for _, getResult := range result {
		if getResult.ParentLabel == "" {
			resultString = append(resultString, getResult.Label)
		}
	}
	return resultString
}

// allChannels - get all software channels labels
//

// param: auth
// return: list off channels labels
func (h *SyncSoftwareChannels) allChannels(auth *_sumanUseCase.AuthParams) []string {
	var resultString []string
	result, err := h.sumanProxyPrim.ChannelListSoftwareChannels(*auth)
	if err != nil {
		h.logger.Fatal("Unable to get list of software channels", zap.Any("error", err))
	}
	for _, getResult := range result {
		resultString = append(resultString, getResult.Label)
	}
	return resultString
}

// addChannels - add channels
//

// param: channels
// param: authPrim
// param: authSec
func (h *SyncSoftwareChannels) addChannels(channels []string, authPrim *_sumanUseCase.AuthParams, authSec *_sumanUseCase.AuthParams) {
	var neededChannels []string
	for _, baseChannel := range channels {
		neededChannels = append(neededChannels, baseChannel)
		childChannelsAll, err := h.sumanProxyPrim.ChannelSoftwareListChildren(*authPrim, baseChannel)
		if err != nil {
			h.logger.Error("Unable to get all childchannels for base channel: " + baseChannel)
			h.logger.Panic("error: " + err.Error())
		}
		for _, cChannel := range childChannelsAll {
			neededChannels = append(neededChannels, cChannel.Label)
		}
	}
	allChannels := h.allChannels(authSec)
	var syncChannels []string
	for _, sChannel := range neededChannels {
		if !utilContains.Contains(allChannels, sChannel) {
			syncChannels = append(syncChannels, sChannel)
		}
	}
	h.logger.Info("Adding channels", zap.Any("channels", syncChannels))
	if len(syncChannels) != 0 {
		for _, syncChannel := range syncChannels {
			cmd := exec.Command("/usr/bin/mgr-inter-sync", "-c", syncChannel)
			err := cmd.Run()
			if err == nil {
				h.logger.Info("Successful added the channel", zap.Any("channels to be added", syncChannel))
			} else {
				h.logger.Error("Adding the channel failed", zap.Any("channels to be added", syncChannel), zap.Any("error", err))
			}
		}
	}
}

// getNeededBaseChannels
//

// param: hostName
// param: authPrim
// param: authSec
func (h *SyncSoftwareChannels) getNeededBaseChannels(hostName string, authPrim *_sumanUseCase.AuthParams, authSec *_sumanUseCase.AuthParams) {
	// read uyuniyaml
	ufile, err := os.ReadFile(filepath.Clean(fileUyuni))
	if err != nil {
		h.logger.Error("Unable to open file: " + fileUyuni)
		h.logger.Panic("error: " + err.Error())
	}
	allBaseChannels := h.GetBaseChannels(authPrim)
	uconfig := _updateGeneralModels.UyunihubYaml{}
	err = yaml.Unmarshal(ufile, &uconfig)
	if err != nil {
		h.logger.Error("Converting yaml failed for: " + fileUyuni)
		h.logger.Panic("error: " + err.Error())
	}
	for key, value := range uconfig.AllObjects {
		if key == "all" || key == hostName {
			for _, channel := range value.BaseChannels {
				if utilContains.SubInString(allBaseChannels, channel) {
					var syncBaseChannels []string
					syncBaseChannels = append(syncBaseChannels, channel)
					h.addChannels(syncBaseChannels, authPrim, authSec)
				}
			}
			for _, channel := range value.CLMProjects {
				if utilContains.SubInString(allBaseChannels, channel) {
					var syncBaseChannels []string
					for _, bc := range allBaseChannels {
						if strings.Contains(bc, channel) {
							syncBaseChannels = append(syncBaseChannels, bc)
						}
					}
					h.addChannels(syncBaseChannels, authPrim, authSec)
				}
			}
		}
	}
}
