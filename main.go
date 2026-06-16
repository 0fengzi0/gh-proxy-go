package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg = loadConfig()

	fetchAssets()

	if cfg.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	r.GET("/*path", Handler)
	r.POST("/*path", Handler)

	addr := cfg.Host + ":" + cfg.Port
	log.Printf("gh-proxy-go starting on %s", addr)
	r.Run(addr)
}
