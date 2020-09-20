package configs

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

func getConfiguration() *ProjectCfg {
	config := &ProjectCfg{}
	if _, err := toml.DecodeFile("config.toml", config); err != nil {
		fmt.Println("\x1b[35m[\x1b[0m\x1b[31mERROR\x1b[0m\x1b[35m]\x1b[0m \x1b[91m>>>\x1b[0m \x1b[32m", err.Error(), "\x1b[0m")
		os.Exit(1)
	}
	return config
}
