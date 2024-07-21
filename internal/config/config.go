package config

type Host struct {
	Ip   string `env:"IP"`
	Port uint16 `env:"PORT"`
}

type ClientConfig struct {
	ClientHost              Host `env-prefix:"CLIENT_"`
	ServerHost              Host `env-prefix:"SERVER_"`
	BufferSize              uint `env:"BUFFER_SIZE" env-default:"2000"`
	HeartBeatTimeInterval   uint `env:"CLIENT_HEART_BEAT_TIME_INTERVAL" env-default:"60"`
	MTU                     uint `env:"MTU" env-default:"1500"`
	ServerConnectionTimeout uint `env:"CONNECTION_TIMEOUT" env-default:"5"`
}

type ServerConfig struct {
	ServerHost            Host   `env-prefix:"SERVER_"`
	Subnet                string `env:"SERVER_SUBNET"`
	BufferSize            uint   `env:"BUFFER_SIZE" env-default:"2000"`
	HeartBeatTimeInterval uint   `env:"HEART_BEAT_TIME_INTERVAL" env-default:"60"`
	MTU                   uint   `env:"MTU" env-default:"1500"`
}
