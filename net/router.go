package net

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Router interface
type Router interface {
	Run(port string)
}

type router struct {
	router *gin.Engine
}

// NewRouter は、Routerを作成して返す。
func NewRouter() Router {
	r := &router{
		router: gin.Default(),
	}
	SetHandler(r.router)
	return r
}

func (r *router) Run(port string) {
	if len(port) == 0 {
		port = "80"
	}
	log.Info("start listen: ", port)

	_, err := strconv.Atoi(port)
	if err != nil {
		log.Errorf("dont number of port: %s", port)
		return
	}
	r.router.Run(fmt.Sprintf(":%s", port))
}
