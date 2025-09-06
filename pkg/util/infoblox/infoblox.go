package infoblox

import (
	_ "encoding/base64"
	"errors"
	"fmt"
	"os"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"

	returncodes "mlmtool/pkg/util/returnCodes"
)

// API Constants
const (
	zoneAuthAPI    = "zone_auth"
	aRecordAPI     = "record:a"
	ptrRecordAPI   = "record:ptr"
	cnameRecordAPI = "record:cname"
	gridAPI        = "grid"
	restartAPI     = "?_function=restartservices"
)

const (
	DefaultDomain      = "a4.telekom.de"
	InfobloxConfigFile = "/srv/pillar/infoblox.sls"
	RetryCount         = 3
)

// DNSRecord data needed for dnsRecord
type DNSRecord struct {
	PodRecord bool
	PtrRecord bool
	IP        string
	DNSRecord string
	CName     []string
}

type InfobloxData struct {
	Infoblox struct {
		IbIpamview string `yaml:"ib_ipamview"`
		IbNsgroup  string `yaml:"ib_nsgroup"`
		IbPassword string `yaml:"ib_password"`
		IbServer   string `yaml:"ib_server"`
		IbUser     string `yaml:"ib_user"`
		IbView     string `yaml:"ib_view"`
		IbPort     string `yaml:"ib_port,omitempty"`
		IbVersion  string `yaml:"ib_version,omitempty"`
	} `yaml:"infoblox"`
}

// Infoblox - Infoblox Structure
type Infoblox struct {
	logger *zap.Logger
	View   string
	Client *Client
}

type DNSProvider interface {
	IBCheckAndCreateARecord(DNSRecord) error
	IBCheckAndCreateZoneAuth(string) error
	IBCheckAndCreatePTRRecord(DNSRecord) error
	IBCheckAndCreateCNameRecord(DNSRecord, string) error
	IBDeleteCNameRecordName(string) error
	IBDeleteCNameRecordCanonical(string) error
	IBDeleteARecord(DNSRecord) error
	IBDeletePTRRecord(DNSRecord) error
}

func New(logger *zap.Logger, configFile string) (DNSProvider, error) {
	slogger := logger.Sugar()
	if configFile == "" {
		slogger.Fatal("configFile: empty string")
	}
	yamlFile, err := os.ReadFile(configFile)
	if err != nil {
		slogger.Errorf("Unable to open or readfile %v", configFile)
		return nil, err
	}
	var result InfobloxData
	err = yaml.Unmarshal(yamlFile, &result)
	if err != nil {
		return nil, errors.New(returncodes.ErrFailedUnMarshalling)
	}
	var infoData Infoblox
	infoData.logger = logger
	infoData.View = result.Infoblox.IbView
	infoData.Client = Create(ProviderArgs{
		Host:     fmt.Sprintf("https://%s", result.Infoblox.IbServer),
		Port:     result.Infoblox.IbPort,
		Username: result.Infoblox.IbUser,
		Password: result.Infoblox.IbPassword,
		Version:  result.Infoblox.IbVersion,
		log:      logger,
	})
	return &infoData, nil
}
