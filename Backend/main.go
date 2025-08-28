package main

import (
	"flag"
)

const (
	defaultConfigPath=""
)

func main(){
	conf := flag.String("conf", defaultConfigPath, "conf file path")
	flag.Parse()
	c, err := config.UnMarshalConfig(*conf)
	if err != nil {
		painc()
	}

	
}