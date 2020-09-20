package configs

// ProjectCfg is the configuration for the entire project
type ProjectCfg struct {
	Debug bool `toml:"debug"`
	Port  int  `toml:"port"`
}
