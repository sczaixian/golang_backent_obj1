package app

import (
	"context"

	"go.uber.org/zap"
	
	"github.com/gin-gonic/gin"
	"github.com/ProjectsTask/Backend/config"
	"github.com/ProjectsTask/Backend/service/svc"

	"github.com/ProjectsTask/SwapBase/logger/xzap"

)

type Platform struct {
	config    *config.Config
	router    *gin.Engine
	serverCtx *svc.ServerCtx
}

func (p *Platform) Start() {
	xzap.WithContext(context.Background()).Info("EasySwap-End run", zap.String("port", p.config.Api.Port))
	if err := p.router.Run(p.config.Api.Port); err != nil {
		panic(err)
	}
}



func NewPlatform(config *config.Config, router *gin.Engine, serverCtx *svc.ServerCtx) (*Platform, error) {
	return &Platform{
		config:    config,
		router:    router,
		serverCtx: serverCtx,
	}, nil

}