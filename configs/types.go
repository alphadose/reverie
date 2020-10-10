package configs

// Mongo is the configuration for mongodb storage
type Mongo struct {
	URL string `toml:"url"`
}

// Admin is the configuration for the default Admin
type Admin struct {
	Email    string `toml:"email"`
	Password string `toml:"password"`
	Username string `toml:"username"`
}

// ProjectCfg is the configuration for the entire project
type ProjectCfg struct {
	Debug bool  `toml:"debug"`
	Port  int   `toml:"port"`
	Admin Admin `toml:"admin"`
	Mongo Mongo `toml:"mongo"`
}
