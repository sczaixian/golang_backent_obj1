package config


import (
	"fmt"
	"strings"

	logging "github.com/ProjectsTask/SwapBase/logger"
	"github.com/ProjectsTask/SwapBase/stores/gdb"

	"github.com/spf13/viper"
)

type Api struct {
	Port string `toml:"port" json:"port"`
	MaxNum int64 `toml:"max_num" json:"max_num"`
}

type ProjectCfg struct {
	Name string `toml:"name" mapstructure:"name" json:"name"`
}

type Config struct{
	Api `toml:"api" json:"api"`
	ProjectCfg * ProjectCfg `toml:"project_cfg" mapstructure:"project_cfg" json:"project_cfg"`
	Log logging.LogConf `toml:"log" json:"log"`
	DB gdb.Config 		`toml:"db" json:"db"`


}

func UnmarshalConfig(configFilePath strng) (*Config, error){

}