package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/reverie/configs"
	"github.com/reverie/utils"
	"golang.org/x/sync/errgroup"
)

func main() {
	var g errgroup.Group
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", configs.Project.Port),
		Handler:      newRouter(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	g.Go(server.ListenAndServe)
	utils.LogInfo("Main-1", "Server running on port %d", configs.Project.Port)
	if err := g.Wait(); err != nil {
		utils.LogError("Main-2", err)
		os.Exit(1)
	}
}
