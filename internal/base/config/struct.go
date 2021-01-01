package config

type MainConfig struct {
	CWD    string       `toml:"cwd"`
	Server ServerConfig `toml:"server"`
	Log    LogConfig    `toml:"log"`
	MySQL  MySQLConfig  `toml:"mysql"`
	Redis  RedisConfig  `toml:"redis"`
}

type ServerConfig struct {
	Name   string `toml:"name"`
	Listen string `toml:"listen"`
	URL    string `toml:"url"`
	Prefix string `toml:"prefix"`
	PProf  string `toml:"pprof"`
}

type LogConfig struct {
	Level string `toml:"level"`
}

type MySQLConfig struct {
	Host           string `toml:"host"`
	Port           int    `toml:"port"`
	Database       string `toml:"database"`
	User           string `toml:"user"`
	Password       string `toml:"password"`
	Charset        string `toml:"charset"`
	Prefix         string `toml:"prefix"`
	Debug          bool   `toml:"debug"`
	MaxConnections int    `toml:"maxConnections"`
}

type RedisConfig struct {
	Addr     string `toml:"addr"`
	Password string `toml:"password"`
	DB       int    `toml:"db"`
	PoolSize int    `toml:"poolSize"`
	Prefix   string `toml:"prefix"`
}
