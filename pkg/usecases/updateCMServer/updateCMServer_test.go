package updatecmserver

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	sumamodels "mlmtool/pkg/models/susemanager"
	_sumanUseCase "mlmtool/pkg/usecases/susemanager"
	susemocks "mlmtool/pkg/usecases/susemanager/mocks"
	logging "mlmtool/pkg/util/logger"
	util "mlmtool/pkg/util/rest"
)

func TestUpdateCMServer_getBaseChannels(t *testing.T) {
	type fields struct {
		sumanProxy           _sumanUseCase.IProxy
		suse                 _sumanUseCase.ISuseManager
		suseOperationTimeout int
		logger               *zap.Logger
		osRelease            string
		updateServer         string
	}
	type args struct {
		auth *_sumanUseCase.AuthParams
	}

	logger := logging.NewTestingLogger(t.Name())
	authParam := _sumanUseCase.AuthParams{
		SessionKey: "test key",
		Host:       "test Hostname",
	}

	allChannels := []sumamodels.ChannelListSoftwareChannels{
		{
			Label: "sm43-230101-r001-channel1",
		},
		{
			Label: "sm43-230102-r001-channel2",
		},
		{
			Label: "sm43-230103-r001-channel3",
		},
	}

	suseManagerMock := new(susemocks.ISuseManager)
	suseManagerMock.On("SuseManagerCall", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&util.HTTPHelperStruct{}, nil)
	suseManagerMockErr := new(susemocks.ISuseManager)
	suseManagerMockErr.On("SuseManagerCall", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&util.HTTPHelperStruct{}, errors.New("failed to logout suse manager"))
	suseAPIProxy := new(susemocks.IProxy)
	suseAPIProxy.On("ChannelListSoftwareChannels", mock.Anything, mock.Anything, mock.Anything).Return(allChannels, nil)
	suseAPIProxyErr := new(susemocks.IProxy)
	suseAPIProxyErr.On("ChannelListSoftwareChannels", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("failed to get channels"))

	positiveResult := fields{
		sumanProxy:           suseAPIProxy,
		suse:                 suseManagerMock,
		suseOperationTimeout: 30,
		logger:               logger,
		osRelease:            "sm43-230102-r001",
		updateServer:         "server1"}

	negativeResult := fields{
		sumanProxy:           suseAPIProxyErr,
		suse:                 suseManagerMock,
		suseOperationTimeout: 30,
		logger:               logger,
		osRelease:            "sm43-230104-r001",
		updateServer:         "server1"}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "valid",
			fields: positiveResult,
			args: args{
				auth: &authParam,
			},
			wantErr: false,
		},
		{
			name:   "invalid",
			fields: negativeResult,
			args: args{
				auth: &authParam,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &UpdateCMServer{
				sumanProxy:           tt.fields.sumanProxy,
				suse:                 tt.fields.suse,
				suseOperationTimeout: tt.fields.suseOperationTimeout,
				logger:               tt.fields.logger,
				osRelease:            tt.fields.osRelease,
				updateServer:         tt.fields.updateServer,
			}
			_, err := h.getBaseChannels(tt.args.auth)
			if (err != nil) != tt.wantErr {
				t.Errorf("getBaseChannels() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUpdateCMServer_checkOsReleaseExists(t *testing.T) {
	type fields struct {
		sumanProxy           _sumanUseCase.IProxy
		suse                 _sumanUseCase.ISuseManager
		suseOperationTimeout int
		logger               *zap.Logger
		osRelease            string
		updateServer         string
	}
	type args struct {
		auth *_sumanUseCase.AuthParams
	}

	logger := logging.NewTestingLogger(t.Name())
	authParam := _sumanUseCase.AuthParams{
		SessionKey: "test key",
		Host:       "test Hostname",
	}
	allChannels := []sumamodels.ChannelListSoftwareChannels{
		{
			Label: "sm43-230101-r001-channel1",
		},
		{
			Label: "sm43-230102-r001-channel2",
		},
		{
			Label: "sm43-230103-r001-channel3",
		},
	}
	suseManagerMock := new(susemocks.ISuseManager)
	suseManagerMock.On("SuseManagerCall", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&util.HTTPHelperStruct{}, nil)
	suseManagerMockErr := new(susemocks.ISuseManager)
	suseManagerMockErr.On("SuseManagerCall", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&util.HTTPHelperStruct{}, errors.New("failed to logout suse manager"))
	suseAPIProxy := new(susemocks.IProxy)
	suseAPIProxy.On("ChannelListSoftwareChannels", mock.Anything, mock.Anything, mock.Anything).Return(allChannels, nil)
	suseAPIProxyErr := new(susemocks.IProxy)
	suseAPIProxyErr.On("ChannelListSoftwareChannels", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("failed to get channels"))
	negativeResult1 := fields{
		sumanProxy:           suseAPIProxyErr,
		suse:                 suseManagerMock,
		suseOperationTimeout: 30,
		logger:               logger,
		osRelease:            "sm43-230104-r001",
		updateServer:         "server1"}
	positiveResult2 := fields{
		sumanProxy:           suseAPIProxy,
		suse:                 suseManagerMock,
		suseOperationTimeout: 30,
		logger:               logger,
		osRelease:            "sm43-230102-r001",
		updateServer:         "server1"}
	negativeResult2 := fields{
		sumanProxy:           suseAPIProxy,
		suse:                 suseManagerMock,
		suseOperationTimeout: 30,
		logger:               logger,
		osRelease:            "sm43-230104-r001",
		updateServer:         "server1"}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "success",
			fields: positiveResult2,
			args: args{
				auth: &authParam,
			},
			wantErr: false,
		},
		{
			name:   "negative2",
			fields: negativeResult2,
			args: args{
				auth: &authParam,
			},
			wantErr: true,
		},
		{
			name:   "negative1",
			fields: negativeResult1,
			args: args{
				auth: &authParam,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &UpdateCMServer{
				sumanProxy:           tt.fields.sumanProxy,
				suse:                 tt.fields.suse,
				suseOperationTimeout: tt.fields.suseOperationTimeout,
				logger:               tt.fields.logger,
				osRelease:            tt.fields.osRelease,
				updateServer:         tt.fields.updateServer,
			}
			if err := h.checkOsReleaseExists(tt.args.auth); (err != nil) != tt.wantErr {
				t.Errorf("checkOsReleaseExists() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateCMServer_checkServer(t *testing.T) {
	type fields struct {
		sumanProxy           _sumanUseCase.IProxy
		suse                 _sumanUseCase.ISuseManager
		suseOperationTimeout int
		logger               *zap.Logger
		osRelease            string
		updateServer         string
	}
	type args struct {
		auth *_sumanUseCase.AuthParams
	}

	logger := logging.NewTestingLogger(t.Name())
	authParam := _sumanUseCase.AuthParams{
		SessionKey: "test key",
		Host:       "test Hostname",
	}

	activeSystems := []sumamodels.ActiveSystem{
		{
			ID:   1000100001,
			Name: "server1",
		},
		{
			ID:   1000100002,
			Name: "server2",
		},
	}

	baseChannelInfo := sumamodels.SubscribedBaseChannel{
		Label: "sm43-230102-r001",
	}

	formData := sumamodels.MgtsSrvFormular{}
	suseManagerMock := new(susemocks.ISuseManager)
	suseManagerMock.On("SuseManagerCall", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&util.HTTPHelperStruct{}, nil)
	suseManagerMockErr := new(susemocks.ISuseManager)
	suseManagerMockErr.On("SuseManagerCall", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&util.HTTPHelperStruct{}, errors.New("failed to logout suse manager"))
	suseAPIProxy := new(susemocks.IProxy)
	suseAPIProxy.On("SystemListActiveSystems", mock.Anything, mock.Anything, mock.Anything).Return(activeSystems, nil)
	suseAPIProxy.On("SystemGetSubscribedBaseChannel", mock.Anything, mock.Anything, mock.Anything).Return(baseChannelInfo, nil)
	suseAPIProxy.On("GetSystemFormulaData", mock.Anything, mock.Anything, mock.Anything).Return(&formData, nil)
	suseAPIProxyErr1 := new(susemocks.IProxy)
	suseAPIProxyErr1.On("SystemListActiveSystems", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("failed to get systems"))
	suseAPIProxyErr1.On("SystemGetSubscribedBaseChannel", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("failed to get basechannel info"))
	suseAPIProxyErr1.On("GetSystemFormulaData", mock.Anything, mock.Anything, mock.Anything).Return(&formData, nil)
	suseAPIProxyErr2 := new(susemocks.IProxy)
	suseAPIProxyErr2.On("SystemListActiveSystems", mock.Anything, mock.Anything, mock.Anything).Return(activeSystems, nil)
	suseAPIProxyErr2.On("SystemGetSubscribedBaseChannel", mock.Anything, mock.Anything, mock.Anything).Return(baseChannelInfo, errors.New("failed to get basechannel info"))
	suseAPIProxyErr2.On("GetSystemFormulaData", mock.Anything, mock.Anything, mock.Anything).Return(&formData, nil)
	suseAPIProxyErr3 := new(susemocks.IProxy)
	suseAPIProxyErr3.On("SystemListActiveSystems", mock.Anything, mock.Anything, mock.Anything).Return(activeSystems, nil)
	suseAPIProxyErr3.On("SystemGetSubscribedBaseChannel", mock.Anything, mock.Anything, mock.Anything).Return(baseChannelInfo, nil)
	suseAPIProxyErr3.On("GetSystemFormulaData", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("failed to get formular info"))

	positiveResult1 := fields{
		sumanProxy:           suseAPIProxy,
		suse:                 suseManagerMock,
		suseOperationTimeout: 30,
		logger:               logger,
		osRelease:            "sm43-230102-r001",
		updateServer:         "server1"}

	negativeResult1 := fields{
		sumanProxy:           suseAPIProxyErr1,
		suse:                 suseManagerMock,
		suseOperationTimeout: 30,
		logger:               logger,
		osRelease:            "sm43-230102-r001",
		updateServer:         "server1"}
	negativeResult2 := fields{
		sumanProxy:           suseAPIProxyErr2,
		suse:                 suseManagerMock,
		suseOperationTimeout: 30,
		logger:               logger,
		osRelease:            "sm43-230102-r001",
		updateServer:         "server1"}
	negativeResult3 := fields{
		sumanProxy:           suseAPIProxyErr3,
		suse:                 suseManagerMock,
		suseOperationTimeout: 30,
		logger:               logger,
		osRelease:            "sm43-230102-r001",
		updateServer:         "server1"}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "success",
			fields: positiveResult1,
			args: args{
				auth: &authParam,
			},
			wantErr: false,
		}, {
			name:   "missing system",
			fields: negativeResult1,
			args: args{
				auth: &authParam,
			},
			wantErr: true,
		}, {
			name:   "missing basechannel",
			fields: negativeResult2,
			args: args{
				auth: &authParam,
			},
			wantErr: true,
		}, {
			name:   "missing formular",
			fields: negativeResult3,
			args: args{
				auth: &authParam,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &UpdateCMServer{
				sumanProxy:           tt.fields.sumanProxy,
				suse:                 tt.fields.suse,
				suseOperationTimeout: tt.fields.suseOperationTimeout,
				logger:               tt.fields.logger,
				osRelease:            tt.fields.osRelease,
				updateServer:         tt.fields.updateServer,
			}
			_, err := h.checkServer(tt.args.auth)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkServer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUpdateCMServer_validateData(t *testing.T) {
	type fields struct {
		sumanProxy           _sumanUseCase.IProxy
		suse                 _sumanUseCase.ISuseManager
		suseOperationTimeout int
		logger               *zap.Logger
		osRelease            string
		updateServer         string
	}
	type args struct {
		auth *_sumanUseCase.AuthParams
	}

	logger := logging.NewTestingLogger(t.Name())
	authParam := _sumanUseCase.AuthParams{
		SessionKey: "test key",
		Host:       "test Hostname",
	}

	activeSystems := []sumamodels.ActiveSystem{
		{
			ID:   1000100001,
			Name: "server1",
		},
		{
			ID:   1000100002,
			Name: "server2",
		},
	}

	baseChannelInfo := sumamodels.SubscribedBaseChannel{
		Label: "sm43-230102-r001",
	}

	allChannels := []sumamodels.ChannelListSoftwareChannels{
		{
			Label: "sm43-230101-r001-channel1",
		},
		{
			Label: "sm43-230102-r001-channel2",
		},
		{
			Label: "sm43-230103-r001-channel3",
		},
	}

	formData := sumamodels.MgtsSrvFormular{}
	suseManagerMock := new(susemocks.ISuseManager)
	suseManagerMock.On("SuseManagerCall", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&util.HTTPHelperStruct{}, nil)
	suseManagerMockErr := new(susemocks.ISuseManager)
	suseManagerMockErr.On("SuseManagerCall", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&util.HTTPHelperStruct{}, errors.New("failed to logout suse manager"))
	suseAPIProxy := new(susemocks.IProxy)
	suseAPIProxy.On("SystemListActiveSystems", mock.Anything, mock.Anything, mock.Anything).Return(activeSystems, nil)
	suseAPIProxy.On("SystemGetSubscribedBaseChannel", mock.Anything, mock.Anything, mock.Anything).Return(baseChannelInfo, nil)
	suseAPIProxy.On("GetSystemFormulaData", mock.Anything, mock.Anything, mock.Anything).Return(&formData, nil)
	suseAPIProxy.On("ChannelListSoftwareChannels", mock.Anything, mock.Anything, mock.Anything).Return(allChannels, nil)
	suseAPIProxyErr1 := new(susemocks.IProxy)
	suseAPIProxyErr1.On("SystemListActiveSystems", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("failed to get systems"))
	suseAPIProxyErr1.On("SystemGetSubscribedBaseChannel", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("failed to get basechannel info"))
	suseAPIProxyErr1.On("GetSystemFormulaData", mock.Anything, mock.Anything, mock.Anything).Return(&formData, nil)
	suseAPIProxyErr1.On("ChannelListSoftwareChannels", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("failed to get channels"))
	suseAPIProxyErr2 := new(susemocks.IProxy)
	suseAPIProxyErr2.On("SystemListActiveSystems", mock.Anything, mock.Anything, mock.Anything).Return(activeSystems, nil)
	suseAPIProxyErr2.On("SystemGetSubscribedBaseChannel", mock.Anything, mock.Anything, mock.Anything).Return(baseChannelInfo, errors.New("failed to get basechannel info"))
	suseAPIProxyErr2.On("GetSystemFormulaData", mock.Anything, mock.Anything, mock.Anything).Return(&formData, nil)
	suseAPIProxyErr2.On("ChannelListSoftwareChannels", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("failed to get channels"))
	suseAPIProxyErr3 := new(susemocks.IProxy)
	suseAPIProxyErr3.On("SystemListActiveSystems", mock.Anything, mock.Anything, mock.Anything).Return(activeSystems, nil)
	suseAPIProxyErr3.On("SystemGetSubscribedBaseChannel", mock.Anything, mock.Anything, mock.Anything).Return(baseChannelInfo, nil)
	suseAPIProxyErr3.On("GetSystemFormulaData", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("failed to get formular info"))
	suseAPIProxyErr3.On("ChannelListSoftwareChannels", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("failed to get channels"))

	positiveResult1 := fields{
		sumanProxy:           suseAPIProxy,
		suse:                 suseManagerMock,
		suseOperationTimeout: 30,
		logger:               logger,
		osRelease:            "sm43-230103-r001",
		updateServer:         "server1"}

	negativeResult1 := fields{
		sumanProxy:           suseAPIProxyErr1,
		suse:                 suseManagerMock,
		suseOperationTimeout: 30,
		logger:               logger,
		osRelease:            "sm43-230102-r001",
		updateServer:         "server1"}
	negativeResult2 := fields{
		sumanProxy:           suseAPIProxyErr2,
		suse:                 suseManagerMock,
		suseOperationTimeout: 30,
		logger:               logger,
		osRelease:            "sm43-230102-r001",
		updateServer:         "server1"}
	negativeResult3 := fields{
		sumanProxy:           suseAPIProxyErr3,
		suse:                 suseManagerMock,
		suseOperationTimeout: 30,
		logger:               logger,
		osRelease:            "sm43-230102-r001",
		updateServer:         "server1"}
	negativeResult4 := fields{
		sumanProxy:           suseAPIProxy,
		suse:                 suseManagerMock,
		suseOperationTimeout: 30,
		logger:               logger,
		osRelease:            "sm43-230108-r001",
		updateServer:         "server1"}
	negativeResult5 := fields{
		sumanProxy:           suseAPIProxyErr1,
		suse:                 suseManagerMock,
		suseOperationTimeout: 30,
		logger:               logger,
		osRelease:            "sm43-230104-r001",
		updateServer:         "server1"}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "check_server success",
			fields: positiveResult1,
			args: args{
				auth: &authParam,
			},
			wantErr: false,
		}, {
			name:   "check_server missing system",
			fields: negativeResult1,
			args: args{
				auth: &authParam,
			},
			wantErr: true,
		}, {
			name:   "check_server missing basechannel",
			fields: negativeResult2,
			args: args{
				auth: &authParam,
			},
			wantErr: true,
		}, {
			name:   "check_server missing formular",
			fields: negativeResult3,
			args: args{
				auth: &authParam,
			},
			wantErr: true,
		},
		{
			name:   "osrelease_exist negative2",
			fields: negativeResult5,
			args: args{
				auth: &authParam,
			},
			wantErr: true,
		},
		{
			name:   "osrelease_exist negative1",
			fields: negativeResult4,
			args: args{
				auth: &authParam,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &UpdateCMServer{
				sumanProxy:           tt.fields.sumanProxy,
				suse:                 tt.fields.suse,
				suseOperationTimeout: tt.fields.suseOperationTimeout,
				logger:               tt.fields.logger,
				osRelease:            tt.fields.osRelease,
				updateServer:         tt.fields.updateServer,
			}
			_, err := h.validateData(tt.args.auth)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUpdateCMServer_setNewChannels(t *testing.T) {
	type fields struct {
		sumanProxy           _sumanUseCase.IProxy
		suse                 _sumanUseCase.ISuseManager
		suseOperationTimeout int
		logger               *zap.Logger
		osRelease            string
		updateServer         string
	}
	type args struct {
		serverInfo ServerHostInfo
		auth       *_sumanUseCase.AuthParams
	}
	logger := logging.NewTestingLogger(t.Name())
	authParam := _sumanUseCase.AuthParams{
		SessionKey: "test key",
		Host:       "test Hostname",
	}

	serverInfo := ServerHostInfo{
		serverID:           100001000,
		currentBaseChannel: "sm43-230102-r001-basechannel",
		serverType:         "sm43",
		testingMode:        true,
	}
	suseManagerMock := new(susemocks.ISuseManager)
	suseManagerMock.On("SuseManagerCall", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&util.HTTPHelperStruct{}, nil)
	suseManagerMock.On("ChangeChannels", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	suseManagerMockErr := new(susemocks.ISuseManager)
	suseManagerMockErr.On("SuseManagerCall", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&util.HTTPHelperStruct{}, errors.New("failed to logout suse manager"))
	suseManagerMockErr.On("ChangeChannels", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("basechannel set failed"))
	suseAPIProxy := new(susemocks.IProxy)

	positiveResult := fields{
		sumanProxy:           suseAPIProxy,
		suse:                 suseManagerMock,
		suseOperationTimeout: 30,
		logger:               logger,
		osRelease:            "sm43-230102-r001",
		updateServer:         "server1"}

	negativeResult := fields{
		sumanProxy:           suseAPIProxy,
		suse:                 suseManagerMockErr,
		suseOperationTimeout: 30,
		logger:               logger,
		osRelease:            "sm43-230104-r001",
		updateServer:         "server1"}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "valid",
			fields: positiveResult,
			args: args{
				auth:       &authParam,
				serverInfo: serverInfo,
			},
			wantErr: false,
		},
		{
			name:   "invalid",
			fields: negativeResult,
			args: args{
				auth:       &authParam,
				serverInfo: serverInfo,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &UpdateCMServer{
				sumanProxy:           tt.fields.sumanProxy,
				suse:                 tt.fields.suse,
				suseOperationTimeout: tt.fields.suseOperationTimeout,
				logger:               tt.fields.logger,
				osRelease:            tt.fields.osRelease,
				updateServer:         tt.fields.updateServer,
			}
			if err := h.setNewChannels(tt.args.serverInfo, tt.args.auth); (err != nil) != tt.wantErr {
				t.Errorf("setNewChannels() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateCMServer_doUpdate(t *testing.T) {
	type fields struct {
		sumanProxy           _sumanUseCase.IProxy
		suse                 _sumanUseCase.ISuseManager
		suseOperationTimeout int
		logger               *zap.Logger
		osRelease            string
		updateServer         string
	}
	type args struct {
		serverInfo ServerHostInfo
		auth       *_sumanUseCase.AuthParams
	}
	logger := logging.NewTestingLogger(t.Name())
	authParam := _sumanUseCase.AuthParams{
		SessionKey: "test key",
		Host:       "test Hostname",
	}

	serverInfo := ServerHostInfo{
		serverID:           100001000,
		currentBaseChannel: "sm43-230102-r001-basechannel",
		serverType:         "sm43",
		currentOsRelease:   "sm43-230102-r001",
		newBaseChannel:     "sm43-230103-r001-basechannel",
		spMig:              false,
		testingMode:        true,
	}
	suseManagerMock := new(susemocks.ISuseManager)
	suseManagerMock.On("SuseManagerCall", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&util.HTTPHelperStruct{}, nil)
	suseManagerMock.On("ChangeChannels", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	suseManagerMockErr := new(susemocks.ISuseManager)
	suseManagerMockErr.On("SuseManagerCall", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&util.HTTPHelperStruct{}, errors.New("failed to logout suse manager"))
	suseManagerMockErr.On("ChangeChannels", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("basechannel set failed"))
	suseAPIProxy := new(susemocks.IProxy)
	suseAPIProxy.On("ScheduleScriptRun", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	suseAPIProxy.On("SystemScheduleReboot", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	suseAPIProxyErr1 := new(susemocks.IProxy)
	suseAPIProxyErr1.On("ScheduleScriptRun", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("script failed"))
	suseAPIProxyErr1.On("SystemScheduleReboot", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	suseAPIProxyErr2 := new(susemocks.IProxy)
	suseAPIProxyErr2.On("ScheduleScriptRun", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	suseAPIProxyErr2.On("SystemScheduleReboot", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("reboot failed"))

	positiveBase := fields{
		sumanProxy:           suseAPIProxy,
		suse:                 suseManagerMock,
		suseOperationTimeout: 30,
		logger:               logger,
		osRelease:            "sm43-230102-r001",
		updateServer:         "server1"}

	negativeBase := fields{
		sumanProxy:           suseAPIProxy,
		suse:                 suseManagerMockErr,
		suseOperationTimeout: 30,
		logger:               logger,
		osRelease:            "sm43-230104-r001",
		updateServer:         "server1"}
	positiveUpdateMi := fields{
		sumanProxy:           suseAPIProxy,
		suse:                 suseManagerMock,
		suseOperationTimeout: 30,
		logger:               logger,
		osRelease:            "mi52-230102-r001",
		updateServer:         "server1"}

	negativeUpdateMi := fields{
		sumanProxy:           suseAPIProxyErr1,
		suse:                 suseManagerMock,
		suseOperationTimeout: 30,
		logger:               logger,
		osRelease:            "mi52-230102-r001",
		updateServer:         "server1"}
	positiveUpdateOther := fields{
		sumanProxy:           suseAPIProxy,
		suse:                 suseManagerMock,
		suseOperationTimeout: 30,
		logger:               logger,
		osRelease:            "sm43-230102-r001",
		updateServer:         "server1"}
	negativeUpdateOther := fields{
		sumanProxy:           suseAPIProxyErr1,
		suse:                 suseManagerMock,
		suseOperationTimeout: 30,
		logger:               logger,
		osRelease:            "sm43-230102-r001",
		updateServer:         "server1"}
	positiveReboot := fields{
		sumanProxy:           suseAPIProxy,
		suse:                 suseManagerMock,
		suseOperationTimeout: 30,
		logger:               logger,
		osRelease:            "sm43-230102-r001",
		updateServer:         "server1"}
	negativeReboot := fields{
		sumanProxy:           suseAPIProxyErr2,
		suse:                 suseManagerMock,
		suseOperationTimeout: 30,
		logger:               logger,
		osRelease:            "sm43-230102-r001",
		updateServer:         "server1"}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "basechannels set successful",
			fields: positiveBase,
			args: args{
				auth:       &authParam,
				serverInfo: serverInfo,
			},
			wantErr: false,
		},
		{
			name:   "basechannels set failed",
			fields: negativeBase,
			args: args{
				auth:       &authParam,
				serverInfo: serverInfo,
			},
			wantErr: true,
		},
		{
			name:   "update mi valid",
			fields: positiveUpdateMi,
			args: args{
				auth:       &authParam,
				serverInfo: serverInfo,
			},
			wantErr: false,
		},
		{
			name:   "update mi failed",
			fields: negativeUpdateMi,
			args: args{
				auth:       &authParam,
				serverInfo: serverInfo,
			},
			wantErr: true,
		},
		{
			name:   "update other valid",
			fields: positiveUpdateOther,
			args: args{
				auth:       &authParam,
				serverInfo: serverInfo,
			},
			wantErr: false,
		},
		{
			name:   "update other failed",
			fields: negativeUpdateOther,
			args: args{
				auth:       &authParam,
				serverInfo: serverInfo,
			},
			wantErr: true,
		},
		{
			name:   "reboot valid",
			fields: positiveReboot,
			args: args{
				auth:       &authParam,
				serverInfo: serverInfo,
			},
			wantErr: false,
		},
		{
			name:   "reboot failed",
			fields: negativeReboot,
			args: args{
				auth:       &authParam,
				serverInfo: serverInfo,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &UpdateCMServer{
				sumanProxy:           tt.fields.sumanProxy,
				suse:                 tt.fields.suse,
				suseOperationTimeout: tt.fields.suseOperationTimeout,
				logger:               tt.fields.logger,
				osRelease:            tt.fields.osRelease,
				updateServer:         tt.fields.updateServer,
			}
			if err := h.doUpdate(tt.args.serverInfo, tt.args.auth); (err != nil) != tt.wantErr {
				t.Errorf("doUpdate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateCMServer_doSPMig(t *testing.T) {
	type fields struct {
		sumanProxy           _sumanUseCase.IProxy
		suse                 _sumanUseCase.ISuseManager
		suseOperationTimeout int
		logger               *zap.Logger
		osRelease            string
		updateServer         string
	}
	type args struct {
		serverInfo ServerHostInfo
		auth       *_sumanUseCase.AuthParams
	}
	logger := logging.NewTestingLogger(t.Name())
	authParam := _sumanUseCase.AuthParams{
		SessionKey: "test key",
		Host:       "test Hostname",
	}

	serverInfo := ServerHostInfo{
		serverID:           100001000,
		currentBaseChannel: "sm43-230102-r001-basechannel",
		serverType:         "sm43",
		currentOsRelease:   "sm43-230102-r001",
		newBaseChannel:     "sm43-230103-r001-basechannel",
		spMig:              false,
		testingMode:        true,
	}
	suseManagerMock := new(susemocks.ISuseManager)
	suseManagerMock.On("SuseManagerCall", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&util.HTTPHelperStruct{}, nil)
	suseManagerMock.On("ChangeChannels", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	suseManagerMockErr := new(susemocks.ISuseManager)
	suseManagerMockErr.On("SuseManagerCall", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&util.HTTPHelperStruct{}, errors.New("failed to logout suse manager"))
	suseManagerMockErr.On("ChangeChannels", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("basechannel set failed"))
	suseAPIProxy := new(susemocks.IProxy)
	suseAPIProxy.On("ScheduleScriptRun", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	suseAPIProxy.On("SystemScheduleReboot", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	suseAPIProxyErr1 := new(susemocks.IProxy)
	suseAPIProxyErr1.On("ScheduleScriptRun", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("script failed"))
	suseAPIProxyErr1.On("SystemScheduleReboot", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	suseAPIProxyErr2 := new(susemocks.IProxy)
	suseAPIProxyErr2.On("ScheduleScriptRun", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	suseAPIProxyErr2.On("SystemScheduleReboot", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("reboot failed"))

	positiveBase := fields{
		sumanProxy:           suseAPIProxy,
		suse:                 suseManagerMock,
		suseOperationTimeout: 30,
		logger:               logger,
		osRelease:            "sm43-230102-r001",
		updateServer:         "server1"}

	negativeBase := fields{
		sumanProxy:           suseAPIProxy,
		suse:                 suseManagerMockErr,
		suseOperationTimeout: 30,
		logger:               logger,
		osRelease:            "sm43-230104-r001",
		updateServer:         "server1"}
	positiveUpdateMi := fields{
		sumanProxy:           suseAPIProxy,
		suse:                 suseManagerMock,
		suseOperationTimeout: 30,
		logger:               logger,
		osRelease:            "mi52-230102-r001",
		updateServer:         "server1"}

	negativeUpdateMi := fields{
		sumanProxy:           suseAPIProxyErr1,
		suse:                 suseManagerMock,
		suseOperationTimeout: 30,
		logger:               logger,
		osRelease:            "mi52-230102-r001",
		updateServer:         "server1"}
	positiveUpdateOther := fields{
		sumanProxy:           suseAPIProxy,
		suse:                 suseManagerMock,
		suseOperationTimeout: 30,
		logger:               logger,
		osRelease:            "sm43-230102-r001",
		updateServer:         "server1"}
	negativeUpdateOther := fields{
		sumanProxy:           suseAPIProxyErr1,
		suse:                 suseManagerMock,
		suseOperationTimeout: 30,
		logger:               logger,
		osRelease:            "sm43-230102-r001",
		updateServer:         "server1"}
	positiveReboot := fields{
		sumanProxy:           suseAPIProxy,
		suse:                 suseManagerMock,
		suseOperationTimeout: 30,
		logger:               logger,
		osRelease:            "sm43-230102-r001",
		updateServer:         "server1"}
	negativeReboot := fields{
		sumanProxy:           suseAPIProxyErr2,
		suse:                 suseManagerMock,
		suseOperationTimeout: 30,
		logger:               logger,
		osRelease:            "sm43-230102-r001",
		updateServer:         "server1"}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "basechannels set successful",
			fields: positiveBase,
			args: args{
				auth:       &authParam,
				serverInfo: serverInfo,
			},
			wantErr: false,
		},
		{
			name:   "basechannels set failed",
			fields: negativeBase,
			args: args{
				auth:       &authParam,
				serverInfo: serverInfo,
			},
			wantErr: true,
		},
		{
			name:   "update mi valid",
			fields: positiveUpdateMi,
			args: args{
				auth:       &authParam,
				serverInfo: serverInfo,
			},
			wantErr: false,
		},
		{
			name:   "update mi failed",
			fields: negativeUpdateMi,
			args: args{
				auth:       &authParam,
				serverInfo: serverInfo,
			},
			wantErr: true,
		},
		{
			name:   "update other valid",
			fields: positiveUpdateOther,
			args: args{
				auth:       &authParam,
				serverInfo: serverInfo,
			},
			wantErr: false,
		},
		{
			name:   "update other failed",
			fields: negativeUpdateOther,
			args: args{
				auth:       &authParam,
				serverInfo: serverInfo,
			},
			wantErr: true,
		},
		{
			name:   "reboot valid",
			fields: positiveReboot,
			args: args{
				auth:       &authParam,
				serverInfo: serverInfo,
			},
			wantErr: false,
		},
		{
			name:   "reboot failed",
			fields: negativeReboot,
			args: args{
				auth:       &authParam,
				serverInfo: serverInfo,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &UpdateCMServer{
				sumanProxy:           tt.fields.sumanProxy,
				suse:                 tt.fields.suse,
				suseOperationTimeout: tt.fields.suseOperationTimeout,
				logger:               tt.fields.logger,
				osRelease:            tt.fields.osRelease,
				updateServer:         tt.fields.updateServer,
			}
			if err := h.doSPMig(tt.args.serverInfo, tt.args.auth); (err != nil) != tt.wantErr {
				t.Errorf("doSPMig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
