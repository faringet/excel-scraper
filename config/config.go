package config

type Config struct {
	PlatformAPI  PlatformAPI  `mapstructure:"PLATFORM_API"`
	Scheduler    Scheduler    `mapstructure:"SCHEDULER"`
	Logger       Logger       `mapstructure:"LOGGER"`
	ScanningOpts ScanningOpts `mapstructure:"SCANNING_OPTS"`
	LocalURL     string       `mapstructure:"THIS_APP_URL"`
}

type PlatformAPI struct {
	URL string `mapstructure:"URL"`
}

type Scheduler struct {
	Update int `mapstructure:"UPDATE"`
}

type Logger struct {
	Production  string `mapstructure:"PRODUCTION"`
	Development string `mapstructure:"DEVELOPMENT"`
}

type ScanningOpts struct {
	Path   string `mapstructure:"PATH"`
	Format string `mapstructure:"FORMAT"`
}
