package config

type StorageType string

//nolint:gochecknoglobals // it is a config package
var GlobalConfig *Config

const (
	StorageTypeRedis StorageType = "redis"
	StorageTypeLocal StorageType = "local"
)

type Config struct {
	StorageType StorageType
	Redis       struct {
		Host     string
		Port     int
		Password string
		DB       int
	}
	Input struct {
		File struct {
			Path string
		}
		URL struct {
			URL string
		}
	}
	WorkerPoolSize   int
	WordsChannelSize int
}
