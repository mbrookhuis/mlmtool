package cmdexecutor

import (
	"os/exec"
	"testing"

	"go.uber.org/zap"

	logging "mlmtool/pkg/util/logger"
)

func TestNewCMDExecutor(t *testing.T) {
	type args struct {
		logger *zap.Logger
	}
	logger := logging.NewTestingLogger(t.Name())
	wantCMDExecutor := NewCMDExecutor(logger)
	tests := []struct {
		name string
		args args
		want ICMDExecutor
	}{
		{
			name: "New OS Use Case",
			args: args{
				logger: logger,
			},
			want: wantCMDExecutor,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewCMDExecutor(tt.args.logger)
			if got == nil {
				t.Errorf("NewOS() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_os_ExecuteCommand(t *testing.T) {
	type fields struct {
		logger      *zap.Logger
		execCommand func(name string, arg ...string) *exec.Cmd
	}
	type args struct {
		name string
		args []string
	}

	execCMDSuccess := func(name string, arg ...string) *exec.Cmd {
		return exec.Command("hostname")
	}

	execCMDErr := func(name string, arg ...string) *exec.Cmd {
		return exec.Command("commandNotPresent")
	}

	execCMDStdErr := func(name string, arg ...string) *exec.Cmd {
		return exec.Command("/usr/bin/hostname", "-rr")
	}

	logger := logging.NewTestingLogger(t.Name())
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Execute Command",
			fields: fields{
				logger:      logger,
				execCommand: execCMDSuccess,
			},
			args: args{
				name: "test binary",
				args: []string{"Test", "Args"},
			},
			//want:    []string{"Test Args", ""},
			wantErr: false,
		},
		{
			name: "Execute Command Negative",
			fields: fields{
				logger:      logger,
				execCommand: execCMDErr,
			},
			args: args{
				name: "test binary",
				args: []string{"Test", "Args"},
			},
			wantErr: true,
		},
		{
			name: "Execute Command StdErr Negative",
			fields: fields{
				logger:      logger,
				execCommand: execCMDStdErr,
			},
			args: args{
				name: "test binary",
				args: []string{"Test", "Args"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os := &cmdUtils{
				logger:      tt.fields.logger,
				execCommand: tt.fields.execCommand,
			}
			_, err := os.ExecuteCommand(tt.args.name, tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("os.ExecuteCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_cmdUtils_CreateDirectory(t *testing.T) {
	type fields struct {
		logger      *zap.Logger
		execCommand func(name string, arg ...string) *exec.Cmd
	}
	type args struct {
		path string
	}
	execCMDSuccess := func(name string, arg ...string) *exec.Cmd {
		return exec.Command("hostname")
	}

	execCMDErr := func(name string, arg ...string) *exec.Cmd {
		return exec.Command("commandNotPresent")
	}

	logger := logging.NewTestingLogger(t.Name())

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Create Dir",
			fields: fields{
				logger:      logger,
				execCommand: execCMDSuccess,
			},
			args: args{
				path: "/tmp/present",
			},
			wantErr: false,
		},
		{
			name: "Execute Command Negative",
			fields: fields{
				logger:      logger,
				execCommand: execCMDErr,
			},
			args: args{
				path: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmdUtils := &cmdUtils{
				logger:      tt.fields.logger,
				execCommand: tt.fields.execCommand,
			}
			if err := cmdUtils.CreateDirectory(tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("CreateDirectory() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
