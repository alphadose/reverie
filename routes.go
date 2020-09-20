package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func newRouter() http.Handler {
	router := gin.Default()
	return router
}
