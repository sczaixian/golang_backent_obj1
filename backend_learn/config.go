package backendlearn

/*
[project_cfg]
name = "EasySwap"
*/

type ProjectCfg struct {
	Name string `toml:"name" mapstructure:"name" json:"name"`
}

/*
[api]
port = ":80"
max_num = 500
*/
type Api struct {
	Port   string `toml:"port" json:"port"`
	MaxNum int64  `toml:"max_num" json:"max_num"`
}

/*
[log]
compress = false
leep_days = 7
level = "info"
mode = "console"
path = "logs/v1-backend"
service_name = "v1-backend"
*/
type LogConf struct {
	ServiceName string `toml:"service_name" mapstructure:"service_name" json:"service_name"`
	Mode        string `toml:"mode" json:"mode"`
	Path        string `toml:"path" json:"path"`
	Level       string `toml:"level" json:"level"`
	Compress    bool   `toml:"compress" json:"compress"`
	KeepDays    int    `toml:"keep_days" mapstructure:"keep_days" json:"keep_days"`
}

/*
[db]
database = "easyswap"
password = "123"
port = 3306 #4000
max_open_conns = 1500
host = "127.0.0.1"
log_level = "info"
max_conn_max_lifetime = 300
user = "root"
max_idle_conns = 10
*/
type Config struct {
	User               string `toml:"user" json:"user"`                                                                        // 用户
	Password           string `toml:"password" json:"password"`                                                                // 密码
	Host               string `toml:"host" json:"host"`                                                                        // 地址
	Port               int    `toml:"port" json:"port"`                                                                        // 端口
	Database           string `toml:"database" json:"database"`                                                                // 数据库
	MaxIdleConns       int    `toml:"max_idle_conns" mapstructure:"max_idle_conns" json:"max_idle_conns"`                      // 最大空闲连接数
	MaxOpenConns       int    `toml:"max_open_conns" mapstructure:"max_open_conns" json:"max_open_conns"`                      // 最大打开连接数
	MaxConnMaxLifetime int64  `toml:"max_conn_max_lifetime" mapstructure:"max_conn_max_lifetime" json:"max_conn_max_lifetime"` // 连接复用时间
	LogLevel           string `toml:"log_level" mapstructure:"log_level" json:"log_level"`                                     // 日志级别，枚举（info、warn、error和silent）
}

/*
[[kv.redis]]
pass = ""
host = "127.0.0.1:6379"
type = "node"
*/
type KvConf struct {
	Redis []*Redis `toml:"redis" mapstructure:"redis" json:"redis"`
}

type Redis struct {
	MasterName string `toml:"master_name" mapstructure:"master_name" json:"master_name"`
	Host       string `toml:"host" json:"host"`
	Type       string `toml:"type" json:"type"`
	Pass       string `toml:"pass" json:"pass"`
}


/*
[metadata_parse]
name_tags = ["name", "title"]
image_tags = ["image", "image_url", "animation_url", "media_url", "image_data", "imageUrl"]
attributes_tags = ["attributes", "properties", "attribute"]
trait_name_tags = ["trait_type"]
trait_value_tags = ["value"]
*/
type MetadataParse struct {
	NameTags       []string `toml:"name_tags" mapstructure:"name_tags" json:"name_tags"`
	ImageTags      []string `toml:"image_tags" mapstructure:"image_tags" json:"image_tags"`
	AttributesTags []string `toml:"attributes_tags" mapstructure:"attributes_tags" json:"attributes_tags"`
	TraitNameTags  []string `toml:"trait_name_tags" mapstructure:"trait_name_tags" json:"trait_name_tags"`
	TraitValueTags []string `toml:"trait_value_tags" mapstructure:"trait_value_tags" json:"trait_value_tags"`
}

/*
[[chain_supported]]
name="sepolia"
chain_id=11155111
endpoint = "https://eth-sepolia.g.alchemy.com/v2/RmsPYhly5O6-XH8UdmqCQ"
*/
type ChainSupported struct {
	Name     string `toml:"name" mapstructure:"name" json:"name"`
	ChainID  int    `toml:"chain_id" mapstructure:"chain_id" json:"chain_id"`
	Endpoint string `toml:"endpoint" mapstructure:"endpoint" json:"endpoint"`
}

type Config struct {
	Api        `toml:"api" json:"api"`
	ProjectCfg *ProjectCfg     `toml:"project_cfg" mapstructure:"project_cfg" json:"project_cfg"`
	Log        logging.LogConf `toml:"log" json:"log"`
	//ImageCfg       *image.Config     `toml:"image_cfg" mapstructure:"image_cfg" json:"image_cfg"`
	DB             gdb.Config        `toml:"db" json:"db"`
	Kv             *KvConf           `toml:"kv" json:"kv"`


	Evm            *erc.NftErc       `toml:"evm" json:"evm"`  <---- TODO:


	MetadataParse  *MetadataParse    `toml:"metadata_parse" mapstructure:"metadata_parse" json:"metadata_parse"`
	ChainSupported []*ChainSupported `toml:"chain_supported" mapstructure:"chain_supported" json:"chain_supported"`
}