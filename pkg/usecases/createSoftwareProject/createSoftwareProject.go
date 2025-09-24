package createSoftwareProject

import (
	"fmt"
	csp "mlmtool/pkg/models/createSoftwareProject"
	"mlmtool/pkg/models/inputfile"
	"strings"

	_ "mlmtool/pkg/models/susemanager"
	_sumanUseCase "mlmtool/pkg/usecases/susemanager"
	log "mlmtool/pkg/util/logger"
	returnCodes "mlmtool/pkg/util/returnCodes"
	"reflect"
	_ "strings"
)

type CreateSoftwareProject struct {
	sumanProxy           _sumanUseCase.IProxy
	suse                 _sumanUseCase.ISuseManager
	suseoperationtimeout int
	genConfig            inputfile.Config
	input                csp.InputData
}

func NewCreateSoftwareProject(sumanProxy _sumanUseCase.IProxy, suse _sumanUseCase.ISuseManager, suseoperationtimeout int, genConfig inputfile.Config, input csp.InputData) *CreateSoftwareProject {
	return &CreateSoftwareProject{
		sumanProxy:           sumanProxy,
		suse:                 suse,
		suseoperationtimeout: suseoperationtimeout,
		genConfig:            genConfig,
		input:                input,
	}
}

// CreateSoftwareProject creates a new software project or updates an existing one based on the provided input parameters.
// It logs in to the Suse Manager, validates the input data, and performs the necessary operations to configure the project.
// Returns an error if any step in the process fails, including login, validation, or project creation/update steps.
func (h *CreateSoftwareProject) CreateSoftwareProject() error {
	log.Debug("CreateSoftwareProject started")

	sessionKey, err := h.sumanProxy.SumanLogin()
	if err != nil {
		log.Error(fmt.Sprintf("%v - error %v", returnCodes.ErrLoginSuseManager, err))
		return err
	}
	var authParm _sumanUseCase.AuthParams
	authParm.Host = h.genConfig.Suman.Server
	authParm.SessionKey = sessionKey
	err = h.validateCreateSoftwareProject(authParm)
	if err != nil {
		return err
	}
	err = h.doCreateSoftwareProject(authParm)
	if err != nil {
		return err
	}
	log.Info("CreateSoftwareProject finished")
	return nil
}

// validateCreateSoftwareProject validates input data for creating a software project based on the provided parameters.
// Ensures mandatory fields (project, basechannel, environment) are present and the basechannel exists.
// Returns an error if validation fails or any required field is missing.
func (h *CreateSoftwareProject) validateCreateSoftwareProject(authParm _sumanUseCase.AuthParams) error {
	log.Debug("DoSystem validateCreateSoftwareProject started")
	if len(h.input.Project) == 0 {
		return fmt.Errorf("project name is mandatory")
	}
	if len(h.input.BaseChannel) == 0 {
		return fmt.Errorf("basechannel is mandatory")
	}
	if len(h.input.Environment) == 0 {
		return fmt.Errorf("environment is mandatory")
	}
	channelPresent, err := h.sumanProxy.ChannelSoftwareIsExisting(authParm, h.input.BaseChannel)
	if err != nil {
		return err
	}
	if !channelPresent {
		return fmt.Errorf("given basechannel doesn't exist")
	}
	log.Debug("DoSystem validateCreateSoftwareProject finished")
	return nil
}

// doCreateSoftwareProject creates a new software project or updates an existing one based on the provided input parameters.
// It performs the necessary operations to configure the project.
// Returns an error if any step in the process fails, including project creation/update steps.
func (h *CreateSoftwareProject) doCreateSoftwareProject(authParm _sumanUseCase.AuthParams) error {
	log.Debug("doCreateSoftwareProject started")
	project, err := h.sumanProxy.ContentManagementLookupProject(authParm, h.input.Project)
	if err != nil {
		return err
	}
	if reflect.ValueOf(project).IsZero() {
		err := h.doSoftwareProject(authParm)
		if err != nil {
			return err
		}
	} else {
		log.Info(fmt.Sprintf("project %v already exists. Only addChannel or deleteChannel are being processed", h.input.Project))
		if len(h.input.AddChannel) > 0 {
			err := h.doChannel(authParm, h.input.AddChannel, "add")
			if err != nil {
				return err
			}
			log.Info(fmt.Sprintf("added channels %v to project %v", h.input.AddChannel, h.input.Project))
		}
		if len(h.input.DeleteChannel) > 0 {
			err := h.doChannel(authParm, h.input.DeleteChannel, "delete")
			if err != nil {
				return err
			}
			log.Info(fmt.Sprintf("deleted channels %v from project %v", h.input.DeleteChannel, h.input.Project))
		}
	}
	log.Debug("doCreateSoftwareProject finished")
	return nil
}

func (h *CreateSoftwareProject) doSoftwareProject(authParm _sumanUseCase.AuthParams) error {
	log.Debug("doSoftwareProject started")
	if len(h.input.Description) == 0 {
		h.input.Description = h.input.Project
	}
	_, err := h.sumanProxy.ContentManagementCreate(authParm, h.input.Project, h.input.Project, h.input.Description)
	if err != nil {
		return err
	}
	preEnv := ""
	for _, env := range strings.Split(h.input.Environment, ",") {
		_, err = h.sumanProxy.ContentManagementCreateEnvironment(authParm, h.input.Project, preEnv, env, env, h.input.Description)
		if err != nil {
			return err
		}
		preEnv = env
	}
	childChannels := h.input.BaseChannel
	if len(h.input.AddChannel) > 0 {
		childChannels = childChannels + "," + h.input.AddChannel
	} else {
		children, err := h.sumanProxy.ChannelSoftwareListChildren(authParm, h.input.BaseChannel)
		if err != nil {
			return err
		}
		for _, child := range children {
			childChannels = childChannels + "," + child.Label
		}
	}
	err = h.doChannel(authParm, childChannels, "add")
	if err != nil {
		return err
	}
	if len(h.input.DeleteChannel) > 0 {
		err = h.doChannel(authParm, h.input.DeleteChannel, "delete")
		if err != nil {
			return err
		}
	}
	log.Debug("doSoftwareProject finished")
	return nil
}

// doChannel manages the addition or removal of provided channels to/from a software project based on the specified action.
// It splits the channels string, checks if each channel exists, and accordingly performs the requested operation (add/delete).
// Returns an error if any step fails, including channel existence check or attach/detach operations.
func (h *CreateSoftwareProject) doChannel(authParam _sumanUseCase.AuthParams, channels string, action string) error {
	log.Debug("doChannel finished")
	for _, channel := range strings.Split(channels, ",") {
		present, err := h.sumanProxy.ChannelSoftwareIsExisting(authParam, channel)
		if err != nil {
			return err
		}
		if present {
			if action == "add" {
				_, err := h.sumanProxy.ContentManagementAttachSource(authParam, h.input.Project, "software", channel)
				if err != nil {
					return err
				}
			}
			if action == "delete" {
				err := h.sumanProxy.ContentManagementDetachSource(authParam, h.input.Project, "software", channel)
				if err != nil {
					return err
				}
			}
		} else {
			log.Debug(fmt.Sprintf("channel %v not found", channel))
		}
	}
	return nil
}
