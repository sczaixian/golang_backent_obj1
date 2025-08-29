package main

import (
	"flag"

	"github.com/ProjectsTask/Backend/config"
	"github.com/ProjectsTask/Backend/service/svc"
	"github.com/ProjectsTask/Backend/api/router"
	"github.com/ProjectsTask/Backend/app"
)

const (
	defaultConfigPath = ""
)

func main() {
	conf := flag.String("conf", defaultConfigPath, "conf file path")
	flag.Parse()
	c, err := config.UnmarshalConfig(*conf)
	if err != nil {
		panic(err)
	}

	for _, chain := range c.ChainSupported {
		if chain.ChainID == 0 || chain.Name == "" {
			panic("invalid chain_suffix config")
		}
	}

	serverCtx, err := svc.NewServiceContext(c)
	if err != nil {
		panic(err)
	}

	r := router.NewRouter(serverCtx)

	app, err := app.NewPlatform(c, r, serverCtx)

	if err != nil {
		panic(err)
	}
	app.Start()
}
