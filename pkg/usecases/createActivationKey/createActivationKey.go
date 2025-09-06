// Package createactivationkey - call createActivationKey
package createactivationkey

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"go.uber.org/zap"

	_sumanUseCase "mlmtool/pkg/usecases/susemanager"
	utilContains "mlmtool/pkg/util/contains"
)

// CreateActivationKey - general def
type CreateActivationKey struct {
	sumanProxy           _sumanUseCase.IProxy
	suse                 _sumanUseCase.ISuseManager
	suseoperationtimeout int
	logger               *zap.Logger
}

// NewCreateActivationKey - call createActivationKey
//
// param: sumanProxy
// param: suse
// param: suseoperationtimeout
// param: logger
// return:
func NewCreateActivationKey(sumanProxy _sumanUseCase.IProxy, suse _sumanUseCase.ISuseManager, suseoperationtimeout int, logger *zap.Logger) ICreateActivationKey {
	return &CreateActivationKey{
		sumanProxy:           sumanProxy,
		suse:                 suse,
		suseoperationtimeout: suseoperationtimeout,
		logger:               logger,
	}
}

// CreateActivationKey  - call createActivationKey
//
// return: error
func (h *CreateActivationKey) CreateActivationKey() error {
	h.logger.Info("Create Activation keys started")

	// login to SUSE Manager.
	sessionKey, err := h.sumanProxy.SumanLogin()
	if err != nil {
		h.logger.Error("Error while login to susemanager Primary", zap.Any("error", err))
		return err
	}
	auth, err := h.suse.GetAuth(sessionKey)
	if err != nil {
		return err
	}
	defer func(sumanProxy _sumanUseCase.IProxy, auth _sumanUseCase.AuthParams) {
		err := sumanProxy.SumanLogout(auth)
		if err != nil {
			h.logger.Error("Error during logout from susemanager Server", zap.Any("error", err))
		}
	}(h.sumanProxy, *auth)
	var actKeys []string
	allKeys, err := h.sumanProxy.ActivationKeyListActivationKeys(*auth)
	if err != nil {
		return err
	}
	for _, akey := range allKeys {
		actKeys = append(actKeys, akey.Key)
	}
	baseChannelsAll, err := h.sumanProxy.ChannelListSoftwareChannels(*auth)
	if err != nil {
		return err
	}
	var baseChannels []string
	for _, channel := range baseChannelsAll {
		if channel.ParentLabel == "" {
			baseChannels = append(baseChannels, channel.Label)
		}
	}
	for _, channel := range baseChannels {
		var activationKeyName string
		if !(strings.HasPrefix(channel, "suse-")) && !(strings.HasPrefix(channel, "sle-pr")) && !(strings.HasPrefix(channel, "custom")) && !(strings.Contains(channel, "repo")) {
			if strings.Contains(channel, "special") {
				activationKeyName = fmt.Sprintf("1-%s-%s-%s-%s", strings.Split(channel, "-")[0], strings.Split(channel, "-")[1], strings.Split(channel, "-")[2], strings.Split(channel, "-")[3])
			} else {
				activationKeyName = fmt.Sprintf("1-%s-%s-%s", strings.Split(channel, "-")[0], strings.Split(channel, "-")[1], strings.Split(channel, "-")[2])
			}
			if !utilContains.Contains(actKeys, activationKeyName) {
				if strings.Contains(channel, "-r00") {
					// create the activation key
					var entitlement []string
					if strings.Contains(channel[0:4], "s15") || strings.Contains(channel[0:4], "sm") {
						entitlement = append(entitlement, "monitoring_entitled")
					}
					result, err := h.sumanProxy.ActivationKeyCreate(*auth, activationKeyName[2:], channel, entitlement)
					if err != nil {
						return err
					}
					if result != activationKeyName {
						h.logger.Fatal("Error creating requested activation key", zap.Any("Wanted activationKey", activationKeyName), zap.Any("Received key", result), zap.Any("error", err))
					}
					// add childchannels of the basechannel software
					childChannelsAll, err := h.sumanProxy.ChannelSoftwareListChildren(*auth, channel)
					if err != nil {
						return err
					}
					var childChannels []string
					for _, cChannel := range childChannelsAll {
						childChannels = append(childChannels, cChannel.Label)
					}
					resultAddChildChannels, err := h.sumanProxy.ActivationKeyAddChildChannels(*auth, activationKeyName, childChannels)
					if err != nil {
						return err
					}
					if resultAddChildChannels != 1 {
						h.logger.Fatal("Error adding childchannels to activation key", zap.Any("AcitvationKey", activationKeyName), zap.Any("childchannels", childChannels), zap.Any("error", err))
					}
					// add the systemgroup general
					var groups []int
					resultSystemGroup, _ := h.sumanProxy.SystemGroupGetDetails(*auth, "general")
					groups = append(groups, resultSystemGroup.ID)
					resultAddGroup, err := h.sumanProxy.ActivationKeyAddServerGroups(*auth, activationKeyName, groups)
					if err != nil {
						return err
					}
					if resultAddGroup != 1 {
						h.logger.Fatal("Error adding systemgroup to activation key", zap.Any("AcitvationKey", activationKeyName), zap.Any("groupId", groups), zap.Any("error", err))
					}
					// remove packages that are added by default
					keyDetails, err := h.sumanProxy.ActivationKeyGetDetails(*auth, activationKeyName)
					if err != nil {
						return err
					}
					packages := keyDetails.Packages
					resultRemovePackeges, err := h.sumanProxy.ActivationKeyRemovePackages(*auth, activationKeyName, packages)
					if err != nil {
						return err
					}
					if resultRemovePackeges != 1 {
						h.logger.Fatal("Error removingpackages from activation key", zap.Any("AcitvationKey", activationKeyName), zap.Any("Packages", packages), zap.Any("error", err))
					}
					h.logger.Info("Activationkey created", zap.Any("AcitvationKey", activationKeyName[2:]))
					// create bootstrap script when it doesn't exist.
					scriptName := fmt.Sprintf("--script=%s.sh", activationKeyName[2:])
					if _, err := os.Stat(scriptName); errors.Is(err, os.ErrNotExist) {
						cmd := exec.Command("/usr/bin/mgr-bootstrap", fmt.Sprintf("--activation-keys=%s", activationKeyName), scriptName)
						err := cmd.Run()
						if err == nil {
							h.logger.Info("Successful created bootstrap script", zap.Any("AcitvationKey", activationKeyName[2:]))
						} else {
							h.logger.Error("Bootstrap script creation failed", zap.Any("AcitvationKey", activationKeyName[2:]), zap.Any("error", err))
						}
						scriptLoc := fmt.Sprintf("/srv/www/htdocs/pub/bootstrap/%s.sh", activationKeyName[2:])
						cmd = exec.Command("/usr/bin/sed", "-i", "s/FORCE_VENV_SALT_MINION=0/FORCE_VENV_SALT_MINION=1/", scriptLoc)
						err = cmd.Run()
						if err == nil {
							h.logger.Info("Successful changed AVOID_VENV_SALT_MINION in bootstrap script", zap.Any("AcitvationKey", activationKeyName[2:]))
						} else {
							h.logger.Error("Failed to change AVOID_VENV_SALT_MINION in bootstrap script", zap.Any("AcitvationKey", activationKeyName[2:]), zap.Any("error", err))
						}
					}
				}
			}
		}
	}
	for _, actKey := range actKeys {
		if !utilContains.SubInString(baseChannels, actKey[2:]) && strings.Contains(actKey, "-r0") {
			resultDeleteKey, err := h.sumanProxy.ActivationKeyDelete(*auth, actKey)
			if err != nil {
				return err
			}
			if resultDeleteKey != 1 {
				h.logger.Error("Unable to delete activationKey", zap.Any("AcitvationKey", actKey), zap.Any("error", err))
			} else {
				h.logger.Info("ActivationKey deleted", zap.Any("AcitvationKey", actKey))
			}
		}
	}
	return nil
}
