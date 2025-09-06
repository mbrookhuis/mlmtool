package updatesusecversion

type SuseC struct { //nolint:maligned
	APIMgmt             bool
	APIVersionPod       string
	APIVersionMgmt      string
	K3sVersion          string
	KeepAlivedVersion   string
	RegistrationVersion string
	CheckFabricVersion  string
	SaltPodVersion      string
	OsPodVersion        string
	EcpSumaVer          string
	Server              string
	SkipBmc             bool
}
