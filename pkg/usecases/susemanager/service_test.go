// Package SUSE Manager - SUSE Manager api call and support functions
package susemanager

import (
	"reflect"
	"testing"

	"go.uber.org/zap"

	logging "mlmtool/pkg/util/logger"
)

func TestNewProxy(t *testing.T) {
	type args struct {
		s      *SumanConfig
		suse   ISuseManagerAPI
		logger *zap.Logger
	}

	config := &SumanConfig{
		Host:     "test host",
		Password: "test",
		Insecure: true,
		Login:    "test",
	}
	logger := logging.NewTestingLogger(t.Name())

	suseManagerMock := new(MockISuseManagerAPI)
	iproxy := NewProxy(config, suseManagerMock, logger, 5)

	tests := []struct {
		name string
		args args
		want IProxy
	}{
		{
			name: "New Proxy",
			args: args{
				s:      config,
				suse:   suseManagerMock,
				logger: logger,
			},
			want: iproxy,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewProxy(tt.args.s, tt.args.suse, tt.args.logger, 5); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewProxy() = %v, want %v", got, tt.want)
			}
		})
	}
}
