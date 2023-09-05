package config

type SenderType int

const (
	MsGraph SenderType = iota
	Dummy
)

type BasicAuthConfig struct {
	Enabled  bool
	Username string
	Password string
}

type MsGraphConfig struct {
	TenantId     string
	ClientId     string
	ClientSecret string
	// The object ID of the user in Azure AD. Will send using that's user email.
	SenderUserId string
}

type SkylineConfig struct {
	Hostname        string
	Port            uint
	MetricsPort     uint
	Debug           bool
	SenderType      SenderType
	MsGraphConfig   *MsGraphConfig
	BasicAuthConfig *BasicAuthConfig
}
