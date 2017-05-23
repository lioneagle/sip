package config

type LocalServerConfig struct {
	transport string
	localAddr string
}

type LocalConnectionConfig struct {
	transport  string
	localAddr  string
	remoteAddr string
}

type TransportConfig struct {
	localServers
}
