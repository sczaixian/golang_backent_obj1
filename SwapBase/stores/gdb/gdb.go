
package gdb

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