package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func newRouter() http.Handler {
	router := gin.Default()
	return router
}
