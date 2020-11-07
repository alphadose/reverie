package configs

import "time"

// Mongo is the configuration for mongodb storage
type Mongo struct {
	URL string `toml:"url"`
}

// Crypto is the configuration for encryption/decryption env variables
type Crypto struct {
	Key   string `toml:"key"`
	Nonce string `toml:"nonce"`
}

// Admin is the configuration for the default Admin
type Admin struct {
	Email    string `toml:"email"`
	Password string `toml:"password"`
	Username string `toml:"username"`
}

// JWT is the configuration for auth token
type JWT struct {
	Timeout    time.Duration `toml:"timeout"`
	MaxRefresh time.Duration `toml:"max_refresh"`
	Secret     string        `toml:"secret"`
}

// ProjectCfg is the configuration for the entire project
type ProjectCfg struct {
	Debug  bool   `toml:"debug"`
	Port   int    `toml:"port"`
	Crypto Crypto `toml:"crypto"`
	Admin  Admin  `toml:"admin"`
	Mongo  Mongo  `toml:"mongo"`
	JWT    JWT    `toml:"jwt"`
}
