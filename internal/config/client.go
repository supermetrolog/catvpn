package config

type Host struct {
	Ip   string `env:"IP"`
	Port uint16 `env:"PORT"`
}

type ClientConfig struct {
	ClientHost              Host `env-prefix:"CLIENT_"`
	ServerHost              Host `env-prefix:"CLIENT_SERVER_"`
	BufferSize              uint `env:"CLIENT_BUFFER_SIZE" env-default:"2000"`
	HeartBeatTimeInterval   uint `env:"CLIENT_HEART_BEAT_TIME_INTERVAL" env-default:"60"`
	MTU                     uint `env:"CLIENT_MTU" env-default:"1501"`
	ServerConnectionTimeout uint `env:"CLIENT_SERVER_CONNECTION_TIMEOUT" env-default:"5"`
}
