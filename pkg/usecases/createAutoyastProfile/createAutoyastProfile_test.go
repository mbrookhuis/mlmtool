package createautoyastprofile

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	sumamodels "mlmtool/pkg/models/susemanager"
	_sumanUseCase "mlmtool/pkg/usecases/susemanager"
	susemocks "mlmtool/pkg/usecases/susemanager/mocks"
	logging "mlmtool/pkg/util/logger"
	util "mlmtool/pkg/util/rest"
)

/*
var data1 = `<?xml version="1.0"?>
<!DOCTYPE profile>
<profile xmlns="http://www.suse.com/1.0/yast2ns" xmlns:config="http://www.suse.com/1.0/configns">

	<scripts>
	 <pre-scripts config:type="list">
	  <script>
	   <interpreter>shell</interpreter>
	   <filename>pre-fetch.sh</filename>
	   <location>http://$SUMAN_SERVER/pub/dt-a4pod/pre-fetch.sh</location>
	   <notification>Please wait while pre-fetch.sh is running...</notification>
	   <debug config:type="boolean">false</debug>
	   <feedback config:type="boolean">false</feedback>
	  </script>
	 </pre-scripts>
	</scripts>
	<software>
	 <products config:type="list">
	  <product>SLES</product>
	 </products>
	</software>

</profile>
`

var fileName = "/tmp/dtag_server.xml"
*/
func TestCreateAutoyastProfile_addProfileVar(t *testing.T) {
	type fields struct {
		sumanProxy           _sumanUseCase.IProxy
		suse                 _sumanUseCase.ISuseManager
		suseOperationTimeout int
		logger               *zap.Logger
		locationXML          string
		profileName          string
		replaceExisting      bool
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
	suseAPIProxy.On("KickstartProfileSetVariables", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(1, nil)
	suseAPIProxyErr1 := new(susemocks.IProxy)
	suseAPIProxyErr1.On("KickstartProfileSetVariables", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(0, errors.New("failed to add"))
	authParam := _sumanUseCase.AuthParams{
		SessionKey: "test key",
		Host:       "test Hostname",
	}

	positiveResult := fields{
		sumanProxy:           suseAPIProxy,
		suse:                 suseManagerMock,
		suseOperationTimeout: 30,
		logger:               logger,
		locationXML:          "/opt/autoyast",
		profileName:          "dtag_server",
		replaceExisting:      false,
	}

	negativeResult1 := fields{
		sumanProxy:           suseAPIProxyErr1,
		suse:                 suseManagerMock,
		suseOperationTimeout: 30,
		logger:               logger,
		locationXML:          "/opt/autoyast",
		profileName:          "dtag_server",
		replaceExisting:      false,
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
		}, {
			name:   "failed",
			fields: negativeResult1,
			args: args{
				auth: &authParam,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &CreateAutoyastProfile{
				sumanProxy:           tt.fields.sumanProxy,
				suse:                 tt.fields.suse,
				suseOperationTimeout: tt.fields.suseOperationTimeout,
				logger:               tt.fields.logger,
				locationXML:          tt.fields.locationXML,
				profileName:          tt.fields.profileName,
				replaceExisting:      tt.fields.replaceExisting,
			}
			if err := h.addProfileVar(tt.args.auth); (err != nil) != tt.wantErr {
				t.Errorf("addProfileVar() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateAutoyastProfile_createProfile(t *testing.T) {
	type fields struct {
		sumanProxy           _sumanUseCase.IProxy
		suse                 _sumanUseCase.ISuseManager
		suseOperationTimeout int
		logger               *zap.Logger
		locationXML          string
		profileName          string
		replaceExisting      bool
	}
	type args struct {
		autoyastXML string
		auth        *_sumanUseCase.AuthParams
	}
	logger := logging.NewTestingLogger(t.Name())
	suseManagerMock := new(susemocks.ISuseManager)
	suseManagerMock.On("SuseManagerCall", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&util.HTTPHelperStruct{}, nil)
	suseManagerMockErr := new(susemocks.ISuseManager)
	suseManagerMockErr.On("SuseManagerCall", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&util.HTTPHelperStruct{}, errors.New("failed to logout suse manager"))
	suseAPIProxy := new(susemocks.IProxy)
	suseAPIProxy.On("KickstartImportRawFile", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(1, nil)
	suseAPIProxyErr1 := new(susemocks.IProxy)
	suseAPIProxyErr1.On("KickstartImportRawFile", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(0, errors.New("failed to add"))
	authParam := _sumanUseCase.AuthParams{
		SessionKey: "test key",
		Host:       "test Hostname",
	}

	positiveResult := fields{
		sumanProxy:           suseAPIProxy,
		suse:                 suseManagerMock,
		suseOperationTimeout: 30,
		logger:               logger,
		locationXML:          "/opt/autoyast",
		profileName:          "dtag_server",
		replaceExisting:      false,
	}

	negativeResult1 := fields{
		sumanProxy:           suseAPIProxyErr1,
		suse:                 suseManagerMock,
		suseOperationTimeout: 30,
		logger:               logger,
		locationXML:          "/opt/autoyast",
		profileName:          "dtag_server",
		replaceExisting:      false,
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
				autoyastXML: "<xml>data</xml>\n<line>1</line>\n",
				auth:        &authParam,
			},
			wantErr: false,
		}, {
			name:   "failed",
			fields: negativeResult1,
			args: args{
				autoyastXML: "<xml>data</xml>\n<line>1</line>\n",
				auth:        &authParam,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &CreateAutoyastProfile{
				sumanProxy:           tt.fields.sumanProxy,
				suse:                 tt.fields.suse,
				suseOperationTimeout: tt.fields.suseOperationTimeout,
				logger:               tt.fields.logger,
				locationXML:          tt.fields.locationXML,
				profileName:          tt.fields.profileName,
				replaceExisting:      tt.fields.replaceExisting,
			}
			if err := h.createProfile(tt.args.autoyastXML, tt.args.auth); (err != nil) != tt.wantErr {
				t.Errorf("createProfile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateAutoyastProfile_deleteProfile(t *testing.T) {
	type fields struct {
		sumanProxy           _sumanUseCase.IProxy
		suse                 _sumanUseCase.ISuseManager
		suseOperationTimeout int
		logger               *zap.Logger
		locationXML          string
		profileName          string
		replaceExisting      bool
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
	suseAPIProxy.On("KickstartDeleteProfile", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(1, nil)
	suseAPIProxyErr1 := new(susemocks.IProxy)
	suseAPIProxyErr1.On("KickstartDeleteProfile", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(0, errors.New("failed to add"))
	authParam := _sumanUseCase.AuthParams{
		SessionKey: "test key",
		Host:       "test Hostname",
	}

	positiveResult := fields{
		sumanProxy:           suseAPIProxy,
		suse:                 suseManagerMock,
		suseOperationTimeout: 30,
		logger:               logger,
		locationXML:          "/opt/autoyast",
		profileName:          "dtag_server",
		replaceExisting:      false,
	}

	negativeResult1 := fields{
		sumanProxy:           suseAPIProxyErr1,
		suse:                 suseManagerMock,
		suseOperationTimeout: 30,
		logger:               logger,
		locationXML:          "/opt/autoyast",
		profileName:          "dtag_server",
		replaceExisting:      false,
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
		}, {
			name:   "failed",
			fields: negativeResult1,
			args: args{
				auth: &authParam,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &CreateAutoyastProfile{
				sumanProxy:           tt.fields.sumanProxy,
				suse:                 tt.fields.suse,
				suseOperationTimeout: tt.fields.suseOperationTimeout,
				logger:               tt.fields.logger,
				locationXML:          tt.fields.locationXML,
				profileName:          tt.fields.profileName,
				replaceExisting:      tt.fields.replaceExisting,
			}
			if err := h.deleteProfile(tt.args.auth); (err != nil) != tt.wantErr {
				t.Errorf("deleteProfile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateAutoyastProfile_checkProfileIsPresent(t *testing.T) {
	type fields struct {
		sumanProxy           _sumanUseCase.IProxy
		suse                 _sumanUseCase.ISuseManager
		suseOperationTimeout int
		logger               *zap.Logger
		locationXML          string
		profileName          string
		replaceExisting      bool
	}
	type args struct {
		auth *_sumanUseCase.AuthParams
	}
	logger := logging.NewTestingLogger(t.Name())

	kickstartList := []sumamodels.KickstartListProfiles{
		{
			Name: "SL_SERVER",
		},
		{
			Name: "dtag_server",
		},
	}

	suseManagerMock := new(susemocks.ISuseManager)
	suseManagerMock.On("SuseManagerCall", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&util.HTTPHelperStruct{}, nil)
	suseManagerMockErr := new(susemocks.ISuseManager)
	suseManagerMockErr.On("SuseManagerCall", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&util.HTTPHelperStruct{}, errors.New("failed to logout suse manager"))
	suseAPIProxy := new(susemocks.IProxy)
	suseAPIProxy.On("KickstartListKickstarts", mock.Anything, mock.Anything).Return(kickstartList, nil)
	suseAPIProxyErr1 := new(susemocks.IProxy)
	suseAPIProxyErr1.On("KickstartListKickstarts", mock.Anything, mock.Anything).Return(nil, errors.New("failed"))
	suseAPIProxyErr2 := new(susemocks.IProxy)
	suseAPIProxyErr2.On("KickstartListKickstarts", mock.Anything, mock.Anything).Return(kickstartList, nil)
	authParam := _sumanUseCase.AuthParams{
		SessionKey: "test key",
		Host:       "test Hostname",
	}

	positiveResult := fields{
		sumanProxy:           suseAPIProxy,
		suse:                 suseManagerMock,
		suseOperationTimeout: 30,
		logger:               logger,
		locationXML:          "/opt/autoyast",
		profileName:          "dtag_server",
		replaceExisting:      false,
	}

	negativeResult1 := fields{
		sumanProxy:           suseAPIProxyErr1,
		suse:                 suseManagerMock,
		suseOperationTimeout: 30,
		logger:               logger,
		locationXML:          "/opt/autoyast",
		profileName:          "dtag_server",
		replaceExisting:      false,
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
		}, {
			name:   "failed",
			fields: negativeResult1,
			args: args{
				auth: &authParam,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &CreateAutoyastProfile{
				sumanProxy:           tt.fields.sumanProxy,
				suse:                 tt.fields.suse,
				suseOperationTimeout: tt.fields.suseOperationTimeout,
				logger:               tt.fields.logger,
				locationXML:          tt.fields.locationXML,
				profileName:          tt.fields.profileName,
				replaceExisting:      tt.fields.replaceExisting,
			}
			_, err := h.checkProfileIsPresent(tt.args.auth)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkProfileIsPresent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

/*
func TestCreateAutoyastProfile_getAutoyastXML(t *testing.T) {
	err := ioutil.WriteFile(fileName, []byte(data1), 0666)
	if err != nil {
		fmt.Printf("Error writing file")
		os.Exit(1)
	}
	type fields struct {
		sumanProxy           _sumanUseCase.IProxy
		suse                 _sumanUseCase.ISuseManager
		suseOperationTimeout int
		logger               *zap.Logger
		locationXML          string
		profileName          string
		replaceExisting      bool
	}
	logger := logging.NewTestingLogger(t.Name())
	suseManagerMock := new(susemocks.ISuseManager)
	suseAPIProxy := new(susemocks.IProxy)

	succesRead := fields{
		sumanProxy:           suseAPIProxy,
		suse:                 suseManagerMock,
		suseOperationTimeout: 30,
		logger:               logger,
		locationXML:          "/tmp",
		profileName:          "dtag_server",
		replaceExisting:      false,
	}

	failedRead := fields{
		sumanProxy:           suseAPIProxy,
		suse:                 suseManagerMock,
		suseOperationTimeout: 30,
		logger:               logger,
		locationXML:          "/tmp",
		profileName:          "does_not_exist",
		replaceExisting:      false,
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "read_success",
			fields: succesRead,
			wantErr: false,
		},
		{
			name:   "read_failed",
			fields: failedRead,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &CreateAutoyastProfile{
				sumanProxy:           tt.fields.sumanProxy,
				suse:                 tt.fields.suse,
				suseOperationTimeout: tt.fields.suseOperationTimeout,
				logger:               tt.fields.logger,
				locationXML:          tt.fields.locationXML,
				profileName:          tt.fields.profileName,
				replaceExisting:      tt.fields.replaceExisting,
			}
			_, err := h.getAutoyastXML()
			if (err != nil) != tt.wantErr {
				t.Errorf("getAutoyastXML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}


*/
