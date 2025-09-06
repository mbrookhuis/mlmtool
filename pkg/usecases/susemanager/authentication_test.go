// Package SUSE Manager - SUSE Manager api call and support functions
package susemanager

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	logging "mlmtool/pkg/util/logger"
	util "mlmtool/pkg/util/rest"
)

func TestProxy_SumanLogin(t *testing.T) {
	type fields struct {
		cfg               *SumanConfig
		contentTypeHeader map[string]string
		suse              ISuseManagerAPI
	}

	header := make(map[string]string)
	header["Content-Type"] = "application/json"

	config := &SumanConfig{
		Host:     "test host",
		Password: "test",
		Insecure: true,
		Login:    "test",
	}

	suseManagerMockErr := new(MockISuseManagerAPI)
	suseManagerMockErr.On("SuseManagerCall", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&util.HTTPHelperStruct{}, errors.New("Failed to login suse manager"))
	var successJSON = `{"success": true,
		"messages": []
	}`
	suseManagerMockSuccess := new(MockISuseManagerAPI)
	suseManagerMockSuccess.On("SuseManagerCall", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&util.HTTPHelperStruct{
		Cookies: []*http.Cookie{{
			Value: "1",
		}, {
			Value: "2",
		}, {
			Value: "3",
		}},
		StatusCode: 200,
		Body:       []byte(successJSON),
	}, nil)
	negativeKey := fields{
		cfg:               config,
		contentTypeHeader: header,
		suse:              suseManagerMockErr,
	}

	postiveKey := fields{
		cfg:               config,
		contentTypeHeader: header,
		suse:              suseManagerMockSuccess,
	}

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "Suman Login Negative",
			fields:  negativeKey,
			wantErr: true,
		},
		{
			name:    "Suman Login Postive",
			fields:  postiveKey,
			wantErr: false,
		},
	}
	logger := logging.NewTestingLogger(t.Name())
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Proxy{
				cfg:               tt.fields.cfg,
				contentTypeHeader: tt.fields.contentTypeHeader,
				suse:              tt.fields.suse,
				logger:            logger,
			}
			_, err := p.SumanLogin()
			if (err != nil) != tt.wantErr {
				t.Errorf("Proxy.SumanLogin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestProxy_SumanLogout(t *testing.T) {
	type fields struct {
		cfg               *SumanConfig
		contentTypeHeader map[string]string
		suse              ISuseManagerAPI
		logger            *zap.Logger
	}
	logger := logging.NewTestingLogger(t.Name())
	type args struct {
		auth AuthParams
	}

	header := make(map[string]string)
	header["Content-Type"] = "application/json"

	config := &SumanConfig{
		Host:     "test host",
		Password: "test",
		Insecure: true,
		Login:    "test",
	}

	suseManagerMock := new(MockISuseManagerAPI)
	suseManagerMock.On("SuseManagerCall", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&util.HTTPHelperStruct{}, nil)

	suseManagerMockErr := new(MockISuseManagerAPI)
	suseManagerMockErr.On("SuseManagerCall", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&util.HTTPHelperStruct{}, errors.New("Failed to logout suse manager"))

	positiveKey := fields{
		cfg:               config,
		contentTypeHeader: header,
		suse:              suseManagerMock,
		logger:            logger,
	}

	negativeKey := fields{
		cfg:               config,
		contentTypeHeader: header,
		suse:              suseManagerMockErr,
		logger:            logger,
	}

	authPara := AuthParams{
		SessionKey: "test key",
		Host:       "test Hostname",
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "Suman Logout Positive",
			fields: positiveKey,
			args: args{
				auth: authPara,
			},
			wantErr: false,
		},
		{
			name:   "Suman Logout Negative",
			fields: negativeKey,
			args: args{
				auth: authPara,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Proxy{
				cfg:               tt.fields.cfg,
				contentTypeHeader: tt.fields.contentTypeHeader,
				suse:              tt.fields.suse,
				logger:            tt.fields.logger,
			}
			if err := p.SumanLogout(tt.args.auth); (err != nil) != tt.wantErr {
				t.Errorf("Proxy.SumanLogout() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
