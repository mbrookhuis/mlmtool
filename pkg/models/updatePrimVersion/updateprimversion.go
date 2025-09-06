// Package updateprimversion models used for this function.
package updateprimversion

// SuseManagerPrimaryFormular - info for SUSE Manager Primary
type SuseManagerPrimaryFormular struct {
	SuseManagerPrimary SuseManagerPrimary `json:"suse_manager_primary"`
}

// SuseManagerPrimary - info for SUSE Manager Primary
type SuseManagerPrimary struct {
	SmasVer       string `json:"smas_ver"`
	EcpSumaVer    string `json:"ecp_suma_ver"`
	EcpSuseCmVer  string `json:"ecp_suse_cm_ver"`
	EcpSusePodVer string `json:"ecp_suse_pod_ver"`
}
