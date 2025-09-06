package updatepod

type PODversion struct { //nolint:maligned
	NegName            string
	CheckFabricVersion string
	SaltPodVersion     string
	SkipSync           bool
}

// DTAGK3SConfig stores K3S config data
type DTAGK3SConfig struct {
	K3SConfig K3SConfig `json:"k3sconfig"`
}

// K3SConfig k3s configuration
type K3SConfig struct {
	// Servers server list
	Servers []K3SServer `json:"k3s_server"`
	// Location neg name
	Location string `json:"k3s_location"`
	// PrimaryIP management IP
	PrimaryIP string `json:"k3s_primary_ip"`
	// Token Network Element ID
	Token string `json:"k3s_token"`
	// VIP virtual IP
	VIP string `json:"k3s_vip"`
	// Type server type
	Version string `json:"k3s_version"`
	// RancherVersion rancher version
	SaltVersion string `json:"salt_pod_version"`
	// OSRelease os release version
	OSRelease string `json:"os_pod_version"`
	// KeepAlivedVersion keep alived version
	KeepAlivedVersion string `json:"keep_alived_version"`
	// Proxy rancher k3s proxy
	Proxy string `json:"k3s_proxy"`
	// Rancher rancher virtual host
	Rancher string `json:"k3s_rancher"`
	// API rancher token
	API string `json:"k3s_api"`
	// Rancher Cluster Registration Image Tag
	ClusterRegistrationImageTag string `json:"registration_version"`
	// Version of Check Fabric Binary
	CheckFabricVersion string `json:"checkfabric_version"`
}

// K3SServer k3s server details
type K3SServer struct {
	// Name hostname of server
	Name string `json:"pod_name"`
	// Role either primary or secondary
	Role string `json:"pod_role"`
	// IP server IP
	IP string `json:"pod_ip_k3s"`
}
