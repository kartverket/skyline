package config

type SenderType int

const (
	MsGraph SenderType = iota
	Dummy
)

type BasicAuthConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type MsGraphConfig struct {
	TenantId     string `mapstructure:"tenant-id"`
	ClientId     string `mapstructure:"client-id"`
	ClientSecret string `mapstructure:"client-secret"`
	// The object ID of the user in Azure AD. Will send using that's user email.
	SenderUserId string `mapstructure:"sender-user-id"`
}

type SkylineConfig struct {
	Hostname        string           `mapstructure:"hostname"`
	Port            uint             `mapstructure:"port"`
	MetricsPort     uint             `mapstructure:"metrics-port"`
	Debug           bool             `mapstructure:"debug"`
	SenderType      SenderType       `mapstructure:"sender-type"`
	MsGraphConfig   *MsGraphConfig   `mapstructure:"ms-graph-config"`
	BasicAuthConfig *BasicAuthConfig `mapstructure:"basic-auth-config"`
}
