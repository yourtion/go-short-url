package config

// 动态配置信息
type DynamicInfo struct {
	// 白名单列表
	WhiteList []string
	// 目标域名列表
	Domains []string
}

// 主要配置
type MainConfig struct {
	// 当前运行目录
	CWD string `toml:"cwd"`
	// 服务器配置
	Server ServerConfig `toml:"server"`
	// 日志配置
	Log LogConfig `toml:"log"`
	// MySQL 配置
	MySQL MySQLConfig `toml:"mysql"`
	// Redis 配置
	Redis RedisConfig `toml:"redis"`
}

// 服务器配置
type ServerConfig struct {
	// 服务名
	Name string `toml:"name"`
	// 监听配置
	Listen string `toml:"listen"`
	// 访问地址
	URL string `toml:"url"`
	// 是否启动 pprof （配置监听路径）
	PProf string `toml:"pprof"`
}

// 日志配置
type LogConfig struct {
	// 日志级别
	Level string `toml:"level"`
}

// 数据库配置
type MySQLConfig struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Database string `toml:"database"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	Charset  string `toml:"charset"`
	// 表前缀
	Prefix string `toml:"prefix"`
	// 是否打开调试
	Debug          bool `toml:"debug"`
	MaxConnections int  `toml:"maxConnections"`
}

// Redis 配置
type RedisConfig struct {
	Addr     string `toml:"addr"`
	Password string `toml:"password"`
	DB       int    `toml:"db"`
	PoolSize int    `toml:"poolSize"`
	// Key 前缀
	Prefix string `toml:"prefix"`
}
