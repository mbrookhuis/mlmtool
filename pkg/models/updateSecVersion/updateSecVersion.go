package updatesecversion

type SuseManagerSecondary struct {
	EcpSumaVer    string
	EcpSusePodVer string
	SecGroup      string
}

/*
Uyunihub formular
*/

type HubForm struct {
	Hub Hub `json:"hub"`
}

type Hub struct {
	EcpSuseCmVer   string        `json:"ecp_suse_cm_ver,omitempty"`
	ServerPassword string        `json:"server_password"`
	ServerUsername string        `json:"server_username"`
	EcpSusePodVer  string        `json:"ecp_suse_pod_ver"`
	EcpSumaVer     string        `json:"ecp_suma_ver"`
	HubOrg         string        `json:"hub_org"`
	MasterRsa      string        `json:"master_rsa"`
	ConfigAll      []ConfigAll   `json:"configchannel,omitempty"`
	ChannelsAll    []ChannelsAll `json:"channels_all,omitempty"`
	ProjectsAll    []ProjectsAll `json:"projects_all,omitempty"`
	Slave          []Slave       `json:"slave,omitempty"`
}

type ConfigAll struct {
	ConfigChannel string `json:"configchannel"`
}

type ChannelsAll struct {
	BaseChannel string `json:"basechannel"`
}

type ProjectsAll struct {
	Project string `json:"project"`
}

type Slave struct {
	Slave    string        `json:"slave"`
	Config   []ConfigAll   `json:"config,omitempty"`
	Channels []ChannelsAll `json:"channels,omitempty"`
	Projects []ProjectsAll `json:"projects,omitempty"`
}
