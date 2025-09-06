// Package createhubform - description needed for creating hub
package createhubform

// HubFormYaml create hub description
type HubFormYaml struct {
	EcpSumaVer    string `yaml:"ecp_suma_ver"`
	EcpSusePodVer string `yaml:"ecp_suse_pod_ver"`
	EcpSuseCMVer  string `yaml:"ecp_suse_cm_ver"`
}
