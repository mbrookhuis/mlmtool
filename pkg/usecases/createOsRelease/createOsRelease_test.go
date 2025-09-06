package createosrelease

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	_osReleaseModels "mlmtool/pkg/models/createOsRelease"
	sumamodels "mlmtool/pkg/models/susemanager"
	_sumanUseCase "mlmtool/pkg/usecases/susemanager"
	susemocks "mlmtool/pkg/usecases/susemanager/mocks"
	"mlmtool/pkg/util/cmdexecutor"
	cmdExecutor "mlmtool/pkg/util/cmdexecutor/mocks"
	logging "mlmtool/pkg/util/logger"
	util "mlmtool/pkg/util/rest"
)

//nolint:funlen
func TestCreateOsRelease_validateGivenOsRelease(t *testing.T) {
	type fields struct {
		sumanProxy           _sumanUseCase.IProxy
		suse                 _sumanUseCase.ISuseManager
		suseoperationtimeout int
		logger               *zap.Logger
		osRelease            string
	}
	type args struct {
		auth *_sumanUseCase.AuthParams
	}

	listProjects := []sumamodels.ContentManagementListProjects{
		{
			Label: "mi52-230317-r001",
		},
	}
	listChannels := []sumamodels.ChannelListSoftwareChannels{
		{
			Label: "mi52-230317-r001",
		},
	}
	kickstartTreeDetails := sumamodels.KickstartTreeGetDetails{
		Label: "mi52-230318-r001",
	}

	listProjectsErr := []sumamodels.ContentManagementListProjects{
		{
			Label: "mi52-230318-r001",
		},
	}
	listChannelsErr := []sumamodels.ChannelListSoftwareChannels{
		{
			Label: "mi52-230318-r001",
		},
	}
	kickstartTreeDetailsErr := sumamodels.KickstartTreeGetDetails{
		Label: "mi52-230317-r001",
	}
	var channelPresent bool
	logger := logging.NewTestingLogger(t.Name())
	suseManagerMock := new(susemocks.ISuseManager)
	suseManagerMock.On("SuseManagerCall", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&util.HTTPHelperStruct{}, nil)
	suseManagerMockErr := new(susemocks.ISuseManager)
	suseManagerMockErr.On("SuseManagerCall", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&util.HTTPHelperStruct{}, errors.New("failed to logout suse manager"))
	suseAPIProxy := new(susemocks.IProxy)
	suseAPIProxy.On("ContentManagementListProjects", mock.Anything, mock.Anything).Return(listProjects, nil)
	suseAPIProxy.On("ChannelListSoftwareChannels", mock.Anything, mock.Anything).Return(listChannels, nil)
	suseAPIProxy.On("KickstartTreeGetDetails", mock.Anything, mock.Anything, mock.Anything).Return(kickstartTreeDetails, errors.New("dist exists"))
	suseAPIProxy.On("ChannelSoftwareIsExisting", mock.Anything, mock.Anything, mock.Anything).Return(channelPresent, nil)
	suseAPIProxy.On("ChannelSoftwareCreate", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(1, nil)
	suseAPIProxyErr := new(susemocks.IProxy)
	suseAPIProxyErr.On("ContentManagementListProjects", mock.Anything, mock.Anything).Return(listProjectsErr, errors.New("failed to get list"))
	suseAPIProxyErr.On("ChannelListSoftwareChannels", mock.Anything, mock.Anything).Return(listChannelsErr, nil)
	suseAPIProxyErr.On("KickstartTreeGetDetails", mock.Anything, mock.Anything, mock.Anything).Return(kickstartTreeDetailsErr, nil)
	suseAPIProxyErr.On("ChannelSoftwareCreate", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(0, errors.New("Failed to create repo"))
	suseAPIProxyErr.On("ChannelSoftwareIsExisting", mock.Anything, mock.Anything, mock.Anything).Return(channelPresent, errors.New("dist exists"))

	authParam := _sumanUseCase.AuthParams{
		SessionKey: "test key",
		Host:       "test Hostname",
	}

	positiveResult := fields{
		sumanProxy:           suseAPIProxy,
		suse:                 suseManagerMock,
		suseoperationtimeout: 30,
		logger:               logger,
		osRelease:            "mi52-230318-r001",
	}

	negativeResult := fields{
		sumanProxy:           suseAPIProxyErr,
		suse:                 suseManagerMock,
		suseoperationtimeout: 30,
		logger:               logger,
		osRelease:            "mi52-www-230318-r001",
	}

	var tests = []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "osRelease valid",
			fields: positiveResult,
			args: args{
				auth: &authParam,
			},
			wantErr: false,
		},
		{
			name:   "osRelease invalid",
			fields: negativeResult,
			args: args{
				auth: &authParam,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &CreateOsRelease{
				sumanProxy:           tt.fields.sumanProxy,
				suse:                 tt.fields.suse,
				suseOperationTimeout: tt.fields.suseoperationtimeout,
				logger:               tt.fields.logger,
				osRelease:            tt.fields.osRelease,
			}
			if err := h.validateGivenOsRelease(tt.args.auth); (err != nil) != tt.wantErr {
				t.Errorf("validateGivenOsRelease() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

//nolint:funlen
func TestCreateOsRelease_CheckOsReleaseCMProjectExists(t *testing.T) {
	type fields struct {
		sumanProxy           _sumanUseCase.IProxy
		suse                 _sumanUseCase.ISuseManager
		suseoperationtimeout int
		logger               *zap.Logger
		osRelease            string
	}
	type args struct {
		auth *_sumanUseCase.AuthParams
	}

	listProjectsErr := []sumamodels.ContentManagementListProjects{
		{
			Label: "mi52-230317-r001",
		},
	}
	listProjects := []sumamodels.ContentManagementListProjects{
		{
			Label: "mi52-230318-r001",
		},
	}
	logger := logging.NewTestingLogger(t.Name())
	suseManagerMock := new(susemocks.ISuseManager)
	suseManagerMock.On("SuseManagerCall", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&util.HTTPHelperStruct{}, nil)
	suseManagerMockErr := new(susemocks.ISuseManager)
	suseManagerMockErr.On("SuseManagerCall", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&util.HTTPHelperStruct{}, errors.New("failed to logout suse manager"))
	suseAPIProxy := new(susemocks.IProxy)
	suseAPIProxy.On("ContentManagementListProjects", mock.Anything, mock.Anything).Return(listProjectsErr, nil)
	suseAPIProxy2 := new(susemocks.IProxy)
	suseAPIProxy2.On("ContentManagementListProjects", mock.Anything, mock.Anything).Return(listProjects, errors.New("CM present"))
	suseAPIProxyErr := new(susemocks.IProxy)
	suseAPIProxyErr.On("ContentManagementListProjects", mock.Anything, mock.Anything).Return(nil, errors.New("failed to get list"))

	authParam := _sumanUseCase.AuthParams{
		SessionKey: "test key",
		Host:       "test Hostname",
	}

	positiveSumaResult := fields{
		sumanProxy:           suseAPIProxy,
		suse:                 suseManagerMock,
		suseoperationtimeout: 30,
		logger:               logger,
		osRelease:            "mi52-230318-r001",
	}

	negativeSumaResult := fields{
		sumanProxy:           suseAPIProxyErr,
		suse:                 suseManagerMock,
		suseoperationtimeout: 30,
		logger:               logger,
		osRelease:            "xxxx-oh-no",
	}

	positiveFindResult := fields{
		sumanProxy:           suseAPIProxy2,
		suse:                 suseManagerMock,
		suseoperationtimeout: 30,
		logger:               logger,
		osRelease:            "mi52-230318-r001",
	}

	negativeFindResult := fields{
		sumanProxy:           suseAPIProxy,
		suse:                 suseManagerMock,
		suseoperationtimeout: 30,
		logger:               logger,
		osRelease:            "mi52-230317-r001",
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "valid",
			fields: positiveSumaResult,
			args: args{
				auth: &authParam,
			},
			wantErr: false,
		},
		{
			name:   "invalid",
			fields: negativeSumaResult,
			args: args{
				auth: &authParam,
			},
			wantErr: true,
		},
		{
			name:   "present",
			fields: positiveFindResult,
			args: args{
				auth: &authParam,
			},
			wantErr: true,
		},
		{
			name:   "not present",
			fields: negativeFindResult,
			args: args{
				auth: &authParam,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &CreateOsRelease{
				sumanProxy:           tt.fields.sumanProxy,
				suse:                 tt.fields.suse,
				suseOperationTimeout: tt.fields.suseoperationtimeout,
				logger:               tt.fields.logger,
				osRelease:            tt.fields.osRelease,
			}
			err := h.CheckOsReleaseCMProjectExists(tt.args.auth)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckOsReleaseCMProjectExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestCreateOsRelease_CheckOsReleaseFormat(t *testing.T) {
	type fields struct {
		sumanProxy           _sumanUseCase.IProxy
		suse                 _sumanUseCase.ISuseManager
		suseoperationtimeout int
		logger               *zap.Logger
		osRelease            string
	}

	logger := logging.NewTestingLogger(t.Name())
	suseManagerMock := new(susemocks.ISuseManager)
	suseAPIProxy := new(susemocks.IProxy)

	positiveResult := fields{
		sumanProxy:           suseAPIProxy,
		suse:                 suseManagerMock,
		suseoperationtimeout: 30,
		logger:               logger,
		osRelease:            "mi52-230318-r001",
	}
	negativeResult := fields{
		sumanProxy:           suseAPIProxy,
		suse:                 suseManagerMock,
		suseoperationtimeout: 30,
		logger:               logger,
		osRelease:            "oh-no-wrong",
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "valid",
			fields:  positiveResult,
			wantErr: false,
		},
		{
			name:    "invalid",
			fields:  negativeResult,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &CreateOsRelease{
				sumanProxy:           tt.fields.sumanProxy,
				suse:                 tt.fields.suse,
				suseOperationTimeout: tt.fields.suseoperationtimeout,
				logger:               tt.fields.logger,
				osRelease:            tt.fields.osRelease,
			}
			if err := h.CheckOsReleaseFormat(); (err != nil) != tt.wantErr {
				t.Errorf("CheckOsReleaseFormat() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

//nolint:funlen
func TestCreateOsRelease_CheckOsReleaseDistroExists(t *testing.T) {
	type fields struct {
		sumanProxy           _sumanUseCase.IProxy
		suse                 _sumanUseCase.ISuseManager
		suseoperationtimeout int
		logger               *zap.Logger
		osRelease            string
	}
	type args struct {
		auth *_sumanUseCase.AuthParams
	}

	kickstartTreeDetails := sumamodels.KickstartTreeGetDetails{
		Label: "mi52-230318-r001",
	}
	kickstartTreeDetailsErr := sumamodels.KickstartTreeGetDetails{
		Label: "mi52-230317-r001",
	}

	logger := logging.NewTestingLogger(t.Name())
	suseManagerMock := new(susemocks.ISuseManager)
	suseAPIProxy := new(susemocks.IProxy)
	suseAPIProxy.On("KickstartTreeGetDetails", mock.Anything, mock.Anything, mock.Anything).Return(kickstartTreeDetails, errors.New("dist exists"))
	suseAPIProxyErr := new(susemocks.IProxy)
	suseAPIProxyErr.On("KickstartTreeGetDetails", mock.Anything, mock.Anything, mock.Anything).Return(kickstartTreeDetailsErr, nil)

	authParam := _sumanUseCase.AuthParams{
		SessionKey: "test key",
		Host:       "test Hostname",
	}

	positiveResult := fields{
		sumanProxy:           suseAPIProxy,
		suse:                 suseManagerMock,
		suseoperationtimeout: 30,
		logger:               logger,
		osRelease:            "mi52-230318-r001",
	}

	negativeResult := fields{
		sumanProxy:           suseAPIProxyErr,
		suse:                 suseManagerMock,
		suseoperationtimeout: 30,
		logger:               logger,
		osRelease:            "mi52-230318-r001",
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "osRelease valid",
			fields: positiveResult,
			args: args{
				auth: &authParam,
			},
			wantErr: false,
		},
		{
			name:   "osRelease invalid",
			fields: negativeResult,
			args: args{
				auth: &authParam,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &CreateOsRelease{
				sumanProxy:           tt.fields.sumanProxy,
				suse:                 tt.fields.suse,
				suseOperationTimeout: tt.fields.suseoperationtimeout,
				logger:               tt.fields.logger,
				osRelease:            tt.fields.osRelease,
			}
			if err := h.CheckOsReleaseDistroExists(tt.args.auth); (err != nil) != tt.wantErr {
				t.Errorf("CheckOsReleaseDistroExists() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateOsRelease_CheckOsReleaseLabel(t *testing.T) {
	type fields struct {
		sumanProxy           _sumanUseCase.IProxy
		suse                 _sumanUseCase.ISuseManager
		suseoperationtimeout int
		logger               *zap.Logger
		osRelease            string
	}
	logger := logging.NewTestingLogger(t.Name())
	suseManagerMock := new(susemocks.ISuseManager)
	suseAPIProxy := new(susemocks.IProxy)

	positiveResult := fields{
		sumanProxy:           suseAPIProxy,
		suse:                 suseManagerMock,
		suseoperationtimeout: 30,
		logger:               logger,
		osRelease:            "mi52-230318-r001",
	}
	negativeResult := fields{
		sumanProxy:           suseAPIProxy,
		suse:                 suseManagerMock,
		suseoperationtimeout: 30,
		logger:               logger,
		osRelease:            "oh-no-wrong",
	}

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "valid",
			fields:  positiveResult,
			wantErr: false,
		},
		{
			name:    "invalid",
			fields:  negativeResult,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &CreateOsRelease{
				sumanProxy:           tt.fields.sumanProxy,
				suse:                 tt.fields.suse,
				suseOperationTimeout: tt.fields.suseoperationtimeout,
				logger:               tt.fields.logger,
				osRelease:            tt.fields.osRelease,
			}
			if err := h.CheckOsReleaseLabel(); (err != nil) != tt.wantErr {
				t.Errorf("CheckOsReleaseLabel() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateOsRelease_CheckOsReleaseDate(t *testing.T) {
	type fields struct {
		sumanProxy           _sumanUseCase.IProxy
		suse                 _sumanUseCase.ISuseManager
		suseoperationtimeout int
		logger               *zap.Logger
		osRelease            string
	}
	logger := logging.NewTestingLogger(t.Name())
	suseManagerMock := new(susemocks.ISuseManager)
	suseAPIProxy := new(susemocks.IProxy)

	positiveResult := fields{
		sumanProxy:           suseAPIProxy,
		suse:                 suseManagerMock,
		suseoperationtimeout: 30,
		logger:               logger,
		osRelease:            "mi52-230318-r001",
	}
	negativeResult := fields{
		sumanProxy:           suseAPIProxy,
		suse:                 suseManagerMock,
		suseoperationtimeout: 30,
		logger:               logger,
		osRelease:            "oh-no-wrong",
	}

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "valid",
			fields:  positiveResult,
			wantErr: false,
		},
		{
			name:    "invalid",
			fields:  negativeResult,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &CreateOsRelease{
				sumanProxy:           tt.fields.sumanProxy,
				suse:                 tt.fields.suse,
				suseOperationTimeout: tt.fields.suseoperationtimeout,
				logger:               tt.fields.logger,
				osRelease:            tt.fields.osRelease,
			}
			if err := h.CheckOsReleaseDate(); (err != nil) != tt.wantErr {
				t.Errorf("CheckOsReleaseDate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateOsRelease_CheckOsReleaseEnv(t *testing.T) {
	type fields struct {
		sumanProxy           _sumanUseCase.IProxy
		suse                 _sumanUseCase.ISuseManager
		suseoperationtimeout int
		logger               *zap.Logger
		osRelease            string
	}
	logger := logging.NewTestingLogger(t.Name())
	suseManagerMock := new(susemocks.ISuseManager)
	suseAPIProxy := new(susemocks.IProxy)

	positiveResult := fields{
		sumanProxy:           suseAPIProxy,
		suse:                 suseManagerMock,
		suseoperationtimeout: 30,
		logger:               logger,
		osRelease:            "mi52-230318-r001",
	}
	negativeResult := fields{
		sumanProxy:           suseAPIProxy,
		suse:                 suseManagerMock,
		suseoperationtimeout: 30,
		logger:               logger,
		osRelease:            "mi52-230318-r002",
	}

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "valid",
			fields:  positiveResult,
			wantErr: false,
		},
		{
			name:    "invalid",
			fields:  negativeResult,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &CreateOsRelease{
				sumanProxy:           tt.fields.sumanProxy,
				suse:                 tt.fields.suse,
				suseOperationTimeout: tt.fields.suseoperationtimeout,
				logger:               tt.fields.logger,
				osRelease:            tt.fields.osRelease,
			}
			if err := h.CheckOsReleaseEnv(); (err != nil) != tt.wantErr {
				t.Errorf("CheckOsReleaseEnv() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateOsRelease_createOsRelease(t *testing.T) {
	type fields struct {
		sumanProxy           _sumanUseCase.IProxy
		suse                 _sumanUseCase.ISuseManager
		cmdExec              cmdexecutor.ICMDExecutor
		suseoperationtimeout int
		logger               *zap.Logger
		osRelease            string
	}
	type args struct {
		auth *_sumanUseCase.AuthParams
	}
	contentManagementDetails := sumamodels.ContentManagementListProjects{
		Label: "mi52-230318-r001",
	}
	childChannels := []sumamodels.ChannelSoftwareListChildren{
		{
			Label: "child 1",
		},
		{
			Label: "child 2",
		},
	}
	childChannel := sumamodels.ChannelSoftwareListChildren{
		Label: "child 1",
	}
	contentManagementSource := sumamodels.ContentManagementSource{
		ContentProjectLabel: "label",
	}
	contentManagementListFilters := []sumamodels.ContentManagementFilter{
		{
			Name: "filter 1",
		},
		{
			Name: "filter 2",
		},
	}
	contentManagementCreateFilter := sumamodels.ContentManagementFilter{
		Name: "filter",
	}
	contentManagementEnv := sumamodels.ContentManagementEnvironment{
		Name: "Environment",
	}
	createRepo := sumamodels.ChannelSoftwareCreateRepo{
		Label: "Customchannel",
	}

	logger := logging.NewTestingLogger(t.Name())
	suseManagerMock := new(susemocks.ISuseManager)
	suseManagerMock.On("SuseManagerCall", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&util.HTTPHelperStruct{}, nil)
	suseManagerMockErr := new(susemocks.ISuseManager)
	suseManagerMockErr.On("SuseManagerCall", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&util.HTTPHelperStruct{}, errors.New("failed to logout suse manager"))
	suseAPIProxy := new(susemocks.IProxy)
	suseAPIProxy.On("ContentManagementCreate", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(contentManagementDetails, nil)
	suseAPIProxy.On("ChannelSoftwareListChildren", mock.Anything, mock.Anything, mock.Anything).Return(childChannels, nil)
	suseAPIProxy.On("ContentManagementAttachSource", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(contentManagementSource, nil)
	suseAPIProxy.On("ContentManagementListFilters", mock.Anything, mock.Anything).Return(contentManagementListFilters, nil)
	suseAPIProxy.On("ContentManagementCreateFilter", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(contentManagementCreateFilter, nil)
	suseAPIProxy.On("ContentManagementAttachFilter", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(contentManagementCreateFilter, nil)
	suseAPIProxy.On("ContentManagementCreateEnvironment", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(contentManagementEnv, nil)
	suseAPIProxy.On("ContentManagementBuildProject", mock.Anything, mock.Anything, mock.Anything).Return(1, nil)
	suseAPIProxy.On("KickstartTreeCreate", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(1, nil)
	suseAPIProxy.On("KickstartTreeCreateKernelOptions", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(0, errors.New("failed"))
	suseAPIProxy.On("ChannelSoftwareCreateRepo", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(createRepo, nil)
	suseAPIProxy.On("ChannelSoftwareCreate", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(1, nil)
	suseAPIProxy.On("ChannelSoftwareAssociateRepo", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(childChannel, nil)
	suseAPIProxy.On("ChannelSoftwareSyncRepo", mock.Anything, mock.Anything, mock.Anything).Return(1, nil)
	suseAPIProxyErr := new(susemocks.IProxy)
	suseAPIProxyErr.On("ContentManagementCreate", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(contentManagementDetails, errors.New("failed"))
	suseAPIProxyErr.On("ChannelSoftwareListChildren", mock.Anything, mock.Anything, mock.Anything).Return(childChannels, errors.New("failed"))
	suseAPIProxyErr.On("ContentManagementAttachSource", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(contentManagementSource, errors.New("failed"))
	suseAPIProxyErr.On("ContentManagementListFilters", mock.Anything, mock.Anything).Return(contentManagementListFilters, errors.New("failed"))
	suseAPIProxyErr.On("ContentManagementCreateFilter", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(contentManagementCreateFilter, errors.New("failed"))
	suseAPIProxyErr.On("ContentManagementAttachFilter", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(contentManagementCreateFilter, errors.New("failed"))
	suseAPIProxyErr.On("ContentManagementCreateEnvironment", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(contentManagementEnv, errors.New("failed"))
	suseAPIProxyErr.On("ContentManagementBuildProject", mock.Anything, mock.Anything, mock.Anything).Return(0, errors.New("failed"))
	suseAPIProxyErr.On("KickstartTreeCreate", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(0, errors.New("failed"))
	suseAPIProxyErr.On("KickstartTreeCreateKernelOptions", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(0, errors.New("failed"))
	suseAPIProxyErr.On("ChannelSoftwareCreateRepo", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(createRepo, errors.New("failed"))
	suseAPIProxyErr.On("ChannelSoftwareAssociateRepo", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(childChannel, errors.New("failed"))
	suseAPIProxyErr.On("ChannelSoftwareCreate", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(0, errors.New("failed"))
	suseAPIProxyErr.On("ChannelListSyncRepo", mock.Anything, mock.Anything).Return(0, errors.New("failed"))
	var resultRun []string
	cmdExec := new(cmdExecutor.ICMDExecutor)
	cmdExec.On("ExecuteCommand", mock.Anything, mock.Anything, mock.Anything).Return(resultRun, nil)
	cmdExec.On("CreateDirectory", mock.Anything, mock.Anything).Return(nil)
	cmdExecErr := new(cmdExecutor.ICMDExecutor)
	cmdExecErr.On("ExecuteCommand", mock.Anything, mock.Anything, mock.Anything).Return(resultRun, errors.New("error running script"))
	cmdExecErr.On("CreateDirectory", mock.Anything, mock.Anything).Return(errors.New("error creating directory"))

	authParam := _sumanUseCase.AuthParams{
		SessionKey: "test key",
		Host:       "test Hostname",
	}

	positiveResult := fields{
		sumanProxy:           suseAPIProxy,
		suse:                 suseManagerMock,
		cmdExec:              cmdExec,
		suseoperationtimeout: 30,
		logger:               logger,
		osRelease:            "mi52-230318-r001",
	}

	negativeResult := fields{
		sumanProxy:           suseAPIProxyErr,
		suse:                 suseManagerMock,
		cmdExec:              cmdExec,
		suseoperationtimeout: 30,
		logger:               logger,
		osRelease:            "mi52-230318-r001",
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "create success",
			fields: positiveResult,
			args: args{
				auth: &authParam,
			},
			wantErr: false,
		},
		{
			name:   "create failed",
			fields: negativeResult,
			args: args{
				auth: &authParam,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &CreateOsRelease{
				sumanProxy:           tt.fields.sumanProxy,
				suse:                 tt.fields.suse,
				cmdExec:              tt.fields.cmdExec,
				suseOperationTimeout: tt.fields.suseoperationtimeout,
				logger:               tt.fields.logger,
				osRelease:            tt.fields.osRelease,
			}
			if err := h.createOsRelease(tt.args.auth); (err != nil) != tt.wantErr {
				t.Errorf("createOsRelease() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateOsRelease_getDataProjectOsRelease(t *testing.T) {
	type fields struct {
		sumanProxy           _sumanUseCase.IProxy
		suse                 _sumanUseCase.ISuseManager
		suseoperationtimeout int
		logger               *zap.Logger
		osRelease            string
	}
	logger := logging.NewTestingLogger(t.Name())
	suseManagerMock := new(susemocks.ISuseManager)
	suseAPIProxy := new(susemocks.IProxy)
	suseAPIProxyErr := new(susemocks.IProxy)

	positiveResult := fields{
		sumanProxy:           suseAPIProxy,
		suse:                 suseManagerMock,
		suseoperationtimeout: 30,
		logger:               logger,
		osRelease:            "mi52-230318-r001",
	}

	negativeResult := fields{
		sumanProxy:           suseAPIProxyErr,
		suse:                 suseManagerMock,
		suseoperationtimeout: 30,
		logger:               logger,
		osRelease:            "mi58-230318-r001",
	}

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "success",
			fields:  positiveResult,
			wantErr: false,
		},
		{
			name:    "failed",
			fields:  negativeResult,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &CreateOsRelease{
				sumanProxy:           tt.fields.sumanProxy,
				suse:                 tt.fields.suse,
				suseOperationTimeout: tt.fields.suseoperationtimeout,
				logger:               tt.fields.logger,
				osRelease:            tt.fields.osRelease,
			}
			_, err := h.getDataProjectOsRelease()
			if (err != nil) != tt.wantErr {
				t.Errorf("getDataProjectOsRelease() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestCreateOsRelease_createProjectOsRelease(t *testing.T) {
	type fields struct {
		sumanProxy           _sumanUseCase.IProxy
		suse                 _sumanUseCase.ISuseManager
		suseoperationtimeout int
		logger               *zap.Logger
		osRelease            string
	}
	type args struct {
		auth *_sumanUseCase.AuthParams
	}
	contentManagementDetails := sumamodels.ContentManagementListProjects{
		Label: "mi52-230318-r001",
	}
	logger := logging.NewTestingLogger(t.Name())
	suseManagerMock := new(susemocks.ISuseManager)
	suseManagerMock.On("SuseManagerCall", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&util.HTTPHelperStruct{}, nil)
	suseManagerMockErr := new(susemocks.ISuseManager)
	suseManagerMockErr.On("SuseManagerCall", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&util.HTTPHelperStruct{}, errors.New("failed to logout suse manager"))
	suseAPIProxy := new(susemocks.IProxy)
	suseAPIProxy.On("ContentManagementCreate", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(contentManagementDetails, nil)
	suseAPIProxyErr := new(susemocks.IProxy)
	suseAPIProxyErr.On("ContentManagementCreate", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(contentManagementDetails, errors.New("failed"))

	authParam := _sumanUseCase.AuthParams{
		SessionKey: "test key",
		Host:       "test Hostname",
	}

	positiveResult := fields{
		sumanProxy:           suseAPIProxy,
		suse:                 suseManagerMock,
		suseoperationtimeout: 30,
		logger:               logger,
		osRelease:            "mi52-230318-r001",
	}

	negativeResult := fields{
		sumanProxy:           suseAPIProxyErr,
		suse:                 suseManagerMock,
		suseoperationtimeout: 30,
		logger:               logger,
		osRelease:            "mi52-230318-r001",
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "create success",
			fields: positiveResult,
			args: args{
				auth: &authParam,
			},
			wantErr: false,
		},
		{
			name:   "create failed",
			fields: negativeResult,
			args: args{
				auth: &authParam,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &CreateOsRelease{
				sumanProxy:           tt.fields.sumanProxy,
				suse:                 tt.fields.suse,
				suseOperationTimeout: tt.fields.suseoperationtimeout,
				logger:               tt.fields.logger,
				osRelease:            tt.fields.osRelease,
			}
			if err := h.createProjectOsRelease(tt.args.auth); (err != nil) != tt.wantErr {
				t.Errorf("createProjectOsRelease() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateOsRelease_addChannelsProjectOsRelease(t *testing.T) {
	type fields struct {
		sumanProxy           _sumanUseCase.IProxy
		suse                 _sumanUseCase.ISuseManager
		suseoperationtimeout int
		logger               *zap.Logger
		osRelease            string
	}
	type args struct {
		auth          *_sumanUseCase.AuthParams
		osReleaseData _osReleaseModels.OsReleaseRecord
	}
	childChannels := []sumamodels.ChannelSoftwareListChildren{
		{
			Label: "child 1",
		},
		{
			Label: "child 2",
		},
	}
	contentManagementSource := sumamodels.ContentManagementSource{
		ContentProjectLabel: "label",
	}

	logger := logging.NewTestingLogger(t.Name())
	suseManagerMock := new(susemocks.ISuseManager)
	suseManagerMock.On("SuseManagerCall", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&util.HTTPHelperStruct{}, nil)
	suseManagerMockErr := new(susemocks.ISuseManager)
	suseManagerMockErr.On("SuseManagerCall", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&util.HTTPHelperStruct{}, errors.New("failed to logout suse manager"))
	suseAPIProxy := new(susemocks.IProxy)
	suseAPIProxy.On("ChannelSoftwareListChildren", mock.Anything, mock.Anything, mock.Anything).Return(childChannels, nil)
	suseAPIProxy.On("ContentManagementAttachSource", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(contentManagementSource, nil)
	suseAPIProxyErr1 := new(susemocks.IProxy)
	suseAPIProxyErr1.On("ChannelSoftwareListChildren", mock.Anything, mock.Anything, mock.Anything).Return(childChannels, errors.New("failed"))
	suseAPIProxyErr2 := new(susemocks.IProxy)
	suseAPIProxyErr2.On("ChannelSoftwareListChildren", mock.Anything, mock.Anything, mock.Anything).Return(childChannels, errors.New("failed"))
	suseAPIProxyErr2.On("ContentManagementAttachSource", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(contentManagementSource, errors.New("failed"))
	suseAPIProxyErr3 := new(susemocks.IProxy)
	suseAPIProxyErr3.On("ChannelSoftwareListChildren", mock.Anything, mock.Anything, mock.Anything).Return(childChannels, errors.New("failed"))
	suseAPIProxyErr3.On("ContentManagementAttachSource", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(contentManagementSource, errors.New("failed"))

	authParam := _sumanUseCase.AuthParams{
		SessionKey: "test key",
		Host:       "test Hostname",
	}
	osRelData := _osReleaseModels.OsReleaseRecord{
		ParentChannel: "suse-microos-5.2-pool-x86_64",
	}

	positiveResult := fields{
		sumanProxy:           suseAPIProxy,
		suse:                 suseManagerMock,
		suseoperationtimeout: 30,
		logger:               logger,
		osRelease:            "mi52-230318-r001",
	}
	negativeResultList := fields{
		sumanProxy:           suseAPIProxyErr1,
		suse:                 suseManagerMock,
		suseoperationtimeout: 30,
		logger:               logger,
		osRelease:            "mi52-230318-r001",
	}
	negativeResultAddP := fields{
		sumanProxy:           suseAPIProxyErr2,
		suse:                 suseManagerMock,
		suseoperationtimeout: 30,
		logger:               logger,
		osRelease:            "mi52-230318-r001",
	}
	negativeResultAddC := fields{
		sumanProxy:           suseAPIProxyErr3,
		suse:                 suseManagerMock,
		suseoperationtimeout: 30,
		logger:               logger,
		osRelease:            "mi52-230318-r001",
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "success list children",
			fields: positiveResult,
			args: args{
				auth:          &authParam,
				osReleaseData: osRelData,
			},
			wantErr: false,
		},
		{
			name:   "failed list children",
			fields: negativeResultList,
			args: args{
				auth:          &authParam,
				osReleaseData: osRelData,
			},
			wantErr: true,
		},
		{
			name:   "failed attach parent",
			fields: negativeResultAddP,
			args: args{
				auth:          &authParam,
				osReleaseData: osRelData,
			},
			wantErr: true,
		},
		{
			name:   "failed attach children",
			fields: negativeResultAddC,
			args: args{
				auth:          &authParam,
				osReleaseData: osRelData,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &CreateOsRelease{
				sumanProxy:           tt.fields.sumanProxy,
				suse:                 tt.fields.suse,
				suseOperationTimeout: tt.fields.suseoperationtimeout,
				logger:               tt.fields.logger,
				osRelease:            tt.fields.osRelease,
			}
			if err := h.addChannelsProjectOsRelease(tt.args.auth, tt.args.osReleaseData); (err != nil) != tt.wantErr {
				t.Errorf("addChannelsProjectOsRelease() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateOsRelease_createFilterOsRelease(t *testing.T) {
	type fields struct {
		sumanProxy           _sumanUseCase.IProxy
		suse                 _sumanUseCase.ISuseManager
		suseoperationtimeout int
		logger               *zap.Logger
		osRelease            string
	}
	type args struct {
		auth *_sumanUseCase.AuthParams
	}
	contentManagementListFilters := []sumamodels.ContentManagementFilter{
		{
			Name: "filter 1",
		},
		{
			Name: "filter 2",
		},
	}
	contentManagementCreateFilter := sumamodels.ContentManagementFilter{
		Name: "filter",
	}

	logger := logging.NewTestingLogger(t.Name())
	suseManagerMock := new(susemocks.ISuseManager)
	suseManagerMock.On("SuseManagerCall", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&util.HTTPHelperStruct{}, nil)
	suseManagerMockErr := new(susemocks.ISuseManager)
	suseManagerMockErr.On("SuseManagerCall", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&util.HTTPHelperStruct{}, errors.New("failed to logout suse manager"))
	suseAPIProxy := new(susemocks.IProxy)
	suseAPIProxy.On("ContentManagementListFilters", mock.Anything, mock.Anything).Return(contentManagementListFilters, nil)
	suseAPIProxy.On("ContentManagementCreateFilter", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(contentManagementCreateFilter, nil)
	suseAPIProxy.On("ContentManagementAttachFilter", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(contentManagementCreateFilter, nil)
	suseAPIProxyErr1 := new(susemocks.IProxy)
	suseAPIProxyErr1.On("ContentManagementListFilters", mock.Anything, mock.Anything).Return(contentManagementListFilters, errors.New("failed"))
	suseAPIProxyErr2 := new(susemocks.IProxy)
	suseAPIProxyErr2.On("ContentManagementListFilters", mock.Anything, mock.Anything).Return(contentManagementListFilters, errors.New("failed"))
	suseAPIProxyErr2.On("ContentManagementCreateFilter", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(contentManagementCreateFilter, errors.New("failed"))
	suseAPIProxyErr3 := new(susemocks.IProxy)
	suseAPIProxyErr3.On("ContentManagementListFilters", mock.Anything, mock.Anything).Return(contentManagementListFilters, errors.New("failed"))
	suseAPIProxyErr3.On("ContentManagementCreateFilter", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(contentManagementCreateFilter, errors.New("failed"))
	suseAPIProxyErr3.On("ContentManagementAttachFilter", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(contentManagementCreateFilter, errors.New("failed"))

	authParam := _sumanUseCase.AuthParams{
		SessionKey: "test key",
		Host:       "test Hostname",
	}

	positiveResult := fields{
		sumanProxy:           suseAPIProxy,
		suse:                 suseManagerMock,
		suseoperationtimeout: 30,
		logger:               logger,
		osRelease:            "mi52-230318-r001",
	}

	negativeResult1 := fields{
		sumanProxy:           suseAPIProxyErr1,
		suse:                 suseManagerMock,
		suseoperationtimeout: 30,
		logger:               logger,
		osRelease:            "mi52-230318-r001",
	}
	negativeResult2 := fields{
		sumanProxy:           suseAPIProxyErr2,
		suse:                 suseManagerMock,
		suseoperationtimeout: 30,
		logger:               logger,
		osRelease:            "mi52-230318-r001",
	}
	negativeResult3 := fields{
		sumanProxy:           suseAPIProxyErr3,
		suse:                 suseManagerMock,
		suseoperationtimeout: 30,
		logger:               logger,
		osRelease:            "mi52-230318-r001",
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "success",
			fields: positiveResult,
			args: args{
				auth: &authParam,
			},
			wantErr: false,
		},
		{
			name:   "failed list",
			fields: negativeResult1,
			args: args{
				auth: &authParam,
			},
			wantErr: true,
		},
		{
			name:   "failed create",
			fields: negativeResult2,
			args: args{
				auth: &authParam,
			},
			wantErr: true,
		},
		{
			name:   "failed attach",
			fields: negativeResult3,
			args: args{
				auth: &authParam,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &CreateOsRelease{
				sumanProxy:           tt.fields.sumanProxy,
				suse:                 tt.fields.suse,
				suseOperationTimeout: tt.fields.suseoperationtimeout,
				logger:               tt.fields.logger,
				osRelease:            tt.fields.osRelease,
			}
			if err := h.createFilterOsRelease(tt.args.auth); (err != nil) != tt.wantErr {
				t.Errorf("createFilterOsRelease() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateOsRelease_createEnvironmentOsRelease(t *testing.T) {
	type fields struct {
		sumanProxy           _sumanUseCase.IProxy
		suse                 _sumanUseCase.ISuseManager
		suseoperationtimeout int
		logger               *zap.Logger
		osRelease            string
	}
	type args struct {
		auth *_sumanUseCase.AuthParams
	}
	contentManagementEnv := sumamodels.ContentManagementEnvironment{
		Name: "Environment",
	}

	logger := logging.NewTestingLogger(t.Name())
	suseManagerMock := new(susemocks.ISuseManager)
	suseManagerMock.On("SuseManagerCall", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&util.HTTPHelperStruct{}, nil)
	suseManagerMockErr := new(susemocks.ISuseManager)
	suseManagerMockErr.On("SuseManagerCall", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&util.HTTPHelperStruct{}, errors.New("failed to logout suse manager"))
	suseAPIProxy := new(susemocks.IProxy)
	suseAPIProxy.On("ContentManagementCreateEnvironment", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(contentManagementEnv, nil)
	suseAPIProxyErr := new(susemocks.IProxy)
	suseAPIProxyErr.On("ContentManagementCreateEnvironment", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(contentManagementEnv, errors.New("failed"))

	authParam := _sumanUseCase.AuthParams{
		SessionKey: "test key",
		Host:       "test Hostname",
	}

	positiveResult := fields{
		sumanProxy:           suseAPIProxy,
		suse:                 suseManagerMock,
		suseoperationtimeout: 30,
		logger:               logger,
		osRelease:            "mi52-230318-r001",
	}

	negativeResult := fields{
		sumanProxy:           suseAPIProxyErr,
		suse:                 suseManagerMock,
		suseoperationtimeout: 30,
		logger:               logger,
		osRelease:            "mi52-230318-r001",
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "create success",
			fields: positiveResult,
			args: args{
				auth: &authParam,
			},
			wantErr: false,
		},
		{
			name:   "create failed",
			fields: negativeResult,
			args: args{
				auth: &authParam,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &CreateOsRelease{
				sumanProxy:           tt.fields.sumanProxy,
				suse:                 tt.fields.suse,
				suseOperationTimeout: tt.fields.suseoperationtimeout,
				logger:               tt.fields.logger,
				osRelease:            tt.fields.osRelease,
			}
			if err := h.createEnvironmentOsRelease(tt.args.auth); (err != nil) != tt.wantErr {
				t.Errorf("createEnvironmentOsRelease() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateOsRelease_syncOsRelease(t *testing.T) {
	type fields struct {
		sumanProxy           _sumanUseCase.IProxy
		suse                 _sumanUseCase.ISuseManager
		suseoperationtimeout int
		logger               *zap.Logger
		osRelease            string
	}
	type args struct {
		auth *_sumanUseCase.AuthParams
	}

	logger := logging.NewTestingLogger(t.Name())
	suseManagerMock := new(susemocks.ISuseManager)
	suseManagerMock.On("SuseManagerCall", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&util.HTTPHelperStruct{}, nil)
	suseManagerMockErr := new(susemocks.ISuseManager)
	suseManagerMockErr.On("SuseManagerCall", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&util.HTTPHelperStruct{}, errors.New("failed to logout suse manager"))
	suseAPIProxy := new(susemocks.IProxy)
	suseAPIProxy.On("ContentManagementBuildProject", mock.Anything, mock.Anything, mock.Anything).Return(1, nil)
	suseAPIProxyErr := new(susemocks.IProxy)
	suseAPIProxyErr.On("ContentManagementBuildProject", mock.Anything, mock.Anything, mock.Anything).Return(0, errors.New("failed"))

	authParam := _sumanUseCase.AuthParams{
		SessionKey: "test key",
		Host:       "test Hostname",
	}

	positiveResult := fields{
		sumanProxy:           suseAPIProxy,
		suse:                 suseManagerMock,
		suseoperationtimeout: 30,
		logger:               logger,
		osRelease:            "mi52-230318-r001",
	}

	negativeResult := fields{
		sumanProxy:           suseAPIProxyErr,
		suse:                 suseManagerMock,
		suseoperationtimeout: 30,
		logger:               logger,
		osRelease:            "mi52-230318-r001",
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "create success",
			fields: positiveResult,
			args: args{
				auth: &authParam,
			},
			wantErr: false,
		},
		{
			name:   "create failed",
			fields: negativeResult,
			args: args{
				auth: &authParam,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &CreateOsRelease{
				sumanProxy:           tt.fields.sumanProxy,
				suse:                 tt.fields.suse,
				suseOperationTimeout: tt.fields.suseoperationtimeout,
				logger:               tt.fields.logger,
				osRelease:            tt.fields.osRelease,
			}
			if err := h.syncOsRelease(tt.args.auth); (err != nil) != tt.wantErr {
				t.Errorf("syncOsRelease() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateOsRelease_createDistributionOsRelease(t *testing.T) {
	type fields struct {
		sumanProxy           _sumanUseCase.IProxy
		suse                 _sumanUseCase.ISuseManager
		suseoperationtimeout int
		logger               *zap.Logger
		osRelease            string
	}
	type args struct {
		auth          *_sumanUseCase.AuthParams
		osReleaseData _osReleaseModels.OsReleaseRecord
	}

	logger := logging.NewTestingLogger(t.Name())
	suseManagerMock := new(susemocks.ISuseManager)
	suseManagerMock.On("SuseManagerCall", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&util.HTTPHelperStruct{}, nil)
	suseManagerMockErr := new(susemocks.ISuseManager)
	suseManagerMockErr.On("SuseManagerCall", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&util.HTTPHelperStruct{}, errors.New("failed to logout suse manager"))
	suseAPIProxy := new(susemocks.IProxy)
	suseAPIProxy.On("KickstartTreeCreate", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(1, nil)
	suseAPIProxy.On("KickstartTreeCreateKernelOptions", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(1, nil)
	suseAPIProxyErr := new(susemocks.IProxy)
	suseAPIProxyErr.On("KickstartTreeCreate", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(0, errors.New("failed"))
	suseAPIProxyErr.On("KickstartTreeCreateKernelOptions", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(0, errors.New("failed"))

	authParam := _sumanUseCase.AuthParams{
		SessionKey: "test key",
		Host:       "test Hostname",
	}

	positiveResult := fields{
		sumanProxy:           suseAPIProxy,
		suse:                 suseManagerMock,
		suseoperationtimeout: 30,
		logger:               logger,
		osRelease:            "mi52-230318-r001",
	}

	negativeResult := fields{
		sumanProxy:           suseAPIProxyErr,
		suse:                 suseManagerMock,
		suseoperationtimeout: 30,
		logger:               logger,
		osRelease:            "mi52-230318-r001",
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "create success",
			fields: positiveResult,
			args: args{
				auth: &authParam,
			},
			wantErr: false,
		},
		{
			name:   "create failed",
			fields: negativeResult,
			args: args{
				auth: &authParam,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &CreateOsRelease{
				sumanProxy:           tt.fields.sumanProxy,
				suse:                 tt.fields.suse,
				suseOperationTimeout: tt.fields.suseoperationtimeout,
				logger:               tt.fields.logger,
				osRelease:            tt.fields.osRelease,
			}
			if err := h.createDistributionOsRelease(tt.args.auth, tt.args.osReleaseData); (err != nil) != tt.wantErr {
				t.Errorf("createDistributionOsRelease() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateOsRelease_GetBaseChannels(t *testing.T) {
	type fields struct {
		sumanProxy           _sumanUseCase.IProxy
		suse                 _sumanUseCase.ISuseManager
		suseoperationtimeout int
		logger               *zap.Logger
		osRelease            string
	}
	type args struct {
		auth *_sumanUseCase.AuthParams
	}
	listChannels := []sumamodels.ChannelListSoftwareChannels{
		{
			Label: "mi52-230317-r001",
		},
	}
	listChannelsErr := []sumamodels.ChannelListSoftwareChannels{
		{
			Label: "mi52-230318-r001",
		},
	}

	logger := logging.NewTestingLogger(t.Name())
	suseManagerMock := new(susemocks.ISuseManager)
	suseManagerMock.On("SuseManagerCall", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&util.HTTPHelperStruct{}, nil)
	suseManagerMockErr := new(susemocks.ISuseManager)
	suseManagerMockErr.On("SuseManagerCall", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&util.HTTPHelperStruct{}, errors.New("failed to logout suse manager"))
	suseAPIProxy := new(susemocks.IProxy)
	suseAPIProxy.On("ChannelListSoftwareChannels", mock.Anything, mock.Anything).Return(listChannels, nil)
	suseAPIProxyErr := new(susemocks.IProxy)
	suseAPIProxyErr.On("ChannelListSoftwareChannels", mock.Anything, mock.Anything).Return(listChannelsErr, errors.New("failed"))

	authParam := _sumanUseCase.AuthParams{
		SessionKey: "test key",
		Host:       "test Hostname",
	}

	positiveResult := fields{
		sumanProxy:           suseAPIProxy,
		suse:                 suseManagerMock,
		suseoperationtimeout: 30,
		logger:               logger,
		osRelease:            "mi52-230318-r001",
	}

	negativeResult := fields{
		sumanProxy:           suseAPIProxyErr,
		suse:                 suseManagerMock,
		suseoperationtimeout: 30,
		logger:               logger,
		osRelease:            "mi52-230318-r001",
	}

	var tests = []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "succes",
			fields: positiveResult,
			args: args{
				auth: &authParam,
			},
			wantErr: false,
		},
		{
			name:   "failed",
			fields: negativeResult,
			args: args{
				auth: &authParam,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &CreateOsRelease{
				sumanProxy:           tt.fields.sumanProxy,
				suse:                 tt.fields.suse,
				suseOperationTimeout: tt.fields.suseoperationtimeout,
				logger:               tt.fields.logger,
				osRelease:            tt.fields.osRelease,
			}
			_, err := h.GetBaseChannels(tt.args.auth)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBaseChannels() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestCreateOsRelease_createExtraRepo(t *testing.T) {
	type fields struct {
		sumanProxy           _sumanUseCase.IProxy
		suse                 _sumanUseCase.ISuseManager
		cmdExec              cmdexecutor.ICMDExecutor
		suseOperationTimeout int
		logger               *zap.Logger
		osRelease            string
	}
	type args struct {
		auth *_sumanUseCase.AuthParams
	}
	logger := logging.NewTestingLogger(t.Name())
	var resultRun []string
	childChannel := sumamodels.ChannelSoftwareListChildren{
		Label: "child 1",
	}
	createRepo := sumamodels.ChannelSoftwareCreateRepo{
		Label: "Customchannel",
	}

	cmdExec1 := new(cmdExecutor.ICMDExecutor)
	cmdExec1.On("CreateDirectory", mock.Anything, mock.Anything).Return(nil)
	cmdExec1.On("ExecuteCommand", mock.Anything, mock.Anything, mock.Anything).Return(resultRun, nil)
	cmdExecErr1 := new(cmdExecutor.ICMDExecutor)
	cmdExecErr1.On("CreateDirectory", mock.Anything, mock.Anything).Return(nil)
	cmdExecErr1.On("ExecuteCommand", mock.Anything, mock.Anything, mock.Anything).Return(resultRun, errors.New("error running script"))
	suseManagerMock := new(susemocks.ISuseManager)
	suseManagerMock.On("SuseManagerCall", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&util.HTTPHelperStruct{}, nil)
	suseManagerMockErr := new(susemocks.ISuseManager)
	suseManagerMockErr.On("SuseManagerCall", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&util.HTTPHelperStruct{}, errors.New("failed to logout suse manager"))
	suseAPIProxy1 := new(susemocks.IProxy)
	suseAPIProxy1.On("ChannelSoftwareCreateRepo", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(createRepo, nil)
	suseAPIProxy1.On("ChannelSoftwareCreate", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(1, nil)
	suseAPIProxy1.On("ChannelSoftwareAssociateRepo", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(childChannel, nil)
	suseAPIProxy1.On("ChannelSoftwareSyncRepo", mock.Anything, mock.Anything, mock.Anything).Return(1, nil)
	suseAPIProxy1.On("ChannelSoftwareIsExisting", mock.Anything, mock.Anything, mock.Anything).Return(1, nil)
	suseAPIProxyErr1 := new(susemocks.IProxy)
	suseAPIProxyErr1.On("ChannelSoftwareCreateRepo", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(createRepo, errors.New("failed"))
	suseAPIProxyErr1.On("ChannelSoftwareCreate", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(0, errors.New("failed"))
	suseAPIProxyErr1.On("ChannelSoftwareAssociateRepo", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(childChannel, errors.New("failed"))
	suseAPIProxyErr1.On("ChannelListSyncRepo", mock.Anything, mock.Anything).Return(0, errors.New("failed"))
	suseAPIProxyErr1.On("ChannelSoftwareIsExisting", mock.Anything, mock.Anything, mock.Anything).Return(0, errors.New("failed"))
	suseAPIProxy2 := new(susemocks.IProxy)
	suseAPIProxy2.On("ChannelSoftwareCreateRepo", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(createRepo, nil)
	suseAPIProxy2.On("ChannelSoftwareCreate", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(1, nil)
	suseAPIProxy2.On("ChannelSoftwareAssociateRepo", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(childChannel, nil)
	suseAPIProxy2.On("ChannelSoftwareSyncRepo", mock.Anything, mock.Anything, mock.Anything).Return(1, nil)
	suseAPIProxy2.On("ChannelSoftwareIsExisting", mock.Anything, mock.Anything, mock.Anything).Return(1, nil)
	suseAPIProxyErr2 := new(susemocks.IProxy)
	suseAPIProxyErr2.On("ChannelSoftwareCreateRepo", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(createRepo, nil)
	suseAPIProxyErr2.On("ChannelSoftwareCreate", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(0, errors.New("failed"))
	suseAPIProxyErr2.On("ChannelSoftwareAssociateRepo", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(childChannel, errors.New("failed"))
	suseAPIProxyErr2.On("ChannelListSyncRepo", mock.Anything, mock.Anything).Return(0, errors.New("failed"))
	suseAPIProxyErr2.On("ChannelSoftwareIsExisting", mock.Anything, mock.Anything, mock.Anything).Return(0, errors.New("failed"))
	suseAPIProxy3 := new(susemocks.IProxy)
	suseAPIProxy3.On("ChannelSoftwareCreateRepo", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(createRepo, nil)
	suseAPIProxy3.On("ChannelSoftwareCreate", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(1, nil)
	suseAPIProxy3.On("ChannelSoftwareAssociateRepo", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(childChannel, nil)
	suseAPIProxy3.On("ChannelSoftwareSyncRepo", mock.Anything, mock.Anything, mock.Anything).Return(1, nil)
	suseAPIProxy3.On("ChannelSoftwareIsExisting", mock.Anything, mock.Anything, mock.Anything).Return(1, nil)
	suseAPIProxyErr3 := new(susemocks.IProxy)
	suseAPIProxyErr3.On("ChannelSoftwareCreateRepo", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(createRepo, nil)
	suseAPIProxyErr3.On("ChannelSoftwareCreate", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(1, nil)
	suseAPIProxyErr3.On("ChannelSoftwareAssociateRepo", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(childChannel, errors.New("failed"))
	suseAPIProxyErr3.On("ChannelListSyncRepo", mock.Anything, mock.Anything).Return(0, errors.New("failed"))
	suseAPIProxyErr3.On("ChannelSoftwareIsExisting", mock.Anything, mock.Anything, mock.Anything).Return(0, errors.New("failed"))
	suseAPIProxy4 := new(susemocks.IProxy)
	suseAPIProxy4.On("ChannelSoftwareCreateRepo", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(createRepo, nil)
	suseAPIProxy4.On("ChannelSoftwareCreate", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(1, nil)
	suseAPIProxy4.On("ChannelSoftwareAssociateRepo", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(childChannel, nil)
	suseAPIProxy4.On("ChannelSoftwareSyncRepo", mock.Anything, mock.Anything, mock.Anything).Return(1, nil)
	suseAPIProxy4.On("ChannelSoftwareIsExisting", mock.Anything, mock.Anything, mock.Anything).Return(1, nil)
	suseAPIProxyErr4 := new(susemocks.IProxy)
	suseAPIProxyErr4.On("ChannelSoftwareCreateRepo", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(createRepo, errors.New("failed"))
	suseAPIProxyErr4.On("ChannelSoftwareCreate", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(1, nil)
	suseAPIProxyErr4.On("ChannelSoftwareAssociateRepo", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(childChannel, errors.New("failed"))
	suseAPIProxyErr4.On("ChannelListSyncRepo", mock.Anything, mock.Anything).Return(0, errors.New("failed"))
	suseAPIProxyErr4.On("ChannelSoftwareIsExisting", mock.Anything, mock.Anything, mock.Anything).Return(0, errors.New("failed"))

	authParam := _sumanUseCase.AuthParams{
		SessionKey: "test key",
		Host:       "test Hostname",
	}
	positiveCreateDir := fields{
		sumanProxy:           suseAPIProxy1,
		suse:                 suseManagerMock,
		cmdExec:              cmdExec1,
		suseOperationTimeout: 30,
		logger:               logger,
		osRelease:            "mi52-230318-r001",
	}
	negativeCreateDir := fields{
		sumanProxy:           suseAPIProxy1,
		suse:                 suseManagerMock,
		cmdExec:              cmdExecErr1,
		suseOperationTimeout: 30,
		logger:               logger,
		osRelease:            "mi52-230318-r001",
	}
	positiveCreateRepo := fields{
		sumanProxy:           suseAPIProxy1,
		suse:                 suseManagerMock,
		cmdExec:              cmdExec1,
		suseOperationTimeout: 30,
		logger:               logger,
		osRelease:            "mi52-230318-r001",
	}
	negativeCreateRepo := fields{
		sumanProxy:           suseAPIProxyErr1,
		suse:                 suseManagerMock,
		cmdExec:              cmdExec1,
		suseOperationTimeout: 30,
		logger:               logger,
		osRelease:            "mi52-230318-r001",
	}
	positiveCreateChannel := fields{
		sumanProxy:           suseAPIProxy2,
		suse:                 suseManagerMock,
		cmdExec:              cmdExec1,
		suseOperationTimeout: 30,
		logger:               logger,
		osRelease:            "mi52-230318-r001",
	}
	negativeCreateChannel := fields{
		sumanProxy:           suseAPIProxy2,
		suse:                 suseManagerMock,
		cmdExec:              cmdExecErr1,
		suseOperationTimeout: 30,
		logger:               logger,
		osRelease:            "mi52-230318-r001",
	}
	positiveAssociateRepo := fields{
		sumanProxy:           suseAPIProxy3,
		suse:                 suseManagerMock,
		cmdExec:              cmdExec1,
		suseOperationTimeout: 30,
		logger:               logger,
		osRelease:            "mi52-230318-r001",
	}
	negativeAssociateRepo := fields{
		sumanProxy:           suseAPIProxyErr3,
		suse:                 suseManagerMock,
		cmdExec:              cmdExec1,
		suseOperationTimeout: 30,
		logger:               logger,
		osRelease:            "mi52-230318-r001",
	}
	positiveSyncRepo := fields{
		sumanProxy:           suseAPIProxy4,
		suse:                 suseManagerMock,
		cmdExec:              cmdExec1,
		suseOperationTimeout: 30,
		logger:               logger,
		osRelease:            "mi52-230318-r001",
	}
	negativeSyncRepo := fields{
		sumanProxy:           suseAPIProxyErr4,
		suse:                 suseManagerMock,
		cmdExec:              cmdExec1,
		suseOperationTimeout: 30,
		logger:               logger,
		osRelease:            "mi52-230318-r001",
	}

	var tests = []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "success create dir",
			fields: positiveCreateDir,
			args: args{
				auth: &authParam,
			},
			wantErr: false,
		},
		{
			name:   "failed create dir",
			fields: negativeCreateDir,
			args: args{
				auth: &authParam,
			},
			wantErr: true,
		},
		{
			name:   "success create repo",
			fields: positiveCreateRepo,
			args: args{
				auth: &authParam,
			},
			wantErr: false,
		},
		{
			name:   "failed create repo",
			fields: negativeCreateRepo,
			args: args{
				auth: &authParam,
			},
			wantErr: true,
		},
		{
			name:   "success create channel",
			fields: positiveCreateChannel,
			args: args{
				auth: &authParam,
			},
			wantErr: false,
		},
		{
			name:   "failed create channel",
			fields: negativeCreateChannel,
			args: args{
				auth: &authParam,
			},
			wantErr: true,
		},
		{
			name:   "success create add repo",
			fields: positiveAssociateRepo,
			args: args{
				auth: &authParam,
			},
			wantErr: false,
		},
		{
			name:   "failed create add repo",
			fields: negativeAssociateRepo,
			args: args{
				auth: &authParam,
			},
			wantErr: true,
		},
		{
			name:   "success sync repo",
			fields: positiveSyncRepo,
			args: args{
				auth: &authParam,
			},
			wantErr: false,
		},
		{
			name:   "failed sync repo",
			fields: negativeSyncRepo,
			args: args{
				auth: &authParam,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &CreateOsRelease{
				sumanProxy:           tt.fields.sumanProxy,
				suse:                 tt.fields.suse,
				cmdExec:              tt.fields.cmdExec,
				suseOperationTimeout: tt.fields.suseOperationTimeout,
				logger:               tt.fields.logger,
				osRelease:            tt.fields.osRelease,
			}
			if err := h.createExtraRepo(tt.args.auth); (err != nil) != tt.wantErr {
				t.Errorf("createExtraRepo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateOsRelease_createExtraRepoFS(t *testing.T) {
	type fields struct {
		sumanProxy           _sumanUseCase.IProxy
		suse                 _sumanUseCase.ISuseManager
		cmdExec              cmdexecutor.ICMDExecutor
		suseOperationTimeout int
		logger               *zap.Logger
		osRelease            string
	}
	var resultRun []string
	suseManagerMock := new(susemocks.ISuseManager)
	suseAPIProxy := new(susemocks.IProxy)
	cmdExec1 := new(cmdExecutor.ICMDExecutor)
	cmdExec1.On("CreateDirectory", mock.Anything, mock.Anything).Return(nil)
	cmdExec1.On("ExecuteCommand", mock.Anything, mock.Anything, mock.Anything).Return(resultRun, nil)
	cmdExec2 := new(cmdExecutor.ICMDExecutor)
	cmdExec2.On("CreateDirectory", mock.Anything, mock.Anything).Return(nil)
	cmdExec2.On("ExecuteCommand", mock.Anything, mock.Anything, mock.Anything).Return(resultRun, nil)
	cmdExecErr1 := new(cmdExecutor.ICMDExecutor)
	cmdExecErr1.On("CreateDirectory", mock.Anything, mock.Anything).Return(errors.New("error creating directory"))
	cmdExecErr2 := new(cmdExecutor.ICMDExecutor)
	cmdExecErr2.On("CreateDirectory", mock.Anything, mock.Anything).Return(nil)
	cmdExecErr2.On("ExecuteCommand", mock.Anything, mock.Anything, mock.Anything).Return(resultRun, errors.New("error running script"))

	logger := logging.NewTestingLogger(t.Name())

	testResult1 := fields{
		sumanProxy:           suseAPIProxy,
		suse:                 suseManagerMock,
		cmdExec:              cmdExec1,
		suseOperationTimeout: 30,
		logger:               logger,
		osRelease:            "mi52-230318-r001",
	}

	testResult2 := fields{
		sumanProxy:           suseAPIProxy,
		suse:                 suseManagerMock,
		cmdExec:              cmdExecErr1,
		suseOperationTimeout: 30,
		logger:               logger,
		osRelease:            "mi52-230318-r001",
	}
	testResult3 := fields{
		sumanProxy:           suseAPIProxy,
		suse:                 suseManagerMock,
		cmdExec:              cmdExec2,
		suseOperationTimeout: 30,
		logger:               logger,
		osRelease:            "mi52-230318-r001",
	}

	testResult4 := fields{
		sumanProxy:           suseAPIProxy,
		suse:                 suseManagerMock,
		cmdExec:              cmdExecErr2,
		suseOperationTimeout: 30,
		logger:               logger,
		osRelease:            "mi52-230318-r001",
	}

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "succes test1",
			fields:  testResult1,
			wantErr: false,
		}, {
			name:    "error test1",
			fields:  testResult2,
			wantErr: true,
		}, {
			name:    "succes test2",
			fields:  testResult3,
			wantErr: false,
		}, {
			name:    "error test2",
			fields:  testResult4,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &CreateOsRelease{
				sumanProxy:           tt.fields.sumanProxy,
				suse:                 tt.fields.suse,
				cmdExec:              tt.fields.cmdExec,
				suseOperationTimeout: tt.fields.suseOperationTimeout,
				logger:               tt.fields.logger,
				osRelease:            tt.fields.osRelease,
			}
			if err := h.createExtraRepoFS(); (err != nil) != tt.wantErr {
				t.Errorf("createExtraRepoFS() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateOsRelease_CheckOsReleaseDefaultExtraChannels(t *testing.T) {
	type fields struct {
		sumanProxy           _sumanUseCase.IProxy
		suse                 _sumanUseCase.ISuseManager
		cmdExec              cmdexecutor.ICMDExecutor
		suseOperationTimeout int
		logger               *zap.Logger
		osRelease            string
	}
	type args struct {
		auth *_sumanUseCase.AuthParams
	}
	logger := logging.NewTestingLogger(t.Name())

	suseManagerMock := new(susemocks.ISuseManager)
	suseManagerMock.On("SuseManagerCall", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&util.HTTPHelperStruct{}, nil)
	suseManagerMockErr := new(susemocks.ISuseManager)
	suseManagerMockErr.On("SuseManagerCall", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&util.HTTPHelperStruct{}, errors.New("failed to logout suse manager"))
	suseAPIProxy1 := new(susemocks.IProxy)
	suseAPIProxy1.On("ChannelSoftwareCreate", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(1, nil)
	suseAPIProxy1.On("ChannelSoftwareIsExisting", mock.Anything, mock.Anything, mock.Anything).Return(true, nil)
	suseAPIProxyErr1 := new(susemocks.IProxy)
	suseAPIProxyErr1.On("ChannelSoftwareCreate", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(1, nil)
	suseAPIProxyErr1.On("ChannelSoftwareIsExisting", mock.Anything, mock.Anything, mock.Anything).Return(false, errors.New("failed"))
	suseAPIProxy2 := new(susemocks.IProxy)
	suseAPIProxy2.On("ChannelSoftwareCreate", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(1, nil)
	suseAPIProxy2.On("ChannelSoftwareIsExisting", mock.Anything, mock.Anything, mock.Anything).Return(true, nil)
	suseAPIProxyErr2 := new(susemocks.IProxy)
	suseAPIProxyErr2.On("ChannelSoftwareCreate", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(0, errors.New("failed"))
	suseAPIProxyErr2.On("ChannelSoftwareIsExisting", mock.Anything, mock.Anything, mock.Anything).Return(true, nil)

	authParam := _sumanUseCase.AuthParams{
		SessionKey: "test key",
		Host:       "test Hostname",
	}
	positiveBasePresent1 := fields{
		sumanProxy:           suseAPIProxy1,
		suse:                 suseManagerMock,
		suseOperationTimeout: 30,
		logger:               logger,
		osRelease:            "mi52-230318-r001",
	}
	negativeBasePresent1 := fields{
		sumanProxy:           suseAPIProxyErr1,
		suse:                 suseManagerMock,
		suseOperationTimeout: 30,
		logger:               logger,
		osRelease:            "mi52-230318-r001",
	}
	positiveBaseCreate := fields{
		sumanProxy:           suseAPIProxy1,
		suse:                 suseManagerMock,
		suseOperationTimeout: 30,
		logger:               logger,
		osRelease:            "mi52-230318-r001",
	}
	negativeBaseCreate := fields{
		sumanProxy:           suseAPIProxyErr1,
		suse:                 suseManagerMock,
		suseOperationTimeout: 30,
		logger:               logger,
		osRelease:            "mi52-230318-r001",
	}

	var tests = []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "base channel present",
			fields: positiveBasePresent1,
			args: args{
				auth: &authParam,
			},
			wantErr: false,
		},
		{
			name:   "base channel present",
			fields: negativeBasePresent1,
			args: args{
				auth: &authParam,
			},
			wantErr: true,
		},
		{
			name:   "base channel create successful",
			fields: positiveBaseCreate,
			args: args{
				auth: &authParam,
			},
			wantErr: false,
		},
		{
			name:   "base channel create failed",
			fields: negativeBaseCreate,
			args: args{
				auth: &authParam,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &CreateOsRelease{
				sumanProxy:           tt.fields.sumanProxy,
				suse:                 tt.fields.suse,
				cmdExec:              tt.fields.cmdExec,
				suseOperationTimeout: tt.fields.suseOperationTimeout,
				logger:               tt.fields.logger,
				osRelease:            tt.fields.osRelease,
			}
			if err := h.CheckOsReleaseDefaultExtraChannels(tt.args.auth); (err != nil) != tt.wantErr {
				t.Errorf("CheckOsReleaseDefaultExtraChannels() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
