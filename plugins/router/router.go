package router

import "github.com/gin-gonic/gin"

type Http struct {
	// Router path
	Router string

	// Router response handle
	Handle func(*gin.Context)
}
