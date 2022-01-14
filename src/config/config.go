package config

type Configuration struct {
	Unifi  UnifiConfiguration
	Backup BackupConfiguration
}

type UnifiConfiguration struct {
	Server            string
	Username          string
	Password          string
	ControllerVersion string `mapstructure:"controller_version"`
}

type BackupConfiguration struct {
	OutputDirectory string `mapstructure:"output_directory"`
	Keep            int
}
