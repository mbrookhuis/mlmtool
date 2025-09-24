package syncStage

import (
	"fmt"
	"mlmtool/pkg/models/inputfile"
	csp "mlmtool/pkg/models/syncStage"
	_sumanUseCase "mlmtool/pkg/usecases/susemanager"
	log "mlmtool/pkg/util/logger"
	returnCodes "mlmtool/pkg/util/returnCodes"
	"reflect"
	"time"
)

type SyncStage struct {
	sumanProxy           _sumanUseCase.IProxy
	suse                 _sumanUseCase.ISuseManager
	suseoperationtimeout int
	genConfig            inputfile.Config
	input                csp.InputData
}

func NewSyncStage(sumanProxy _sumanUseCase.IProxy, suse _sumanUseCase.ISuseManager, suseoperationtimeout int, genConfig inputfile.Config, input csp.InputData) *SyncStage {
	return &SyncStage{
		sumanProxy:           sumanProxy,
		suse:                 suse,
		suseoperationtimeout: suseoperationtimeout,
		genConfig:            genConfig,
		input:                input,
	}
}

func (h *SyncStage) SyncStage() error {
	log.Debug("SyncStage started")
	sessionKey, err := h.sumanProxy.SumanLogin()
	if err != nil {
		log.Error(fmt.Sprintf("%v - error %v", returnCodes.ErrLoginSuseManager, err))
		return err
	}
	var authParm _sumanUseCase.AuthParams
	authParm.Host = h.genConfig.Suman.Server
	authParm.SessionKey = sessionKey
	err = h.validateSyncStage(authParm)
	if err != nil {
		return err
	}
	err = h.doSyncStage(authParm)
	if err != nil {
		return err
	}
	if h.input.Wait {
		err = h.waitUntilFinished(authParm)
		if err != nil {
			return err
		}
	}
	log.Info("SyncStage finished")
	return nil
}

func (h *SyncStage) validateSyncStage(authParm _sumanUseCase.AuthParams) error {
	log.Debug("syncStage validateSyncStage started")
	if len(h.input.Project) == 0 {
		return fmt.Errorf("project name is mandatory")
	}
	project, err := h.sumanProxy.ContentManagementLookupProject(authParm, h.input.Project)
	if err != nil {
		return err
	}
	if reflect.ValueOf(project).IsZero() {
		return fmt.Errorf("project %v does not exist", h.input.Project)
	}
	if len(h.input.Environment) == 0 {
		return fmt.Errorf("environment is mandatory")
	}
	environment, err := h.sumanProxy.ContentManagementLookupEnvironment(authParm, h.input.Project, h.input.Environment)
	if err != nil {
		return err
	}
	if reflect.ValueOf(environment).IsZero() {
		return fmt.Errorf("project %v environment %v does not exist", h.input.Project, h.input.Environment)
	}
	if !reflect.ValueOf(environment.PreviousEnvironmentLabel).IsZero() {
		previousEnvironment, err := h.sumanProxy.ContentManagementLookupEnvironment(authParm, h.input.Project, environment.PreviousEnvironmentLabel)
		if err != nil {
			return err
		}
		if previousEnvironment.Status == "unknown" {
			return fmt.Errorf("previous environment %v has never been build. Please build first", environment.PreviousEnvironmentLabel)
		}
		if previousEnvironment.Status == "building" || previousEnvironment.Status == "generating_repodata" {
			return fmt.Errorf("previous environment %v still being build. Please wait untill finished", environment.PreviousEnvironmentLabel)
		}
	}
	log.Debug("syncStage validateSyncStage finished")
	return nil
}

func (h *SyncStage) doSyncStage(authParm _sumanUseCase.AuthParams) error {
	log.Debug("doSyncStage started")
	environment, err := h.sumanProxy.ContentManagementLookupEnvironment(authParm, h.input.Project, h.input.Environment)
	if err != nil {
		return err
	}
	if reflect.ValueOf(environment.PreviousEnvironmentLabel).IsZero() {
		_, err := h.sumanProxy.ContentManagementBuildProject(authParm, h.input.Project)
		if err != nil {
			return err
		}
	} else {
		_, err := h.sumanProxy.ContentManagementPromoteProject(authParm, h.input.Project, environment.PreviousEnvironmentLabel)
		if err != nil {
			return err
		}
	}
	log.Debug("doCreateSoftwareProject finished")
	return nil
}

func (h *SyncStage) waitUntilFinished(authParm _sumanUseCase.AuthParams) error {
	log.Debug("waitUntilFinished started")
	for {
		environment, err := h.sumanProxy.ContentManagementLookupEnvironment(authParm, h.input.Project, h.input.Environment)
		if err != nil {
			return err
		}
		if environment.Status == "built" {
			break
		}
		log.Info(fmt.Sprintf("waiting for environment %v to be built", environment.Label))
		time.Sleep(time.Second * 30)
	}
	log.Debug("waitUntilFinished finished")
	return nil
}
