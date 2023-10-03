package config

type SenderType int

const (
	MsGraph SenderType = iota
	Dummy
)

func (s SenderType) IsValid() bool {
	switch s {
	case MsGraph, Dummy:
		return true
	}
	return false
}

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
	Debug           bool             `mapstructure:"debug"`
	MetricsPort     uint             `mapstructure:"metrics-port"`
	SenderType      SenderType       `mapstructure:"sender-type"`
	MsGraphConfig   *MsGraphConfig   `mapstructure:"ms-graph-config"`
	BasicAuthConfig *BasicAuthConfig `mapstructure:"basic-auth-config"`
}
