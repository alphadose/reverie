package main

import (
	"fmt"

	"github.com/reverie/configs"
	"github.com/reverie/utils"
)

func main() {
	utils.LogInfo("Main-1", "Server running on port %d", configs.Project.Port)
	newRouter().Listen(fmt.Sprintf(":%d", configs.Project.Port))
}
