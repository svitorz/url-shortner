package config

type RedisConfig struct {
	RedisHost     string
	RedisPort     int64
	RedisPassword string
	RedisDb       int64
	RedisProtocol int64
}
